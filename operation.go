package wopan

func (w *WoClient) InitPhone() error {
	if w.Phone == "" {
		_resp, err := w.AppQueryUser()
		if err != nil {
			return err
		}
		// The client sends the phone number as phoneNum; the account API only
		// exposes userId here, which is what this call has always used.
		w.Phone = _resp.UserId
	}
	return nil
}

func (w *WoClient) InitClassifyRule() error {
	if w.ClassifyRuleData != nil {
		return nil
	}
	data, err := w.ClassifyRule()
	if err != nil {
		return err
	}
	w.ClassifyRuleData = data
	return nil
}

func (w *WoClient) InitZoneURL() error {
	if w.ZoneURL != "" {
		return nil
	}
	data, err := w.GetZoneInfo()
	if err != nil {
		return err
	}
	w.ZoneURL = data.Url
	return nil
}

func (w *WoClient) RefreshToken() error {
	if w.refreshToken == "" {
		return ErrInvalidRefreshToken
	}
	resp, err := w.AppRefreshToken()
	if err != nil {
		return err
	}
	if resp.AccessToken == "" {
		return ErrInvalidAccessToken
	}
	w.SetTokens(resp.AccessToken, resp.RefreshToken)
	w.notifyTokensChanged(resp.AccessToken, resp.RefreshToken)
	return nil
}

func (w *WoClient) InitData() error {
	if w.accessToken == "" && w.refreshToken != "" {
		if err := w.RefreshToken(); err != nil {
			return err
		}
	}
	if err := w.InitPhone(); err != nil {
		return err
	}
	if err := w.InitClassifyRule(); err != nil {
		return err
	}
	if err := w.InitZoneURL(); err != nil {
		return err
	}
	return nil
}
