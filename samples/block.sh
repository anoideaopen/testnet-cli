#!/bin/bash

set -xe

#./cli  --organization middleeast --connection ./bh-dev/bh-connection-dev.yaml query ba02 metadata

#./cli channelHeight atz029olp005xx dev-peer-middleeast-001.dev.bh.ledger.n-t.io --organization middleeast --connection ./bh-dev/bh-connection-dev.yaml

#./cli -s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr invoke currub healthCheck--organization middleeast --connection ./bh-dev/bh-connection-dev.yaml

#b4848ee647b79b632aaa9cd27acc292e32478023fcfcf608c88a064eb39408fc

./cli block otf 355 dev-peer-middleeast-001.dev.bh.ledger.n-t.io --organization middleeast --connection ./bh-dev/bh-connection-dev.yaml
./cli block otf 354 dev-peer-middleeast-001.dev.bh.ledger.n-t.io --organization middleeast --connection ./bh-dev/bh-connection-dev.yaml

#./cli tx atz029olp005xx b4848ee647b79b632aaa9cd27acc292e32478023fcfcf608c88a064eb39408fc dev-peer-middleeast-001.dev.bh.ledger.n-t.io --organization middleeast --connection ./bh-dev/bh-connection-dev.yaml
#./cli tx currub 7e15edc5af2ab15b7ea558a21f00802c27341c515b42957d1ff9301759d8bf66 dev-peer-middleeast-001.dev.bh.ledger.n-t.io --organization middleeast --connection ./bh-dev/bh-connection-dev.yaml
#./cli block currub 17 dev-peer-middleeast-001.dev.bh.ledger.n-t.io --organization middleeast --connection ./bh-dev/bh-connection-dev.yaml


#a091d643e35984270e103957f0630c0b00cc6d148ca94de044e148681f25f455


#./cli tx fra46046 a091d643e35984270e103957f0630c0b00cc6d148ca94de044e148681f25f455 peer0.eternyze.atmz-ch-dev-v2.ledger.n-t.io --connection connection-ch-dev-v2.yaml --organization Eternyze
#./cli block fra46046 72 peer0.eternyze.atmz-ch-dev-v2.ledger.n-t.io --connection connection-ch-dev-v2.yaml --organization Eternyze
#time ./cli -s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr invoke currub healthCheck --organization middleeast --connection ./connection.json ^C
time ./cli -s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr invoke currub healthCheck  --connection connection-ch-dev-v2.yaml --organization Eternyze
13.8
4.5
4.3
14
2.6
4.1
4.3
4.0



time ./cli -s 71684a0e25c11632a11977eea27aa3107bbf128c8425f3438054042257f85aaabdf91a67cd6d6669c0c05c33955397c11e4ac1390025ae89bad195225bb6e3ba  --organization middleeast --connection ./bh-dev/bh-connection-dev.yaml invoke curusd emit D4PHY9uDSZ7axgKWhJ82wfLCUVR625FzHZ1iCHvFMUoVhdmaV 1
15
4.3
4.5
4.1
4.2

time ./cli -s 71684a0e25c11632a11977eea27aa3107bbf128c8425f3438054042257f85aaabdf91a67cd6d6669c0c05c33955397c11e4ac1390025ae89bad195225bb6e3ba  --connection connection-ch-dev-v2.yaml --organization Eternyze invoke curusd emit D4PHY9uDSZ7axgKWhJ82wfLCUVR625FzHZ1iCHvFMUoVhdmaV 1
11
4.7
5.1
12
4.5
2.4

#

#./cli tx currencytoken 4c898cd7fc9ae4a8534570a18b3a35c962374a7cb4ca8441969fdbcde1676724 peer0.eternyze.atmz-ch-dev-v2.ledger.n-t.io --connection connection-ch-dev-v2.yaml --organization Eternyze
#./cli block currencytoken 775 peer0.eternyze.atmz-ch-dev-v2.ledger.n-t.io --connection connection-ch-dev-v2.yaml --organization Eternyze
#./cli block currencytoken 776 peer0.eternyze.atmz-ch-dev-v2.ledger.n-t.io --connection connection-ch-dev-v2.yaml --organization Eternyze


userPrivateKeyArg=$(./cli privkey)
userPublicKey=$(./cli pubkey -s $userPrivateKeyArg)
userAddress=$(./cli address -s $userPrivateKeyArg)
time ./cli --organization org0 -f connection.yaml invoke acl addUser $userPublicKey 123 testuser true | true