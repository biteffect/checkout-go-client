package liqpay

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/biteffect/go.gm-fin"
	"net/http"
	"net/url"
)

type apiRequest map[string]interface{}
type apiResponse map[string]interface{}

type Client struct {
	url       *url.URL
	publicKey PublicKey
	secret    []byte
}

func (c *Client) Refund(orderId string, amount gmfin.Amount) {
	req := c.getBaseRequest("refund")
	req["order_id"] = orderId
	req["amount"] = amount
}

func (c *Client) OrderStatus(orderId string) {
	req := c.getBaseRequest("status")
	req["order_id"] = orderId
}

func (c *Client) OperationStatus(operationId string) {

}

func (c *Client) getBaseRequest(action string) apiRequest {
	return map[string]interface{}{
		"version":    3,
		"public_key": c.publicKey.String(),
		"action":     action,
	}
}

func (c *Client) callApi(req apiRequest, res interface{}) error {

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	data = []byte(base64.StdEncoding.EncodeToString(data))
	form := url.Values{
		"data":      {string(data)},
		"signature": {c.sign(data)},
	}

	httpRes, err := http.Post(c.url.String(), "application/x-www-form-urlencoded",
		bytes.NewBufferString(form.Encode()))
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode != 200 {
		return fmt.Errorf("bad response status code : %d", httpRes.StatusCode)
	}

	apiRes := apiResponse{}

	if err := json.NewDecoder(httpRes.Body).Decode(&apiRes); err != nil {
		return err
	}

	if apiRes["status"] == "error" {
		errMsg, ok := apiRes["err_description"].(string)
		if ok {
			return errors.New(errMsg)
		}
		return errors.New("response body has status error but didn't get error description")
	}

	return nil
}

func (c *Client) sign(data []byte) string {
	hasher := sha1.New()
	hasher.Write(c.secret)
	hasher.Write(data)
	hasher.Write(c.secret)
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}
