package utils

import "encoding/json"

func ObjectMapper(in, out any) error {
	b, err := json.Marshal(in)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &out)
	if err != nil {
		return err
	}

	return err
}
