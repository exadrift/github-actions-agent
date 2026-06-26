package version

import (
	"fmt"
	"strings"
)

type TextVersion struct {
	data string
}

func NewTextVersion(b []byte) (*TextVersion, error) {
	return &TextVersion{
		data: string(b),
	}, nil
}

func (v *TextVersion) Get(key string) (string, error) {
	return strings.Trim(v.data, "\r\n "), nil
}

func (v *TextVersion) Set(key string, value string) error {
	v.data = fmt.Sprintf("%s\n", value)
	return nil
}

func (v *TextVersion) Marshal() ([]byte, error) {
	return []byte(v.data), nil
}
