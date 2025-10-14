package service

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/anoideaopen/cartridge"
	"github.com/anoideaopen/cartridge/manager"
	"github.com/anoideaopen/foundation/keys"
	"github.com/anoideaopen/foundation/proto"
	"github.com/anoideaopen/testnet-cli/logger"
	pb "github.com/golang/protobuf/proto" //nolint:staticcheck
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	contextApi "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	core2 "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"go.uber.org/zap"
)

const (
	BatchExecuteEvent          = "batchExecute"
	defaultTimeoutEventSeconds = 30
)

type HLFClient struct {
	sdk                     *fabsdk.FabricSDK
	channels                map[string]*ChannelConnection
	NotifierChaincodeEvents map[string]<-chan *fab.CCEvent
	afterInvokeHandler      HlfAfterInvokeHandler
	beforeInvokeHandler     HlfBeforeInvokeHandler
	ContextOptions          []fabsdk.ContextOption
}

func (hlf *HLFClient) AddConnectionConfig(connectionConfigPath string) (core2.ConfigProvider, []core2.ConfigBackend, error) {
	configProvider := config.FromFile(connectionConfigPath)
	configBackends, err := configProvider()
	return configProvider, configBackends, err
}

func (hlf *HLFClient) AddFabsdk(configProvider core2.ConfigProvider, opts ...fabsdk.Option) error {
	sdk, err := fabsdk.New(configProvider, opts...)
	if err != nil {
		msg := "failed to create new SDK"
		logger.Error(msg, zap.Error(err))
		return fmt.Errorf("%s: %w", msg, err)
	}

	hlf.sdk = sdk
	return nil
}

func NewHLFClientFile(connectionConfigPath string, username string, organization string) (*HLFClient, error) {
	hlfClient := newHLFClient()

	configProvider, _, err := hlfClient.AddConnectionConfig(connectionConfigPath)
	if err != nil {
		logger.Error("Failed to AddConnectionConfig", zap.Error(err))
		return nil, err
	}

	err = hlfClient.AddFabsdk(configProvider)
	if err != nil {
		logger.Error("Failed to create new channel client", zap.Error(err))
		return nil, err
	}
	if hlfClient == nil {
		logger.Error("Failed to create new channel client")
		return nil, err
	}

	hlfClient.ContextOptions = append(hlfClient.ContextOptions, fabsdk.WithUser(username))
	hlfClient.ContextOptions = append(hlfClient.ContextOptions, fabsdk.WithOrg(organization))

	return hlfClient, nil
}

type VaultConfig struct {
	Token string `mapstructure:"token,omitempty"`
	Path  string `mapstructure:"path,omitempty"`
	// CAcert       string `mapstructure:"ca_cert,omitempty"`
	Address      string `mapstructure:"address,omitempty"`
	UserCertName string `mapstructure:"user_cert_name,omitempty"`
	UserOrgMspID string `mapstructure:"user_org_msp_id,omitempty"`
}

func NewHLFClient(connectionConfigPath string, username string, organization string, vaultConfig *VaultConfig) (*HLFClient, error) {
	var hlfClient *HLFClient
	var err error
	if vaultConfig != nil && len(vaultConfig.Address) != 0 && len(vaultConfig.Token) != 0 && len(vaultConfig.Path) != 0 &&
		len(vaultConfig.UserCertName) != 0 && len(vaultConfig.UserOrgMspID) != 0 {
		hlfClient, err = NewHLFClientVault(connectionConfigPath, organization, vaultConfig)
		if err != nil {
			return nil, err
		}
	} else {
		hlfClient, err = NewHLFClientFile(connectionConfigPath, username, organization)
		if err != nil {
			return nil, err
		}
	}
	return hlfClient, nil
}

func NewHLFClientVault(connectionConfigPath string, organization string, vaultConfig *VaultConfig) (*HLFClient, error) {
	hlfClient := newHLFClient()

	configProvider, backends, err := hlfClient.AddConnectionConfig(connectionConfigPath)
	if err != nil {
		logger.Error("Failed to AddConnectionConfig", zap.Error(err))
		return nil, err
	}

	vaultManager, err := manager.NewVaultManager(
		vaultConfig.UserOrgMspID,
		vaultConfig.UserCertName,
		vaultConfig.Address,
		vaultConfig.Token,
		vaultConfig.Path,
	)
	if err != nil {
		logger.Error("Failed to create new vault manager", zap.Error(err))
		return nil, err
	}

	var connectOpts []fabsdk.Option

	connector := cartridge.NewConnector(vaultManager, cartridge.NewVaultConnectProvider(backends...))
	if connector != nil {
		connectOpts, err = connector.Opts()
		if err != nil {
			logger.Error("Failed to get connector Opts", zap.Error(err))
			return nil, err
		}
	}

	err = hlfClient.AddFabsdk(configProvider, connectOpts...)
	if err != nil {
		logger.Error("Failed to create new channel client", zap.Error(err))
		return nil, err
	}
	if hlfClient == nil {
		logger.Error("Failed to create new channel client")
		return nil, err
	}

	hlfClient.ContextOptions = append(hlfClient.ContextOptions, fabsdk.WithIdentity(vaultManager.SigningIdentity()))
	hlfClient.ContextOptions = append(hlfClient.ContextOptions, fabsdk.WithOrg(organization))

	return hlfClient, nil
}

func newHLFClient() *HLFClient {
	hlf := &HLFClient{
		channels:                make(map[string]*ChannelConnection),
		NotifierChaincodeEvents: map[string]<-chan *fab.CCEvent{},
	}

	return hlf
}

func (hlf *HLFClient) AddChannel(channelID string) error {
	channelConnection, err := hlf.GetChannelConnection(channelID)
	if err != nil {
		return fmt.Errorf("failed to GetChannelConnection: %w", err)
	}
	if channelConnection == nil {
		return errors.New("channelConnection can't be nil")
	}
	return nil
}

func (hlf *HLFClient) AddBeforeInvokeHandler(beforeInvokeHandler HlfBeforeInvokeHandler) {
	hlf.beforeInvokeHandler = beforeInvokeHandler
}

func (hlf *HLFClient) AddAfterInvokeHandler(afterInvokeHandler HlfAfterInvokeHandler) {
	hlf.afterInvokeHandler = afterInvokeHandler
}

func (hlf *HLFClient) createChannelClient(channelID string, options ...fabsdk.ContextOption) (*channel.Client, contextApi.ChannelProvider, error) {
	// prepare channel client context using client context
	clientChannelContext := hlf.sdk.ChannelContext(channelID, options...)
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, nil, err
	}

	return client, clientChannelContext, nil
}

type ChannelConnection struct {
	ChannelClient   *channel.Client
	ChannelProvider contextApi.ChannelProvider
}

func (hlf *HLFClient) GetChannelConnection(channelID string) (*ChannelConnection, error) {
	var result *ChannelConnection

	val, isExists := hlf.channels[channelID]
	if !isExists {
		channelClient, channelProvider, err := hlf.createChannelClient(channelID, hlf.ContextOptions...)
		if err != nil {
			return nil, err
		}
		result = &ChannelConnection{
			ChannelClient:   channelClient,
			ChannelProvider: channelProvider,
		}
		hlf.channels[channelID] = result
	} else {
		result = val
	}

	return result, nil
}

func (hlf *HLFClient) Query(channelID string, chaincodeName string, methodName string, methodArgs []string, options ...channel.RequestOption) (*channel.Response, error) {
	channelConnection, err := hlf.GetChannelConnection(channelID)
	if err != nil {
		return nil, err
	}

	return hlf.RequestChaincode(chaincodeName, methodName, methodArgs,
		channelConnection.ChannelClient.Query, options)
}

// InvokeWithSecretKey - method to sign arguments and send invoke request to hlf
// methodArgs []string -
// secretKey string - private key ed25519 - in base58check, or hex or base58
// chaincodeName string - chaincode name for invoke
// methodName string - chaincode method name for invoke
// noBatch bool - if wait batchTransaction set 'true'
// peers string - peer0.testnet
func (hlf *HLFClient) InvokeWithSecretKey(waitBatch bool, channelID string, chaincodeName string, methodName string, methodArgs []string, secretKey string, keyType proto.KeyType, requestOptions ...channel.RequestOption) (*channel.Response, error) {
	if len(secretKey) != 0 {
		k, err := GetKeys(secretKey, keyType)
		if err != nil {
			logger.Error("failed getPrivateKey", zap.Error(err))
			return nil, err
		}
		methodArgs, err = hlf.SignArgs(channelID, chaincodeName, methodName, methodArgs, k)
		if err != nil {
			logger.Error("failed signArgs", zap.Error(err))
			return nil, err
		}
	}

	return hlf.Invoke(waitBatch, channelID, chaincodeName, methodName, methodArgs, requestOptions...)
}

// InvokeWithPublicAndPrivateKey - method to sign arguments and send invoke request to hlf
// privateKey string - private key in ed25519
// publicKey string - private key in ed25519
// channelID string - channel name for invoke
// chaincodeName string - chaincode name for invoke
// methodName string - chaincode method name for invoke
// methodArgs []string -
// noBatch bool - if wait batchTransaction set 'true'
// peers string - peer0.testnet
func (hlf *HLFClient) InvokeWithPublicAndPrivateKey(waitBatch bool, k *keys.Keys, channelID string, chaincodeName string, methodName string, methodArgs []string, requestOptions ...channel.RequestOption) (*channel.Response, error) {
	methodArgs, err := hlf.SignArgs(channelID, chaincodeName, methodName, methodArgs, k)
	if err != nil {
		logger.Error("failed signArgs", zap.Error(err))
		return nil, err
	}

	return hlf.Invoke(waitBatch, channelID, chaincodeName, methodName, methodArgs, requestOptions...)
}

// Invoke - method to sign arguments and send invoke request to hlf
// channelID string - channel name for invoke
// chaincodeName string - chaincode name for invoke
// methodName string - chaincode method name for invoke
// methodArgs []string -
// peers string - target peer for invoke, if empty use default peer count by policy
func (hlf *HLFClient) Invoke(waitBatch bool, channelID string, chaincodeName string, methodName string, methodArgs []string, options ...channel.RequestOption) (*channel.Response, error) {
	channelConnection, err := hlf.GetChannelConnection(channelID)
	if err != nil {
		return nil, err
	}

	var beforeInvokeData interface{}
	if hlf.beforeInvokeHandler != nil {
		beforeInvokeData, err = hlf.beforeInvokeHandler(channelID, chaincodeName, methodName, methodArgs)
		if err != nil {
			return nil, err
		}
	}

	var notifier <-chan *fab.CCEvent

	if waitBatch {
		notifier, err = hlf.GetCCEventNotifier(channelConnection.ChannelClient, chaincodeName, BatchExecuteEvent)
		if err != nil {
			logger.Error("failed RegisterChaincodeEvent", zap.Error(err))
			return nil, err
		}
	}

	response, err := hlf.RequestChaincode(
		chaincodeName,
		methodName,
		methodArgs,
		channelConnection.ChannelClient.Execute,
		options,
	)

	if waitBatch { //nolint:nestif
		// start event start search batch event
		select {
		case ccEvent := <-notifier:
			logger.GetLogger().Info("- found cc event:")
			logger.GetLogger().Info("ccEvent.TxID: %s " + ccEvent.TxID)
			// logger.GetLogger().Info("ccEvent.BlockNumber: %s " + s(ccEvent.BlockNumber))
			logger.GetLogger().Info("ccEvent.ChaincodeID: %s " + ccEvent.ChaincodeID)
			logger.GetLogger().Info("ccEvent.EventName: %s " + ccEvent.EventName)
			logger.GetLogger().Info("ccEvent.SourceURL: %s " + ccEvent.SourceURL)
			logger.GetLogger().Error("ccEvent.EventName", zap.String("ccEvent.EventName", ccEvent.EventName))
			logger.Debug("ccEvent.EventName", zap.String("ccEvent.EventName", ccEvent.EventName))

			if ccEvent.EventName == BatchExecuteEvent && ccEvent.ChaincodeID == channelID {
				logger.Debug("ccEvent.Payload", zap.ByteString("ccEvent.Payload", ccEvent.Payload))
				batchEvent := &proto.BatchEvent{}
				if err = pb.Unmarshal(ccEvent.Payload, batchEvent); err != nil {
					logger.Error("err", zap.Error(err))
					return nil, err
				}

				for _, event := range batchEvent.GetEvents() {
					if event.GetError() != nil {
						logger.Error("err",
							zap.String("event.Id", hex.EncodeToString(event.GetId())),
							zap.String("event.Error.Error", event.GetError().GetError()),
							zap.Int32("event.Error.Code", event.GetError().GetCode()),
							zap.Error(err),
						)

						continue
					}
				}
			}

			logger.Debug("payload", zap.ByteString("payload", ccEvent.Payload))
		case <-time.After(time.Second * defaultTimeoutEventSeconds):
			logger.Debug(fmt.Sprintf("Did NOT receive CC for eventId(%s)\n", BatchExecuteEvent))
		}
	}

	if hlf.afterInvokeHandler != nil {
		err = hlf.afterInvokeHandler(beforeInvokeData, channelID, chaincodeName, methodName, methodArgs, response, err)
		if err != nil {
			return nil, err
		}
	}

	return response, err
}

type HlfAfterInvokeHandler func(beforeInvokeData interface{}, channelID string, chaincodeName string, methodName string, methodArgs []string, response *channel.Response, err error) error

type HlfBeforeInvokeHandler func(channelID string, chaincodeName string, methodName string, methodArgs []string) (interface{}, error)

func (hlf *HLFClient) RequestChaincode(
	chaincodeName string, methodName string, methodArgs []string,
	requestFunc func(channel.Request, ...channel.RequestOption) (channel.Response, error),
	options []channel.RequestOption,
) (*channel.Response, error) {
	logger.Debug("chaincodeName")
	logger.Debug(fmt.Sprintf("%v\n", chaincodeName))
	logger.Debug("methodName")
	logger.Debug(fmt.Sprintf("%v\n", methodName))
	logger.Debug("methodArgs", zap.Strings("methodArgs", methodArgs))

	channelRequest := channel.Request{
		ChaincodeID: chaincodeName,
		Fcn:         methodName,
		Args:        AsBytes(methodArgs),
	}

	logger.Debug("channelRequest",
		zap.String("ChaincodeID", channelRequest.ChaincodeID),
		zap.String("Fcn", channelRequest.Fcn),
		zap.Strings("Args", methodArgs))

	response, err := requestFunc(
		channelRequest,
		options...,
	)
	if err != nil {
		logger.Error("error", zap.Error(err))
		return nil, err
	}

	logger.Debug("payload", zap.ByteString("payload", response.Payload))

	return &response, nil
}

func (hlf *HLFClient) GetCCEventNotifier(client *channel.Client, chaincodeName string, event string) (<-chan *fab.CCEvent, error) {
	key := chaincodeName + event
	notifier := hlf.NotifierChaincodeEvents[key]
	if notifier == nil {
		var err error
		// var reg fab.Registration

		// Register chaincode event (pass in channel which receives event details when the event is complete)
		_, notifier, err = client.RegisterChaincodeEvent(chaincodeName, event)
		if err != nil {
			logger.Error("failed RegisterChaincodeEvent", zap.Error(err))
			return nil, err
		}
		if notifier == nil {
			logger.Error("failed RegisterChaincodeEvent notifier can't be nil")
			return nil, err
		}

		// defer client.UnregisterChaincodeEvent(reg)

		hlf.NotifierChaincodeEvents[key] = notifier
	}

	return notifier, nil
}

func (hlf *HLFClient) SignArgs(channelID string, chaincodeName string, methodName string, methodArgs []string, keys *keys.Keys) ([]string, error) {
	signedMessage, _, err := Sign(keys, channelID, chaincodeName, methodName, methodArgs)
	if err != nil {
		return nil, err
	}
	return signedMessage, nil
}

func (hlf *HLFClient) QueryBlockByTxID(channelID string, transactionID string, peer string) (*common.Block, error) {
	channelConnection, err := hlf.GetChannelConnection(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to GetChannelConnection: %w", err)
	}
	ledgerClient, err := ledger.New(channelConnection.ChannelProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create new ledger client: %w", err)
	}

	return ledgerClient.QueryBlockByTxID(fab.TransactionID(transactionID), ledger.WithTargetEndpoints(peer))
}

// ChaincodeVersion - only for admin user, get version for chaincode, firstly try to get version for 1.4, secondly try to get for 2.3 lifecycle
func (hlf *HLFClient) ChaincodeVersion(chaincode string, peer string) (string, error) {
	clientProvider := hlf.sdk.Context(hlf.ContextOptions...)
	client, err := resmgmt.New(clientProvider)
	if err != nil {
		return "", err
	}

	chaincodeVersion14, err := hlf.ChaincodeVersion14(client, chaincode, peer)
	if err == nil {
		return chaincodeVersion14, nil
	}

	chaincodeVersion23, err := hlf.ChaincodeVersion23Lifecycle(client, chaincode, peer)
	if err == nil {
		return chaincodeVersion23, nil
	}

	return "", fmt.Errorf("chaincode %s in channel %s not found", chaincode, chaincode)
}

// ChaincodeVersion14 - only for admin user, get version for chaincode for 1.4 hlf
func (hlf *HLFClient) ChaincodeVersion14(client *resmgmt.Client, chaincode string, peer string) (string, error) {
	chaincodeQueryResponse, err := client.QueryInstalledChaincodes(resmgmt.WithTargetEndpoints(peer))
	if err != nil {
		return "", err
	}

	for _, a := range chaincodeQueryResponse.GetChaincodes() {
		if a.GetName() == chaincode {
			return a.GetVersion(), nil
		}
	}

	return "", fmt.Errorf("chaincode %s in channel %s not found", chaincode, chaincode)
}

// ChaincodeVersion23Lifecycle - only for admin user, get version for Committed chaincode for 2.3 hlf - lifecycle
func (hlf *HLFClient) ChaincodeVersion23Lifecycle(client *resmgmt.Client, chaincode string, peer string) (string, error) {
	lifecycleQueryCommittedCC, err := client.LifecycleQueryCommittedCC(chaincode, resmgmt.LifecycleQueryCommittedCCRequest{Name: chaincode}, resmgmt.WithTargetEndpoints(peer))
	if err != nil {
		return "", err
	}

	// TODO return json array with chaincode info in channel, add test after chaincode update
	for _, a := range lifecycleQueryCommittedCC {
		if a.Name == chaincode {
			return a.Version, nil
		}
	}

	return "", fmt.Errorf("chaincode %s in channel %s not found", chaincode, chaincode)
}

// GetPeerInfo - return channel height
func (hlf *HLFClient) GetPeerInfo(channelID string, peer string) (*fab.BlockchainInfoResponse, error) {
	channelConnection, err := hlf.GetChannelConnection(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to GetChannelConnection: %w", err)
	}
	ledgerClient, err := ledger.New(channelConnection.ChannelProvider)
	if err != nil {
		return nil, err
	}
	blockchainInfoResponse, err := ledgerClient.QueryInfo(ledger.WithTargetEndpoints(peer))
	if err != nil {
		return nil, err
	}

	return blockchainInfoResponse, err
}

// GetTransactionByID - return transaction
func (hlf *HLFClient) GetTransactionByID(channelID string, transactionID string, peer string) (*peer.ProcessedTransaction, error) {
	channelConnection, err := hlf.GetChannelConnection(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to GetChannelConnection: %w", err)
	}
	ledgerClient, err := ledger.New(channelConnection.ChannelProvider)
	if err != nil {
		return nil, err
	}
	processedTransaction, err := ledgerClient.QueryTransaction(fab.TransactionID(transactionID), ledger.WithTargetEndpoints(peer))
	if err != nil {
		return nil, err
	}

	return processedTransaction, err
}

func (hlf *HLFClient) QueryBlock(channelID string, blockID string, endpoints string) (*common.Block, error) {
	channelConnection, err := hlf.GetChannelConnection(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to GetChannelConnection: %w", err)
	}
	ledgerClient, err := ledger.New(channelConnection.ChannelProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create new ledger client: %w", err)
	}

	blockIDUint, err := strconv.ParseUint(blockID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse blockID: %w", err)
	}
	return ledgerClient.QueryBlock(blockIDUint, ledger.WithTargetEndpoints(endpoints))
}
