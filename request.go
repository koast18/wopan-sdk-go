package wopan

import (
	"fmt"
	"strings"
)

// Response codes the client has to act on rather than merely report.
//
//	9999   generic failure; what the edge has historically answered on an
//	       expired session for the web/PC client ids
//	1001   RESULT_RSP_CODE_INVALID_TOKEN — the Android app's own token-expired
//	       signal (BaseResponse.isTokenInvalid, acted on by its interceptor)
//	10002  the user was removed from the family group
//	10004  the family group was dismissed
const (
	rspCodeTokenInvalid      = "1001"
	rspCodeGenericFailure    = "9999"
	rspCodeRemovedFromFamily = "10002"
	rspCodeFamilyDismissed   = "10004"
	rspCodeOK                = "0000"
)

func (w *WoClient) request(channel string, key string, param, other Json, resp interface{}, retry bool, opts ...RestyOption) ([]byte, error) {
	return w.requestURL(fmt.Sprintf("%s/%s/dispatcher", DefaultBaseURL, channel),
		channel, key, param, other, resp, retry, opts...)
}

// RequestPath sends the dispatcher envelope to an arbitrary path instead of
// {base}/{channel}/dispatcher.
//
// A few endpoints share the very same envelope (signed header + AES encrypted
// body) but are not reachable through the dispatcher — for example
// /wohome/v1/displayMembershipCard, which returns the membership level and the
// per-file upload limit.
func (w *WoClient) RequestPath(path, channel, key string, param, other Json, resp interface{}, opts ...RestyOption) ([]byte, error) {
	return w.requestURL(DefaultBaseURL+path, channel, key, param, other, resp, true, opts...)
}

func (w *WoClient) requestURL(url, channel, key string, param, other Json, resp interface{}, retry bool, opts ...RestyOption) ([]byte, error) {
	req := w.NewRequest()
	req.SetHeaders(map[string]string{
		"Origin":  "https://pan.wo.cn",
		"Referer": "https://pan.wo.cn/",
	})
	if w.accessToken != "" {
		req.SetHeader("Accesstoken", w.accessToken)
	}
	header := calHeader(channel, key)
	body, err := w.NewBody(channel, param, other)
	if err != nil {
		return nil, err
	}
	var _resp Resp
	req.SetBody(Req[interface{}]{
		Header: header,
		Body:   body,
	}).SetResult(&_resp)
	for _, opt := range opts {
		opt(req)
	}
	res, err := req.Post(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return res.Body(), fmt.Errorf("request failed with status: %s", res.Status())
	}
	if _resp.Status != "200" {
		return res.Body(), fmt.Errorf("request failed with status: %s, msg: %s", _resp.Status, _resp.Msg)
	}
	if _resp.Rsp.RspCode != rspCodeOK {
		switch _resp.Rsp.RspCode {
		case rspCodeRemovedFromFamily:
			return res.Body(), ErrRemovedFromFamily
		case rspCodeFamilyDismissed:
			return res.Body(), ErrFamilyDismissed
		}
		// An expired access token is refreshed once and replayed. The api-user
		// channel is excluded because refreshing is exactly what it does.
		if channel != ChannelAPIUser && retry &&
			(_resp.Rsp.RspCode == rspCodeGenericFailure || _resp.Rsp.RspCode == rspCodeTokenInvalid) &&
			w.refreshToken != "" {
			if err := w.RefreshToken(); err != nil {
				return res.Body(), err
			}
			return w.requestURL(url, channel, key, param, other, resp, false, opts...)
		}
		return res.Body(), fmt.Errorf("request failed with rsp_code: %s,rep_desc: %s", _resp.Rsp.RspCode, _resp.Rsp.RspDesc)
	}
	if resp != nil {
		data := string(_resp.Rsp.Data)
		if strings.HasSuffix(data, "\"") && strings.HasPrefix(data, "\"") {
			data, err = w.crypto.Decrypt(data[1:len(data)-1], channel)
			if err != nil {
				return res.Body(), err
			}
		}
		err = w.jsonUnmarshalFunc([]byte(data), resp)
		if err != nil {
			return res.Body(), err
		}
	}
	return res.Body(), nil
}

func (w *WoClient) Request(channel string, key string, param, other Json, resp interface{}, opts ...RestyOption) ([]byte, error) {
	return w.request(channel, key, param, other, resp, true, opts...)
}

func (w *WoClient) RequestApiUser(key string, param, other Json, resp interface{}, opts ...RestyOption) ([]byte, error) {
	return w.Request(ChannelAPIUser, key, param, other, resp, opts...)
}

func (w *WoClient) RequestWoHome(key string, param, other Json, resp interface{}, opts ...RestyOption) ([]byte, error) {
	return w.Request(ChannelWoHome, key, param, other, resp, opts...)
}
