package zaif

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"github.com/google/go-querystring/query"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const privateEndPointURL = "https://api.zaif.jp/tapi"

var initialNonce = 0

// PrivateAPI API有効にした際の キー,シークレットキー を設定
type PrivateAPI struct {
	Key    string //
	Secret string //
}

// BodyActiveOrders Params of ActiveOrders
type BodyActiveOrders struct {
	From         int    `url:"from,omitempty"`
	Count        int    `url:"count,omitempty"`
	FromID       int    `url:"from_id,omitempty"`
	EndID        int    `url:"end_id,omitempty"`
	Order        string `url:"order,omitempty"`
	Since        int    `url:"since,omitempty"`
	End          int    `url:"end,omitempty"`
	CurrencyPair string `url:"currency_pair,omitempty"`
}

// BodyTrade Params of Trade
type BodyTrade struct {
	CurrencyPair string  `url:"currency_pair"`
	Action       string  `url:"action"`
	Price        int     `url:"price"`
	Amount       float32 `url:"amount"`
	Limit        bool    `url:"limit,omitempty"`
}

// BodyCancel Params of Cancel
type BodyCancel struct {
	OrderID int `url:"order_id"`
}

// BodyWithdraw Params of Withdraw
type BodyWithdraw struct {
	Currency string  `url:"currency"`
	Address  string  `url:"address"`
	Amount   float32 `url:"amount"`
	OptFee   float32 `url:"opt_fee,omitempty"`
}

// BodyTradeHistory Params of TradeHistory
type BodyTradeHistory struct {
	From         int    `url:"from,omitempty"`
	Count        int    `url:"count,omitempty"`
	FromID       int    `url:"from_id,omitempty"`
	EndID        int    `url:"end_id,omitempty"`
	Order        string `url:"order,omitempty"`
	Since        int    `url:"since,omitempty"`
	End          int    `url:"end,omitempty"`
	CurrencyPair string `url:"currency_pair,omitempty"`
}

// BodyDepositHistory Params of DepositHistory()
type BodyDepositHistory struct {
	From     int    `url:"from,omitempty"`
	Count    int    `url:"count,omitempty"`
	FromID   int    `url:"from_id,omitempty"`
	EndID    int    `url:"end_id,omitempty"`
	Order    string `url:"order,omitempty"`
	Since    int    `url:"since,omitempty"`
	End      int    `url:"end,omitempty"`
	Currency string `url:"currency"`
}

// BodyWithdrawHistory Params of WithdrawHistory()
type BodyWithdrawHistory struct {
	From     int    `url:"from,omitempty"`
	Count    int    `url:"count,omitempty"`
	FromID   int    `url:"from_id,omitempty"`
	EndID    int    `url:"end_id,omitempty"`
	Order    string `url:"order,omitempty"`
	Since    int    `url:"since,omitempty"`
	End      int    `url:"end,omitempty"`
	Currency string `url:"currency"`
}

// NewPrivateAPI To use PrivateAPI
func NewPrivateAPI(key string, secret string) *PrivateAPI {
	return &PrivateAPI{
		Key:    key,
		Secret: secret,
	}
}

// TradingParam 全リクエストで使用するparams
type TradingParam struct {
	Method string `url:"method"`
	Nonce  int    `url:"nonce"`
}

// NewTradingParam To make TradingParam
func NewTradingParam(method string) TradingParam {
	return TradingParam{
		Method: method,
		Nonce:  getNonceWithIncr(),
	}
}

// MakeSign messageをHMAC-SHA512で署名
func MakeSign(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha512.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// getNonceWithIncr nonce取得
func getNonceWithIncr() int {
	initialNonce++
	return (int(time.Now().Unix())-1420070400)*30 + (initialNonce % 30)
}

// Post API送信
func (api *PrivateAPI) Post(tradingParam TradingParam, parameter string) (string, error) {
	v, err := query.Values(tradingParam)
	if err != nil {
		return "", err
	}
	encodedParams := v.Encode()
	if parameter != "" {
		encodedParams += "&" + parameter
	}

	req, err := http.NewRequest("POST", privateEndPointURL, strings.NewReader(encodedParams))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", api.Key)
	req.Header.Add("Sign", MakeSign(encodedParams, api.Secret))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// GetInfo 現在の残高（余力および残高）、APIキーの権限、過去のトレード数、アクティブな注文数、サーバーのタイムスタンプを取得
func (api *PrivateAPI) GetInfo() (string, error) {

	return api.Post(
		NewTradingParam("get_info"),
		"",
	)
}

// ActiveOrders 現在有効な注文一覧を取得
func (api *PrivateAPI) ActiveOrders(body BodyActiveOrders) (string, error) {
	v, err := query.Values(body)
	if err != nil {
		return "", err
	}
	return api.Post(
		NewTradingParam("active_orders"),
		v.Encode(),
	)
}

// Trade 注文
func (api *PrivateAPI) Trade(body BodyTrade) (string, error) {
	v, err := query.Values(body)
	if err != nil {
		return "", err
	}
	return api.Post(
		NewTradingParam("trade"),
		v.Encode(),
	)
}

// Cancel 注文キャンセル
func (api *PrivateAPI) Cancel(body BodyCancel) (string, error) {
	v, err := query.Values(body)
	if err != nil {
		return "", err
	}
	return api.Post(
		NewTradingParam("cancel"),
		v.Encode(),
	)
}

// Withdraw 資金の引き出しリクエストを送信
func (api *PrivateAPI) Withdraw(body BodyWithdraw) (string, error) {
	v, err := query.Values(body)
	if err != nil {
		return "", err
	}
	return api.Post(
		NewTradingParam("withdraw"),
		v.Encode(),
	)
}

// DepositHistory 入金履歴を取得
func (api *PrivateAPI) DepositHistory(body BodyDepositHistory) (string, error) {
	v, err := query.Values(body)
	if err != nil {
		return "", err
	}
	return api.Post(
		NewTradingParam("deposit_history"),
		v.Encode(),
	)
}

// WithdrawHistory 出金履歴を取得
func (api *PrivateAPI) WithdrawHistory(body BodyWithdrawHistory) (string, error) {
	v, err := query.Values(body)
	if err != nil {
		return "", err
	}
	return api.Post(
		NewTradingParam("withdraw_history"),
		v.Encode(),
	)
}
