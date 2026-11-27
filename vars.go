package wopan

import "errors"

var (
	// JsonSecret marks a wohome request whose param is encrypted.
	JsonSecret = Json{
		"secret": true,
	}

	// JsonClientIDSecret is kept for source compatibility only. It is bound to
	// the historical DefaultClientID, so requests built with it cannot follow a
	// configured client identity; use WoClient.jsonClientIDSecret instead.
	//
	// Deprecated: build the fragment through the client.
	JsonClientIDSecret = Json{
		"clientId": DefaultClientID,
		"secret":   true,
	}
)

var (
	ErrInvalidAccessToken  = errors.New("invalid access token")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrInvalidPsToken      = errors.New("invalid psToken")

	// RSP_CODE 10002 / 10004 mean the configured family group is gone. Retrying
	// or refreshing the token cannot help, the caller has to pick another group.
	ErrRemovedFromFamily = errors.New("the current user has been removed from this family group")
	ErrFamilyDismissed   = errors.New("this family group has been dismissed")

	// ErrVerificationCodeRequired is returned by the login helpers when the
	// server asks for the graphical captcha (AppLoginBean.isNeedVerfyCode != 0).
	ErrVerificationCodeRequired = errors.New("the server requires the graphical verification code: fetch it with GetVerifyCode(uuid) and retry the login with uuid+code")
)
