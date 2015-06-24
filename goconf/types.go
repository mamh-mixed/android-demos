package goconf

import (
	"fmt"
	"strconv"
	"time"
)

type Duration time.Duration

func (d *Duration) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(time.Duration(*d).String())), nil
}

func (d *Duration) UnmarshalJSON(v []byte) error {
	dv, err := strconv.Unquote(string(v))
	if err != nil {
		return err
	}

	x, err := time.ParseDuration(dv)
	if err != nil {
		fmt.Println(err)
		return err
	}
	*d = Duration(x)
	return nil
}
