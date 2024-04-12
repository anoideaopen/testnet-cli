#!/bin/bash

set -xe

#searchFreeBarsParts 9e7fca93baaad9ad0b386e2a8f8e495fd36dc3c3c8635e507cd61d32b956d17c
export connection_file="./connection.yaml"


connectionConfig=" -f ./connection.yaml --organization org0 "
./cli $connectionConfig query hermitage allowedBalanceOf $issuerAddress FIAT

# QueryMetadata
 ./cli --organization org0 -f $connection_file query hermitage metadata
 ./cli --organization org0 -f $connection_file query hermitage emissionsList "{}"
#
export signAdmin=" -s DGgvqCBCSw84d7KywEhpXMFEbyMdL7WcXaEiCgHKBAoPyJaD2CCzCvfHWxgnL8iTcFzwtrpx76KpFE2RWXKVUgMwKJ2vk "
export signIssuer=" -s DGgvqCBCSw84d7KywEhpXMFEbyMdL7WcXaEiCgHKBAoPyJaD2CCzCvfHWxgnL8iTcFzwtrpx76KpFE2RWXKVUgMwKJ2vk "

adminPublicKey=$(./cli pubkey $signAdmin)
issuerPublicKey=$(./cli pubkey $signIssuer)

./cli $signIssuer --organization org0 -f $connection_file invoke hermitage createEmissionApp "{\"ticker_suffix\":\"6\",\"name\":\"Name\",\"description\":\"Description\",\"price\":\"Price\",\"redemption_price\":\"1\",\"currency\":\"FIAT\",\"recipient\":\"Recipient\",\"schedule\":{\"investment_start_date\":\"0001-01-01T00:00:00Z\",\"investment_end_date\":\"0001-01-01T00:00:00Z\",\"redemption_date\":\"0001-01-01T00:00:00Z\",\"trade_start_date\":\"0001-01-01T00:00:00Z\"},\"files\":{\"photo_hash\":\"PhotoHash\",\"photo_url\":\"PhotoURL\",\"video_hash\":\"VideoHash\",\"video_url\":\"VideoURL\",\"docs_list_hash\":[\"test\"]}}"

#./cli $signIssuer --organization org0 -f $connection_file invoke hermitage acceptEmissionApp 6
#./cli $signIssuer --organization org0 -f $connection_file invoke hermitage startInvestment 6
#./cli $signIssuer --organization org0 -f $connection_file invoke hermitage sendInvestmentApp 6 FIAT
#
./cli --organization org0 -f $connection_file query hermitage emissionDetails 6
./cli --organization org0 -f $connection_file query hermitage lockedBalanceOf $issuerAddress FIAT
#
#./cli $signIssuer --organization org0 -f $connection_file invoke hermitage endInvestment 6
#./cli $signIssuer --organization org0 -f $connection_file invoke hermitage revokeInvestmentApp a0853217ab40c906c245ca7148ee8c379e3cca38e233664e97bc8ef36ea7f875
#./cli $signIssuer --organization org0 -f $connection_file invoke hermitage revokeInvestmentApp 68b0dc60d5a7c9523c68db09d6568248ef383155e44825ddf5f38d492ad973b0
#
#./cli --organization org0 -f $connection_file query hermitage emissionsList "{}"
#
#./cli --organization org0 -f $connection_file query hermitage investmentAppsList "{}"
#
#./cli --organization org0 -f $connection_file $signAdmin invoke hermitage rejectEmissionApp TickerSuffix1 reason
#./cli --organization org0 -f $connection_file query hermitage emissionsList "{\"statuses\":[2]}"
#

#
#./cli -s 697384ca952cb3a1a90c5686f1f48527d5ab63456847fdb45bd3a6179616d21e4eb2032b53c57ac6197774dda66f946fc17fd08b6e0f6f3c57616886c0cac111 -f ru-dev-23.yaml invoke hermitage createEmissionApp "{\"ticker_suffix\":\"TickerSuffix1\",\"name\":\"Name\",\"description\":\"Description\",\"price\":\"Price\",\"redemption_price\":\"RedemptionPrice\",\"currency\":\"Currency\",\"recipient\":\"Recipient\",\"schedule\":{\"investment_start_date\":\"0001-01-01T00:00:00Z\",\"investment_end_date\":\"0001-01-01T00:00:00Z\",\"redemption_date\":\"0001-01-01T00:00:00Z\",\"trade_start_date\":\"0001-01-01T00:00:00Z\"},\"files\":{\"photo_hash\":\"PhotoHash\",\"photo_url\":\"PhotoURL\",\"video_hash\":\"VideoHash\",\"video_url\":\"VideoURL\",\"docs_list_hash\":[\"test\"]}}"
#
#./cli -f ru-dev-23.yaml query hermitage emissionsList "{}"
#
#./cli -f ru-dev-23.yaml block hermitage 10 peer0.hlf.org1.dev-23.testnet-ru.ledger.n-t.io
#./cli -f ru-dev-23.yaml readBlockFile hermitage 10 peer0.hlf.org1.dev-23.testnet-ru.ledger.n-t.io
#./cli -f ru-dev-23.yaml channelHeight hermitage peer0.hlf.org1.dev-23.testnet-ru.ledger.n-t.io

#issuerUserNewSecretKey="AG1XotNqaRbtZtuWNbT9nNjbndqhrFvRCgERuGz9vCHc"
#issuerUserNewPublicKey="DGgvqCBCSw84d7KywEhpXMFEbyMdL7WcXaEiCgHKBAoPyJaD2CCzCvfHWxgnL8iTcFzwtrpx76KpFE2RWXKVUgMwKJ2vk"
#
#./cli --organization org0 -f c.yaml invoke acl addUser ${issuerUserNewSecretKey} 123 testuser true
#
#./cli --organization org0 -f c.yaml query hermitage metadata
#
./cli -s 697384ca952cb3a1a90c5686f1f48527d5ab63456847fdb45bd3a6179616d21e4eb2032b53c57ac6197774dda66f946fc17fd08b6e0f6f3c57616886c0cac111 -f ru-dev-23.yaml invoke hermitage createEmissionApp "{\"ticker_suffix\":\"TickerSuffix3\",\"name\":\"Name\",\"description\":\"Description\",\"price\":\"Price\",\"redemption_price\":\"RedemptionPrice\",\"currency\":\"Currency\",\"recipient\":\"Recipient\",\"schedule\":{\"investment_start_date\":\"0001-01-01T00:00:00Z\",\"investment_end_date\":\"0001-01-01T00:00:00Z\",\"redemption_date\":\"0001-01-01T00:00:00Z\",\"trade_start_date\":\"0001-01-01T00:00:00Z\"},\"files\":{\"photo_hash\":\"PhotoHash\",\"photo_url\":\"PhotoURL\",\"video_hash\":\"VideoHash\",\"video_url\":\"VideoURL\",\"docs_list_hash\":[\"test\"]}}"
#./cli --organization org0 -f c.yaml query hermitage emissionsList "{}"
#
./cli -f ru-dev-23.yaml -s 6fb7f9ad0c307d8fa80a5e9918002c9dbb066eb14e7175fde647cd0e58a8a5de974a32f42be7b72d735d80843106d87add11c5b107b6e2429dea43a1250d4a2b invoke hermitage acceptEmissionApp TickerSuffix3
./cli -f ru-dev-23.yaml query hermitage emissionsList "{\"statuses\":[1]}"


./cli -s 6fb7f9ad0c307d8fa80a5e9918002c9dbb066eb14e7175fde647cd0e58a8a5de974a32f42be7b72d735d80843106d87add11c5b107b6e2429dea43a1250d4a2b -f ru-dev-23.yaml query hermitage deleteEmissionApp TickerSuffix3
