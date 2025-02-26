package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/line/lbm-sdk/client"
	codectypes "github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/genutil/types"
	slashing "github.com/line/lbm-sdk/x/slashing/types"
	staking "github.com/line/lbm-sdk/x/staking/types"
	cryptocodec "github.com/line/ostracon/crypto/encoding"
	octypes "github.com/line/ostracon/types"
	"github.com/pkg/errors"
)

type replacementConfigs []replacementConfig

func (r *replacementConfigs) isReplacedValidator(validatorAddress string) (int, replacementConfig) {
	for i, replacement := range *r {
		if replacement.ValidatorAddress == validatorAddress {
			return i, replacement
		}
	}

	return -1, replacementConfig{}
}

type replacementConfig struct {
	Name             string `json:"validator_name"`
	ValidatorAddress string `json:"validator_address"`
	ConsensusPubkey  string `json:"stargate_consensus_public_key"`
}

func loadKeydataFromFile(clientCtx client.Context, replacementrJSON string, genDoc *octypes.GenesisDoc) *octypes.GenesisDoc {
	jsonReplacementBlob, err := ioutil.ReadFile(replacementrJSON)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "failed to read replacement keys from file %s", replacementrJSON))
	}

	var replacementKeys replacementConfigs

	err = json.Unmarshal(jsonReplacementBlob, &replacementKeys)

	if err != nil {
		log.Fatal("Could not unmarshal replacement keys ")
	}

	var state types.AppMap
	if err := json.Unmarshal(genDoc.AppState, &state); err != nil {
		log.Fatal(errors.Wrap(err, "failed to JSON unmarshal initial genesis state"))
	}

	var stakingGenesis staking.GenesisState
	var slashingGenesis slashing.GenesisState

	clientCtx.JSONMarshaler.MustUnmarshalJSON(state[staking.ModuleName], &stakingGenesis)
	clientCtx.JSONMarshaler.MustUnmarshalJSON(state[slashing.ModuleName], &slashingGenesis)

	for i, val := range stakingGenesis.Validators {
		idx, replacement := replacementKeys.isReplacedValidator(val.OperatorAddress)

		if idx != -1 {
			toReplaceValConsAddress, _ := val.GetConsAddr() // nolint: errcheck

			consPubKey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, replacement.ConsensusPubkey)

			if err != nil {
				log.Fatal(fmt.Errorf("failed to decode key:%s %w", consPubKey, err))
			}

			val.ConsensusPubkey, err = codectypes.NewAnyWithValue(consPubKey)
			if err != nil {
				log.Fatal(fmt.Errorf("failed to decode key:%s %w", consPubKey, err))
			}

			replaceValConsAddress, _ := val.GetConsAddr()                                     // nolint: errcheck
			protoReplaceValConsPubKey, _ := val.TmConsPublicKey()                             // nolint: errcheck
			replaceValConsPubKey, _ := cryptocodec.PubKeyFromProto(protoReplaceValConsPubKey) // nolint: errcheck

			for i, signingInfo := range slashingGenesis.SigningInfos {
				if signingInfo.Address == toReplaceValConsAddress.String() {
					slashingGenesis.SigningInfos[i].Address = replaceValConsAddress.String()
					slashingGenesis.SigningInfos[i].ValidatorSigningInfo.Address = replaceValConsAddress.String()
				}
			}

			for i, missedInfo := range slashingGenesis.MissedBlocks {
				if missedInfo.Address == toReplaceValConsAddress.String() {
					slashingGenesis.MissedBlocks[i].Address = replaceValConsAddress.String()
				}
			}

			for tmIdx, tmval := range genDoc.Validators {
				if tmval.Address.String() == replaceValConsAddress.String() {
					genDoc.Validators[tmIdx].Address = replaceValConsAddress.Bytes()
					genDoc.Validators[tmIdx].PubKey = replaceValConsPubKey
				}
			}

			stakingGenesis.Validators[i] = val
		}
	}

	state[staking.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(&stakingGenesis)
	state[slashing.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(&slashingGenesis)

	genDoc.AppState, err = json.Marshal(state)

	if err != nil {
		log.Fatal("Could not marshal App State")
	}
	return genDoc
}
