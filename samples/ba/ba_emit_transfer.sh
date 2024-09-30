#!/bin/bash

set -xe

fromCC="ba"
fromCC_uppercase=${fromCC^^}

userPrivateKeyArg=" -s AkTSLPVKaYvKWHvb415YrJr6Vg2o3oundd6pqtDF8PVW1DqvjRycJ5PGAFxTvHAedk398tv2gdcgwmmr9hpBVVK4Le4E4"
userPrivateKeyArg=$(./cli privkey)
userPublicKey=$(./cli pubkey -s $userPrivateKeyArg)
userAddress=$(./cli address -s $userPrivateKeyArg)

issuerSKeyArg=" -s VLPf7EoxF7eTnABWx3pweTVGEbfTVAg4aUpeVhJ2dCgj4BtP11TmphZ49u7bfbiZh5Z5btUi5BMbxcpSrW395bwXgFxCC"
issuerPublicKey=$(./cli --organization org0 -f connection.yaml pubkey $issuerSKeyArg)
issuerAddress=$(./cli --organization org0 -f connection.yaml address $issuerSKeyArg)

#// otf usd
issuerSKeyArg=" -s FcjAE6ov6ooqQa3yzVke22ZmMVRrH2yXZyaAxBQBhNUqpC9EwyHGPThPE4rtp8pT5GZHjLdKx2AZ2eBBuT3Mz3sCpefK8"
issuerPublicKey=$(./cli --organization org0 -f connection.yaml pubkey $issuerSKeyArg)
issuerAddress=$(./cli --organization org0 -f connection.yaml address $issuerSKeyArg)

#// ba
issuerSKeyArg=" -s VLPf7EoxF7eTnABWx3pweTVGEbfTVAg4aUpeVhJ2dCgj4BtP11TmphZ49u7bfbiZh5Z5btUi5BMbxcpSrW395bwXgFxCC"
issuerPublicKey=$(./cli --organization org0 -f connection.yaml pubkey $issuerSKeyArg)
issuerAddress=$(./cli --organization org0 -f connection.yaml address $issuerSKeyArg)

echo "=================================================================="
echo "==================  acl addUser  ================================="
echo "=================================================================="
./cli --organization org0 -f connection.yaml invoke acl addUser $issuerPublicKey 123 testuser true | true
./cli --organization org0 -f connection.yaml invoke acl addUser $userPublicKey 123 testuser true | true

echo "=================================================================="
echo "==================  metadata  ================================="
echo "=================================================================="
./cli --organization org0 -f connection.yaml query $fromCC metadata

echo "=================================================================="
echo "==================  add money  ================================="
echo "=================================================================="
num=`date +%N`
a="A$num"
b="B$num"
c="C$num"
./cli --organization org0 -f connection.yaml $issuerSKeyArg invoke $fromCC emitTokensFromBars '{"bars":[{"id":"","group_name":"","serial":"'$a'","amountToTokenization":"1.111","underlying_asset":"gold","unit_of_measure":"oz","gross_weight":"400.512","fine_weight":"399.911","calc_method":"fine","refiner":"Valcambi SA - Suisse","delivery_form":"bar","custodian":"Moscow, Diamond Fund","year":"2019","price":"1000","note":"test1","is_hold":false},{"id":"","group_name":"","serial":"'$b'","amountToTokenization":"1.111","underlying_asset":"silver","unit_of_measure":"oz","gross_weight":"404.055","fine_weight":"403.449","calc_method":"fine","refiner":"Valcambi SA - Suisse","delivery_form":"bar","custodian":"","year":"","price":"100","note":"test1","is_hold":false},{"id":"","group_name":"","serial":"'$c'","amountToTokenization":"1.111","underlying_asset":"copper","unit_of_measure":"MT","gross_weight":"404.055","fine_weight":"403.449","calc_method":"gross","refiner":"ASAHI REF CANADA LTD","delivery_form":"plate","custodian":"","year":"","price":"500","note":"test1","is_hold":false}]}' '{"docs":[{"id":"Doc1","hash":"Hash1"},{"id":"Doc2","hash":"Hash2"}]}'
./cli --organization org0 -f connection.yaml $issuerSKeyArg invoke $fromCC emitTokensFromBars '{"bars":[]}' '{"docs":[{"id":"Doc1","hash":"Hash1"},{"id":"Doc2","hash":"Hash2"}]}'
sleep 5
./cli --organization org0 -f connection.yaml query $fromCC bAAllBalancesOf $issuerAddress
./cli --organization org0 -f connection.yaml query $fromCC bABalanceOf $issuerAddress $fromCC_uppercase"_"$a"goldbar.1"
./cli --organization org0 -f connection.yaml query $fromCC bABalanceOf $issuerAddress $fromCC_uppercase"_"$b"silverbar.1"
./cli --organization org0 -f connection.yaml query $fromCC bABalanceOf $issuerAddress $fromCC_uppercase"_"$c"copperplate.1"
./cli --organization org0 -f connection.yaml $issuerSKeyArg invoke $fromCC transfer $userAddress '{"bars":["'$fromCC_uppercase'_'$a'goldbar.1"]}' "transfer"
./cli --organization org0 -f connection.yaml $issuerSKeyArg invoke $fromCC transfer $userAddress '{"bars":["'$fromCC_uppercase'_'$b'silverbar.1"]}' "transfer"
./cli --organization org0 -f connection.yaml $issuerSKeyArg invoke $fromCC transfer $userAddress '{"bars":["'$fromCC_uppercase'_'$c'copperplate.1"]}' "transfer"
sleep 7
./cli --organization org0 -f connection.yaml query $fromCC bAAllBalancesOf $userAddress
./cli --organization org0 -f connection.yaml query $fromCC bABalanceOf $userAddress $fromCC_uppercase"_"$a"goldbar.1"
./cli --organization org0 -f connection.yaml query $fromCC bABalanceOf $userAddress $fromCC_uppercase"_"$b"silverbar.1"
./cli --organization org0 -f connection.yaml query $fromCC bABalanceOf $userAddress $fromCC_uppercase"_"$c"copperplate.1"

./cli --organization org0 -f connection.yaml $userPrivateKeyArg invoke $fromCC createRedeemRequest '{"bars":["BA_A728916511goldbar.1"]}' "Redeem test"

./cli --organization org0 -f connection.yaml $userPrivateKeyArg invoke $fromCC denyAllRedeemRequest

./cli --organization org0 -f connection.yaml invoke acl addMultisig A4JdE9iZRzU9NEiVDNxYKKWymHeBxHR7mA8AetFrg8m4 123 testuser true

export validatorPublicKeys="A4JdE9iZRzU9NEiVDNxYKKWymHeBxHR7mA8AetFrg8m4,\
5Tevazf8xxwyyKGku4VCCSVMDN56mU3mm2WsnENk1zv5,\
6qFz88dv2R8sXmyzWPjvzN6jafv7t1kNUHztYKjH1Rd4"


./cli --organization org0 -f connection.yaml invoke acl addUser A4JdE9iZRzU9NEiVDNxYKKWymHeBxHR7mA8AetFrg8m4 123 testuser true
./cli --organization org0 -f connection.yaml invoke acl addUser 5Tevazf8xxwyyKGku4VCCSVMDN56mU3mm2WsnENk1zv5 123 testuser true
./cli --organization org0 -f connection.yaml invoke acl addUser 6qFz88dv2R8sXmyzWPjvzN6jafv7t1kNUHztYKjH1Rd4 123 testuser true



# create a multisig signature

./cli generateMessage acl addMultisig $validatorPublicKeys 2
./cli -s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr signMessage
./cli -s 5D2BpuHZwik9zPFuaqba4zbvNP8TB7PQ6usZke5bufPbKf8xG6ZMHReBqwKw9aDfpTaNfaRsg1j2zVZWrX8hg18D signMessage
./cli -s 3sK2wHWxU58kzAeFtShDMsPm5Qh74NAWgfwCmdKyzvp4npivEDDEp14WgQpg7KGaVNF7qWyyMvkKPzGddVkxagNN signMessage
./cli --organization org0 -f connection.yaml sendRequest acl addMultisig $validatorPublicKeys

./cli --organization org0 -f connection.yaml $userPrivateKeyArg invoke $fromCC createRedeemRequest '{"bars":["BA_A728916511goldbar.1"]}' "Redeem test"

./cli --organization org0 -f connection.yaml query $fromCC redeemRequestsList
./cli generateMessage $fromCC denyAllRedeemRequest  $validatorPublicKeys
./cli -s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr signMessage
./cli -s 5D2BpuHZwik9zPFuaqba4zbvNP8TB7PQ6usZke5bufPbKf8xG6ZMHReBqwKw9aDfpTaNfaRsg1j2zVZWrX8hg18D signMessage
./cli -s 3sK2wHWxU58kzAeFtShDMsPm5Qh74NAWgfwCmdKyzvp4npivEDDEp14WgQpg7KGaVNF7qWyyMvkKPzGddVkxagNN signMessage
./cli --organization org0 -f connection.yaml sendRequest $fromCC denyAllRedeemRequest $validatorPublicKeys
./cli --organization org0 -f connection.yaml query $fromCC redeemRequestsList
