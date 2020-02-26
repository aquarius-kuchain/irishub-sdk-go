package service

import (
	"errors"
	"fmt"

	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils/json"
)

var (
	_ sdk.Msg = MsgDefineService{}
	_ sdk.Msg = MsgBindService{}
	_ sdk.Msg = MsgUpdateServiceBinding{}
	_ sdk.Msg = MsgSetWithdrawAddress{}
	_ sdk.Msg = MsgDisableService{}
	_ sdk.Msg = MsgEnableService{}
	_ sdk.Msg = MsgRefundServiceDeposit{}
	_ sdk.Msg = MsgRequestService{}
	_ sdk.Msg = MsgRespondService{}
	_ sdk.Msg = MsgPauseRequestContext{}
	_ sdk.Msg = MsgStartRequestContext{}
	_ sdk.Msg = MsgKillRequestContext{}
	_ sdk.Msg = MsgUpdateRequestContext{}
	_ sdk.Msg = MsgWithdrawEarnedFees{}
	_ sdk.Msg = MsgWithdrawTax{}

	TagServiceName      = "service-name"
	TagProvider         = "provider"
	TagConsumer         = "consumer"
	TagRequestID        = "request-id"
	TagRespondService   = "respond_service"
	TagRequestContextID = "request-context-id"

	RequestContextIDLen = 32 // length of the request context ID in bytes

	cdc = sdk.NewAminoCodec()
)

func init() {
	RegisterCodec(cdc)
}

// MsgDefineService defines a message to define a service
type MsgDefineService struct {
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	Tags              []string       `json:"tags"`
	Author            sdk.AccAddress `json:"author"`
	AuthorDescription string         `json:"author_description"`
	Schemas           string         `json:"schemas"`
}

func (msg MsgDefineService) Type() string {
	return "define_service"
}

func (msg MsgDefineService) ValidateBasic() error {
	if len(msg.Author) == 0 {
		return errors.New("author missing")
	}

	if len(msg.Name) == 0 {
		return errors.New("author missing")
	}

	if len(msg.Schemas) == 0 {
		return errors.New("schemas missing")
	}

	return nil
}

func (msg MsgDefineService) GetSignBytes() []byte {
	if len(msg.Tags) == 0 {
		msg.Tags = nil
	}

	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

func (msg MsgDefineService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Author}
}

// MsgBindService defines a message to bind a service
type MsgBindService struct {
	ServiceName     string         `json:"service_name"`
	Provider        sdk.AccAddress `json:"provider"`
	Deposit         sdk.Coins      `json:"deposit"`
	Pricing         string         `json:"pricing"`
	WithdrawAddress sdk.AccAddress `json:"withdraw_address"`
}

func (msg MsgBindService) Type() string {
	return "bind_service"
}

func (msg MsgBindService) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("serviceName missing")
	}

	if len(msg.Pricing) == 0 {
		return errors.New("pricing missing")
	}
	return nil
}

func (msg MsgBindService) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

func (msg MsgBindService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

// MsgRequestService defines a message to request a service
type MsgRequestService struct {
	ServiceName       string           `json:"service_name"`
	Providers         []sdk.AccAddress `json:"providers"`
	Consumer          sdk.AccAddress   `json:"consumer"`
	Input             string           `json:"input"`
	ServiceFeeCap     sdk.Coins        `json:"service_fee_cap"`
	Timeout           int64            `json:"timeout"`
	SuperMode         bool             `json:"super_mode"`
	Repeated          bool             `json:"repeated"`
	RepeatedFrequency uint64           `json:"repeated_frequency"`
	RepeatedTotal     int64            `json:"repeated_total"`
}

func (msg MsgRequestService) Type() string {
	return "request_service"
}

func (msg MsgRequestService) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}
	if len(msg.Providers) == 0 {
		return errors.New("providers missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("serviceName missing")
	}

	if len(msg.Input) == 0 {
		return errors.New("input missing")
	}
	return nil
}

func (msg MsgRequestService) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

func (msg MsgRequestService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

// MsgRespondService defines a message to respond a service request
type MsgRespondService struct {
	RequestID string         `json:"request_id"`
	Provider  sdk.AccAddress `json:"provider"`
	Output    string         `json:"output"`
	Error     string         `json:"error"`
}

func (msg MsgRespondService) Type() string {
	return "respond_service"
}

func (msg MsgRespondService) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.Output) == 0 && len(msg.Error) == 0 {
		return errors.New("either output or error should be specified, but neither was provided")
	}

	if len(msg.Output) > 0 && len(msg.Error) > 0 {
		return errors.New("either output or error should be specified, but both were provided")
	}
	return nil
}

func (msg MsgRespondService) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

func (msg MsgRespondService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgUpdateServiceBinding defines a message to update a service binding
type MsgUpdateServiceBinding struct {
	ServiceName string         `json:"service_name"`
	Provider    sdk.AccAddress `json:"provider"`
	Deposit     sdk.Coins      `json:"deposit"`
	Pricing     string         `json:"pricing"`
}

// Type implements Msg.
func (msg MsgUpdateServiceBinding) Type() string { return "update_service_binding" }

// GetSignBytes implements Msg.
func (msg MsgUpdateServiceBinding) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgUpdateServiceBinding) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("service name missing")
	}

	if !msg.Deposit.Empty() {
		return errors.New(fmt.Sprintf("invalid deposit: %s", msg.Deposit))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgUpdateServiceBinding) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSetWithdrawAddress defines a message to set the withdrawal address for a service binding
type MsgSetWithdrawAddress struct {
	ServiceName     string         `json:"service_name"`
	Provider        sdk.AccAddress `json:"provider"`
	WithdrawAddress sdk.AccAddress `json:"withdraw_address"`
}

// Type implements Msg.
func (msg MsgSetWithdrawAddress) Type() string { return "set_withdraw_address" }

// GetSignBytes implements Msg.
func (msg MsgSetWithdrawAddress) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgSetWithdrawAddress) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("service name missing")
	}

	if len(msg.WithdrawAddress) == 0 {
		return errors.New("withdrawal address missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgSetWithdrawAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgDisableService defines a message to disable a service binding
type MsgDisableService struct {
	ServiceName string         `json:"service_name"`
	Provider    sdk.AccAddress `json:"provider"`
}

// Type implements Msg.
func (msg MsgDisableService) Type() string { return "disable_service" }

// GetSignBytes implements Msg.
func (msg MsgDisableService) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgDisableService) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("service name missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgDisableService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgEnableService defines a message to enable a service binding
type MsgEnableService struct {
	ServiceName string         `json:"service_name"`
	Provider    sdk.AccAddress `json:"provider"`
	Deposit     sdk.Coins      `json:"deposit"`
}

// Type implements Msg.
func (msg MsgEnableService) Type() string { return "enable_service" }

// GetSignBytes implements Msg.
func (msg MsgEnableService) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgEnableService) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("service name missing")
	}

	if !msg.Deposit.Empty() {
		return errors.New(fmt.Sprintf("invalid deposit: %s", msg.Deposit))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgEnableService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgRefundServiceDeposit defines a message to refund deposit from a service binding
type MsgRefundServiceDeposit struct {
	ServiceName string         `json:"service_name"`
	Provider    sdk.AccAddress `json:"provider"`
}

// Type implements Msg.
func (msg MsgRefundServiceDeposit) Type() string { return "refund_service_deposit" }

// GetSignBytes implements Msg.
func (msg MsgRefundServiceDeposit) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgRefundServiceDeposit) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("service name missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgRefundServiceDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgPauseRequestContext defines a message to suspend a request context
type MsgPauseRequestContext struct {
	RequestContextID []byte         `json:"request_context_id"`
	Consumer         sdk.AccAddress `json:"consumer"`
}

// Type implements Msg.
func (msg MsgPauseRequestContext) Type() string { return "pause_request_context" }

// GetSignBytes implements Msg.
func (msg MsgPauseRequestContext) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgPauseRequestContext) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}

	if len(msg.RequestContextID) != RequestContextIDLen {
		return errors.New(fmt.Sprintf("length of the request context ID must be %d in bytes", RequestContextIDLen))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgPauseRequestContext) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgStartRequestContext defines a message to resume a request context
type MsgStartRequestContext struct {
	RequestContextID []byte         `json:"request_context_id"`
	Consumer         sdk.AccAddress `json:"consumer"`
}

// Type implements Msg.
func (msg MsgStartRequestContext) Type() string { return "start_request_context" }

// GetSignBytes implements Msg.
func (msg MsgStartRequestContext) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgStartRequestContext) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}

	if len(msg.RequestContextID) != RequestContextIDLen {
		return errors.New(fmt.Sprintf("length of the request context ID must be %d in bytes", RequestContextIDLen))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgStartRequestContext) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgKillRequestContext defines a message to terminate a request context
type MsgKillRequestContext struct {
	RequestContextID []byte         `json:"request_context_id"`
	Consumer         sdk.AccAddress `json:"consumer"`
}

// Type implements Msg.
func (msg MsgKillRequestContext) Type() string { return "kill_request_context" }

// GetSignBytes implements Msg.
func (msg MsgKillRequestContext) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgKillRequestContext) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}

	if len(msg.RequestContextID) != RequestContextIDLen {
		return errors.New(fmt.Sprintf("length of the request context ID must be %d in bytes", RequestContextIDLen))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgKillRequestContext) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgUpdateRequestContext defines a message to update a request context
type MsgUpdateRequestContext struct {
	RequestContextID  []byte           `json:"request_context_id"`
	Providers         []sdk.AccAddress `json:"providers"`
	ServiceFeeCap     sdk.Coins        `json:"service_fee_cap"`
	Timeout           int64            `json:"timeout"`
	RepeatedFrequency uint64           `json:"repeated_frequency"`
	RepeatedTotal     int64            `json:"repeated_total"`
	Consumer          sdk.AccAddress   `json:"consumer"`
}

// Type implements Msg.
func (msg MsgUpdateRequestContext) Type() string { return "update_request_context" }

// GetSignBytes implements Msg.
func (msg MsgUpdateRequestContext) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgUpdateRequestContext) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}

	if len(msg.RequestContextID) != RequestContextIDLen {
		return errors.New(fmt.Sprintf("length of the request context ID must be %d in bytes", RequestContextIDLen))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgUpdateRequestContext) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgWithdrawEarnedFees defines a message to withdraw the fees earned by the provider
type MsgWithdrawEarnedFees struct {
	Provider sdk.AccAddress `json:"provider"`
}

// Type implements Msg.
func (msg MsgWithdrawEarnedFees) Type() string { return "withdraw_earned_fees" }

// GetSignBytes implements Msg.
func (msg MsgWithdrawEarnedFees) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgWithdrawEarnedFees) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgWithdrawEarnedFees) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgWithdrawTax defines a message to withdraw the service tax
type MsgWithdrawTax struct {
	Trustee     sdk.AccAddress `json:"trustee"`
	DestAddress sdk.AccAddress `json:"dest_address"`
	Amount      sdk.Coins      `json:"amount"`
}

// Type implements Msg.
func (msg MsgWithdrawTax) Type() string { return "withdraw_tax" }

// GetSignBytes implements Msg.
func (msg MsgWithdrawTax) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgWithdrawTax) ValidateBasic() error {
	if len(msg.Trustee) == 0 {
		return errors.New("trustee missing")
	}

	if len(msg.DestAddress) == 0 {
		return errors.New("destination address missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgWithdrawTax) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Trustee}
}

func RegisterCodec(cdc sdk.Codec) {
	cdc.RegisterConcrete(MsgDefineService{}, "irishub/service/MsgDefineService")
	cdc.RegisterConcrete(MsgBindService{}, "irishub/service/MsgBindService")
	cdc.RegisterConcrete(MsgUpdateServiceBinding{}, "irishub/service/MsgUpdateServiceBinding")
	cdc.RegisterConcrete(MsgSetWithdrawAddress{}, "irishub/service/MsgSetWithdrawAddress")
	cdc.RegisterConcrete(MsgDisableService{}, "irishub/service/MsgDisableService")
	cdc.RegisterConcrete(MsgEnableService{}, "irishub/service/MsgEnableService")
	cdc.RegisterConcrete(MsgRefundServiceDeposit{}, "irishub/service/MsgRefundServiceDeposit")
	cdc.RegisterConcrete(MsgRequestService{}, "irishub/service/MsgRequestService")
	cdc.RegisterConcrete(MsgRespondService{}, "irishub/service/MsgRespondService")
	cdc.RegisterConcrete(MsgPauseRequestContext{}, "irishub/service/MsgPauseRequestContext")
	cdc.RegisterConcrete(MsgStartRequestContext{}, "irishub/service/MsgStartRequestContext")
	cdc.RegisterConcrete(MsgKillRequestContext{}, "irishub/service/MsgKillRequestContext")
	cdc.RegisterConcrete(MsgUpdateRequestContext{}, "irishub/service/MsgUpdateRequestContext")
	cdc.RegisterConcrete(MsgWithdrawEarnedFees{}, "irishub/service/MsgWithdrawEarnedFees")
	cdc.RegisterConcrete(MsgWithdrawTax{}, "irishub/service/MsgWithdrawTax")
}