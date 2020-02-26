package service

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type serviceClient struct {
	sdk.AbstractClient
}

func New(ac sdk.AbstractClient) sdk.Service {
	return serviceClient{
		AbstractClient: ac,
	}
}

//DefineService is responsible for creating a new service definition
func (s serviceClient) DefineService(request sdk.ServiceDefinitionRequest) (sdk.Result, error) {
	author, err := s.QueryAddress(request.From, request.Password)
	if err != nil {
		return nil, err
	}
	msg := MsgDefineService{
		Name:              request.ServiceName,
		Description:       request.Description,
		Tags:              request.Tags,
		Author:            author,
		AuthorDescription: request.AuthorDescription,
		Schemas:           request.Schemas,
	}
	return s.Broadcast(request.BaseTx, []sdk.Msg{msg})
}

//BindService is responsible for binding a new service definition
func (s serviceClient) BindService(request sdk.ServiceBindingRequest) (sdk.Result, error) {
	provider, err := s.QueryAddress(request.From, request.Password)
	if err != nil {
		return nil, err
	}

	var withdrawAddress sdk.AccAddress
	if len(request.WithdrawAddr) > 0 {
		withdrawAddress, err = sdk.AccAddressFromBech32(request.WithdrawAddr)
		if err != nil {
			return nil, err
		}
	}

	msg := MsgBindService{
		ServiceName:     request.ServiceName,
		Provider:        provider,
		Deposit:         request.Deposit,
		Pricing:         request.Pricing,
		WithdrawAddress: withdrawAddress,
	}
	return s.Broadcast(request.BaseTx, []sdk.Msg{msg})
}

//UpdateServiceBinding updates the specified service binding
func (s serviceClient) UpdateServiceBinding(request sdk.UpdateServiceBindingRequest) (sdk.Result, error) {
	provider, err := s.QueryAddress(request.From, request.Password)
	if err != nil {
		return nil, err
	}
	msg := MsgUpdateServiceBinding{
		ServiceName: request.ServiceName,
		Provider:    provider,
		Deposit:     request.Deposit,
		Pricing:     request.Pricing,
	}
	return s.Broadcast(request.BaseTx, []sdk.Msg{msg})
}

// DisableService disables the specified service binding
func (s serviceClient) DisableService(serviceName string, baseTx sdk.BaseTx) (sdk.Result, error) {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}
	msg := MsgDisableService{
		ServiceName: serviceName,
		Provider:    provider,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

// EnableService enables the specified service binding
func (s serviceClient) EnableService(serviceName string, deposit sdk.Coins, baseTx sdk.BaseTx) (sdk.Result, error) {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}
	msg := MsgEnableService{
		ServiceName: serviceName,
		Provider:    provider,
		Deposit:     deposit,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

//InvokeService is responsible for invoke a new service and callback `callback`
func (s serviceClient) InvokeService(request sdk.ServiceInvocationRequest,
	callback sdk.ServiceInvokeHandler) (string, error) {
	consumer, err := s.QueryAddress(request.From, request.Password)
	if err != nil {
		return "", err
	}

	var providers []sdk.AccAddress
	for _, provider := range request.Providers {
		provider := sdk.MustAccAddressFromBech32(provider)
		providers = append(providers, provider)
	}

	msg := MsgRequestService{
		ServiceName:       request.ServiceName,
		Providers:         providers,
		Consumer:          consumer,
		Input:             request.Input,
		ServiceFeeCap:     request.ServiceFeeCap,
		Timeout:           request.Timeout,
		SuperMode:         request.SuperMode,
		Repeated:          request.Repeated,
		RepeatedFrequency: request.RepeatedFrequency,
		RepeatedTotal:     request.RepeatedTotal,
	}

	result, err := s.Broadcast(request.BaseTx, []sdk.Msg{msg})
	if err != nil {
		return "", err
	}

	requestContextID := result.GetTags().GetValue(TagRequestContextID)
	builder := sdk.NewEventQueryBuilder().
		AddCondition(sdk.ActionKey, sdk.EventValue(TagRespondService)).
		AddCondition(sdk.EventKey(TagConsumer), sdk.EventValue(consumer.String())).
		AddCondition(sdk.EventKey(TagServiceName), sdk.EventValue(request.ServiceName))

	var subscription sdk.Subscription
	subscription, err = s.SubscribeTx(builder, func(tx sdk.EventDataTx) {
		for _, msg := range tx.Tx.Msgs {
			msg, ok := msg.(MsgRespondService)
			if ok {
				callback(requestContextID, msg.Output)
			}
		}
		//cancel subscription
		if !request.Repeated {
			_ = s.Unscribe(subscription)
		}
	})
	if err != nil {
		return "", err
	}

	return requestContextID, nil
}

// SetWithdrawAddress sets a new withdrawal address for the specified service binding
func (s serviceClient) SetWithdrawAddress(serviceName string, withdrawAddress string, baseTx sdk.BaseTx) (sdk.Result, error) {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}
	withdrawAddr := sdk.MustAccAddressFromBech32(withdrawAddress)
	msg := MsgSetWithdrawAddress{
		ServiceName:     serviceName,
		Provider:        provider,
		WithdrawAddress: withdrawAddr,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

// RefundServiceDeposit refunds the deposit from the specified service binding
func (s serviceClient) RefundServiceDeposit(serviceName string, baseTx sdk.BaseTx) (sdk.Result, error) {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}
	msg := MsgRefundServiceDeposit{
		ServiceName: serviceName,
		Provider:    provider,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

// StartRequestContext starts the specified request context
func (s serviceClient) StartRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.Result, error) {
	consumer, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}
	msg := MsgStartRequestContext{
		RequestContextID: []byte(requestContextID),
		Consumer:         consumer,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

// PauseRequestContext suspends the specified request context
func (s serviceClient) PauseRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.Result, error) {
	consumer, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}
	msg := MsgPauseRequestContext{
		RequestContextID: []byte(requestContextID),
		Consumer:         consumer,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

// KillRequestContext terminates the specified request context
func (s serviceClient) KillRequestContext(requestContextID string, baseTx sdk.BaseTx) (sdk.Result, error) {
	consumer, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}
	msg := MsgKillRequestContext{
		RequestContextID: []byte(requestContextID),
		Consumer:         consumer,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

// UpdateRequestContext updates the specified request context
func (s serviceClient) UpdateRequestContext(request sdk.UpdateContextRequest) (sdk.Result, error) {
	consumer, err := s.QueryAddress(request.From, request.Password)
	if err != nil {
		return nil, err
	}

	var providers []sdk.AccAddress
	for _, p := range request.Providers {
		provider := sdk.MustAccAddressFromBech32(p)
		providers = append(providers, provider)
	}

	msg := MsgUpdateRequestContext{
		RequestContextID:  []byte(request.RequestContextID),
		Providers:         providers,
		ServiceFeeCap:     request.ServiceFeeCap,
		Timeout:           request.Timeout,
		RepeatedFrequency: request.RepeatedFrequency,
		RepeatedTotal:     request.RepeatedTotal,
		Consumer:          consumer,
	}
	return s.Broadcast(request.BaseTx, []sdk.Msg{msg})
}

// WithdrawEarnedFees withdraws the earned fees to the specified provider
func (s serviceClient) WithdrawEarnedFees(baseTx sdk.BaseTx) (sdk.Result, error) {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}

	msg := MsgWithdrawEarnedFees{
		Provider: provider,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

// WithdrawTax withdraws the service tax to the speicified destination address by the trustee
func (s serviceClient) WithdrawTax(destAddress string, amount sdk.Coins, baseTx sdk.BaseTx) (sdk.Result, error) {
	trustee, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}

	receipt := sdk.MustAccAddressFromBech32(destAddress)
	msg := MsgWithdrawTax{
		Trustee:     trustee,
		DestAddress: receipt,
		Amount:      amount,
	}
	return s.Broadcast(baseTx, []sdk.Msg{msg})
}

//RegisterInvocationListener is responsible for registering a group of service handler
func (s serviceClient) RegisterInvocationListener(serviceRouter sdk.ServiceRouter, baseTx sdk.BaseTx) error {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	_, err = s.SubscribeNewBlock(func(block sdk.EventDataNewBlock) {
		reqIDs := block.ResultEndBlock.Tags.GetValues(TagRequestID)
		for _, reqID := range reqIDs {
			request, err := s.QueryRequest(reqID)
			if err != nil {
				//TODO
				continue
			}
			if handler, ok := serviceRouter[request.ServiceName]; ok && provider.Equals(request.Provider) {
				output, errMsg := handler(request.Input)
				msg := MsgRespondService{
					RequestID: reqID,
					Provider:  provider,
					Output:    output,
					Error:     errMsg,
				}
				go func() {
					if _, err = s.Broadcast(baseTx, []sdk.Msg{msg}); err != nil {
						panic(err)
					}
				}()
			}
		}
	})
	return err
}

//RegisterSingleInvocationListener is responsible for registering a single service handler
func (s serviceClient) RegisterSingleInvocationListener(serviceName string,
	respondHandler sdk.ServiceRespondHandler,
	baseTx sdk.BaseTx) error {
	provider, err := s.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	_, err = s.SubscribeNewBlock(func(block sdk.EventDataNewBlock) {
		reqIDs := block.ResultEndBlock.Tags.GetValues(TagRequestID)
		for _, reqID := range reqIDs {
			request, err := s.QueryRequest(reqID)
			if err != nil {
				//TODO
				continue
			}
			if provider.Equals(request.Provider) && request.ServiceName == serviceName {
				output, errMsg := respondHandler(request.Input)
				msg := MsgRespondService{
					RequestID: reqID,
					Provider:  provider,
					Output:    output,
					Error:     errMsg,
				}
				go func() {
					if _, err = s.Broadcast(baseTx, []sdk.Msg{msg}); err != nil {
						panic(err)
					}
				}()
			}
		}
	})
	return err
}

// QueryDefinition return a service definition of the specified name
func (s serviceClient) QueryDefinition(serviceName string) (result sdk.ServiceDefinition, err error) {
	param := struct {
		ServiceName string
	}{
		ServiceName: serviceName,
	}
	err = s.Query("custom/service/definition", param, &result)
	if err != nil {
		return result, err
	}
	return
}

// QueryBinding return the specified service binding
func (s serviceClient) QueryBinding(serviceName string, provider sdk.AccAddress) (result sdk.ServiceBinding, err error) {
	param := struct {
		ServiceName string
		Provider    sdk.AccAddress
	}{
		ServiceName: serviceName,
		Provider:    provider,
	}
	err = s.Query("custom/service/binding", param, &result)
	if err != nil {
		return result, err
	}
	return
}

// QueryBindings returns all bindings of the specified service
func (s serviceClient) QueryBindings(serviceName string) (result []sdk.ServiceBinding, err error) {
	param := struct {
		ServiceName string
	}{
		ServiceName: serviceName,
	}
	err = s.Query("custom/service/bindings", param, &result)
	if err != nil {
		return result, err
	}
	return
}

// QueryRequest returns  the active request of the specified requestID
func (s serviceClient) QueryRequest(requestID string) (request sdk.Request, err error) {
	param := struct {
		RequestID string
	}{
		RequestID: requestID,
	}
	err = s.Query("custom/service/request", param, &request)
	if err != nil {
		return request, err
	}
	return
}

// QueryRequest returns all the active requests of the specified service binding
func (s serviceClient) QueryRequests(serviceName string, provider sdk.AccAddress) (result []sdk.Request, err error) {
	param := struct {
		ServiceName string
		Provider    sdk.AccAddress
	}{
		ServiceName: serviceName,
		Provider:    provider,
	}
	err = s.Query("custom/service/requests", param, &result)
	if err != nil {
		return result, err
	}
	return
}

// QueryRequestsByReqCtx returns all requests of the specified request context ID and batch counter
func (s serviceClient) QueryRequestsByReqCtx(requestContextID string, batchCounter uint64) (result []sdk.Request, err error) {
	param := struct {
		RequestContextID []byte
		BatchCounter     uint64
	}{
		RequestContextID: []byte(requestContextID),
		BatchCounter:     batchCounter,
	}
	err = s.Query("custom/service/requests_by_ctx", param, &result)
	if err != nil {
		return result, err
	}
	return
}

// QueryResponse returns a response with the speicified request ID
func (s serviceClient) QueryResponse(requestID string) (result sdk.Response, err error) {
	param := struct {
		RequestID string
	}{
		RequestID: requestID,
	}
	err = s.Query("custom/service/response", param, &result)
	if err != nil {
		return result, err
	}
	return
}

// QueryResponses returns all responses of the specified request context and batch counter
func (s serviceClient) QueryResponses(requestContextID string, batchCounter uint64) (result []sdk.Response, err error) {
	param := struct {
		RequestContextID []byte
		BatchCounter     uint64
	}{
		RequestContextID: []byte(requestContextID),
		BatchCounter:     batchCounter,
	}
	err = s.Query("custom/service/responses", param, &result)
	if err != nil {
		return result, err
	}
	return
}

// QueryRequestContext return the specified request context
func (s serviceClient) QueryRequestContext(requestContextID string) (result sdk.RequestContext, err error) {
	param := struct {
		RequestContextID []byte
	}{
		RequestContextID: []byte(requestContextID),
	}
	err = s.Query("custom/service/context", param, &result)
	if err != nil {
		return result, err
	}
	return
}

//QueryFees return the earned fees for a provider
func (s serviceClient) QueryFees(provider sdk.AccAddress) (result sdk.EarnedFees, err error) {
	param := struct {
		Address sdk.AccAddress
	}{
		Address: provider,
	}
	err = s.Query("custom/service/fees", param, &result)
	if err != nil {
		return result, err
	}
	return
}