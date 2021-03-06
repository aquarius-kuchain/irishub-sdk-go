package gov

import (
	"time"

	"github.com/irisnet/irishub-sdk-go/rpc"

	sdk "github.com/irisnet/irishub-sdk-go/types"
)

var (
	_ proposal = (*BasicProposal)(nil)
	_ proposal = (*plainTextProposal)(nil)
	_ proposal = (*parameterProposal)(nil)
	_ proposal = (*communityTaxUsageProposal)(nil)
	_ proposal = (*softwareUpgradeProposal)(nil)
)

// proposal interface
type proposal interface {
	GetProposalID() uint64
	GetTitle() string
	GetDescription() string
	GetProposalType() string
	GetStatus() string
	GetTallyResult() tallyResult
	GetSubmitTime() time.Time
	GetDepositEndTime() time.Time
	GetTotalDeposit() sdk.Coins
	GetVotingStartTime() time.Time
	GetVotingEndTime() time.Time
	GetProposer() sdk.AccAddress
	sdk.Response
}

type proposals []proposal

func (ps proposals) Convert() interface{} {
	//should not implement
	return nil
}

// Basic proposals
type BasicProposal struct {
	ProposalID      uint64         `json:"proposal_id"`       //  ID of the proposal
	Title           string         `json:"title"`             //  Title of the proposal
	Description     string         `json:"description"`       //  Description of the proposal
	ProposalType    string         `json:"proposal_type"`     //  Type of proposal. Initial set {plainTextProposal, softwareUpgradeProposal}
	Status          string         `json:"proposal_status"`   //  Status of the proposal {Pending, Active, Passed, Rejected}
	TallyResult     tallyResult    `json:"tally_result"`      //  Result of Tallys
	SubmitTime      time.Time      `json:"submit_time"`       //  Time of the block where TxGovSubmitProposal was included
	DepositEndTime  time.Time      `json:"deposit_end_time"`  // Time that the proposal would expire if deposit amount isn't met
	TotalDeposit    sdk.Coins      `json:"total_deposit"`     //  Current deposit on this proposal. Initial value is set at InitialDeposit
	VotingStartTime time.Time      `json:"voting_start_time"` //  Time of the block where MinDeposit was reached. -1 if MinDeposit is not reached
	VotingEndTime   time.Time      `json:"voting_end_time"`   // Time that the VotingPeriod for this proposal will end and votes will be tallied
	Proposer        sdk.AccAddress `json:"proposer"`
}

func (b BasicProposal) GetTitle() string {
	return b.Title
}

func (b BasicProposal) GetDescription() string {
	return b.Description
}

func (b BasicProposal) GetProposalType() string {
	return b.ProposalType
}

func (b BasicProposal) GetStatus() string {
	return b.Status
}

func (b BasicProposal) GetTallyResult() tallyResult {
	return b.TallyResult
}

func (b BasicProposal) GetSubmitTime() time.Time {
	return b.SubmitTime
}

func (b BasicProposal) GetDepositEndTime() time.Time {
	return b.DepositEndTime
}

func (b BasicProposal) GetTotalDeposit() sdk.Coins {
	return b.TotalDeposit
}

func (b BasicProposal) GetVotingStartTime() time.Time {
	return b.VotingStartTime
}

func (b BasicProposal) GetVotingEndTime() time.Time {
	return b.VotingEndTime
}

func (b BasicProposal) GetProposer() sdk.AccAddress {
	return b.Proposer
}

func (b BasicProposal) Convert() interface{} {
	return rpc.BasicProposal{
		Title:          b.Title,
		Description:    b.Description,
		ProposalID:     b.ProposalID,
		ProposalStatus: b.Status,
		ProposalType:   b.ProposalType,
		TallyResult: rpc.TallyResult{
			Yes:               b.TallyResult.Yes,
			Abstain:           b.TallyResult.Abstain,
			No:                b.TallyResult.No,
			NoWithVeto:        b.TallyResult.NoWithVeto,
			SystemVotingPower: b.TallyResult.SystemVotingPower,
		},
		SubmitTime:      b.SubmitTime,
		DepositEndTime:  b.DepositEndTime,
		TotalDeposit:    b.TotalDeposit,
		VotingStartTime: b.VotingStartTime,
		VotingEndTime:   b.VotingEndTime,
		Proposer:        b.Proposer.String(),
	}
}

func (b BasicProposal) GetProposalID() uint64 {
	return b.ProposalID
}

type plainTextProposal struct {
	BasicProposal
}

func (b plainTextProposal) Convert() interface{} {
	return rpc.PlainTextProposal{
		Proposal: b.BasicProposal.Convert().(rpc.BasicProposal),
	}
}

type param struct {
	Subspace string `json:"subspace"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

type params []param

// Implements proposal Interface
type parameterProposal struct {
	BasicProposal
	Params params `json:"params"`
}

func (b parameterProposal) Convert() interface{} {
	var params []rpc.Param
	for _, p := range b.Params {
		params = append(params, rpc.Param{
			Subspace: "", //TODO
			Key:      p.Key,
			SubKey:   "", //TODO
			Value:    p.Value,
		})
	}
	return rpc.ParameterProposal{
		Proposal: b.BasicProposal.Convert().(rpc.BasicProposal),
		Params:   params,
	}
}

// Implements proposal Interface
type taxUsage struct {
	Usage       string         `json:"usage"`
	DestAddress sdk.AccAddress `json:"dest_address"`
	Percent     string         `json:"percent"`
	Amount      sdk.Coins      `json:"amount"`
}

type communityTaxUsageProposal struct {
	BasicProposal
	TaxUsage taxUsage `json:"tax_usage"`
}

func (b communityTaxUsageProposal) Convert() interface{} {
	return rpc.CommunityTaxUsageProposal{
		Proposal: b.BasicProposal.Convert().(rpc.BasicProposal),
		TaxUsage: rpc.TaxUsage{
			Usage:       b.TaxUsage.Usage,
			DestAddress: b.TaxUsage.DestAddress.String(),
			Percent:     b.TaxUsage.Percent,
		},
	}
}

type softwareUpgradeProposal struct {
	BasicProposal
	ProtocolDefinition protocolDefinition `json:"protocol_definition"`
}

type protocolDefinition struct {
	Version   uint64 `json:"version"`
	Software  string `json:"software"`
	Height    uint64 `json:"height"`
	Threshold string `json:"threshold"`
}

func (b softwareUpgradeProposal) Convert() interface{} {
	return rpc.SoftwareUpgradeProposal{
		Proposal: b.BasicProposal.Convert().(rpc.BasicProposal),
		ProtocolDefinition: rpc.ProtocolDefinition{
			Version:   b.ProtocolDefinition.Version,
			Software:  b.ProtocolDefinition.Software,
			Height:    b.ProtocolDefinition.Height,
			Threshold: b.ProtocolDefinition.Threshold,
		},
	}
}
