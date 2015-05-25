package lib

import (
	"encoding/json"
)

func (bigv *BigVClient) CreateGroup(name GroupName) error {
	path := BuildUrl("/accounts/%s/groups", name.Account)

	obj := map[string]string{
		"name": name.Group,
	}

	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, _, err = bigv.Request(true, "POST", path, string(bytes))
	return err
}