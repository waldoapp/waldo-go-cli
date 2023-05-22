package tpw

import (
	"gopkg.in/yaml.v3"
)

func DecodeFromYAML(in []byte, out any) error {
	return yaml.Unmarshal(in, out)
}

func EncodeToYAML(in any) ([]byte, error) {
	return yaml.Marshal(in)
}
