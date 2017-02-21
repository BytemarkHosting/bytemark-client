package brain

import "fmt"

// VLAN is a representation of a VLAN, as used by admin endpoints
type VLAN struct {
	ID        int    `json:"id"`
	Num       int    `json:"num"`
	UsageType string `json:"usage_type"`
}

// String serialises a VLAN to easily be output
func (v *VLAN) String() string {
	return fmt.Sprintf("%d: %s (Num: %d)", v.ID, v.UsageType, v.ID)
}
