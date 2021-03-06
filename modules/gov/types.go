package gov

import (
	json2 "encoding/json"
	"errors"
	"fmt"

	"github.com/irisnet/irishub-sdk-go/rpc"

	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils/json"
)

const (
	ModuleName = "gov"
)

var (
	_ sdk.Msg = MsgDeposit{}
	_ sdk.Msg = MsgVote{}

	cdc = sdk.NewAminoCodec()
)

func init() {
	registerCodecForProposal(cdc)
}

// Type that represents Proposal Type as a byte
type ProposalKind byte
type UsageType byte

type Param struct {
	Subspace string `json:"subspace"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

type Params []Param

//-----------------------------------------------------------
// MsgSubmitProposal
type MsgSubmitProposal struct {
	Title          string         `json:"title"`           //  Title of the proposal
	Description    string         `json:"description"`     //  Description of the proposal
	ProposalType   ProposalKind   `json:"proposal_type"`   //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Proposer       sdk.AccAddress `json:"proposer"`        //  Address of the proposer
	InitialDeposit sdk.Coins      `json:"initial_deposit"` //  Initial deposit paid by sender. Must be strictly positive.
	Params         Params         `json:"params"`
}

//nolint
func (msg MsgSubmitProposal) Type() string { return "submit_proposal" }

// Implements Msg.
func (msg MsgSubmitProposal) ValidateBasic() error {
	return nil
}

// Implements Msg.
func (msg MsgSubmitProposal) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// Implements Msg.
func (msg MsgSubmitProposal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proposer}
}

type MsgSubmitSoftwareUpgradeProposal struct {
	MsgSubmitProposal
	Version      uint64  `json:"version"`
	Software     string  `json:"software"`
	SwitchHeight uint64  `json:"switch_height"`
	Threshold    sdk.Dec `json:"threshold"`
}

func (msg MsgSubmitSoftwareUpgradeProposal) ValidateBasic() error {
	return nil
}

func (msg MsgSubmitSoftwareUpgradeProposal) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

type MsgSubmitCommunityTaxUsageProposal struct {
	MsgSubmitProposal
	Usage       UsageType      `json:"usage"`
	DestAddress sdk.AccAddress `json:"dest_address"`
	Amount      sdk.Coins      `json:"amount"`
}

func (msg MsgSubmitCommunityTaxUsageProposal) ValidateBasic() error {
	return nil
}

func (msg MsgSubmitCommunityTaxUsageProposal) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

//-----------------------------------------------------------
// MsgDeposit
type MsgDeposit struct {
	ProposalID uint64         `json:"proposal_id"` // ID of the proposal
	Depositor  sdk.AccAddress `json:"depositor"`   // Address of the depositor
	Amount     sdk.Coins      `json:"amount"`      // Coins to add to the proposal's deposit
}

func (msg MsgDeposit) Route() string { return ModuleName }

// Implements Msg.
// nolint
func (msg MsgDeposit) Type() string { return "deposit" }

// Implements Msg.
func (msg MsgDeposit) ValidateBasic() error {
	if len(msg.Depositor) == 0 {
		return errors.New("depositor is empty")
	}
	if msg.ProposalID < 0 {
		return errors.New("invalid proposalID")
	}
	return nil
}

// Implements Msg.
func (msg MsgDeposit) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// Implements Msg.
func (msg MsgDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Depositor}
}

//-----------------------------------------------------------
// MsgVote
type MsgVote struct {
	ProposalID uint64         `json:"proposal_id"` // ID of the proposal
	Voter      sdk.AccAddress `json:"voter"`       //  address of the voter
	Option     VoteOption     `json:"option"`      //  option from OptionSet chosen by the voter
}

func (msg MsgVote) Route() string { return ModuleName }

// Implements Msg.
// nolint
func (msg MsgVote) Type() string { return "vote" }

// Implements Msg.
func (msg MsgVote) ValidateBasic() error {
	if len(msg.Voter) == 0 {
		return errors.New("voter is empty")
	}
	if msg.ProposalID < 0 {
		return errors.New("invalid proposalID")
	}
	if msg.Option != OptionYes &&
		msg.Option != OptionNo &&
		msg.Option != OptionNoWithVeto &&
		msg.Option != OptionAbstain {
		return errors.New("invalid option")
	}
	return nil
}

// Implements Msg.
func (msg MsgVote) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// Implements Msg.
func (msg MsgVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Voter}
}

// Type that represents VoteOption as a byte
type VoteOption byte

//nolint
const (
	OptionEmpty      VoteOption = 0x00
	OptionYes        VoteOption = 0x01
	OptionAbstain    VoteOption = 0x02
	OptionNo         VoteOption = 0x03
	OptionNoWithVeto VoteOption = 0x04
)

// String to proposalType byte.  Returns ff if invalid.
func VoteOptionFromString(option rpc.VoteOption) (VoteOption, sdk.Error) {
	switch option {
	case rpc.Yes:
		return OptionYes, nil
	case rpc.Abstain:
		return OptionAbstain, nil
	case rpc.No:
		return OptionNo, nil
	case rpc.NoWithVeto:
		return OptionNoWithVeto, nil
	default:
		return OptionEmpty, sdk.Wrapf(fmt.Sprintf("'%s' is not a valid vote option", option))
	}
}

// Marshals to JSON using string
func (vo VoteOption) MarshalJSON() ([]byte, error) {
	return json2.Marshal(vo.String())
}

// Turns VoteOption byte to String
func (vo VoteOption) String() string {
	switch vo {
	case OptionYes:
		return "Yes"
	case OptionAbstain:
		return "Abstain"
	case OptionNo:
		return "No"
	case OptionNoWithVeto:
		return "NoWithVeto"
	default:
		return ""
	}
}

// Tally Results
type tallyResult struct {
	Yes               string `json:"yes"`
	Abstain           string `json:"abstain"`
	No                string `json:"no"`
	NoWithVeto        string `json:"no_with_veto"`
	SystemVotingPower string `json:"system_voting_power"`
}

func (t tallyResult) Convert() interface{} {
	return rpc.TallyResult{
		Yes:               t.Yes,
		Abstain:           t.Abstain,
		No:                t.No,
		NoWithVeto:        t.NoWithVeto,
		SystemVotingPower: t.SystemVotingPower,
	}
}

//for query
type vote struct {
	Voter      sdk.AccAddress `json:"voter"`       //  address of the voter
	ProposalID uint64         `json:"proposal_id"` //  proposalID of the proposal
	Option     string         `json:"option"`      //  option from OptionSet chosen by the voter
}

func (v vote) Convert() interface{} {
	return rpc.Vote{
		Voter:      v.Voter.String(),
		ProposalID: v.ProposalID,
		Option:     v.Option,
	}
}

type votes []vote

func (vs votes) Convert() interface{} {
	votes := make([]rpc.Vote, len(vs))
	for _, v := range vs {
		votes = append(votes, v.Convert().(rpc.Vote))
	}
	return votes
}

// deposit
type deposit struct {
	Depositor  sdk.AccAddress `json:"depositor"`   //  Address of the depositor
	ProposalID uint64         `json:"proposal_id"` //  proposalID of the proposal
	Amount     sdk.Coins      `json:"amount"`      //  deposit amount
}

func (d deposit) Convert() interface{} {
	return rpc.Deposit{
		Depositor:  d.Depositor.String(),
		ProposalID: d.ProposalID,
		Amount:     d.Amount,
	}
}

type deposits []deposit

func (ds deposits) Convert() interface{} {
	deposits := make([]rpc.Deposit, len(ds))
	for _, d := range ds {
		deposits = append(deposits, d.Convert().(rpc.Deposit))
	}
	return deposits
}

func registerCodec(cdc sdk.Codec) {
	cdc.RegisterConcrete(MsgSubmitProposal{}, "irishub/gov/MsgSubmitProposal")
	cdc.RegisterConcrete(MsgSubmitCommunityTaxUsageProposal{}, "irishub/gov/MsgSubmitCommunityTaxUsageProposal")
	cdc.RegisterConcrete(MsgSubmitSoftwareUpgradeProposal{}, "irishub/gov/MsgSubmitSoftwareUpgradeProposal")
	cdc.RegisterConcrete(MsgDeposit{}, "irishub/gov/MsgDeposit")
	cdc.RegisterConcrete(MsgVote{}, "irishub/gov/MsgVote")

	registerCodecForProposal(cdc)
	cdc.RegisterConcrete(&vote{}, "irishub/gov/Vote")
}

func registerCodecForProposal(cdc sdk.Codec) {
	cdc.RegisterInterface((*proposal)(nil))
	cdc.RegisterConcrete(&BasicProposal{}, "irishub/gov/BasicProposal")
	cdc.RegisterConcrete(&parameterProposal{}, "irishub/gov/ParameterProposal")
	cdc.RegisterConcrete(&plainTextProposal{}, "irishub/gov/PlainTextProposal")
	cdc.RegisterConcrete(&softwareUpgradeProposal{}, "irishub/gov/SoftwareUpgradeProposal")
	cdc.RegisterConcrete(&communityTaxUsageProposal{}, "irishub/gov/CommunityTaxUsageProposal")
}
