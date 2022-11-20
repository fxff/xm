package auth

import (
	"encoding/base64"
)

type (
	Key []byte

	Config struct {
		Key Key
	}
)

func (pk Key) MarshalYAML() (interface{}, error) {
	return base64.StdEncoding.EncodeToString(pk), nil
}

//nolint:wrapcheck // unnecessary
func (pk *Key) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	ba, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}

	*pk = ba
	return nil
}
