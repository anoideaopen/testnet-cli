#!/bin/bash

set -xe

itOwnerKey=" -s VLPf7EoxF7eTnABWx3pweTVGEbfTVAg4aUpeVhJ2dCgj4BtP11TmphZ49u7bfbiZh5Z5btUi5BMbxcpSrW395bwXgFxCC"
user1Key=" -s AkTSLPVKaYvKWHvb415YrJr6Vg2o3oundd6pqtDF8PVW1DqvjRycJ5PGAFxTvHAedk398tv2gdcgwmmr9hpBVVK4Le4E4"
user2Key=" -s Bi39jBAHakJ3yvvSyVui5otU2U8j8dpj6Gf8pK1LN5oV6k14zyupThWQTSodhmtCTpQp6TdqdcnhQHMUafxPBXu98n4Ff"

itOwnerPublicKey=$(./cli pubkey $itOwnerKey)
user1PublicKey=$(./cli pubkey $user1Key)
user2PublicKey=$(./cli pubkey $user2Key)

itOwnerAddress=$(./cli address $itOwnerKey)
user1Address=$(./cli address $user1Key)
user2Address=$(./cli address $user2Key)

echo "============================================================================"
echo "==================  check installed tokens  ================================="
echo "============================================================================"
./cli query ittest metadata
#./cli query otf metadata unknown method

## Initialize token groups

./cli invoke ittest initialize --noBatch $itOwnerKey

./cli query ittest industrialBalanceOf $itOwnerAddress
#
./cli invoke ittest transferIndustrial $user2Address "30112025" "10" "for lunch" $itOwnerKey
sleep 5
./cli query ittest industrialBalanceOf $user1Address
./cli query ittest industrialBalanceOf $user2Address

./cli invoke ittest transferIndustrial $user1Address "31122025" "1" "for lunch" $user2Key
sleep 5
./cli query ittest industrialBalanceOf $user1Address
./cli query ittest industrialBalanceOf $user2Address

REDEEM_REQUEST_1=$(./cli invoke ittest createRedeemRequest "30112025" "1" "ref1" -r tx $user1Key)
REDEEM_REQUEST_2=$(./cli invoke ittest createRedeemRequest "30112025" "2" "ref2" -r tx $user1Key)

echo "REDEEM_REQUEST_1"
echo $REDEEM_REQUEST_1
echo "REDEEM_REQUEST_2"
echo $REDEEM_REQUEST_2

./cli query ittest industrialBalanceOf $user1Address

./cli query ittest redeemRequestsList

./cli invoke ittest denyRedeemRequest $REDEEM_REQUEST_1 $user1Key
./cli invoke ittest acceptRedeemRequest $REDEEM_REQUEST_2 "1" "acceptRef2" $user1Key

./cli query ittest industrialBalanceOf $user1Address
./cli query ittest industrialBalanceOf $user2Address

./cli query ittest metadata
