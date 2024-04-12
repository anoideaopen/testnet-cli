#!/bin/bash

fromCC="rub"
toCC="hermitage"

fromCC_uppercase=${fromCC^^}
toCC_uppercase=${toCC^^}

userPrivateKeyArg=" -s f92d1d8dc55645d94da7c2cf11fbf40303bae3f33d06f534eb92acce0b0ac40f84a9fbce96a6f7af760effb15b5a14f51b8f08517f9549ff84dc9a52921f4798"
userPublicKey=$(./cli pubkey $userPrivateKeyArg)
userAddress=$(./cli address $userPrivateKeyArg)

swapHash="7d4e3eec80026719639ed4dba68916eb94c7a49a053e05c8f9578fe4e5a3d7ea"
swapKey="12345"

echo "==================  swapBegin  ================================="
TX_ID=`./cli $connection_file $userPrivateKeyArg -r tx invoke $fromCC swapBegin $fromCC_uppercase $toCC_uppercase 100100 "$swapHash"`
echo $TX_ID
sleep 10
echo "==================  swapGet  ================================="
./cli $connection_file query $fromCC swapGet $TX_ID
./cli $connection_file query $toCC swapGet $TX_ID
echo "==================  swapDone  ================================="
./cli $connection_file invoke $toCC swapDone $TX_ID $swapKey
echo "================== balance ================================="
./cli $connection_file query $fromCC balanceOf $userAddress
./cli $connection_file query $toCC allowedBalanceOf $userAddress $fromCC_uppercase
