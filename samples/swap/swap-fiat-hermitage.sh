#!/bin/bash

set -xe

connectionConfig=" -f ./connection.yaml --organization org0 "

fromCC="fiat"
toCC="hermitage"

fromCC_uppercase=${fromCC^^}
toCC_uppercase=${toCC^^}

issuerSKeyArg=" -s DGgvqCBCSw84d7KywEhpXMFEbyMdL7WcXaEiCgHKBAoPyJaD2CCzCvfHWxgnL8iTcFzwtrpx76KpFE2RWXKVUgMwKJ2vk"
issuerPublicKey=$(./cli pubkey $issuerSKeyArg)
issuerAddress=$(./cli address $issuerSKeyArg)

userPrivateKeyArg=$issuerSKeyArg
userPublicKey=$issuerPublicKey
userAddress=$issuerAddress

swapHash="7d4e3eec80026719639ed4dba68916eb94c7a49a053e05c8f9578fe4e5a3d7ea"
swapKey="12345"

echo "==================  swapBegin  ================================="
TX_ID=`./cli $connectionConfig $userPrivateKeyArg -r tx invoke $fromCC swapBegin $fromCC_uppercase $toCC_uppercase 100 "$swapHash"`
echo $TX_ID
sleep 10
echo "==================  swapGet  ================================="
./cli $connectionConfig query $fromCC swapGet $TX_ID
./cli $connectionConfig query $toCC swapGet $TX_ID
echo "==================  swapDone  ================================="
./cli $connectionConfig invoke $toCC swapDone $TX_ID $swapKey
echo "================== balance ================================="
./cli $connectionConfig query $fromCC balanceOf $userAddress
./cli $connectionConfig query $toCC allowedBalanceOf $userAddress $fromCC_uppercase
