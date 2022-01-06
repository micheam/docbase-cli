package docbasecli

import (
	"bytes"
	"errors"
	"io"

	"github.com/pelletier/go-toml"
)

type ConfigMap map[string]Config

type Config struct {
	AccessToken string
	Domain      string
	UserID      string
	Editor      string
}

// LoadConfig は、 Default 設定を読み込みます。
func LoadConfig(r io.Reader) (*Config, error) {
	const key = "default"
	var (
		confMap = ConfigMap{}
		err     error
	)
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(buf.Bytes(), &confMap)
	if err != nil {
		return nil, err
	}
	found, ok := confMap[key]
	if !ok {
		return nil, errors.New("'default' profile not found")
	}
	return &found, nil
}
