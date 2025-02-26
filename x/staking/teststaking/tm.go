package teststaking

import (
	occrypto "github.com/line/ostracon/crypto"
	octypes "github.com/line/ostracon/types"

	cryptocodec "github.com/line/lbm-sdk/crypto/codec"
	"github.com/line/lbm-sdk/x/staking/types"
)

// GetOcConsPubKey gets the validator's public key as an occrypto.PubKey.
func GetOcConsPubKey(v types.Validator) (occrypto.PubKey, error) {
	pk, err := v.ConsPubKey()
	if err != nil {
		return nil, err
	}

	return cryptocodec.ToOcPubKeyInterface(pk)
}

// ToOcValidator casts an SDK validator to an ostracon type Validator.
func ToOcValidator(v types.Validator) (*octypes.Validator, error) {
	tmPk, err := GetOcConsPubKey(v)
	if err != nil {
		return nil, err
	}

	return octypes.NewValidator(tmPk, v.ConsensusPower()), nil
}

// ToOcValidators casts all validators to the corresponding tendermint type.
func ToOcValidators(v types.Validators) ([]*octypes.Validator, error) {
	validators := make([]*octypes.Validator, len(v))
	var err error
	for i, val := range v {
		validators[i], err = ToOcValidator(val)
		if err != nil {
			return nil, err
		}
	}

	return validators, nil
}
