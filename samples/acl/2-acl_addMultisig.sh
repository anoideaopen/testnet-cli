#!/bin/bash

#set -xe

export validatorPublicKeys="A4JdE9iZRzU9NEiVDNxYKKWymHeBxHR7mA8AetFrg8m4,\
5Tevazf8xxwyyKGku4VCCSVMDN56mU3mm2WsnENk1zv5,\
6qFz88dv2R8sXmyzWPjvzN6jafv7t1kNUHztYKjH1Rd4"


validatorPrivateKey1Arg="-s 3aDebSkgXq37VPrzThboaV8oMMbYXrRAt7hnGrod4PNMnGfXjh14TY7cQs8eVT46C4RK4ZyNKLrBmyD5CYZiFmkr"
validatorPrivateKey2Arg="-s 5D2BpuHZwik9zPFuaqba4zbvNP8TB7PQ6usZke5bufPbKf8xG6ZMHReBqwKw9aDfpTaNfaRsg1j2zVZWrX8hg18D"
validatorPrivateKey3Arg="-s 3sK2wHWxU58kzAeFtShDMsPm5Qh74NAWgfwCmdKyzvp4npivEDDEp14WgQpg7KGaVNF7qWyyMvkKPzGddVkxagNN"

validatorPublicKey1=$(./cli pubkey $validatorPrivateKey1Arg)
validatorPublicKey2=$(./cli pubkey $validatorPrivateKey2Arg)
validatorPublicKey3=$(./cli pubkey $validatorPrivateKey3Arg)

echo "=================================================================="
echo "==================  acl addUser  ================================="
echo "=================================================================="
./cli --connection connection.yaml invoke acl addUser $validatorPublicKey1 123 testuser true | true
./cli --connection connection.yaml invoke acl addUser $validatorPublicKey2 123 testuser true | true
./cli --connection connection.yaml invoke acl addUser $validatorPublicKey3 123 testuser true | true


echo "=================================================================="
echo "==================  acl addMultisig  ================================="
echo "=================================================================="
./cli generateMessage acl addMultisig $validatorPublicKeys "3"

./cli $validatorPrivateKey1Arg signMessage
./cli $validatorPrivateKey2Arg signMessage
./cli $validatorPrivateKey3Arg signMessage

./cli --connection connection.yaml sendRequest acl addMultisig $validatorPublicKeys
