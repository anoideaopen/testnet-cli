#!/bin/bash

set -xe

fromCC="ba02"
toCC="ba"

fromCC_uppercase=${fromCC^^}
toCC_uppercase=${toCC^^}

issuerSKeyArg=" -s 4d14d1db55d4612ac2e0d2369e75120b1b46dde5b9a3b8ede4231f63e37d3ca30f531508888a7bbf2c6a14df441a14401d8192df3f7c93b05d08115298565189"
issuerPublicKey=$(./cli pubkey $issuerSKeyArg)
issuerAddress=$(./cli address $issuerSKeyArg)
echo "$issuerAddress"
echo "$issuerPublicKey"

userSecretKey="AkTSLPVKaYvKWHvb415YrJr6Vg2o3oundd6pqtDF8PVW1DqvjRycJ5PGAFxTvHAedk398tv2gdcgwmmr9hpBVVK4Le4E4"
#userSecretKey=$(./cli privkey)
userPrivateKeyArg=" -s $userSecretKey"
userPublicKey=$(./cli pubkey $userPrivateKeyArg)
userAddress=$(./cli address $userPrivateKeyArg)

swapHash="7d4e3eec80026719639ed4dba68916eb94c7a49a053e05c8f9578fe4e5a3d7ea"
swapKey="12345"

echo "=================================================================="
echo "==================  acl addUser  ================================="
echo "=================================================================="
./cli invoke acl addUser $issuerPublicKey 123 testuser true | true
./cli invoke acl addUser $userPublicKey 123 testuser true | true

echo "=================================================================="
echo "==================  acl checkKeys  ================================="
echo "=================================================================="
./cli query acl checkKeys $userPublicKey

echo "=================================================================="
echo "==================  add money  ================================="
echo "=================================================================="
num=`date +%N`
a="A$num"
b="B$num"
c="C$num"
./cli $issuerSKeyArg invoke $fromCC emitTokensFromBars '{"bars":[{"id":"","group_name":"","serial":"'$a'","amountToTokenization":"1.111","underlying_asset":"gold","unit_of_measure":"oz","gross_weight":"400.512","fine_weight":"399.911","calc_method":"fine","refiner":"Valcambi SA - Suisse","delivery_form":"bar","custodian":"Moscow, Diamond Fund","year":"2019","price":"1000","note":"test1","is_hold":false},{"id":"","group_name":"","serial":"'$b'","amountToTokenization":"1.111","underlying_asset":"silver","unit_of_measure":"oz","gross_weight":"404.055","fine_weight":"403.449","calc_method":"fine","refiner":"Valcambi SA - Suisse","delivery_form":"bar","custodian":"","year":"","price":"100","note":"test1","is_hold":false},{"id":"","group_name":"","serial":"'$c'","amountToTokenization":"1.111","underlying_asset":"copper","unit_of_measure":"MT","gross_weight":"404.055","fine_weight":"403.449","calc_method":"gross","refiner":"ASAHI REF CANADA LTD","delivery_form":"plate","custodian":"","year":"","price":"500","note":"test1","is_hold":false}]}' '{"docs":[{"id":"Doc1","hash":"Hash1"},{"id":"Doc2","hash":"Hash2"}]}'
sleep 5
./cli query $fromCC bAAllBalancesOf $issuerAddress
./cli query $fromCC bABalanceOf $issuerAddress $fromCC_uppercase"_"$a"goldbar.1"
./cli query $fromCC bABalanceOf $issuerAddress $fromCC_uppercase"_"$b"silverbar.1"
./cli query $fromCC bABalanceOf $issuerAddress $fromCC_uppercase"_"$c"copperplate.1"
./cli $issuerSKeyArg invoke $fromCC transfer $userAddress '{"bars":["'$fromCC_uppercase'_'$a'goldbar.1"]}' "transfer"
./cli $issuerSKeyArg invoke $fromCC transfer $userAddress '{"bars":["'$fromCC_uppercase'_'$b'silverbar.1"]}' "transfer"
./cli $issuerSKeyArg invoke $fromCC transfer $userAddress '{"bars":["'$fromCC_uppercase'_'$c'copperplate.1"]}' "transfer"
sleep 7
./cli query $fromCC bAAllBalancesOf $userAddress
./cli query $fromCC bABalanceOf $userAddress $fromCC_uppercase"_"$a"goldbar.1"
./cli query $fromCC bABalanceOf $userAddress $fromCC_uppercase"_"$b"silverbar.1"
./cli query $fromCC bABalanceOf $userAddress $fromCC_uppercase"_"$c"copperplate.1"

echo "=================================================================="
echo "==================  multiSwapBegin from balance to allowedBalance ================================="
echo "=================================================================="
TX_ID=`./cli $userPrivateKeyArg -r tx invoke $fromCC multiSwapBegin $fromCC_uppercase '{"Assets":[{"group":"'$fromCC_uppercase'_'$a'goldbar.1","amount":"1"},{"group":"'$fromCC_uppercase'_'$b'silverbar.1","amount":"1"}]}' $toCC_uppercase "$swapHash"`
echo $TX_ID
sleep 10
echo "=================================================================="
echo "==================  multiSwapGet  ================================="
echo "=================================================================="
./cli query $fromCC multiSwapGet $TX_ID
./cli query $toCC multiSwapGet $TX_ID
echo "=================================================================="
echo "==================  multiSwapDone  ================================="
echo "=================================================================="
./cli invoke $toCC multiSwapDone $TX_ID $swapKey
echo "=================================================================="
echo "================== balance $fromCC_uppercase ================================="
echo "=================================================================="
./cli query $fromCC bAAllBalancesOf $userAddress
./cli query $fromCC bABalanceOf $userAddress $fromCC_uppercase"_"$a"goldbar.1"
./cli query $fromCC bABalanceOf $userAddress $fromCC_uppercase"_"$b"silverbar.1"
./cli query $fromCC bABalanceOf $userAddress $fromCC_uppercase"_"$c"copperplate.1"
echo "=================================================================="
echo "================== balance $toCC_uppercase ================================="
echo "=================================================================="
./cli query $toCC bAAllBalancesOf $userAddress
./cli query $toCC bABalanceOf $userAddress $fromCC_uppercase"_"$a"goldbar.1"
./cli query $toCC bABalanceOf $userAddress $fromCC_uppercase"_"$b"silverbar.1"
./cli query $toCC bABalanceOf $userAddress $fromCC_uppercase"_"$c"copperplate.1"
./cli query $toCC allowedBalanceOf $userAddress $fromCC_uppercase"_"$a"goldbar.1"
./cli query $toCC allowedBalanceOf $userAddress $fromCC_uppercase"_"$b"silverbar.1"
./cli query $toCC allowedBalanceOf $userAddress $fromCC_uppercase"_"$c"copperplate.1"

echo "=================================================================="
echo "==================  multiSwapBegin from allowedBalance to balance ================================="
echo "=================================================================="
TX_ID_BACK=`./cli $userPrivateKeyArg -r tx invoke $toCC multiSwapBegin $toCC_uppercase '{"Assets":[{"group":"'$fromCC_uppercase'_'$a'goldbar.1","amount":"1"},{"group":"'$fromCC_uppercase'_'$b'silverbar.1","amount":"1"}]}' $fromCC_uppercase "$swapHash"`
echo $TX_ID_BACK
sleep 10
echo "=================================================================="
echo "==================  multiSwapGet  ================================="
echo "=================================================================="
./cli query $toCC multiSwapGet $TX_ID_BACK
./cli query $fromCC multiSwapGet $TX_ID_BACK
echo "=================================================================="
echo "==================  multiSwapDone  ================================="
echo "=================================================================="
./cli invoke $fromCC multiSwapDone $TX_ID_BACK $swapKey
echo "=================================================================="
echo "================== balance $fromCC_uppercase ================================="
echo "=================================================================="
./cli query $fromCC bAAllBalancesOf $userAddress
./cli query $fromCC bABalanceOf $userAddress $fromCC_uppercase"_"$a"goldbar.1"
./cli query $fromCC bABalanceOf $userAddress $fromCC_uppercase"_"$b"silverbar.1"
./cli query $fromCC bABalanceOf $userAddress $fromCC_uppercase"_"$c"copperplate.1"
echo "=================================================================="
echo "================== balance $toCC_uppercase ================================="
echo "=================================================================="
./cli query $toCC bAAllBalancesOf $userAddress
./cli query $toCC bABalanceOf $userAddress $fromCC_uppercase"_"$a"goldbar.1"
./cli query $toCC bABalanceOf $userAddress $fromCC_uppercase"_"$b"silverbar.1"
./cli query $toCC bABalanceOf $userAddress $fromCC_uppercase"_"$c"copperplate.1"
./cli query $toCC allowedBalanceOf $userAddress $fromCC_uppercase"_"$a"goldbar.1"
./cli query $toCC allowedBalanceOf $userAddress $fromCC_uppercase"_"$b"silverbar.1"
./cli query $toCC allowedBalanceOf $userAddress $fromCC_uppercase"_"$c"copperplate.1"
