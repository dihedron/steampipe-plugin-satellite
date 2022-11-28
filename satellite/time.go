package satellite

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type Time time.Time

// 2020-06-10 10:03:19 UTC
const layout = "2006-01-02 15:04:05 MST"

// MarshalJSON implements the json.Marshaler interface.
func (t Time) MarshalJSON() ([]byte, error) {
	if y := time.Time(t).Year(); y < 0 || y >= 10000 {
		// Assuming years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(layout)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, layout)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Time) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	s, err := time.Parse(`"`+layout+`"`, string(data))
	*t = Time(s)
	return err
}

func (t *Time) String() string {
	if t == nil || (time.Time(*t)).UnixNano() == 0 {
		return ""
	}
	return time.Time(*t).Format(layout)
}

func (t *Time) IsZero() bool {
	if t == nil {
		return true
	}
	return time.Time(*t).IsZero()
}

func ToTime(ctx context.Context, d *transform.TransformData) (any, error) {
	var err error
	switch t := d.Value.(type) {
	case *Time:
		if t == nil || t.IsZero() {
			return nil, nil
		}
		return t.String(), nil
	case Time:
		if t.IsZero() {
			return nil, nil
		}
		return t.String(), nil
	case time.Time:
		if t.IsZero() {
			return nil, nil
		}
		return t.String(), nil
	default:
		err = fmt.Errorf("invalid type: %T", d.Value)
	}
	return nil, err
}
