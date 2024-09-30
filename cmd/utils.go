package cmd

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/anoideaopen/testnet-cli/service"
	"github.com/golang/protobuf/proto" //nolint:staticcheck
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func handlerArgs(args []string) (string, string, []string) {
	channelID := args[0]
	methodName := args[1]
	methodArgs := args[2:]

	if len(config.ChaincodeName) == 0 {
		config.ChaincodeName = channelID
	}

	return channelID, methodName, methodArgs
}

func FatalError(errorMessage string, err error) {
	if err != nil {
		logger.Error(errorMessage, zap.Error(err))
	} else {
		logger.Error(errorMessage)
	}

	panic(err)
}

var HlfClient *service.HLFClient

func initHlfClient() {
	var err error
	HlfClient, err = service.NewHLFClient(config.Connection, config.User, config.Organization, nil)
	if err != nil {
		FatalError("Failed to create new channel client", err)
	}
}

// UnmarshalBlock unmarshals bytes to a Block
func UnmarshalBlock(encoded []byte) (*common.Block, error) {
	block := &common.Block{}
	err := proto.Unmarshal(encoded, block)
	return block, errors.Wrap(err, "error unmarshalling Block")
}

// UnmarshalTransaction unmarshals bytes to a Transaction
func UnmarshalTransaction(txBytes []byte) (*peer.Transaction, error) {
	tx := &peer.Transaction{}
	err := proto.Unmarshal(txBytes, tx)
	return tx, errors.Wrap(err, "error unmarshalling Transaction")
}

// UnmarshalChaincodeEvents unmarshals bytes to a ChaincodeEvent
func UnmarshalChaincodeEvents(eBytes []byte) (*peer.ChaincodeEvent, error) {
	chaincodeEvent := &peer.ChaincodeEvent{}
	err := proto.Unmarshal(eBytes, chaincodeEvent)
	return chaincodeEvent, errors.Wrap(err, "error unmarshalling ChaicnodeEvent")
}

// GetOrComputeTxIDFromEnvelope gets the txID present in a given transaction
// envelope. If the txID is empty, it constructs the txID from nonce and
// creator fields in the envelope.
func GetOrComputeTxIDFromEnvelope(txEnvelopBytes []byte) (string, error) {
	txEnvelope, err := UnmarshalEnvelope(txEnvelopBytes)
	if err != nil {
		return "", errors.WithMessage(err, "error getting txID from envelope")
	}

	txPayload, err := UnmarshalPayload(txEnvelope.GetPayload())
	if err != nil {
		return "", errors.WithMessage(err, "error getting txID from payload")
	}

	if txPayload.GetHeader() == nil {
		return "", errors.New("error getting txID from header: payload header is nil")
	}

	chdr, err := UnmarshalChannelHeader(txPayload.GetHeader().GetChannelHeader())
	if err != nil {
		return "", errors.WithMessage(err, "error getting txID from channel header")
	}

	if chdr.GetTxId() != "" {
		return chdr.GetTxId(), nil
	}

	sighdr, err := UnmarshalSignatureHeader(txPayload.GetHeader().GetSignatureHeader())
	if err != nil {
		return "", errors.WithMessage(err, "error getting nonce and creator for computing txID")
	}

	txid := ComputeTxID(sighdr.GetNonce(), sighdr.GetCreator())
	return txid, nil
}

// UnmarshalEnvelope unmarshals bytes to a Envelope
func UnmarshalEnvelope(encoded []byte) (*common.Envelope, error) {
	envelope := &common.Envelope{}
	err := proto.Unmarshal(encoded, envelope)
	return envelope, errors.Wrap(err, "error unmarshalling Envelope")
}

// UnmarshalPayload unmarshals bytes to a Payload
func UnmarshalPayload(encoded []byte) (*common.Payload, error) {
	payload := &common.Payload{}
	err := proto.Unmarshal(encoded, payload)
	return payload, errors.Wrap(err, "error unmarshalling Payload")
}

// UnmarshalChannelHeader unmarshals bytes to a ChannelHeader
func UnmarshalChannelHeader(bytes []byte) (*common.ChannelHeader, error) {
	chdr := &common.ChannelHeader{}
	err := proto.Unmarshal(bytes, chdr)
	return chdr, errors.Wrap(err, "error unmarshalling ChannelHeader")
}

// UnmarshalSignatureHeader unmarshals bytes to a SignatureHeader
func UnmarshalSignatureHeader(bytes []byte) (*common.SignatureHeader, error) {
	sh := &common.SignatureHeader{}
	err := proto.Unmarshal(bytes, sh)
	return sh, errors.Wrap(err, "error unmarshalling SignatureHeader")
}

// ComputeTxID computes TxID as the Hash computed
// over the concatenation of nonce and creator.
func ComputeTxID(nonce, creator []byte) string {
	// TODO: Get the Hash function to be used from
	// channel configuration
	hasher := sha256.New()
	hasher.Write(nonce)
	hasher.Write(creator)
	return hex.EncodeToString(hasher.Sum(nil))
}

// GetEnvelopeFromBlock gets an envelope from a block's Data field.
func GetEnvelopeFromBlock(data []byte) (*common.Envelope, error) {
	// Block always begins with an envelope
	var err error
	env := &common.Envelope{}
	if err = proto.Unmarshal(data, env); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling Envelope")
	}

	return env, nil
}

func UnmarshalSerializedIdentity(bytes []byte) (*msp.SerializedIdentity, error) {
	sid := &msp.SerializedIdentity{}
	err := proto.Unmarshal(bytes, sid)
	return sid, errors.Wrap(err, "error unmarshalling SerializedIdentity")
}
