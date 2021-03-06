package brain

import (
	"fmt"
	"io"
	"strings"

	"github.com/BytemarkHosting/bytemark-client/lib/output"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
)

// PrivilegeLevel is a type to represent different privilege levels.
// since privilege levels in the brain are just strings, they're just a string type here too.
type PrivilegeLevel string

const (
	// ClusterAdminPrivilege allows a user to administer the cluster managed by the brain, and do things like create/delete VMs on accounts they have no explicit right on, grant others AccountAdminPrivilege, and set disc iops_limit
	ClusterAdminPrivilege PrivilegeLevel = "cluster_admin"
	// AccountAdminPrivilege allows a user to create, modify & delete groups and servers in an account.
	AccountAdminPrivilege = "account_admin"
	// GroupAdminPrivilege allows a user to create, modify & delete servers in a specific group.
	GroupAdminPrivilege = "group_admin"
	// VMAdminPrivilege allows a user to modify & administer a server, including increasing the performance (and hence the price on the uk0 cluster) and accessing the console.
	VMAdminPrivilege = "vm_admin"
	// VMConsolePrivilege allows a user to access the console for a particular server.
	VMConsolePrivilege = "vm_console"
)

// String returns the privilege level cast to a string.
func (pl PrivilegeLevel) String() string {
	return string(pl)
}

const (
	// PrivilegeTargetTypeVM is the prefix for all privilege levels that affect VMs
	PrivilegeTargetTypeVM = "vm"
	// PrivilegeTargetTypeGroup is the prefix for all privilege levels that affect Groups
	PrivilegeTargetTypeGroup = "group"
	// PrivilegeTargetTypeAccount is the prefix for all privilege levels that affect Accounts
	PrivilegeTargetTypeAccount = "account"
	// PrivilegeTargetTypeCluster is the prefix for all privilege levels that affect the whole cluster.
	PrivilegeTargetTypeCluster = "cluster"
)

// Privilege represents a privilege on the brain.
// A user may have multiple privileges, and multiple privileges may be granted on the same object.
// At the moment we're not worried about the extra fields that privileges have on the brain (IP restrictions) because they're unused
type Privilege struct {
	// ID is the numeric ID used mostly by the brain
	ID int `json:"id,omitempty"`
	// Username is the user who the privilege is granted to
	Username string `json:"username,omitempty"`
	// VirtualMachineID is the ID of the virtual machine the privilege is granted on
	VirtualMachineID int `json:"virtual_machine_id,omitempty"`
	// VirtualMachineName is the name of the virtual machine the privilege is granted on
	VirtualMachineName string `json:"virtual_machine_name,omitempty"`
	// AccountID is the ID of the account the privilege is granted on
	AccountID int `json:"account_id,omitempty"`
	// AccountID is the name of the account the privilege is granted on
	AccountName string `json:"account_name,omitempty"`
	// GroupID is the ID of the group the privilege is granted on
	GroupID int `json:"group_id,omitempty"`
	// GroupID is the name of the group the privilege is granted on
	GroupName string `json:"group_name,omitempty"`
	// APIKeyID is the ID of the api key required in order to use this Privilege. If not set, privilege relates to normal logins using auth.
	APIKeyID int `json:"api_key_id,omitempty"`
	// Level is the PrivilegeLevel they have
	Level PrivilegeLevel `json:"level,omitempty"`
	// PasswordRequired is true for all privileges requiring a password (all except api key privileges)
	PasswordRequired bool `json:"password_required"`
	// YubikeyRequired is true if the user should have to authenticate with a yubikey in order to use this privilege. Only set it to true if you're sure the user has a yubikey set up on their account, and that they know where it is!
	YubikeyRequired bool `json:"yubikey_required"`
	// YubikeyOTPMaxAge should set how long (in seconds) a yubikey one-time-password would be accepted for, but it might not be used?
	YubikeyOTPMaxAge int `json:"yubikey_otp_max_age,omitempty"`
}

// Target returns a formatted string containing the target type and its ID.
func (p Privilege) Target() string {
	switch p.TargetType() {
	case PrivilegeTargetTypeVM:
		return fmt.Sprintf("server %s", p.VirtualMachineName)
	case PrivilegeTargetTypeGroup:
		return fmt.Sprintf("group %s", p.GroupName)
	case PrivilegeTargetTypeAccount:
		return fmt.Sprintf("account %s", p.AccountName)
	}
	return ""

}

// TargetType returns the prefix of the PrivilegeLevel, which should be one of the PrivilegeTargetType* constants.
func (p Privilege) TargetType() string {
	return strings.Split(string(p.Level), "_")[0]
}

// String returns a string representation of the Privilege in English.
// Privileges are a little tricky to represent in English because the Privilege itself doesn't know if it exists on a user or if it has just been removed from a user, nor does it now anything about the target it's been granted on/revoked from other than a numerical ID. So we do the best we can.
func (p Privilege) String() string {
	requiresYubikey := ""
	if p.YubikeyRequired {
		requiresYubikey = " (requires yubikey)"
	}
	switch p.TargetType() {
	case PrivilegeTargetTypeVM:
		return fmt.Sprintf("%s on VM %s for %s%s", p.Level, p.VirtualMachineName, p.Username, requiresYubikey)
	case PrivilegeTargetTypeGroup:
		return fmt.Sprintf("%s on group %s for %s%s", p.Level, p.GroupName, p.Username, requiresYubikey)
	case PrivilegeTargetTypeAccount:
		return fmt.Sprintf("%s on account %s for %s%s", p.Level, p.AccountName, p.Username, requiresYubikey)
	}
	return fmt.Sprintf("%s for %s%s", p.Level, p.Username, requiresYubikey)
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type.
func (p Privilege) DefaultFields(f output.Format) string {
	switch f {
	case output.List:
		return "ID, Username, Level, APIKeyID, PasswordRequired, YubikeyRequired, Target"
	}
	return "ID, Username, Level, Target, APIKeyID, PasswordRequired, YubikeyRequired"
}

// PrettyPrint nicely formats the Privilege and sends it to the given writer.
// At the moment, the detail parameter is ignored.
func (p Privilege) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) (err error) {
	_, err = wr.Write([]byte(p.String()))
	return
}

// Privileges is used to allow API consumers to use IndexOf on the array of privileges.
type Privileges []Privilege

// IndexOf finds the privilege given in the list of privileges, ignoring the Privilege ID and returns the index. If it couldn't find it, returns -1.
func (ps Privileges) IndexOf(priv Privilege) int {
	if priv.Username == "" || priv.Level == "" {
		return -1
	}
	for i, p := range ps {
		if p.VirtualMachineID == priv.VirtualMachineID &&
			p.GroupID == priv.GroupID && p.AccountID == priv.AccountID &&
			p.YubikeyRequired == priv.YubikeyRequired &&
			p.Level == priv.Level && p.Username == priv.Username {
			return i
		}
	}
	return -1
}

// DefaultFields returns the list of default fields to feed to github.com/BytemarkHosting/row.From for this type, which is the same as Privilege.DefaultFields.
func (ps Privileges) DefaultFields(f output.Format) string {
	return (Privilege{}).DefaultFields(f)
}

// PrettyPrint writes a human-readable summary of the privileges to writer at the given detail level.
func (ps Privileges) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	privilegesTpl := `
{{ define "privileges_sgl" }}{{ len . }} servers{{ end }}

{{ define "privileges_medium" -}}
{{- range . -}}
{{- prettysprint . "_sgl" }}
{{ end -}}
{{- end }}

{{ define "privileges_full" }}{{ template "privileges_medium" . }}{{ end }}
`
	return prettyprint.Run(wr, privilegesTpl, "privileges"+string(detail), ps)
}
