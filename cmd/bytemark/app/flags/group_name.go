package flags

import (
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
)

// GroupNameFlag is used for all --group flags, including the global one.
type GroupNameFlag struct {
	GroupName pathers.GroupName
	Value     string
}

// Set runs lib.ParseGroupName to make sure we have a valid group name
func (name *GroupNameFlag) Set(value string) error {
	name.Value = value
	return nil
}

// Preprocess defaults the value of this flag to the default group from the
// config attached to the context and then runs lib.ParseGroupName
// This is an implementation of `app.Preprocessor`, which is detected and
// called automatically by actions created with `app.Action`
func (name *GroupNameFlag) Preprocess(c *app.Context) (err error) {
	if name.Value == "" {
		return
	}
	groupName := lib.ParseGroupName(name.Value, c.Config().GetGroup())
	name.GroupName = groupName
	return
}

// String returns the GroupName as a string.
func (name GroupNameFlag) String() string {
	return name.GroupName.String()
}
