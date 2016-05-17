package lib

import (
	"bytes"
	"encoding/json"
)

// CreateGroup sends a request to the API server to create a group with the given name.
func (c *bytemarkClient) CreateGroup(name *GroupName) (err error) {
	err = c.validateGroupName(name)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("POST", EP_BRAIN, "/accounts/%s/groups", name.Account)
	if err != nil {
		return
	}

	obj := map[string]string{
		"name": name.Group,
	}

	js, err := json.Marshal(obj)
	if err != nil {
		return
	}
	_, _, err = r.Run(bytes.NewBuffer(js), nil)
	return
}

// DeleteGroup requests that a given group be deleted. Will return an error if there are VMs in the group.
func (c *bytemarkClient) DeleteGroup(name *GroupName) (err error) {
	err = c.validateGroupName(name)
	if err != nil {
		return
	}
	r, err := c.BuildRequest("DELETE", EP_BRAIN, "/accounts/%s/groups/%s", name.Account, name.Group)
	if err != nil {
		return
	}
	_, _, err = r.Run(nil, nil)
	return
}

// GetGroup requests an overview of the group with the given name
func (c *bytemarkClient) GetGroup(name *GroupName) (group *Group, err error) {
	group = new(Group)
	err = c.validateGroupName(name)
	if err != nil {
		return
	}

	r, err := c.BuildRequest("GET", EP_BRAIN, "/accounts/%s/groups/%s?view=overview&include_deleted=true", name.Account, name.Group)
	if err != nil {
		return
	}

	_, _, err = r.Run(nil, group)
	if err != nil {
		return
	}
	return
}