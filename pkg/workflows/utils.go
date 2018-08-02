package workflows

import "encoding/json"

// bind params uses json serializing and reflect package that is underneath
// to avoid direct access to map for getting appropriate field values.
func bindParams(params map[string]string, object interface{}) error {
	data, err := json.Marshal(params)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, object)

	if err != nil {
		return err
	}

	return nil
}
