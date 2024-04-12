#!/bin/bash
#set -xe


user1Key=$(./cli privkey)
user1KeyArg=" -s $user1Key"

user1PublicKey=$(./cli pubkey $user1KeyArg)

user1Address=$(./cli address $user1KeyArg)

echo "=================================================================="
echo "==================  acl addUser  ================================="
echo "=================================================================="
#./cli invoke acl addUser $itOwnerPublicKey 123 testuser true
./cli invoke acl addUser $user1PublicKey 123 testuser true

echo "=================================================================="
echo "==================  acl checkKeys  ================================="
echo "=================================================================="
#./cli query acl checkKeys $itOwnerPublicKey
./cli query acl checkKeys $user1PublicKey

echo $user1Address
./cli query acl getAddresses 100 ""
