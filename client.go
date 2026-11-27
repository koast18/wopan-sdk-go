package wopan

import (
	"encoding/json"
	"net/http"
	"path"
	"sync"

	"github.com/go-resty/resty/v2"
)

type WoClient struct {
	accessToken  string
	refreshToken string
	psToken      string

	// clientID/clientSecret identify the client that minted the tokens. They
	// must stay consistent: the api-user channel encrypts with clientSecret and
	// every body carries clientID.
	clientID     string
	clientSecret string
	appID        string

	client            *resty.Client
	crypto            *Crypto
	ua                string
	jsonMarshalFunc   func(v interface{}) ([]byte, error)
	jsonUnmarshalFunc func(data []byte, v interface{}) error

	Phone            string
	ZoneURL          string
	zoneURLOnce      sync.Once
	ClassifyRuleData *ClassifyRuleData

	onRefreshToken func(accessToken, refreshToken string)
}

func New(opts ...Option) *WoClient {
	w := &WoClient{
		client:            resty.New(),
		crypto:            NewCrypto(),
		jsonMarshalFunc:   json.Marshal,
		jsonUnmarshalFunc: json.Unmarshal,
		clientID:          DefaultClientID,
		clientSecret:      DefaultClientSecret,
		appID:             DefaultAppID,
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func DefaultWithAccessToken(accessToken string) *WoClient {
	w := Default()
	w.SetAccessToken(accessToken)
	return w
}

func DefaultWithRefreshToken(refreshToken string) *WoClient {
	w := Default()
	w.SetRefreshToken(refreshToken)
	return w
}

// DefaultWithAccessAndPsToken keeps its historical (buggy) signature for source
// compatibility: the first argument has always been set as the access token.
// Prefer DefaultWithRefreshToken followed by SetPsToken.
func DefaultWithAccessAndPsToken(accessToken, psToken string) *WoClient {
	w := Default()
	w.SetAccessToken(accessToken)
	w.SetPsToken(psToken)
	return w
}

func Default() *WoClient {
	return New(WithUA(DefaultUA))
}

func (w *WoClient) SetUA(ua string) {
	w.ua = ua
}

func (w *WoClient) SetJsonMarshalFunc(f func(v interface{}) ([]byte, error)) {
	w.jsonMarshalFunc = f
}

func (w *WoClient) SetJsonUnmarshalFunc(f func(data []byte, v interface{}) error) {
	w.jsonUnmarshalFunc = f
}

// SetClientCredentials switches the client identity. Call it before the first
// request; tokens already set stay untouched but can only be refreshed by the
// client that minted them.
func (w *WoClient) SetClientCredentials(clientID, clientSecret string) {
	if clientID != "" {
		w.clientID = clientID
	}
	if clientSecret != "" {
		w.clientSecret = clientSecret
		w.crypto.SetClientSecret(clientSecret)
	}
}

// SetClientID sets only the clientId that travels inside the bodies.
func (w *WoClient) SetClientID(clientID string) {
	if clientID != "" {
		w.clientID = clientID
	}
}

// SetClientSecret sets the client secret, which doubles as the AES key of the
// api-user channel.
func (w *WoClient) SetClientSecret(secret string) {
	if secret != "" {
		w.clientSecret = secret
		w.crypto.SetClientSecret(secret)
	}
}

func (w *WoClient) ClientID() string {
	return w.clientID
}

func (w *WoClient) ClientSecret() string {
	return w.clientSecret
}

func (w *WoClient) SetAppID(appID string) {
	if appID != "" {
		w.appID = appID
	}
}

func (w *WoClient) AppID() string {
	return w.appID
}

func (w *WoClient) SetAccessToken(token string) {
	w.accessToken = token
	_ = w.crypto.SetAccessToken(token)
}

func (w *WoClient) SetRefreshToken(token string) {
	w.refreshToken = token
}

func (w *WoClient) SetPsToken(psToken string) {
	w.psToken = psToken
}

func (w *WoClient) PsToken() string {
	return w.psToken
}

func (w *WoClient) SetPhone(phone string) {
	w.Phone = phone
}

// SetTokens installs a token pair as returned by the login/refresh endpoints.
func (w *WoClient) SetTokens(accessToken, refreshToken string) {
	if accessToken != "" {
		w.SetAccessToken(accessToken)
	}
	if refreshToken != "" {
		w.SetRefreshToken(refreshToken)
	}
}

func (w *WoClient) GetToken() (string, string) {
	return w.accessToken, w.refreshToken
}

func (w *WoClient) SetHttpClient(httpClient *http.Client) *WoClient {
	w.client = resty.NewWithClient(httpClient)
	return w
}

func (w *WoClient) SetUserAgent(userAgent string) *WoClient {
	w.client.SetHeader("User-Agent", userAgent)
	return w
}

func (w *WoClient) SetDebug(d bool) *WoClient {
	w.client.SetDebug(d)
	return w
}

func (w *WoClient) EnableTrace() *WoClient {
	w.client.EnableTrace()
	return w
}

func (w *WoClient) SetProxy(proxy string) *WoClient {
	w.client.SetProxy(proxy)
	return w
}

func (w *WoClient) NewRequest() *resty.Request {
	return w.client.R()
}

// jsonClientIDSecret is the body fragment of api-user requests: the clientId
// plus the "secret" marker which tells the edge the param is encrypted.
func (w *WoClient) jsonClientIDSecret() Json {
	return Json{
		"clientId": w.clientID,
		"secret":   true,
	}
}

func (w *WoClient) GetFileType(filename string) string {
	ext := path.Ext(filename)
	if ext == "" {
		return "5"
	}
	ext = ext[1:]
	err := w.InitClassifyRule()
	if err != nil {
		return "5"
	}
	if _type, ok := w.ClassifyRuleData.FileTypes[ext]; ok {
		return _type.Type
	}
	return "5"
}

func (w *WoClient) OnRefreshToken(f func(accessToken, refreshToken string)) {
	w.onRefreshToken = f
}

// notifyTokensChanged reports a new token pair to the owner. The callback is
// optional: the previous revision dereferenced it unconditionally, which
// panicked for every caller that never installed one.
func (w *WoClient) notifyTokensChanged(accessToken, refreshToken string) {
	if w.onRefreshToken != nil {
		w.onRefreshToken(accessToken, refreshToken)
	}
}
