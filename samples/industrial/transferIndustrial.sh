#!/bin/bash

set -xe

itOwnerKey=" -s VLPf7EoxF7eTnABWx3pweTVGEbfTVAg4aUpeVhJ2dCgj4BtP11TmphZ49u7bfbiZh5Z5btUi5BMbxcpSrW395bwXgFxCC"
user1Key=" -s AkTSLPVKaYvKWHvb415YrJr6Vg2o3oundd6pqtDF8PVW1DqvjRycJ5PGAFxTvHAedk398tv2gdcgwmmr9hpBVVK4Le4E4"
user2Key=" -s Bi39jBAHakJ3yvvSyVui5otU2U8j8dpj6Gf8pK1LN5oV6k14zyupThWQTSodhmtCTpQp6TdqdcnhQHMUafxPBXu98n4Ff"

itOwnerPublicKey=$(./cli -f connection.yaml pubkey $itOwnerKey)
user1PublicKey=$(./cli -f connection.yaml pubkey $user1Key)
user2PublicKey=$(./cli -f connection.yaml pubkey $user2Key)

itOwnerAddress=$(./cli -f connection.yaml address $itOwnerKey)
user1Address=$(./cli -f connection.yaml address $user1Key)
user2Address=$(./cli -f connection.yaml address $user2Key)

echo "============================================================================"
echo "==================  check installed tokens  ================================="
echo "============================================================================"
./cli -f connection.yaml query gf28iln060 metadata
#./cli -f connection.yaml query otf metadata unknown method

## Initialize token groups

#./cli -f connection.yaml invoke gf28iln060 initialize --noBatch $itOwnerKey

./cli -f connection.yaml query gf28iln060 industrialBalanceOf $itOwnerAddress


./cli -f connection.yaml invoke acl addUser $user1PublicKey 123 testuser true
./cli -f connection.yaml invoke acl addUser $user2PublicKey 123 testuser true
./cli -f connection.yaml invoke acl addUser $itOwnerPublicKey 123 testuser true

./cli -f connection.yaml invoke gf28iln060 transferIndustrial $user2Address "28022021" "75000000000" "for lunch" $itOwnerKey
sleep 5
./cli -f connection.yaml query gf28iln060 industrialBalanceOf $user1Address
./cli -f connection.yaml query gf28iln060 industrialBalanceOf $user2Address

./cli -f connection.yaml invoke gf28iln060 transferIndustrial $user1Address "31122020" "1" "for lunch" $user2Key
sleep 5
./cli -f connection.yaml query gf28iln060 industrialBalanceOf $user1Address
./cli -f connection.yaml query gf28iln060 industrialBalanceOf $user2Address

REDEEM_REQUEST_3=$(./cli -f connection.yaml invoke gf28iln060 createRedeemRequest "28022021" "100000000" "ref3" -r tx $user2Key)
REDEEM_REQUEST_2=$(./cli -f connection.yaml invoke gf28iln060 createRedeemRequest "31122020" "2" "ref2" -r tx $user1Key)

echo "REDEEM_REQUEST_1"
echo $REDEEM_REQUEST_1
echo "REDEEM_REQUEST_2"
echo $REDEEM_REQUEST_2

./cli -f connection.yaml query gf28iln060 industrialBalanceOf $user1Address

./cli -f connection.yaml query gf28iln060 redeemRequestsList

./cli -f connection.yaml invoke gf28iln060 denyRedeemRequest $REDEEM_REQUEST_1 $user1Key
./cli -f connection.yaml invoke gf28iln060 acceptRedeemRequest $REDEEM_REQUEST_3 "1" "acceptRef3" $user2Key

./cli -f connection.yaml query gf28iln060 industrialBalanceOf $user1Address
./cli -f connection.yaml query gf28iln060 industrialBalanceOf $user2Address

./cli -f connection.yaml query gf28iln060 metadata
