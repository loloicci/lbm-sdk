package legacytx

import (
	"github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
)

var _ types.UnpackInterfacesMessage = StdSignMsg{}

// StdSignMsg is a convenience structure for passing along a Msg with the other
// requirements for a StdSignDoc before it is signed. For use in the CLI.
type StdSignMsg struct {
	ChainID        string    `json:"chain_id" yaml:"chain_id"`
	SigBlockHeight uint64    `json:"sign_block_height" yaml:"sign_block_height"`
	Sequence       uint64    `json:"sequence" yaml:"sequence"`
	TimeoutHeight  uint64    `json:"timeout_height" yaml:"timeout_height"`
	Fee            StdFee    `json:"fee" yaml:"fee"`
	Msgs           []sdk.Msg `json:"msgs" yaml:"msgs"`
	Memo           string    `json:"memo" yaml:"memo"`
}

// get message bytes
func (msg StdSignMsg) Bytes() []byte {
	return StdSignBytes(msg.ChainID, msg.SigBlockHeight, msg.Sequence, msg.TimeoutHeight, msg.Fee, msg.Msgs, msg.Memo)
}

func (msg StdSignMsg) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	for _, m := range msg.Msgs {
		err := types.UnpackInterfaces(m, unpacker)
		if err != nil {
			return err
		}
	}

	return nil
}
