package liqpay

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/go-pg/pg/v10/types"
	"strconv"
	"strings"
	"time"
)

const jsonDateFormat = "2006-01-02"
const jsonTimeFormat = "15:04:05Z" // "2006-01-02T15:04:05.999Z"
const jsonFullFormat = jsonDateFormat + "T" + jsonTimeFormat

// spetioal type for universal json parse unix time stamp or time string
type LiqPayTime time.Time

func (t LiqPayTime) Value() (driver.Value, error) {
	return time.Time(t).Format(time.RFC3339Nano), nil
}

// Scan implements the sql.Scanner interface.
func (t *LiqPayTime) Scan(src interface{}) error {
	lt, err := types.ParseTime(src.([]byte))
	if err == nil {
		*t = LiqPayTime(lt)
	}
	return err
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *LiqPayTime) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		return nil
	}
	str := strings.Trim(string(bytes), "\"")
	fmt := ""
	switch len(str) {
	case 0:
		return nil
	case 13:
		i, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			*t = LiqPayTime(time.Unix(0, i*int64(time.Millisecond)))
		}
		return err
	case len(jsonDateFormat):
		str += "T00:00:00Z"
		fmt = time.RFC3339
	case len(jsonFullFormat):
		str = str[:10] + "T" + str[11:]
		fmt = time.RFC3339
	case len(time.RFC3339):
		fmt = time.RFC3339
	case len(time.RFC3339Nano):
		fmt = time.RFC3339Nano
	default:
		return errors.New("time must be in RFC3339 or timestamp format")
	}
	nt, err := time.Parse(fmt, str)
	if err == nil {
		*t = LiqPayTime(nt)
	}
	return err
}

// MarshalJSON implements the json.Marshaler interface.
func (t LiqPayTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(t).Format(time.RFC3339))), nil
}
