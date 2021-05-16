package liqpay

import "net/url"

func NewClient(url *url.URL, publicKey string, secret []byte) (*Client, error) {
	key := PublicKey{}
	if err := key.UnmarshalJSON([]byte(publicKey)); err != nil {
		return nil, err
	}
	if url == nil {
		url, _ = url.Parse(CheckoutApiUrl)
	}
	return &Client{
		url:       url,
		publicKey: key,
		secret:    secret,
	}, nil
}
