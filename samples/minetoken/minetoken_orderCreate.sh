./cli -f connaction-preprod-23-testnet-ru.yaml query minetoken metadata

./cli -f connaction-preprod-23-testnet-ru.yaml invoke minetoken orderCreate -s 6evUikkpSoibhWc3wkP8h1bnAUb11L58RNjQZGXVUcofW3w1nuPhqAt9cH9VxYTUb7BMz7F5pFVzyVfgr8Y92AAZ9wcFo "2QQhdyVb5zKv6T1bZ6zvNggxqzZ5sGFvGe2ms84Lu5H56Q2T6i" "18" "1000000" "101000000" "gCe5cYuGTBsRrULiacRDcSsA2As9HkiMMSDBRyQzhGKtTsLZN" "{\"initId\":\"\",\"type\":\"p2p\",\"direction\":\"SELL\",\"side\":\"SELLER\",\"conditions\":{\"baseAsset\":\"MINETOKEN_18\",\"baseAmount\":\"1000000\",\"quoteAsset\":\"RUB\",\"quoteAmount\":\"1000000000\"}}"

./cli channelHeight minetoken peer0.hlf.testnet.preprod-23.testnet-ru.ledger.n-t.io -f connaction-preprod-23-testnet-ru.yaml
./cli block minetoken 561483 peer0.hlf.testnet.preprod-23.testnet-ru.ledger.n-t.io -f connaction-preprod-23-testnet-ru.yaml
./cli tx minetoken 03ce1f0f2a36537222897750306d62310f8ae5983fef5c5253e3cbfe12626c6c peer0.hlf.testnet.preprod-23.testnet-ru.ledger.n-t.io -f connaction-preprod-23-testnet-ru.yaml

tx channelID transactionID peerUrl

#03ce1f0f2a36537222897750306d62310f8ae5983fef5c5253e3cbfe12626c6c

#2B37303030303133303939362B37303030303133303939362B37303030303133C747175FFA520E9F2F79602B165E9E1D74D644C4D0C3B91A17511C7E086CA0D7
#
#
#6evUikkpSoibhWc3wkP8h1bnAUb11L58RNjQZGXVUcofW3w1nuPhqAt9cH9VxYTUb7BMz7F5pFVzyVfgr8Y92AAZ9wcFo
