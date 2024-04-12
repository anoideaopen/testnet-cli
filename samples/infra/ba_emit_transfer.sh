#!/bin/bash

set -xe

fromCC="ba02"
fromCC_uppercase=${fromCC^^}

userPrivateKeyArg=" -s AkTSLPVKaYvKWHvb415YrJr6Vg2o3oundd6pqtDF8PVW1DqvjRycJ5PGAFxTvHAedk398tv2gdcgwmmr9hpBVVK4Le4E4"
userPublicKey=$(./cli pubkey $userPrivateKeyArg)
userAddress=$(./cli address $userPrivateKeyArg)

issuerSKeyArg=" -s 4d14d1db55d4612ac2e0d2369e75120b1b46dde5b9a3b8ede4231f63e37d3ca30f531508888a7bbf2c6a14df441a14401d8192df3f7c93b05d08115298565189"
issuerPublicKey=$(./cli pubkey $issuerSKeyArg)
issuerAddress=$(./cli address $issuerSKeyArg)

echo "=================================================================="
echo "==================  acl addUser  ================================="
echo "=================================================================="
./cli invoke acl addUser $issuerPublicKey 123 testuser true | true
./cli invoke acl addUser $userPublicKey 123 testuser true | true

echo "=================================================================="
echo "==================  metadata  ================================="
echo "=================================================================="
./cli query $fromCC metadata

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
