#!/bin/bash

set -xe

export config_cli="./bh-dev/cli.yaml"

export userPrivateKeyArg=$(./cli privkey)
export userPublicKey=$(./cli pubkey $userPrivateKeyArg)
export userAddress=$(./cli address $userPrivateKeyArg)

export userID=$(uuidgen)

./cli --config $config_cli \
  invoke acl addUser $userPublicKey 123 $(userID) true