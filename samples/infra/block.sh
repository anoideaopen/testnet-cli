#!/bin/bash

set -xe

channelName="acl"
blockId=1

./cli block $channelName $blockId peer0.testnet

ls *.block

blockFileName="./"$channelName"_"$blockId".block"
cp $blockFileName /test/
#cat $blockFileName
