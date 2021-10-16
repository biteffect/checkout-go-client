package liqpay

import (
	"encoding/json"
	"fmt"
	gmfin "github.com/biteffect/go.gm-fin"
	"net/url"
)

type CommonStatus struct {
	Version     int             `json:"version"`
	PublicKey   PublicKey       `json:"public_key"`
	Action      string          `json:"action,omitempty"`
	Status      string          `json:"status,omitempty"`
	Amount      gmfin.Amount    `json:"amount" validate:"gte=0"`
	AmountPaid  gmfin.Amount    `json:"amount_paid,omitempty" validate:"gte=0"`
	Currency    *gmfin.Currency `json:"currency"`
	Description string          `json:"description,omitempty"`
	CreateDate  LiqPayTime      `json:"create_date,omitempty"`
	EndDate     *LiqPayTime     `json:"end_date,omitempty"`
}

type OrderStatus struct {
	CommonStatus
	PayType           string `json:"paytype,omitempty"`
	OrderId           string `json:"order_id,omitempty"`
	OrderData         string `json:"order_data,omitempty"`
	LiqpayOrderId     string `json:"liqpay_order_id,omitempty"`
	PaymentId         string `json:"payment_id,omitempty"`
	Info              string `json:"info,omitempty"`
	Ip                string `json:"ip,omitempty"`
	StatusDescription string `json:"status_description,omitempty"`
	SenderCardMask2   string `json:"sender_card_mask2,omitempty"`

	// Custtom GlobalMoney propertioes
	QrCode      []string `json:"qr_code,omitempty"`
	PayShortUrl string   `json:"checkout_short_url,omitempty"`
}

type OffsetStatus struct {
	CommonStatus
	Info          string `json:"info,omitempty"`
	OffsetId      string `json:"offset_id,omitempty"`
	OffsetData    string `json:"offset_data,omitempty"`
	LiqpayOrderId string `json:"liqpay_order_id,omitempty"`
}

type PayRequestOptions struct {
	RequestOptions
	ResultUrl      string      `json:"result_url,omitempty"`
	resultUrlDelay int         `json:"result_url_delay"`
	Splits         []SplitRule `json:"split_rules,omitempty"`
	CartItems      []CartItem  `json:"cart_items,omitempty"`
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
	ExternalId  string
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

type CartItem struct {
	Label  string       `json:"label"`
	Amount gmfin.Amount `json:"amount"`
}

type CartItems []*CartItem

func (i *CartItems) HasItems() bool {
	return i != nil && len(*i) > 0
}

func (i *CartItems) Sum() gmfin.Amount {
	out := gmfin.AmountFromCents(0)
	if i != nil {
		for _, v := range *i {
			out = out.Add(v.Amount)
		}
	}
	return out
}

func (i *CartItems) Items() []*CartItem {
	if i == nil {
		return make([]*CartItem, 0)
	}
	return []*CartItem(*i)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (i *CartItems) UnmarshalJSON(bytes []byte) error {
	//str := strings.Trim(string(bytes), "\" ")
	if len(bytes) < 2 {
		return nil
	}
	v := ""
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return err
	}
	list := make([]*CartItem, 0)
	err = json.Unmarshal([]byte(v), &list)
	if err != nil {
		return err
	}
	for _, v := range list {
		if v.Amount < 0 {
			return fmt.Errorf("cart item %s is negative", v.Label)
		}
	}
	*i = list
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (i *CartItems) MarshalJSON() (out []byte, err error) {
	out = []byte{}
	if !i.HasItems() {
		return out, nil
	}
	out, err = json.Marshal(i)
	if err != nil {
		return
	}
	out = []byte("\"" + string(out) + "\"")
	return
}
