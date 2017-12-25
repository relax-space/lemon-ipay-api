package core

import "github.com/mitchellh/mapstructure"

// Decode takes a map and uses reflection to convert it into the
// given Go native structure. val must be a pointer to a struct.
func Decode(m interface{}, rawVal interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   rawVal,
		TagName:  "mapstruct",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(m)
}
