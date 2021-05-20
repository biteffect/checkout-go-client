package liqpay

import (
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/biteffect/go.gm-fin"
	"net/http"
	"net/url"
	"time"
)

type apiRequest map[string]interface{}
type apiResponse map[string]interface{}

func (res *apiResponse) HasError() bool {
	if res == nil {
		return true
	}
	if s, ok := (*res)["status"]; ok && s != "error" {
		return false
	}
	return true
}

func (res *apiResponse) Error() error {
	if !res.HasError() {
		return nil
	}
	if errMsg, ok := (*res)["err_description"].(string); ok {
		return errors.New(errMsg)
	}
	return errors.New("response body has status error but didn't get error description")
}

type Client struct {
	httpClient *http.Client
	url        *url.URL
	publicKey  PublicKey
	secret     []byte
}

func (c *Client) Refund(orderId string, amount gmfin.Amount, opt PayRequestOptions) (*OrderStatus, error) {
	req := c.getBaseRequest("refund")
	req["order_id"] = orderId
	req["amount"] = amount
	if len(opt.Description) > 0 {
		req["description"] = opt.Description
	}
	res := new(OrderStatus)
	if err := c.callApi(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) Offset(amount gmfin.Amount, currency gmfin.Currency, destinations []SplitRule, opt OffsetRequestOptions) (*OffsetStatus, error) {
	req := c.getBaseRequest("offset")
	req["amount"] = amount
	req["currency"] = currency
	if data, err := json.Marshal(destinations); err != nil {
		return nil, err
	} else {
		req["split_rules"] = string(data)
	}
	if opt.BalanceKey != nil {
		req["balance_key"] = opt.BalanceKey.String()
	}
	if len(opt.Description) > 0 {
		req["description"] = opt.Description
	}
	if len(opt.OrderId) > 0 {
		req["order_id"] = opt.OrderId
	}
	res := new(OffsetStatus)
	if err := c.callApi(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) OrderStatus(orderId string) (*OrderStatus, error) {
	req := c.getBaseRequest("status")
	req["order_id"] = orderId
	res := new(OrderStatus)
	if err := c.callApi(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) OffsetStatus(offsetId string) (*OffsetStatus, error) {
	req := c.getBaseRequest("status")
	req["offset_id"] = offsetId
	res := new(OffsetStatus)
	if err := c.callApi(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) getBaseRequest(action string) apiRequest {
	return map[string]interface{}{
		"version":    3,
		"public_key": c.publicKey.String(),
		"action":     action,
	}
}

func (c *Client) callApi(req apiRequest, res interface{}) error {

	if c.httpClient == nil {

		c.httpClient = &http.Client{
			Timeout: 120 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Renegotiation:      tls.RenegotiateOnceAsClient,
					Certificates:       []tls.Certificate{},
					InsecureSkipVerify: true,
				},
			}}

	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	data = []byte(base64.StdEncoding.EncodeToString(data))
	form := url.Values{
		"data":      {string(data)},
		"signature": {c.sign(data)},
	}

	httpRes, err := c.httpClient.Post(c.url.String(), "application/x-www-form-urlencoded",
		bytes.NewBufferString(form.Encode()))
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode < 300 {
		return json.NewDecoder(httpRes.Body).Decode(res)
	}

	errRes := struct {
		ErrDescription string `json:"err_description"`
	}{}

	if err := json.NewDecoder(httpRes.Body).Decode(&errRes); err != nil {
		return err
	}

	if len(errRes.ErrDescription) > 0 {
		return errors.New(errRes.ErrDescription)
	}

	return errors.New("response body has status error but didn't get error description")
}

func (c *Client) sign(data []byte) string {
	hasher := sha1.New()
	hasher.Write(c.secret)
	hasher.Write(data)
	hasher.Write(c.secret)
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}
