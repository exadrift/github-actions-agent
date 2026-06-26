package version

import (
	"encoding/json"
	"fmt"
	"strings"
)

type JsonVersion struct {
	mapping map[string]any
}

func NewJsonVersion(b []byte) (*JsonVersion, error) {
	jv := &JsonVersion{
		mapping: map[string]any{},
	}
	err := json.Unmarshal(b, &jv.mapping)
	if err != nil {
		return nil, err
	}

	return jv, nil
}

func (v *JsonVersion) Get(key string) (string, error) {
	keyParts := strings.Split(key, ".")
	m := v.mapping
	for i, part := range keyParts {
		next, ok := m[part]
		if !ok {
			return "", fmt.Errorf("unable to locate %s in json file", key)
		}

		// if final part in key, must be a string
		if i == len(keyParts)-1 {
			value, ok := next.(string)
			if !ok {
				return "", fmt.Errorf("json key did not yield a string value")
			}

			return value, nil
		}

		m, ok = next.(map[string]any)
		if !ok {
			return "", fmt.Errorf("unexpected type in json file while iterating key %s", key)
		}
	}

	return "", nil
}

func (v *JsonVersion) Set(key string, value string) error {
	keyParts := strings.Split(key, ".")
	m := v.mapping
	for i, part := range keyParts {
		// if final part in key, must be a string
		if i == len(keyParts)-1 {
			m[part] = value
			return nil
		}

		next, ok := m[part]
		if !ok {
			return fmt.Errorf("unable to locate %s in json file", key)
		}

		m, ok = next.(map[string]any)
		if !ok {
			return fmt.Errorf("unexpected type in json file while iterating key %s", key)
		}
	}

	return nil
}

func (v *JsonVersion) Marshal() ([]byte, error) {
	return json.Marshal(v.mapping)
}
