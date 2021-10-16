package liqpay

import (
	gmfin "github.com/biteffect/go.gm-fin"
	"net/url"
)

type OrderStatus struct {
	Version           int             `json:"version"`
	PublicKey         PublicKey       `json:"public_key"`
	Action            string          `json:"action,omitempty"`
	PayType           string          `json:"paytype,omitempty"`
	Amount            gmfin.Amount    `json:"amount" validate:"gte=0"`
	AmountPaid        gmfin.Amount    `json:"amount_paid,omitempty" validate:"gte=0"`
	Currency          *gmfin.Currency `json:"currency"`
	Description       string          `json:"description,omitempty"`
	OrderId           string          `json:"order_id,omitempty"`
	OrderData         string          `json:"order_data,omitempty"`
	LiqpayOrderId     string          `json:"liqpay_order_id,omitempty"`
	PaymentId         string          `json:"payment_id,omitempty"`
	Info              string          `json:"info,omitempty"`
	CreateDate        LiqPayTime      `json:"create_date,omitempty"`
	EndDate           *LiqPayTime     `json:"end_date,omitempty"`
	Ip                string          `json:"ip,omitempty"`
	Status            string          `json:"status,omitempty"`
	StatusDescription string          `json:"status_description,omitempty"`
	SenderCardMask2   string          `json:"sender_card_mask2,omitempty"`

	// Custtom GlobalMoney propertioes
	QrCode      []string `json:"qr_code,omitempty"`
	PayShortUrl string   `json:"checkout_short_url,omitempty"`
}

type OffsetStatus struct {
	Version       int             `json:"version"`
	PublicKey     PublicKey       `json:"public_key"`
	Action        string          `json:"action,omitempty"`
	PayType       string          `json:"paytype,omitempty"`
	Status        string          `json:"status,omitempty"`
	Amount        gmfin.Amount    `json:"amount" validate:"gte=0"`
	Currency      *gmfin.Currency `json:"currency"`
	Info          string          `json:"info,omitempty"`
	CreateDate    LiqPayTime      `json:"create_date,omitempty"`
	EndDate       *LiqPayTime     `json:"end_date,omitempty"`
	OrderId       string          `json:"order_id,omitempty"`
	OrderData     string          `json:"order_data,omitempty"`
	LiqpayOrderId string          `json:"liqpay_order_id,omitempty"`
}

type PayRequestOptions struct {
	RequestOptions
	ResultUrl      string `json:"result_url,omitempty"`
	resultUrlDelay int    `json:"result_url_delay"`
	Splits         []SplitRule
}

func (r *PayRequestOptions) ReturnDelay(delay int) {
	r.resultUrlDelay = delay
}

func (r *PayRequestOptions) Fill(in map[string]interface{}) map[string]interface{} {
	in = r.RequestOptions.Fill(in)
	if len(r.ResultUrl) > 0 {
		in["result_url"] = r.ResultUrl
		if r.resultUrlDelay >= 0 {
			in["result_url_delay"] = r.resultUrlDelay
		}
	}
	return in
}

type OffsetRequestOptions struct {
	RequestOptions
	Splits []SplitRule
}

type RequestOptions struct {
	BalanceKey  *PublicKey
	OrderData   string
	Info        string
	Description string
	ServerUrl   string
	Language    string
	Date        *LiqPayTime
}

func (r *RequestOptions) Fill(in map[string]interface{}) map[string]interface{} {
	if len(r.Description) > 0 {
		in["description"] = r.Description
	}
	if r.BalanceKey != nil {
		in["balance_key"] = r.BalanceKey.String()
	}
	if r.Date != nil {
		in["date"] = r.Date
	}
	if len(r.Language) == 2 {
		in["language"] = r.Language
	}
	return in
}

type SplitRule struct {
	Amount     gmfin.Amount `json:"amount"`
	BalanceKey PublicKey    `json:"balance_key"`
	ServerUrl  *url.URL     `json:"server_url,omitempty"`
	Info       string       `json:"info,omitempty"`
}

type ItemRule struct {
	Title  string       `json:"title,omitempty"`
	Amount gmfin.Amount `json:"amount"`
}
