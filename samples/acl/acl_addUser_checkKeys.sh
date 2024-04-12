#!/bin/bash

set -xe

userSecretKey=$(./cli privkey)
userPrivateKeyArg=" -s $userSecretKey"

userPublicKey=$(./cli pubkey $userPrivateKeyArg)

echo "=================================================================="
echo "==================  acl addUser  ================================="
echo "=================================================================="
./cli invoke acl addUser $userPublicKey 123 testuser true

echo "=================================================================="
echo "==================  acl checkKeys  ================================="
echo "=================================================================="
./cli query acl checkKeys $userPublicKey