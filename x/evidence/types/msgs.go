package types

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	"github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/evidence/exported"
)

// Message types for the evidence module
const (
	TypeMsgSubmitEvidence = "submit_evidence"
)

var (
	_ sdk.Msg                       = &MsgSubmitEvidence{}
	_ types.UnpackInterfacesMessage = MsgSubmitEvidence{}
	_ exported.MsgSubmitEvidenceI   = &MsgSubmitEvidence{}
)

// NewMsgSubmitEvidence returns a new MsgSubmitEvidence with a signer/submitter.
//nolint:interfacer
func NewMsgSubmitEvidence(s sdk.AccAddress, evi exported.Evidence) (*MsgSubmitEvidence, error) {
	msg, ok := evi.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("cannot proto marshal %T", evi)
	}
	any, err := types.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}
	return &MsgSubmitEvidence{Submitter: s.String(), Evidence: any}, nil
}

// Route returns the MsgSubmitEvidence's route.
func (m MsgSubmitEvidence) Route() string { return RouterKey }

// Type returns the MsgSubmitEvidence's type.
func (m MsgSubmitEvidence) Type() string { return TypeMsgSubmitEvidence }

// ValidateBasic performs basic (non-state-dependant) validation on a MsgSubmitEvidence.
func (m MsgSubmitEvidence) ValidateBasic() error {
	if m.Submitter == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Submitter)
	}
	if err := sdk.ValidateAccAddress(m.Submitter); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid submitter address (%s)", err)
	}

	evi := m.GetEvidence()
	if evi == nil {
		return sdkerrors.Wrap(ErrInvalidEvidence, "missing evidence")
	}
	if err := evi.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

// GetSignBytes returns the raw bytes a signer is expected to sign when submitting
// a MsgSubmitEvidence message.
func (m MsgSubmitEvidence) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners returns the single expected signer for a MsgSubmitEvidence.
func (m MsgSubmitEvidence) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(m.Submitter)}
}

func (m MsgSubmitEvidence) GetEvidence() exported.Evidence {
	evi, ok := m.Evidence.GetCachedValue().(exported.Evidence)
	if !ok {
		return nil
	}
	return evi
}

func (m MsgSubmitEvidence) GetSubmitter() sdk.AccAddress {
	return sdk.AccAddress(m.Submitter)
}

func (m MsgSubmitEvidence) UnpackInterfaces(ctx types.AnyUnpacker) error {
	var evi exported.Evidence
	return ctx.UnpackAny(m.Evidence, &evi)
}
