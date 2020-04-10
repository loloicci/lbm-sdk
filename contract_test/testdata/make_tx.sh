#!/usr/bin/env bash
source "./contract_test/testdata/common.sh"

SEND_TX_HASH=$(./build/linkcli --keyring-backend=test tx send jack $(./build/linkcli --keyring-backend=test --home /tmp/contract_test/.linkcli keys show alice -a) 100link --home /tmp/contract_test/.linkcli --chain-id lcd --yes  | awk '/txhash.*/{print $2}')
echo ${SEND_TX_HASH}
