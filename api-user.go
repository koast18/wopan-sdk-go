package wopan

// AppLoginBean mirrors the login/refresh response of the Android client
// (AppLoginBean.java:18-32). The field names are the wire names.
type AppLoginBean struct {
	AccessToken     string `json:"access_token"`
	TokenType       string `json:"token_type"`
	RefreshToken    string `json:"refresh_token"`
	ExpiresIn       int    `json:"expires_in"`
	Scope           string `json:"scope"`
	Phone           string `json:"phone"`
	IsSetIdInfo     string `json:"isSetIdInfo"`
	IsNeedVerfyCode int    `json:"isNeedVerfyCode"`
}

// PcWebLoginData no encrypt
type PcWebLoginData struct {
	NeedSmsCode string `json:"needSmsCode"`
}

// PcWebLogin is the PC/Web password login. It never returns tokens directly: a
// non-empty needSmsCode means the account is registered with SMS verification,
// and PcLoginVerifyCode has to be called with the received code.
func (w *WoClient) PcWebLogin(phone, password string, opts ...RestyOption) (*PcWebLoginData, error) {
	var resp PcWebLoginData
	_, err := w.RequestApiUser(KeyPcWebLogin, Json{
		"phone":        phone,
		"password":     password,
		"uuid":         "",
		"verifyCode":   "",
		"clientSecret": w.clientSecret,
	}, w.jsonClientIDSecret(), &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// PcLoginVerifyCode completes PcWebLogin with the SMS code.
type PcLoginVerifyCodeData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func (w *WoClient) PcLoginVerifyCode(phone, password, messageCode string, opts ...RestyOption) (*PcLoginVerifyCodeData, error) {
	var resp PcLoginVerifyCodeData
	_, err := w.RequestApiUser(KeyPcLoginVerifyCode, Json{
		"phone":        phone,
		"messageCode":  messageCode,
		"verifyCode":   nil,
		"uuid":         nil,
		"clientSecret": w.clientSecret,
		"password":     password,
	}, w.jsonClientIDSecret(), &resp, opts...)
	if err != nil {
		return nil, err
	}
	if resp.AccessToken != "" {
		w.SetTokens(resp.AccessToken, resp.RefreshToken)
		w.SetPhone(phone)
		w.notifyTokensChanged(resp.AccessToken, resp.RefreshToken)
	}
	return &resp, nil
}

// AppLoginAccount is the Android client's username/password login
// (AppLoginAccount, AppLoginRequest.java:34-51).
//
// uuid/verifyCode: the graphical captcha. Pass empty strings on the first
// attempt; when the server answers isNeedVerfyCode != 0 the caller has to fetch
// the image with GetVerifyCode(uuid) and retry with the same uuid plus the code
// the user typed. The pair is only sent when both are set, exactly like the
// client does (it checks `uuid != null && verifyCode != null`).
func (w *WoClient) AppLoginAccount(phone, password, uuid, verifyCode string, opts ...RestyOption) (*AppLoginBean, error) {
	param := Json{
		"phone":        phone,
		"password":     password,
		"clientSecret": w.clientSecret,
	}
	if uuid != "" && verifyCode != "" {
		param["verifyCode"] = verifyCode
		param["uuid"] = uuid
	}
	var resp AppLoginBean
	_, err := w.RequestApiUser(KeyAppLoginAccount, param, w.jsonClientIDSecret(), &resp, opts...)
	if err != nil {
		return nil, err
	}
	w.adoptLoginBean(&resp, phone)
	if resp.IsNeedVerfyCode != 0 && resp.AccessToken == "" {
		return &resp, ErrVerificationCodeRequired
	}
	return &resp, nil
}

// AppLoginByMobile is the SMS-code login (AppLoginByMobile); when password is
// non-empty the client switches to AppLoginDoubleAuth, which needs both the SMS
// code and the password (AppLoginByMobileRequest.java:34-53).
func (w *WoClient) AppLoginByMobile(phone, smsCode, password string, opts ...RestyOption) (*AppLoginBean, error) {
	key := KeyAppLoginByMobile
	param := Json{
		"phone":        phone,
		"smsCode":      smsCode,
		"clientSecret": w.clientSecret,
	}
	if password != "" {
		key = KeyAppLoginDoubleAuth
		delete(param, "smsCode")
		param["messageCode"] = smsCode
		param["password"] = password
	}
	var resp AppLoginBean
	_, err := w.RequestApiUser(key, param, w.jsonClientIDSecret(), &resp, opts...)
	if err != nil {
		return nil, err
	}
	w.adoptLoginBean(&resp, phone)
	if resp.IsNeedVerfyCode != 0 && resp.AccessToken == "" {
		return &resp, ErrVerificationCodeRequired
	}
	return &resp, nil
}

// adoptLoginBean installs the tokens of a successful login/refresh response.
func (w *WoClient) adoptLoginBean(resp *AppLoginBean, phone string) {
	if resp == nil || resp.AccessToken == "" {
		return
	}
	w.SetTokens(resp.AccessToken, resp.RefreshToken)
	if phone != "" {
		w.SetPhone(phone)
	} else if resp.Phone != "" {
		w.SetPhone(resp.Phone)
	}
	w.notifyTokensChanged(resp.AccessToken, resp.RefreshToken)
}

// GetVerifyCode fetches the graphical captcha image bytes for uuid
// (GET /api-user/getverifycode?uuid=..., r1.java:78-80). The uuid is generated
// by the caller and must be echoed back in the login request.
func (w *WoClient) GetVerifyCode(uuid string, opts ...RestyOption) ([]byte, error) {
	req := w.NewRequest().SetQueryParam("uuid", uuid)
	for _, opt := range opts {
		opt(req)
	}
	res, err := req.Get(DefaultBaseURL + "/api-user/getverifycode")
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, errorFromStatus("fetch verification code", res.Status())
	}
	return res.Body(), nil
}

// SendMessageCode asks the server to send an SMS code
// (POST /api-user/sendMessageCodeBase).
//
// This endpoint does not use the dispatcher envelope: the body is
// {func, clientId, param} with param encrypted by the fixed app secret and no
// signature (DataEncryptBodyRequest.java:51-81).
func (w *WoClient) SendMessageCode(phone string, operateType int, fn string, opts ...RestyOption) error {
	if fn == "" {
		fn = "app_send"
	}
	enc, err := w.EncryptParam(ChannelAPIUser, Json{
		"phone":       phone,
		"operateType": operateType,
	})
	if err != nil {
		return err
	}
	req := w.NewRequest().SetBody(Json{
		"func":     fn,
		"clientId": w.clientID,
		"param":    enc,
	})
	for _, opt := range opts {
		opt(req)
	}
	var resp Resp
	res, err := req.SetResult(&resp).Post(DefaultBaseURL + "/api-user/sendMessageCodeBase")
	if err != nil {
		return err
	}
	if res.IsError() {
		return errorFromStatus("send message code", res.Status())
	}
	if resp.Rsp.RspCode != rspCodeOK && resp.Status != "" {
		return errorFromRsp("send message code", resp.Rsp.RspCode, resp.Rsp.RspDesc)
	}
	return nil
}

type AppQueryUserData struct {
	UserId        string `json:"userId"`
	HeadUrl       string `json:"headUrl"`
	UserName      string `json:"userName"`
	Sex           string `json:"sex"`
	Birthday      string `json:"birthday"`
	IsModify      string `json:"isModify"`
	IsHeadModify  string `json:"isHeadModify"`
	IsSetPassword string `json:"isSetPassword"`
	RegisterTime  string `json:"registerTime"`
}

func (w *WoClient) AppQueryUser(opts ...RestyOption) (*AppQueryUserData, error) {
	var resp AppQueryUserData
	_, err := w.RequestApiUser(KeyAppQueryUser, Json{
		"accessToken": w.accessToken,
	}, w.jsonClientIDSecret(), &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// AppRefreshTokenData no encrypt
type AppRefreshTokenData struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

// AppRefreshToken exchanges the refresh token for a new pair, using the client
// secret the pair was minted with.
func (w *WoClient) AppRefreshToken(opts ...RestyOption) (*AppRefreshTokenData, error) {
	var resp AppRefreshTokenData
	_, err := w.RequestApiUser(KeyAppRefreshToken, Json{
		"refreshToken": w.refreshToken,
		"clientSecret": w.clientSecret,
	}, w.jsonClientIDSecret(), &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, err
}

func (w *WoClient) AppLogout(opts ...RestyOption) error {
	_, err := w.RequestApiUser(KeyAppLogout, Json{
		"accessToken": w.accessToken,
	}, w.jsonClientIDSecret(), nil, opts...)
	return err
}
