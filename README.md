# testnet-cli

## TOC

- [testnet-cli](#testnet-cli)
  - [TOC](#toc)
  - [Description](#description)
  - [Build](#build)
    - [Go for current platform](#go-for-current-platform)
    - [Go for different platform](#go-for-different-platform)
    - [Docker](#docker)
  - [Configuration yaml file](#configuration-yaml-file)
  - [Samples](#samples)
  - [Debug](#debug)
  - [Commands](#commands)
    - [Help](#help)
    - [Version](#version)
    - [Generate private key](#generate-private-key)
    - [Get public key by private key](#get-public-key-by-private-key)
    - [Get address by private key](#get-address-by-private-key)
    - [Get address by public key](#get-address-by-public-key)
    - [Get ski by private key](#get-ski-by-private-key)
    - [Channel height](#channel-height)
      - [Request:](#request)
      - [Response:](#response)
    - [Show tx ids in block file](#show-tx-ids-in-block-file)
    - [Save block from peer by block number](#save-block-from-peer-by-block-number)
      - [Request:](#request-1)
      - [Response:](#response-1)
    - [Save block from peer by transactionId](#save-block-from-peer-by-transactionid)
      - [Request:](#request-2)
      - [Response:](#response-2)
    - [script](#script)
    - [batchExecute](#batchexecute)
    - [Query](#query)
      - [Request:](#request-3)
      - [Response:](#response-3)
    - [Invoke](#invoke)
      - [Invoke with signed args](#invoke-with-signed-args)
      - [Request:](#request-4)
      - [Response:](#response-4)
      - [Invoke without signed args](#invoke-without-signed-args)
      - [Wait batch execute signed invoke](#wait-batch-execute-signed-invoke)
    - [Get batch execute result from observer](#get-batch-execute-result-from-observer)
    - [Convert](#convert)
    - [Performance test](#performance-test)
  - [License](#license)
  - [Links](#links)
  - [Issues:](#issues)

## Description

Утилита предназначена для выполнения следующих задач:
- Создание криптографии/пользовательских данных специфичной для платформы testnet
- Выполнения запросов к чейнкодам с учётом особенностей платформы testnet
- Выполнения запросов к fabric-peer для получения данных для дебага

Основные пользователи этой утилиты:
- тестировщики воспроизведение сценариев ошибок на стенде заказчика
- devops для получения информации
- developers для отладки работы как в sandbox так и на стендах заказчика

Пользователь может не разбираться в особенностях платформы Hyperledger Fabric,
но должен хорошо понимать логику работы для платформы testnet.

#hlf#tool#cli#go#

------
## Build

### Go for current platform
```shell
make build
```

### Go for different platform

build for:
- linux GOARCH=386
- linux GOARCH=amd64
- darwin GOARCH=amd64
- windows GOARCH=386
- windows GOARCH=amd64

```
Windows 32 bit cli-windows-386.exe
Windows 64 bit cli-windows-amd64.exe
Mac 64 bit cli-darwin-amd64
Linux 32 bit cli-linux-386
Linux 64 bit cli-linux-amd64
```
        
Execute:
```shell
make build-all
```

Result:

```
GOOS=linux GOARCH=386 go build -o "./output/cli-linux-386" -mod=vendor -ldflags "-X main.version=load -X main.commit=aaf7b3e35597c7c43d5082e0ed5763e152431bf1 -X main.date=1690909646"
GOOS=linux GOARCH=amd64 go build -o "./output/cli-linux-amd64" -mod=vendor -ldflags "-X main.version=load -X main.commit=aaf7b3e35597c7c43d5082e0ed5763e152431bf1 -X main.date=1690909646"
GOOS=darwin GOARCH=amd64 go build -o "./output/cli-darwin-amd64" -mod=vendor -ldflags "-X main.version=load -X main.commit=aaf7b3e35597c7c43d5082e0ed5763e152431bf1 -X main.date=1690909646"
GOOS=windows GOARCH=386 go build -o "./output/cli-windows-386.exe" -mod=vendor -ldflags "-X main.version=load -X main.commit=aaf7b3e35597c7c43d5082e0ed5763e152431bf1 -X main.date=1690909646"
GOOS=windows GOARCH=amd64 go build -o "./output/cli-windows-amd64.exe" -mod=vendor -ldflags "-X main.version=load -X main.commit=aaf7b3e35597c7c43d5082e0ed5763e152431bf1 -X main.date=1690909646
```

Результат можно посмотреть в директории output

```shell
ls ./output/
```

### Docker
```shell
docker build -t $(IMAGE_NAME):$(VERSION) .
```

------
## Configuration yaml file

Конфигурационный файл cli обязательный.
Пример как задается конфигурационный файл.

```shell
./testnet-cli --config ./bh-dev/cli.yaml ...
```

Пример с описанием параметров:

```yaml
# connection. Обязательно для заполнения. По умолчанию: не задано.
# путь к конфигурационному файлу подключения к стенду
connection: ./bh-dev/bh-dev-connection.yaml

# organization. Опционально. По умолчанию: testnet.
# Название организации из файла указанного в 'connection'
organization: middleeast

# username. Опционально. По умолчанию: backend.
# Название пользователя
username: backend

# waitBatch. Опционально. По умолчанию: false - не ждем событие о выполнении батча
# определяет нужно ли ждать событие "batchExecute" после выполнения запроса.
waitBatch: true

# responseType. Опционально. По умолчанию: resp.
# Тип отражения результатов выполнения запроса
responseType: resp

# observer. Опционально. По умолчанию не задан.
# Указываем подключения к observer для получения информации о выполнении батча
observer:
  # поле опционально, если не заполнено, значит не включен basic auth
  username: "root"
  # поле опционально, если не заполнено значит пароль пустой
  password: "gBm8sPqEGLFAm8v4y7"
  # обязательное поле, если в конфигурации объявлена секция 'observer:'
  url: "https://observer.dev.bh.ledger.n-t.io/api"
```
------

## Samples

Образцы выполнения запросов с помощью этой утилиты размещены в директории [samples](samples)```samples``` 

- [acl. примеры создания адреса пользователя на основе публичного ключа](samples/acl)
- [swap. примеры выполнения операции swap и multiswap](samples/swap)
- Пример выполнения запросов для токенов
  - [fiat](samples/fiat)
  - [hermitage](samples/hermitage)
  - [ba](samples/ba)
  - [industrial](samples/industrial)
  - [minetoken](samples/minetoken)

## Debug

How to run with loglevel DEBUG?

**logger:** https://github.com/sirupsen/logrus

**default log level:** ERROR

```shell
export LOG_LEVEL=debug
```

## Commands

### Help

Просмотр версии для определения списка поддерживаемых команд

Execute:
```shell
./cli -h
./cli --help
```

Return:
```json
{"Version":"load","Commit":"aaf7b3e35597c7c43d5082e0ed5763e152431bf1","Date":"1690909578"}
```


### Version

Просмотр версии для определения списка поддерживаемых команд

Execute:
```shell
./testnet-cli version
```

Return:
```json
{"Version":"load","Commit":"aaf7b3e35597c7c43d5082e0ed5763e152431bf1","Date":"1690909578"}
```

### Generate private key

Generate private key.

**Request:**

```shell
./testnet-cli privkey
```

**Response:**

Return private key in format: `private key -> base58.CheckEncode`

Example:
```
AkTSLPVKaYvKWHvb415YrJr6Vg2o3oundd6pqtDF8PVW1DqvjRycJ5PGAFxTvHAedk398tv2gdcgwmmr9hpBVVK4Le4E4
```

### Get public key by private key

Get public key by private key.

**Request:**

```shell
userPrivateKeyArg="AkTSLPVKaYvKWHvb415YrJr6Vg2o3oundd6pqtDF8PVW1DqvjRycJ5PGAFxTvHAedk398tv2gdcgwmmr9hpBVVK4Le4E4"
./testnet-cli pubkey -s $userPrivateKeyArg
```

**Response:**

Return public key in format: `public key -> base58.Check`

Example:
```
8CjV5L5KKyZR3N9zorKTR6ENNcXcskQdxAgiVKm449Tp
```

### Get address by private key

**Request:**

args:
- [command] address
- [private key] AkTSLPVKaYvKWHvb415YrJr6Vg2o3oundd6pqtDF8PVW1DqvjRycJ5PGAFxTvHAedk398tv2gdcgwmmr9hpBVVK4Le4E4

```shell
./testnet-cli address AkTSLPVKaYvKWHvb415YrJr6Vg2o3oundd6pqtDF8PVW1DqvjRycJ5PGAFxTvHAedk398tv2gdcgwmmr9hpBVVK4Le4E4
```

**Response:**

Return address key in format: (public key -> sha3.Sum256 -> base58.CheckEncode)

```
naBqaB46uCQxNQgLbCpMrVrHS694G9iLw78LFwvsM6duEzpAK
```

### Get address by public key

**Request:**

args:
- [command] address
- [public key] 6bUesd2PwAtCbRZmAU8um34D2WieE6Qsvf3uj5ZqH3B7

```shell
./testnet-cli address 6bUesd2PwAtCbRZmAU8um34D2WieE6Qsvf3uj5ZqH3B7
```

**Response:**

Return address key in format: (public key -> sha3.Sum256 -> base58.CheckEncode)

```
naBqaB46uCQxNQgLbCpMrVrHS694G9iLw78LFwvsM6duEzpAK
```

### Get ski by private key

### Channel height

Receive channel height from peer. Required connection to hlf.

#### Request:

Attributes connection config:
```
  --config ./bh-dev/cli.yaml
```

Args:
1. [command] - channelHeight
2. [channel] - acl
3. [peer url] - dev-peer-middleeast-001.dev.bh.ledger.n-t.io

```shell
./testnet-cli --config ./bh-dev/cli.yaml channelHeight acl dev-peer-middleeast-001.dev.bh.ledger.n-t.io
```

#### Response:

Return channel height on peer.

Example:
```
33233
```

### Show tx ids in block file

```shell
./testnet-cli --config ./bh-dev/cli.yaml getTxIDFromBlock ba02_1.block
```

### Save block from peer by block number

#### Request:

Attributes connection config:
```
  --config ./bh-dev/cli.yaml
```

Args:
1. [command] - block
2. [channel] - acl
3. [block number] - 273 (last block channelHeight - 1)
4. [peer url] - dev-peer-middleeast-001.dev.bh.ledger.n-t.io

```shell
./testnet-cli --config ./bh-dev/cli.yaml block acl 273 dev-peer-middleeast-001.dev.bh.ledger.n-t.io
```

```shell
for ((i=42430;i<42446;i+=1)); do ./testnet-cli --config ./ru-dev23/cli.yaml block dc ${i} peer0.hlf.testnet.dev-23.testnet-ru.ledger.n-t.io; done
```

#### Response:

Проверьте директорию с ```./testnet-cli``` там будет файл с блоков ```[acl_273.block]```

### Save block from peer by transactionId

#### Request:

Attributes connection config:
```
  --config ./bh-dev/cli.yaml
```

Args:
1. [command] - tx
2. [channel] - acl
3. [transactionID] - a67350845287b5d6ae74e5786d0d238fd192dd434ad6d23bb1ae92c79dc3d202
4. [peer url] - dev-peer-middleeast-001.dev.bh.ledger.n-t.io

```shell
./testnet-cli --config ./bh-dev/cli.yaml tx acl a67350845287b5d6ae74e5786d0d238fd192dd434ad6d23bb1ae92c79dc3d202 dev-peer-middleeast-001.dev.bh.ledger.n-t.io
```

#### Response:

Проверьте директорию с ```./testnet-cli``` там будет файл с блоков ```[acl_273.block]```

### script

Запускает переданный набор комманд в json формате

Пример файла ./script.json

```json
{
  "commands": [
    {
      "channel": "ba",
      "chaincode": "ba",
      "method": "createRedeemRequest",
      "args": ["{\"bars\":[\"BA_A728916511goldbar.1\"]}", "Redeem test"],
      "signerPrivateKey": "6fb7f9ad0c307d8fa80a5e9918002c9dbb066eb14e7175fde647cd0e58a8a5de974a32f42be7b72d735d80843106d87add11c5b107b6e2429dea43a1250d4a2b",
      "waitBatch": false
    }
  ]
}
```

```shell
./testnet-cli --config ./internal/cli.yaml script ./script.json
```

### batchExecute

Выполняет query batchExecute для указанного набора txID выполняя их в одном батче

```shell
./testnet-cli --config ./internal/cli.yaml batchExecute testbagrossmt txID
```

Пример выполнения на 1 пире для проверки таймаута:

```shell
./testnet-cli --peers prod-org0-peer-001.internal.org0.prod.core.n-t.io \
    --config ./internal/cli.yaml batchExecute testbagrossmt \
    e6b0c3c471b4518e47baae314a22275b420acd82a00ad7a5b7f8a7548535531b \
    1262b6692604bcd7b132ac7a3a22465796d33bd894aca9f4b9a0e4e7d7d1dab7 \
    9cd8acdec83d15213d0a3c308b737cbe240cb09e8622368c9894ecf74c02e2a6 \
    accfb659ea67206fbfecfd72fbe44e4df8b5b29e2136cfbc99e7c117e2d22ba4 \
    be8eef3bb11351f262e7cd271840eced37c37738293ad58934caf3c2c4cca3e5
```

### Query

Для отправки запросов на конкретный пир/пиры нужно указать параметр --peers "url1,url2,url3"
Через ``,`` перечислить список адресов пиров

#### Request:

Attributes connection config:
```
  --config ./bh-dev/cli.yaml
```

Args:
1. [command] - query
2. [channel] - ba02
3. [method] - metadata
4. [args] - далее через пробел передаем аргументы

```shell
./testnet-cli --config ./bh-dev/cli.yaml query ba02 metadata | jq
```

#### Response:

Return query result

Example:
```json
{
  "name": "bar allocated token",
  "symbol": "BA02",
  "issuer": "bQDgUMUhi2CsE9A2F3CoPuQTn1H7Dn4dEfj2MYhZvARYy4g8r",
  "redeemer": "24PLjfuJQ1GtBBPdxtYuhmMxF61onL9o6WRGdYzDnqxGdsppQW",
  "methods": [
    "acceptRedeemRequest",
    "addDocs",
    "allowedBalanceOf",
    "bAAllBalancesOf",
    "bABalanceOf",
    "barsGroupList",
    "barsPrices",
    "buildInfo",
    "buyBack",
    "buyTokens",
    "cancelCCTransferFrom",
    "channelTransferByAdmin",
    "channelTransferByCustomer",
    "channelTransferFrom",
    "channelTransferTo",
    "channelTransfersFrom",
    "commitCCTransferFrom",
    "coreChaincodeIDName",
    "createCCTransferTo",
    "createRedeemRequest",
    "deleteCCTransferFrom",
    "deleteCCTransferTo",
    "deleteDoc",
    "deleteRate",
    "denyAllRedeemRequest",
    "denyRedeemRequest",
    "documentsList",
    "emitTokensFromBars",
    "getLockedAllowedBalance",
    "getLockedTokenBalance",
    "getNonce",
    "groupBalanceOf",
    "healthCheck",
    "lockAllowedBalance",
    "lockTokenBalance",
    "metadata",
    "multiSwapBegin",
    "multiSwapCancel",
    "multiSwapGet",
    "nameOfFiles",
    "redeemRequestsList",
    "setRate",
    "srcFile",
    "srcPartFile",
    "swapBegin",
    "swapCancel",
    "swapGet",
    "systemEnv",
    "totalEmission",
    "transfer",
    "unlockAllowedBalance",
    "unlockTokenBalance"
  ],
  "bar_groups": [
    "gol523179052",
    "gol523279052",
    "gol523379052",
    "gol523479052",
    "gol523579052",
    "gol523679052",
    "gol523779052",
    "gol523129035",
    "gol523229035",
    "gol523329052",
    "gol523429052",
    "gol523529052",
    "gol523629052",
    "gol523729052",
    "cop523129035",
    "cop523229035",
    "cop523329052",
    "cop523429052",
    "cop523529052",
    "cop523629052",
    "cop523729052",
    "cop523829035",
    "cop523929035",
    "cop523329035",
    "cop523429035",
    "cop523529035",
    "cop523629035",
    "cop523729035",
    "44714279052",
    "74507279052",
    "95061979052",
    "BarPal99029055",
    "BarPal99129055",
    "BarPal99229055",
    "BarPal99329055",
    "BarPal99429055",
    "BarPal99529055",
    "BarPal99629055",
    "BarPal99729055",
    "BarPal99829055",
    "BarPal99929055",
    "BarPa200029055",
    "BarPa200129055",
    "BarKris99629055",
    "BarKris99729055",
    "BarKris99829055",
    "BarKris99929055",
    "BarKris100029055",
    "BarKris100129055",
    "coptest129035",
    "coptest1129035",
    "coptest1329035",
    "coptest1429035",
    "cop524029035",
    "cop524129035",
    "cop524229035",
    "cop524329035",
    "cop524429035",
    "cop524529035",
    "cop524629035",
    "068085179052",
    "331446179052",
    "105015179052",
    "coptest14229035",
    "coptest14329035",
    "coptest14429035",
    "coptest14529035",
    "BarGol00079049",
    "BarGol00179049",
    "BarGol00279049",
    "BarGol00379049",
    "BarGol00479049",
    "BarGol00579049",
    "BarGol00679049",
    "BarGol00779049",
    "BarGol00879049",
    "BarGol00979049"
  ],
  "total_emission": "7700000000",
  "rates": [
    {
      "deal_type": "buyTokens",
      "underlying_asset": "79",
      "delivery_form": "052",
      "currency": "AT99USD",
      "rate": "2"
    },
    {
      "deal_type": "buyBack",
      "underlying_asset": "79",
      "delivery_form": "052",
      "currency": "AT99USD",
      "rate": "3"
    },
    {
      "deal_type": "buyTokens",
      "underlying_asset": "29",
      "delivery_form": "035",
      "currency": "AT99USD",
      "rate": "4"
    },
    {
      "deal_type": "buyBack",
      "underlying_asset": "29",
      "delivery_form": "035",
      "currency": "AT99USD",
      "rate": "3"
    },
    {
      "deal_type": "buyTokens",
      "underlying_asset": "47",
      "delivery_form": "053",
      "currency": "CURRENCYTOKEN",
      "rate": "1"
    },
    {
      "deal_type": "buyBack",
      "underlying_asset": "47",
      "delivery_form": "053",
      "currency": "CURRENCYTOKEN",
      "rate": "0.66"
    },
    {
      "deal_type": "buyTokens",
      "underlying_asset": "29",
      "delivery_form": "035",
      "currency": "CURRENCYTOKEN",
      "rate": "1"
    },
    {
      "deal_type": "buyBack",
      "underlying_asset": "29",
      "delivery_form": "035",
      "currency": "CURRENCYTOKEN",
      "rate": "0.66"
    },
    {
      "deal_type": "buyTokens",
      "underlying_asset": "29",
      "delivery_form": "035",
      "currency": "CURUSD",
      "rate": "1"
    },
    {
      "deal_type": "buyBack",
      "underlying_asset": "29",
      "delivery_form": "035",
      "currency": "CURUSD",
      "rate": "1"
    }
  ]
}
```

### Invoke

Для отправки запросов на конкретный пир/пиры нужно указать параметр --peers "url1,url2,url3"
Через ``,`` перечислить список адресов пиров

#### Invoke with signed args

Sign args and send invoke to hlf. Required connection to hlf.

signed args
- [method]
- [requestId]
- [channel]
- [chaincode]
- [args...]
- [nonce]
- [pubkey]
- [signature]

send args
- [requestId]
- [channel]
- [chaincode]
- [args...]
- [nonce]
- [pubkey]
- [signature]

#### Request:

Attributes connection config:
```
  --config ./bh-dev/cli.yaml
```

Args:
1. [command] - invoke
2. [channel] - ba
3. [args..] - createRedeemRequest '{"bars":["BA_A728916511goldbar.1"]}' "Redeem test"

```
-s "private key"
```

```shell
./testnet-cli --config ./bh-dev/cli.yaml -s 6fb7f9ad0c307d8fa80a5e9918002c9dbb066eb14e7175fde647cd0e58a8a5de974a32f42be7b72d735d80843106d87add11c5b107b6e2429dea43a1250d4a2b invoke ba createRedeemRequest '{"bars":["BA_A728916511goldbar.1"]}' "Redeem test"
```

#### Response:

Return channel height on peer.

Example:
```
33233
```

#### Invoke without signed args

Doesn't sign args and send invoke to hlf. Required connection to hlf.

Request.

```shell
./testnet-cli --config ./bh-dev/cli.yaml invoke acl addUser "6bUesd2PwAtCbRZmAU8um34D2WieE6Qsvf3uj5ZqH3B7" "unknown" "testUser4" "true"
```

Response SUCCESS:

```
TransactionID:
f60f3007b2973cefc858a333277e50802cc059f76124dc00f7c1d96bf5a07d53
TxValidationCode:
VALID
BlockNumber:
278



```

Response ERROR - connection is in TRANSIENT_FAILURE:
```
{"level":"error","ts":1690913989.0852187,"caller":"logger/logger.go:37","msg":"error","error":"CreateAndSendTransaction failed: SendTransaction failed: calling orderer 'stage-orderer-middleeast-004.stage.bh.ledger.n-t.io:12336' failed: Orderer Client Status Code: (2) CONNECTION_FAILED. Description: dialing connection on target [stage-orderer-middleeast-004.stage.bh.ledger.n-t.io:12336]: connection is in TRANSIENT_FAILURE","errorVerbose":"Orderer Client Status Code: (2) CONNECTION_FAILED. Description: dialing connection on target [stage-orderer-middleeast-004.stage.bh.ledger.n-t.io:12336]: connection is in TRANSIENT_FAILURE\ncalling orderer 'stage-orderer-middleeast-004.stage.bh.ledger.n-t.io:12336' failed\ngithub.com/hyperledger/fabric-sdk-go/pkg/fab/txn.sendBroadcast\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/fab/txn/txn.go:284\ngithub.com/hyperledger/fabric-sdk-go/pkg/fab/txn.broadcastEnvelope\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/fab/txn/txn.go:209\ngithub.com/hyperledger/fabric-sdk-go/pkg/fab/txn.BroadcastPayload\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/fab/txn/txn.go:185\ngithub.com/hyperledger/fabric-sdk-go/pkg/fab/txn.Send\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/fab/txn/txn.go:151\ngithub.com/hyperledger/fabric-sdk-go/pkg/fab/channel.(*Transactor).SendTransaction\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/fab/channel/transactor.go:187\ngithub.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke.createAndSendTransaction\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke/txnhandler.go:296\ngithub.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke.(*CommitTxHandler).Handle\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke/txnhandler.go:204\ngithub.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke.(*SignatureValidationHandler).Handle\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke/signature.go:37\ngithub.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke.(*EndorsementValidationHandler).Handle\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke/txnhandler.go:161\ngithub.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke.(*SelectAndEndorseHandler).Handle\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke/selectendorsehandler.go:90\ngithub.com/hyperledger/fabric-sdk-go/pkg/client/channel.(*Client).InvokeHandler.func2.1\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/client/channel/chclient.go:191\ngithub.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry.(*RetryableInvoker).Invoke\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry/invoker.go:63\ngithub.com/hyperledger/fabric-sdk-go/pkg/client/channel.(*Client).InvokeHandler.func2\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/client/channel/chclient.go:189\nruntime.goexit\n\t/home/yury/go/go1.18/src/runtime/asm_amd64.s:1571\nSendTransaction failed\nCreateAndSendTransaction failed\ngithub.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke.(*CommitTxHandler).Handle\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke/txnhandler.go:206\ngithub.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke.(*SignatureValidationHandler).Handle\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke/signature.go:37\ngithub.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke.(*EndorsementValidationHandler).Handle\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke/txnhandler.go:161\ngithub.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke.(*SelectAndEndorseHandler).Handle\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke/selectendorsehandler.go:90\ngithub.com/hyperledger/fabric-sdk-go/pkg/client/channel.(*Client).InvokeHandler.func2.1\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/client/channel/chclient.go:191\ngithub.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry.(*RetryableInvoker).Invoke\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry/invoker.go:63\ngithub.com/hyperledger/fabric-sdk-go/pkg/client/channel.(*Client).InvokeHandler.func2\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/vendor/github.com/hyperledger/fabric-sdk-go/pkg/client/channel/chclient.go:189\nruntime.goexit\n\t/home/yury/go/go1.18/src/runtime/asm_amd64.s:1571","stacktrace":"github.com/anoideaopen/testnet-cli/logger.Error\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/logger/logger.go:37\ngithub.com/anoideaopen/testnet-cli/service.(*HLFClient).RequestChaincode\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/service/fabric.go:412\ngithub.com/anoideaopen/testnet-cli/service.(*HLFClient).Invoke\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/service/fabric.go:320\ngithub.com/anoideaopen/testnet-cli/cmd.glob..func11.2\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/cmd/invokeCmd.go:99"}
Invoke error: CreateAndSendTransaction failed: SendTransaction failed: calling orderer 'stage-orderer-middleeast-004.stage.bh.ledger.n-t.io:12336' failed: Orderer Client Status Code: (2) CONNECTION_FAILED. Description: dialing connection on target [stage-orderer-middleeast-004.stage.bh.ledger.n-t.io:12336]: connection is in TRANSIENT_FAILURE
```

Response ERROR - already exists:
```
{"level":"error","ts":1690914482.816641,"caller":"logger/logger.go:37","msg":"error","error":"Multiple errors occurred: - Transaction processing for endorser [dev-peer-org3-001.dev.bh.ledger.n-t.io:16323]: Chaincode status Code: (500) UNKNOWN. Description: The address 22UpsSKQXh57NAtcj5CZjow8Q88Hc4j46kKjk8huJBDenyCFYf associated with key 870e58cade839037c1bc341b88f02d6bbd3886af7afd133c67798879a9f327d9 already exists - Transaction processing for endorser [dev-peer-middleeast-001.dev.bh.ledger.n-t.io:16536]: Chaincode status Code: (500) UNKNOWN. Description: The address 22UpsSKQXh57NAtcj5CZjow8Q88Hc4j46kKjk8huJBDenyCFYf associated with key 870e58cade839037c1bc341b88f02d6bbd3886af7afd133c67798879a9f327d9 already exists - Transaction processing for endorser [dev-peer-org1-001.dev.bh.ledger.n-t.io:16552]: Chaincode status Code: (500) UNKNOWN. Description: The address 22UpsSKQXh57NAtcj5CZjow8Q88Hc4j46kKjk8huJBDenyCFYf associated with key 870e58cade839037c1bc341b88f02d6bbd3886af7afd133c67798879a9f327d9 already exists","stacktrace":"github.com/anoideaopen/testnet-cli/logger.Error\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/logger/logger.go:37\ngithub.com/anoideaopen/testnet-cli/service.(*HLFClient).RequestChaincode\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/service/fabric.go:412\ngithub.com/anoideaopen/testnet-cli/service.(*HLFClient).Invoke\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/service/fabric.go:320\ngithub.com/anoideaopen/testnet-cli/cmd.glob..func11.2\n\t/home/yury/go/src/github.com/anoideaopen/testnet-cli/cmd/invokeCmd.go:99"}
Invoke error: Multiple errors occurred: - Transaction processing for endorser [dev-peer-org3-001.dev.bh.ledger.n-t.io:16323]: Chaincode status Code: (500) UNKNOWN. Description: The address 22UpsSKQXh57NAtcj5CZjow8Q88Hc4j46kKjk8huJBDenyCFYf associated with key 870e58cade839037c1bc341b88f02d6bbd3886af7afd133c67798879a9f327d9 already exists - Transaction processing for endorser [dev-peer-middleeast-001.dev.bh.ledger.n-t.io:16536]: Chaincode status Code: (500) UNKNOWN. Description: The address 22UpsSKQXh57NAtcj5CZjow8Q88Hc4j46kKjk8huJBDenyCFYf associated with key 870e58cade839037c1bc341b88f02d6bbd3886af7afd133c67798879a9f327d9 already exists - Transaction processing for endorser [dev-peer-org1-001.dev.bh.ledger.n-t.io:16552]: Chaincode status Code: (500) UNKNOWN. Description: The address 22UpsSKQXh57NAtcj5CZjow8Q88Hc4j46kKjk8huJBDenyCFYf associated with key 870e58cade839037c1bc341b88f02d6bbd3886af7afd133c67798879a9f327d9 already exists
```

#### Wait batch execute signed invoke

На данный момент оформлена ошибка [подробнее](#Issues)

Опционально. По умолчанию не ожидаем событий batch execute.

```yaml
waitBatch: true
```

```shell
./testnet-cli --waitBatch true
```

```shell
./testnet-cli -w true
```

### Get batch execute result from observer

Get batch execute result from observer by transaction id. Required connection to observer.

**Request:**

args
- [command] - status
- [transactionID] - cc84fcb8934d326fac095cccd6848ca21b167757636ea93a553d15f392cbf4ac

```shell
./testnet-cli --config ./bh-dev/cli.yaml status cc84fcb8934d326fac095cccd6848ca21b167757636ea93a553d15f392cbf4ac
```

**Response:**
```
-------- Batch tx found in observer:
Request:
TxID: cc84fcb8934d326fac095cccd6848ca21b167757636ea93a553d15f392cbf4ac
CreatedAt: 2023-08-01T21:11:12.690896135+03:00

Batch:
TxID: a6525ddb9b5d2199b8dcfad7d0da6744d3ad881d81e389e13ee6a1f510ae21d0
BlockNumber: 15
CreatedAt: 2023-08-01T21:11:12.69089623+03:00
BatchErrorMsg: 
BatchValidationCode: 0
```

### Convert

Для конвертации строки из одной кодировки в другую.

**Комбинации конвертаций которые поддерживаются на данный момент:**
- из base58 в hex
- из base58 в sum3hex
- из base58 в base58check
- из base58 в Sum256base58CheckEncode
- из base58check в base58
- из base58check в base64
- из hex в base58
- из hex в base58check
- из str в base58
- из str в hex
- из str в hex
- из str в base58check

Request:

```shell
./cli convert base58 hex FmUXc1fudiUREQSpqc5pgi5MZYH6XHaAGoRSDoeB2QpT
./cli convert hex base58 db684c558b4e1dcfdc98b0b629bc572742361b42a5c7ac1709d68bc126cdbc64
./cli convert hex base58check db684c558b4e1dcfdc98b0b629bc572742361b42a5c7ac1709d68bc126cdbc64
./cli convert str base58 строка
./cli convert str hex строка
./cli convert str base58check строка
```

### Performance test

Параметры для метода invoke которые можно использовать для нагрузочного тестирования.
Параметры доступны для методов `./testnet-cli` 'invoke', 'query'

- `-t 50` или `--requestsPerSecond 50` указав кол-во параллельно запущенных goroutine за одну секунду
- общее кол-во запросов
  - `-n 1000` или `--numberRequest 1000` указать максимальное кол-во запросов
  - `-n 0` или `--numberRequest 0` задать чтобы нагрузка выполнялась бесконечно
- во время нагрузки не нужно ждать событий по этому параметр `waitBatch` не должен быть `true` ни в конфиге не в env ни в параметрах запроса. По умолчанию этот параметр false.

**Request:**

```shell
./testnet-cli -n 0 -t 50 --config ./bh-stage/cli.yaml -s 71684a0e25c11632a11977eea27aa3107bbf128c8425f3438054042257f85aaabdf91a67cd6d6669c0c05c33955397c11e4ac1390025ae89bad195225bb6e3ba invoke atz029olp005xx healthCheck
```

**Response:**

Для удобного чтения вывода нагрузки можно сделать табличный вывод указав `-r table` или же записать информацию о нагрузки в бд postgres задав параметр `-r postgres`

**table report**

```shell
| tx | 17abb9210b2c662e408589aca82dc00d3d54b0af8544fa7a2678ff43e7c5c84b | block | 99570 | start | 2023-08-01 21:25:23.422819 | end | 2023-08-01 21:25:31.096039 | dur | 7.673220 |
| tx | 5456fdea088738283911020e0b5148e41121a9536e5d809c9ca3a131310216f1 | block | 99570 | start | 2023-08-01 21:25:22.422870 | end | 2023-08-01 21:25:31.096050 | dur | 8.673181 |
| tx | 2ef4005be9e693bc0b507b157562dca8eb6cf07835eb3248d25bb24fb7d4bcbb | block | 99570 | start | 2023-08-01 21:25:22.922426 | end | 2023-08-01 21:25:31.096072 | dur | 8.173646 |
| tx | 2a4be06ea7d447a443d20ff66b6e10156f364754b45ccc53ba8f44ce05debf37 | block | 99571 | start | 2023-08-01 21:25:23.922324 | end | 2023-08-01 21:25:31.394234 | dur | 7.471909 |
| tx | 27aa6af19aa7d3a61c4573a011e4555095467734193da0827c4dca35ee5b0012 | block | 99572 | start | 2023-08-01 21:25:26.921981 | end | 2023-08-01 21:25:31.512033 | dur | 4.590051 |
| tx | b7636379acaa4160b50892a9ff4426ee9ee67f89df70c6cd73d4400e49985d0d | block | 99572 | start | 2023-08-01 21:25:26.422679 | end | 2023-08-01 21:25:31.512040 | dur | 5.089361 |
| tx | 73fc771635867f30c6f5e58419d6e5ba46469a81f5cec0cf61cc37099a4aecfe | block | 99572 | start | 2023-08-01 21:25:24.422048 | end | 2023-08-01 21:25:31.512045 | dur | 7.089997 |
| tx | aef939de4dbd35009fe346f2f91c4ec5af963ea0784e10a6d293250b2e750d02 | block | 99572 | start | 2023-08-01 21:25:24.922561 | end | 2023-08-01 21:25:31.512049 | dur | 6.589488 |
| tx | 86b41fa94aff151c3ce9c81da925ad2e54ce66a387dba5f5ab3a431e69c6a2e1 | block | 99573 | start | 2023-08-01 21:25:25.422436 | end | 2023-08-01 21:25:31.697953 | dur | 6.275517 |
```

## License

Apache-2.0

## Links

* [origin](https://github.com/anoideaopen/migration-manager)
* [Процесс обновления публичного ключа](CHANGE_PUBLIC_KEY.md)

## Issues:
- "cli не находит ошибку в event о создании батча" https://nwty.atlassian.net/browse/ATMCORE-6634
