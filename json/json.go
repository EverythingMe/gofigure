package json

import (
	"io"
	"strings"

	"encoding/json"
)

// Decoder can take configurations encoded as json dictionaries and decode them to
// config structs
type Decoder struct{}

// Decode just wraps using a json decoder to unmarshal into config, which is a pointer to a struct
func (d Decoder) Decode(r io.Reader, config interface{}) error {

	dec := json.NewDecoder(r)

	return dec.Decode(config)

}

// CanDecode returns true if this is a json file
func (d Decoder) CanDecode(path string) bool {
	return strings.HasSuffix(path, ".json")
}
