package liqpay

import (
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"github.com/satori/go.uuid"
	"strings"
)

var keysEncoder = base64.NewEncoding(encodeStd).WithPadding(base64.NoPadding)

type PublicKey uuid.UUID

func (k PublicKey) Id() uuid.UUID {
	return uuid.UUID(k)
}

func (k PublicKey) String() string {
	return keysEncoder.EncodeToString(uuid.UUID(k).Bytes())
}

func (k PublicKey) Value() (driver.Value, error) {
	return uuid.UUID(k).Value()
}

// Scan implements the sql.Scanner interface.
func (k *PublicKey) Scan(src interface{}) error {
	u := uuid.UUID{}
	err := u.Scan(src)
	if err == nil {
		*k = PublicKey(u)
	}
	return err
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (k *PublicKey) UnmarshalJSON(bytes []byte) error {
	d, err := keysEncoder.DecodeString(strings.Trim(string(bytes), "\" "))
	if err != nil {
		return errors.New("public key has bad format")
	}
	id, err := uuid.FromBytes(d)
	if err != nil {
		return err
	}
	*k = PublicKey(id)
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (k PublicKey) MarshalJSON() ([]byte, error) {
	return []byte("\"" + k.String() + "\""), nil
}
