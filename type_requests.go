package liqpay

import (
	"fmt"
	gmfin "github.com/biteffect/go.gm-fin"
)

type PayRequest struct {
	Amount         gmfin.Amount    `json:"amount" validate:"gte=0"`
	Currency       *gmfin.Currency `json:"currency"`
	Card           string          `json:"card,omitempty"`
	CardCvv        string          `json:"card_cvv,omitempty"`
	CardExpMonth   string          `json:"card_exp_month,omitempty"`
	CardExpYear    string          `json:"card_exp_year,omitempty"`
	CardToken      *PublicKey      `json:"card_token,omitempty"`
	Description    string          `json:"description"`
	Ip             string          `json:"ip,omitempty"`
	OrderId        string          `json:"order_id"`
	OrderData      string          `json:"order_data"`
	Language       string          `json:"language,omitempty"`
	ResultUrl      string          `json:"result_url,omitempty"`
	ResultUrlDelay int             `json:"result_url_delay,omitempty"`
	ServerUrl      string          `json:"server_url,omitempty"`
	Customer       string          `json:"customer,omitempty"`
	Info           string          `json:"info,omitempty"`
	Date           *LiqPayTime     `json:"date,omitempty"`

	// params for info_merchant info_balance reports register
	BalanceKey *PublicKey `json:"balance_key,omitempty"`
}

func (pr *PayRequest) Validate() error {
	if pr.Amount <= 0 {
		return fmt.Errorf("amount (%v) must be greater then 0", pr.Amount)
	}
	return nil
}

func (pr *PayRequest) Map(action string, key *PublicKey) map[string]interface{} {
	req := map[string]interface{}{
		"version":    3,
		"public_key": key.String(),
		"action":     action,
	}
	req["amount"] = pr.Amount
	req["currency"] = pr.Currency
	if pr.BalanceKey != nil {
		req["balance_key"] = pr.BalanceKey.String()
	}
	if len(pr.Description) > 0 {
		req["description"] = pr.Description
	}
	if len(pr.OrderId) > 0 {
		req["order_id"] = pr.OrderId
	}
	if pr.Date != nil {
		req["date"] = pr.Date
	}
	return req
}

type PayRequestEnvironment struct {
	AcceptHeader string `json:"accept_header,omitempty"`
	Lang         string `json:"lang,omitempty"`
	ColorDepth   int    `json:"color_depth,omitempty"`
	ScreenHeight int    `json:"screen_height,omitempty"`
	ScreenWidth  int    `json:"screen_width,omitempty"`
	UserAgent    string `json:"user_agent,omitempty"`
	Fingerprint  string `json:"fingerprint,omitempty"`
	Tz           int    `json:"tz,omitempty"`
}
