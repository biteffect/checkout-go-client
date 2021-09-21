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
	QrCode           []string `json:"qr_code,omitempty"`
	CheckoutShortUrl string   `json:"checkout_short_url,omitempty"`
	CheckoutUrl      string   `json:"checkout_url,omitempty"`
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
	Splits []SplitRule
}

type OffsetRequestOptions struct {
	RequestOptions
}

type RequestOptions struct {
	BalanceKey  *PublicKey
	OrderId     string
	OrderData   string
	Info        string
	Description string
	ServerUrl   string
	Date        *LiqPayTime
}

type ThreeDsOptions struct {
	AcceptHeader string `json:"accept_header,omitempty"`
	Lang         string `json:"lang,omitempty"`
	ColorDepth   string `json:"color_depth,omitempty"`
	ScreenHeight string `json:"screen_height,omitempty"`
	ScreenWidth  string `json:"screen_width,omitempty"`
	TzUserAagent string `json:"user_agent,omitempty"`
	Fingerprint  string `json:"fingerprint,omitempty"`
}

type SplitRule struct {
	Amount     gmfin.Amount `json:"amount"`
	BalanceKey PublicKey    `json:"balance_key"`
	ServerUrl  *url.URL     `json:"server_url,omitempty"`
	Info       string       `json:"info,omitempty"`
}
