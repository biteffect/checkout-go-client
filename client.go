package liqpay

import (
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/biteffect/go.gm-fin"
	"io"
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

func (c *Client) PayWithCard(orderId string, amount gmfin.Amount, card gmfin.CreditCard, opt *PayRequestOptions) (*OrderStatus, error) {
	if err := card.Validate(); err != nil {
		return nil, err
	}
	req := c.getBaseRequest("pay")
	req["amount"] = amount
	if opt != nil {
		req = opt.Fill(req)
	}
	req["order_id"] = orderId
	req["card"] = card.NumberString()
	req["card_exp_month"] = fmt.Sprintf("%v", int(card.ExpiryMonth))
	req["card_exp_year"] = fmt.Sprintf("%v", card.ExpiryYear)
	req["card_cvv"] = fmt.Sprintf("%v", card.CardSecurityCode)
	res := new(OrderStatus)
	if err := c.callApi(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) Refund(orderId string, amount gmfin.Amount, opt *PayRequestOptions) (*OrderStatus, error) {
	req := c.getBaseRequest("refund")
	req["amount"] = amount
	if opt != nil {
		req = opt.Fill(req)
	}
	req["order_id"] = orderId
	res := new(OrderStatus)
	if err := c.callApi(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) PayQr(orderId string, amount gmfin.Amount, currency gmfin.Currency, opt *PayRequestOptions) (*OrderStatus, error) {
	req := c.getBaseRequest("payqr")
	req["amount"] = amount
	req["currency"] = currency
	if opt != nil {
		req = opt.Fill(req)
	}
	req["order_id"] = orderId
	res := new(OrderStatus)
	if err := c.callApi(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) Offset(amount gmfin.Amount, currency gmfin.Currency, opt *OffsetRequestOptions) (*OffsetStatus, error) {
	req := c.getBaseRequest("offset")
	req["amount"] = amount
	req["currency"] = currency
	if opt != nil {
		req = opt.Fill(req)
		if data, err := json.Marshal(opt.Splits); err != nil {
			return nil, err
		} else {
			req["split_rules"] = string(data)
		}
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
		if err.Error() == "unknown order id" {
			err = nil
		}
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

func (c *Client) InfoMerchant() (*OffsetStatus, error) {
	req := c.getBaseRequest("agent_info_merchant")
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
	dataStr := base64.StdEncoding.EncodeToString(data)
	signStr := c.sign([]byte(dataStr))
	form := url.Values{
		"data":      {dataStr},
		"signature": {signStr},
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

	rawRes, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(rawRes, &errRes); err != nil {
		return err
	}

	if len(errRes.ErrDescription) > 0 {
		return errors.New(errRes.ErrDescription)
	}

	return errors.New("response body has status error but didn't get error description")
}

func (c *Client) sign(data []byte) string {
	h := sha1.New()
	h.Write(c.secret)
	h.Write(data)
	h.Write(c.secret)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
