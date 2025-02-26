package types

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"

	"github.com/gogo/protobuf/proto"

	"github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

// Governance message types and routes
const (
	TypeMsgDeposit        = "deposit"
	TypeMsgVote           = "vote"
	TypeMsgVoteWeighted   = "weighted_vote"
	TypeMsgSubmitProposal = "submit_proposal"
)

var (
	_, _, _ sdk.Msg                       = &MsgSubmitProposal{}, &MsgDeposit{}, &MsgVoteWeighted{}
	_       types.UnpackInterfacesMessage = &MsgSubmitProposal{}
)

// NewMsgSubmitProposal creates a new MsgSubmitProposal.
//nolint:interfacer
func NewMsgSubmitProposal(content Content, initialDeposit sdk.Coins, proposer sdk.AccAddress) (*MsgSubmitProposal, error) {
	m := &MsgSubmitProposal{
		InitialDeposit: initialDeposit,
		Proposer:       proposer.String(),
	}
	err := m.SetContent(content)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *MsgSubmitProposal) GetInitialDeposit() sdk.Coins { return m.InitialDeposit }

func (m *MsgSubmitProposal) GetProposer() sdk.AccAddress {
	return sdk.AccAddress(m.Proposer)
}

func (m *MsgSubmitProposal) GetContent() Content {
	content, ok := m.Content.GetCachedValue().(Content)
	if !ok {
		return nil
	}
	return content
}

func (m *MsgSubmitProposal) SetInitialDeposit(coins sdk.Coins) {
	m.InitialDeposit = coins
}

func (m *MsgSubmitProposal) SetProposer(address fmt.Stringer) {
	m.Proposer = address.String()
}

func (m *MsgSubmitProposal) SetContent(content Content) error {
	msg, ok := content.(proto.Message)
	if !ok {
		return fmt.Errorf("can't proto marshal %T", msg)
	}
	any, err := types.NewAnyWithValue(msg)
	if err != nil {
		return err
	}
	m.Content = any
	return nil
}

// Route implements Msg
func (m MsgSubmitProposal) Route() string { return RouterKey }

// Type implements Msg
func (m MsgSubmitProposal) Type() string { return TypeMsgSubmitProposal }

// ValidateBasic implements Msg
func (m MsgSubmitProposal) ValidateBasic() error {
	if m.Proposer == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Proposer)
	}
	if !m.InitialDeposit.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.InitialDeposit.String())
	}
	if m.InitialDeposit.IsAnyNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.InitialDeposit.String())
	}
	if err := sdk.ValidateAccAddress(m.Proposer); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid proposer address (%s)", err)
	}

	content := m.GetContent()
	if content == nil {
		return sdkerrors.Wrap(ErrInvalidProposalContent, "missing content")
	}
	if !IsValidProposalType(content.ProposalType()) {
		return sdkerrors.Wrap(ErrInvalidProposalType, content.ProposalType())
	}
	if err := content.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

// GetSignBytes implements Msg
func (m MsgSubmitProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (m MsgSubmitProposal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(m.Proposer)}
}

// String implements the Stringer interface
func (m MsgSubmitProposal) String() string {
	out, _ := yaml.Marshal(m)
	return string(out)
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m MsgSubmitProposal) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var content Content
	return unpacker.UnpackAny(m.Content, &content)
}

// NewMsgDeposit creates a new MsgDeposit instance
//nolint:interfacer
func NewMsgDeposit(depositor sdk.AccAddress, proposalID uint64, amount sdk.Coins) *MsgDeposit {
	return &MsgDeposit{proposalID, depositor.String(), amount}
}

// Route implements Msg
func (msg MsgDeposit) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgDeposit) Type() string { return TypeMsgDeposit }

// ValidateBasic implements Msg
func (msg MsgDeposit) ValidateBasic() error {
	if msg.Depositor == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Depositor)
	}
	if err := sdk.ValidateAccAddress(msg.Depositor); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid depositor address (%s)", err)
	}
	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}
	if msg.Amount.IsAnyNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}

// String implements the Stringer interface
func (msg MsgDeposit) String() string {
	out, _ := yaml.Marshal(msg)
	return string(out)
}

// GetSignBytes implements Msg
func (msg MsgDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Depositor)}
}

// NewMsgVote creates a message to cast a vote on an active proposal
//nolint:interfacer
func NewMsgVote(voter sdk.AccAddress, proposalID uint64, option VoteOption) *MsgVote {
	return &MsgVote{proposalID, voter.String(), option}
}

// Route implements Msg
func (msg MsgVote) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgVote) Type() string { return TypeMsgVote }

// ValidateBasic implements Msg
func (msg MsgVote) ValidateBasic() error {
	if msg.Voter == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Voter)
	}
	if err := sdk.ValidateAccAddress(msg.Voter); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid voter address (%s)", err)
	}
	if !ValidVoteOption(msg.Option) {
		return sdkerrors.Wrap(ErrInvalidVote, msg.Option.String())
	}

	return nil
}

// String implements the Stringer interface
func (msg MsgVote) String() string {
	out, _ := yaml.Marshal(msg)
	return string(out)
}

// GetSignBytes implements Msg
func (msg MsgVote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Voter)}
}

// NewMsgVoteWeighted creates a message to cast a vote on an active proposal
//nolint:interfacer
func NewMsgVoteWeighted(voter sdk.AccAddress, proposalID uint64, options WeightedVoteOptions) *MsgVoteWeighted {
	return &MsgVoteWeighted{proposalID, voter.String(), options}
}

// Route implements Msg
func (msg MsgVoteWeighted) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgVoteWeighted) Type() string { return TypeMsgVoteWeighted }

// ValidateBasic implements Msg
func (msg MsgVoteWeighted) ValidateBasic() error {
	if msg.Voter == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Voter)
	}

	if len(msg.Options) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, WeightedVoteOptions(msg.Options).String())
	}

	totalWeight := sdk.NewDec(0)
	usedOptions := make(map[VoteOption]bool)
	for _, option := range msg.Options {
		if !ValidWeightedVoteOption(option) {
			return sdkerrors.Wrap(ErrInvalidVote, option.String())
		}
		totalWeight = totalWeight.Add(option.Weight)
		if usedOptions[option.Option] {
			return sdkerrors.Wrap(ErrInvalidVote, "Duplicated vote option")
		}
		usedOptions[option.Option] = true
	}

	if totalWeight.GT(sdk.NewDec(1)) {
		return sdkerrors.Wrap(ErrInvalidVote, "Total weight overflow 1.00")
	}

	if totalWeight.LT(sdk.NewDec(1)) {
		return sdkerrors.Wrap(ErrInvalidVote, "Total weight lower than 1.00")
	}

	return nil
}

// String implements the Stringer interface
func (msg MsgVoteWeighted) String() string {
	out, _ := yaml.Marshal(msg)
	return string(out)
}

// GetSignBytes implements Msg
func (msg MsgVoteWeighted) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgVoteWeighted) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Voter)}
}
