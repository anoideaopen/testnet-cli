#!/bin/bash

set -xe

userSecretKey=$(./cli privkey)
userPublicKey=$(./cli pubkey -s $userSecretKey)

echo "=================================================================="
echo "==================  acl addUser  ================================="
echo "=================================================================="
./cli --connection connection.yaml invoke acl addUser $userPublicKey 123 testuser true

echo "=================================================================="
echo "==================  acl checkKeys  ================================="
echo "=================================================================="
./cli --connection connection.yaml query acl checkKeys $userPublicKey


#++ ./cli privkey
#+ userSecretKey=FyRLfVyA2QEjy2CwnNoN5eAZKhDxfG5nznrajLowNScNWPtL4WtSa2SHFyE2pbVi9sscKwSnpUbew8VCT8PNQ8SV6iwQh
#+ userPrivateKeyArg=' -s FyRLfVyA2QEjy2CwnNoN5eAZKhDxfG5nznrajLowNScNWPtL4WtSa2SHFyE2pbVi9sscKwSnpUbew8VCT8PNQ8SV6iwQh'
#++ ./cli pubkey -s FyRLfVyA2QEjy2CwnNoN5eAZKhDxfG5nznrajLowNScNWPtL4WtSa2SHFyE2pbVi9sscKwSnpUbew8VCT8PNQ8SV6iwQh
#+ userPublicKey=DN38uNW8QtChsZqorkq3KpKUuE9BpqgE3NgPcwKTJbGu

#./cli --connection connection.yaml invoke acl addUser DN38uNW8QtChsZqorkq3KpKUuE9BpqgE3NgPcwKTJbGu 123 testuser true
#./cli --connection connection.yaml query acl checkKeys DN38uNW8QtChsZqorkq3KpKUuE9BpqgE3NgPcwKTJbGu