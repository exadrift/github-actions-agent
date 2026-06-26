package version

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
)

type YamlVersion struct {
	mapping map[string]any
}

func NewYamlVersion(b []byte) (*YamlVersion, error) {
	yv := &YamlVersion{
		mapping: map[string]any{},
	}
	err := yaml.Unmarshal(b, &yv.mapping)
	if err != nil {
		return nil, err
	}

	return yv, nil
}

func (v *YamlVersion) Get(key string) (string, error) {
	keyParts := strings.Split(key, ".")
	m := v.mapping
	for i, part := range keyParts {
		next, ok := m[part]
		if !ok {
			return "", fmt.Errorf("unable to locate %s in yaml file", key)
		}

		// if final part in key, must be a string
		if i == len(keyParts)-1 {
			value, ok := next.(string)
			if !ok {
				return "", fmt.Errorf("yaml key did not yield a string value")
			}

			return value, nil
		}

		m, ok = next.(map[string]any)
		if !ok {
			return "", fmt.Errorf("unexpected type in yaml file while iterating key %s", key)
		}
	}

	return "", nil
}

func (v *YamlVersion) Set(key string, value string) error {
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
			return fmt.Errorf("unable to locate %s in yaml file", key)
		}

		m, ok = next.(map[string]any)
		if !ok {
			return fmt.Errorf("unexpected type in yaml file while iterating key %s", key)
		}
	}

	return nil
}

func (v *YamlVersion) Marshal() ([]byte, error) {
	return yaml.Marshal(v.mapping)
}
