#!/bin/bash

#set -xe

userSecretKey=$(./cli privkey)
userPrivateKeyArg=" -s $userSecretKey"

userPublicKey=$(./cli pubkey $userPrivateKeyArg)

userAddress=$(./cli address $userPrivateKeyArg)

echo "=================================================================="
echo "==================  acl addUser  ================================="
echo "=================================================================="
./cli invoke acl addUser $userPublicKey 123 testuser true

echo "=================================================================="
echo "==================  acl checkKeys  ================================="
echo "=================================================================="
./cli query acl checkKeys $userPublicKey


export validatorPublicKeys="A4JdE9iZRzU9NEiVDNxYKKWymHeBxHR7mA8AetFrg8m4,\
5Tevazf8xxwyyKGku4VCCSVMDN56mU3mm2WsnENk1zv5,\
6qFz88dv2R8sXmyzWPjvzN6jafv7t1kNUHztYKjH1Rd4"


userNewSecretKey=$(./cli privkey)
userNewPublicKey=$(./cli pubkey $userPrivateKeyArg)

## change public key
export changedAddr="$userAddress"
export reason="lost_key"
export reasonId="2"
export newPkey="$userNewPublicKey"

./cli generateMessage acl changePublicKey $validatorPublicKeys $changedAddr $reason $reasonId $newPkey

./cli -s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr signMessage
./cli -s 5D2BpuHZwik9zPFuaqba4zbvNP8TB7PQ6usZke5bufPbKf8xG6ZMHReBqwKw9aDfpTaNfaRsg1j2zVZWrX8hg18D signMessage
./cli -s 3sK2wHWxU58kzAeFtShDMsPm5Qh74NAWgfwCmdKyzvp4npivEDDEp14WgQpg7KGaVNF7qWyyMvkKPzGddVkxagNN signMessage

./cli sendRequest acl changePublicKey $validatorPublicKeys

./cli query acl checkKeys $userNewPublicKey
./cli query acl checkKeys $userPublicKey
