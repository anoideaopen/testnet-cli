#!/bin/bash

set -xe


issuerSKeyArg=" -s FcjAE6ov6ooqQa3yzVke22ZmMVRrH2yXZyaAxBQBhNUqpC9EwyHGPThPE4rtp8pT5GZHjLdKx2AZ2eBBuT3Mz3sCpefK8"
issuerPublicKey=$(./cli  pubkey $issuerSKeyArg)
issuerAddress=$(./cli address $issuerSKeyArg)

# add issuer
./cli  $connection_file invoke acl addUser $issuerPublicKey 123 testuser true

export investorAddress="U1ErYgqxePVF9P5dQTH9ZL4fYDoJE3s3fpvoi7GUZQTwhwR2d"

# emit
./cli  $connection_file $issuerSKeyArg invoke rub emit $investorAddress 100000

# balanceOf
./cli $connection_file query rub balanceOf $investorAddress
