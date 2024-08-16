package tg

import (
	"github.com/pkg/errors"
	"github.com/zelenin/go-tdlib/client"
)

var (
	_ client.AuthorizationStateHandler = (*authHandler)(nil)
	_ AuthHandler                      = (*authHandler)(nil)
)

type AuthHandler interface {
	StateType() string
	SetPhone(phone string)
	SetCode(code string)
	SetPassword(password string)
}

type authHandler struct {
	params    *client.SetTdlibParametersRequest
	stateType string
	phone     string
	password  string
	code      string
	flow      chan struct{}
}

func newAuthHandler(
	params *client.SetTdlibParametersRequest,
) *authHandler {
	return &authHandler{
		params: params,
		flow:   make(chan struct{}),
	}
}

func (a *authHandler) Close() {
	close(a.flow)
	a.params = nil
	a.stateType = ""
	a.phone = ""
	a.password = ""
	a.code = ""
}

func (a *authHandler) Handle(
	clientInstance *client.Client,
	state client.AuthorizationState,
) error {
	a.stateType = state.AuthorizationStateType()

	var err error

	switch state.AuthorizationStateType() {
	case client.TypeAuthorizationStateWaitTdlibParameters:
		_, err = clientInstance.SetTdlibParameters(a.params)

	case client.TypeAuthorizationStateWaitPhoneNumber:
		<-a.flow

		_, err = clientInstance.SetAuthenticationPhoneNumber(&client.SetAuthenticationPhoneNumberRequest{
			PhoneNumber: a.phone,
			Settings: &client.PhoneNumberAuthenticationSettings{
				AllowFlashCall:       false,
				IsCurrentPhoneNumber: false,
				AllowSmsRetrieverApi: false,
			},
		})

	case client.TypeAuthorizationStateWaitCode:
		<-a.flow

		_, err = clientInstance.CheckAuthenticationCode(&client.CheckAuthenticationCodeRequest{
			Code: a.code,
		})

	case client.TypeAuthorizationStateWaitPassword:
		<-a.flow

		_, err = clientInstance.CheckAuthenticationPassword(&client.CheckAuthenticationPasswordRequest{
			Password: a.password,
		})

	case client.TypeAuthorizationStateReady:
		return nil

	case client.TypeAuthorizationStateClosing:
		return nil

	case client.TypeAuthorizationStateClosed:
		return nil

	default:
		err = errors.WithMessage(client.ErrNotSupportedAuthorizationState, a.stateType)
	}

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *authHandler) StateType() string {
	return a.stateType
}

func (a *authHandler) SetPassword(password string) {
	if a.stateType != client.TypeAuthorizationStateWaitPassword {
		return
	}

	a.password = password
	a.flow <- struct{}{}
}

func (a *authHandler) SetPhone(phone string) {
	if a.stateType != client.TypeAuthorizationStateWaitPhoneNumber {
		return
	}

	a.phone = phone
	a.flow <- struct{}{}
}

func (a *authHandler) SetCode(code string) {
	if a.stateType != client.TypeAuthorizationStateWaitCode {
		return
	}

	a.code = code
	a.flow <- struct{}{}
}
