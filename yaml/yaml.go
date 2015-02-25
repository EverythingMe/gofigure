package yaml

import (
	"io"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Decoder struct{}

func (d Decoder) Decode(r io.Reader, config interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}

func (d Decoder) CanDecode(path string) bool {
	return strings.HasSuffix(path, ".yaml")
}
