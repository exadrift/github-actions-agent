package version

import (
	"fmt"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type TomlVersion struct {
	mapping map[string]any
}

func NewTomlVersion(b []byte) (*TomlVersion, error) {
	tv := &TomlVersion{
		mapping: map[string]any{},
	}
	err := toml.Unmarshal(b, &tv.mapping)
	if err != nil {
		return nil, err
	}

	return tv, nil
}

func (v *TomlVersion) Get(key string) (string, error) {
	keyParts := strings.Split(key, ".")
	m := v.mapping
	for i, part := range keyParts {
		next, ok := m[part]
		if !ok {
			return "", fmt.Errorf("unable to locate %s in toml file", key)
		}

		// if final part in key, must be a string
		if i == len(keyParts)-1 {
			value, ok := next.(string)
			if !ok {
				return "", fmt.Errorf("toml key did not yield a string value")
			}

			return value, nil
		}

		m, ok = next.(map[string]any)
		if !ok {
			return "", fmt.Errorf("unexpected type in toml file while iterating key %s", key)
		}
	}

	return "", nil
}

func (v *TomlVersion) Set(key string, value string) error {
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
			return fmt.Errorf("unable to locate %s in toml file", key)
		}

		m, ok = next.(map[string]any)
		if !ok {
			return fmt.Errorf("unexpected type in toml file while iterating key %s", key)
		}
	}

	return nil
}

func (v *TomlVersion) Marshal() ([]byte, error) {
	return toml.Marshal(v.mapping)
}
