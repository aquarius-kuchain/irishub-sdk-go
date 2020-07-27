package distribution

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type distributionClient struct {
	sdk.BaseClient
	*log.Logger
}

func (d distributionClient) RegisterCodec(cdc sdk.Codec) {
	registerCodec(cdc)
}

func (d distributionClient) Name() string {
	return ModuleName
}

func Create(ac sdk.BaseClient) rpc.Distribution {
	return distributionClient{
		BaseClient: ac,
		Logger:     ac.Logger(),
	}
}

func (d distributionClient) QueryRewards(delegator string) (rpc.Rewards, sdk.Error) {
	address, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return rpc.Rewards{}, sdk.Wrap(err)
	}

	param := struct {
		DelegatorAddress sdk.AccAddress `json:"delegator_address"`
	}{
		DelegatorAddress: address,
	}

	var rewards rewardsResponse
	if err := d.QueryWithResponse("custom/distribution/delegator_total_rewards", param, &rewards); err != nil {
		return rpc.Rewards{}, sdk.Wrap(err)
	}
	return rewards.Convert().(rpc.Rewards), nil
}

func (d distributionClient) QueryWithdrawAddr(delegator string) (string, sdk.Error) {
	address, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return "", sdk.Wrap(err)
	}

	param := struct {
		DelegatorAddress sdk.AccAddress `json:"delegator_address"`
	}{
		DelegatorAddress: address,
	}

	res, newErr := d.Query("custom/distribution/withdraw_addr", param)
	if newErr != nil {
		return "", sdk.Wrap(err)
	}
	return string(res), nil
}

func (d distributionClient) QueryCommission(validator string) (rpc.ValidatorAccumulatedCommission, sdk.Error) {
	address, err := sdk.ValAddressFromBech32(validator)
	if err != nil {
		return rpc.ValidatorAccumulatedCommission{}, sdk.Wrap(err)
	}

	param := struct {
		ValidatorAddress sdk.ValAddress `json:"validator_address"`
	}{
		ValidatorAddress: address,
	}

	var commission validatorAccumulatedCommission
	if err := d.QueryWithResponse("custom/distribution/validator_commission", param, &commission); err != nil {
		return rpc.ValidatorAccumulatedCommission{}, sdk.Wrap(err)
	}

	return commission.Convert().(rpc.ValidatorAccumulatedCommission), nil
}

func (d distributionClient) SetWithdrawAddr(withdrawAddr string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	delegator, err := d.QueryAddress(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	withdraw, err := sdk.AccAddressFromBech32(withdrawAddr)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := MsgSetWithdrawAddress{
		DelegatorAddr: delegator,
		WithdrawAddr:  withdraw,
	}
	d.Info().Str("delegator", delegator.String()).
		Str("withdrawAddr", withdrawAddr).
		Msg("execute setWithdrawAddr transaction")
	return d.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (d distributionClient) WithdrawRewards(isValidator bool, onlyFromValidator string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	delegator, err := d.QueryAddress(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	var msgs []sdk.Msg
	switch {
	case isValidator:
		msgs = append(msgs, MsgWithdrawValidatorRewardsAll{
			ValidatorAddr: sdk.ValAddress(delegator.Bytes()),
		})

		d.Info().Str("delegator", delegator.String()).
			Msg("execute withdrawValidatorRewardsAll transaction")
		break
	case onlyFromValidator != "":
		valAddr, err := sdk.ValAddressFromBech32(onlyFromValidator)
		if err != nil {
			return sdk.ResultTx{}, sdk.Wrap(err)
		}
		msgs = append(msgs, MsgWithdrawDelegatorReward{
			ValidatorAddr: valAddr,
			DelegatorAddr: delegator,
		})

		d.Info().Str("delegator", delegator.String()).
			Str("validator", onlyFromValidator).
			Msg("execute withdrawDelegatorReward transaction")
		break
	default:
		msgs = append(msgs, MsgWithdrawDelegatorRewardsAll{
			DelegatorAddr: delegator,
		})

		d.Info().Str("delegator", delegator.String()).
			Str("validator", onlyFromValidator).
			Msg("execute withdrawDelegatorRewardsAll transaction")
		break
	}
	return d.BuildAndSend(msgs, baseTx)
}
