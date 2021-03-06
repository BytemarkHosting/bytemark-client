package mocks

import (
	"fmt"
	"net/http"

	auth3 "github.com/BytemarkHosting/auth-client"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/billing"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/pathers"
	"github.com/BytemarkHosting/bytemark-client/lib/spp"
	"github.com/maraino/go-mock"
)

type Client struct {
	mock.Mock
	MockRequest *Request
}

func (c *Client) AllowInsecureRequests() {
	c.Called()
}
func (c *Client) BuildRequestNoAuth(method string, endpoint lib.Endpoint, path string, parts ...string) (lib.Request, error) {
	if c.MockRequest == nil {
		r := c.Called(method, endpoint, path, parts)
		req, _ := r.Get(0).(lib.Request)
		return req, r.Error(1)
	} else {
		return c.MockRequest, nil
	}
}
func (c *Client) BuildRequest(method string, endpoint lib.Endpoint, path string, parts ...string) (lib.Request, error) {
	if c.MockRequest == nil {
		r := c.Called(method, endpoint, path, parts)
		req, _ := r.Get(0).(lib.Request)
		return req, r.Error(1)
	} else {
		return c.MockRequest, nil
	}
}

func (c *Client) Impersonate(user string) error {
	r := c.Called(user)
	return r.Error(0)
}

func (c *Client) GetEndpoint() string {
	r := c.Called()
	return r.String(0)
}
func (c *Client) GetSessionToken() string {
	r := c.Called()
	return r.String(0)
}
func (c *Client) GetSessionUser() string {
	r := c.Called()
	return r.String(0)
}
func (c *Client) GetSessionFactors() []string {
	r := c.Called()
	ar := r.Get(0)
	return ar.([]string)
}
func (c *Client) GetSPPToken(cc spp.CreditCard, owner billing.Person) (string, error) {
	r := c.Called(cc, owner)
	return r.String(0), r.Error(1)
}
func (c *Client) SetDebugLevel(level int) {
	c.Called(level)
}
func (c *Client) AuthWithToken(token string) error {
	r := c.Called(token)
	return r.Error(0)
}
func (c *Client) AuthWithCredentials(credents auth3.Credentials) error {
	r := c.Called(credents)
	return r.Error(0)
}
func (c *Client) RequestAndUnmarshal(auth bool, method, path, requestBody string, output interface{}) error {
	r := c.Called(auth, method, path, requestBody, output)
	return r.Error(0)
}
func (c *Client) RequestAndRead(auth bool, method, path, requestBody string) (responseBody []byte, err error) {
	r := c.Called(auth, method, path, requestBody)
	return r.Bytes(0), r.Error(1)
}

func (c *Client) Request(auth bool, method string, location string, requestBody string) (req *http.Request, res *http.Response, err error) {
	r := c.Called(auth, method, location, requestBody)
	req, ok := r.Get(0).(*http.Request)
	if !ok {
		panic(fmt.Sprintf("Couldn't typecast request object because it was a %t", r.Get(0)))
	}
	res, ok = r.Get(1).(*http.Response)
	if !ok {
		panic(fmt.Sprintf("Couldn't typecast response object because it was a %t", r.Get(1)))
	}
	return req, res, r.Error(2)
}

func (c *Client) ReadDefinitions() (lib.Definitions, error) {
	r := c.Called()
	defs, _ := r.Get(0).(lib.Definitions)
	return defs, r.Error(1)
}

func (c *Client) AddIP(name pathers.VirtualMachineName, spec brain.IPCreateRequest) (brain.IPs, error) {
	r := c.Called(name, spec)
	ips, _ := r.Get(0).(brain.IPs)
	return ips, r.Error(1)
}

func (c *Client) AddUserAuthorizedKey(name, key string) error {
	r := c.Called(name, key)
	return r.Error(0)
}

func (c *Client) DeleteUserAuthorizedKey(name, key string) error {
	r := c.Called(name, key)
	return r.Error(0)
}

func (c *Client) GetUser(name string) (brain.User, error) {
	r := c.Called(name)
	u, _ := r.Get(0).(brain.User)
	return u, r.Error(1)
}

func (c *Client) CreateCreditCard(cc spp.CreditCard) (string, error) {
	r := c.Called(cc)
	return r.String(0), r.Error(1)
}
func (c *Client) CreateCreditCardWithToken(cc spp.CreditCard, token string) (string, error) {
	r := c.Called(cc, token)
	return r.String(0), r.Error(1)
}
func (c *Client) CreateAccount(acc lib.Account) (lib.Account, error) {
	r := c.Called(acc)
	a, _ := r.Get(0).(lib.Account)
	return a, r.Error(1)
}

func (c *Client) RegisterNewAccount(acc lib.Account) (lib.Account, error) {
	r := c.Called(acc)
	a, _ := r.Get(0).(lib.Account)
	return a, r.Error(1)
}

func (c *Client) GetAccount(name string) (account lib.Account, err error) {
	r := c.Called(name)
	acc, _ := r.Get(0).(lib.Account)
	return acc, r.Error(1)
}

func (c *Client) GetDefaultAccount() (account lib.Account, err error) {
	r := c.Called()
	acc, _ := r.Get(0).(lib.Account)
	return acc, r.Error(1)
}

func (c *Client) GetAccounts() (accounts lib.Accounts, err error) {
	r := c.Called()
	acc, _ := r.Get(0).([]lib.Account)
	return acc, r.Error(1)
}

func (c *Client) CreateDisc(name pathers.VirtualMachineName, disc brain.Disc) error {
	r := c.Called(name, disc)
	return r.Error(0)
}

func (c *Client) GetDisc(name pathers.VirtualMachineName, discId string) (disc brain.Disc, err error) {
	r := c.Called(name, discId)
	disc, _ = r.Get(0).(brain.Disc)
	return disc, r.Error(1)
}

func (c *Client) GetDiscByID(id int) (disc brain.Disc, err error) {
	r := c.Called(id)
	disc, _ = r.Get(0).(brain.Disc)
	return disc, r.Error(1)
}

func (c *Client) CreateGroup(name pathers.GroupName) error {
	r := c.Called(name)
	return r.Error(0)
}

func (c *Client) GetGroup(name pathers.GroupName) (brain.Group, error) {
	r := c.Called(name)
	group, _ := r.Get(0).(brain.Group)
	return group, r.Error(1)
}

func (c *Client) DeleteDisc(name pathers.VirtualMachineName, disc string) error {
	r := c.Called(name, disc)
	return r.Error(0)
}

func (c *Client) DeleteGroup(name pathers.GroupName) error {
	r := c.Called(name)
	return r.Error(0)
}

func (c *Client) DeleteVirtualMachine(name pathers.VirtualMachineName, purge bool) error {
	r := c.Called(name, purge)
	return r.Error(0)
}

func (c *Client) CreateVirtualMachine(group pathers.GroupName, vm brain.VirtualMachineSpec) (brain.VirtualMachine, error) {
	r := c.Called(group, vm)
	rvm, _ := r.Get(0).(brain.VirtualMachine)
	return rvm, r.Error(1)
}

func (c *Client) GetVirtualMachine(name pathers.VirtualMachineName) (vm brain.VirtualMachine, err error) {
	r := c.Called(name)
	vm, _ = r.Get(0).(brain.VirtualMachine)
	return vm, r.Error(1)
}

func (c *Client) MoveVirtualMachine(oldName pathers.VirtualMachineName, newName pathers.VirtualMachineName) error {
	r := c.Called(oldName, newName)
	return r.Error(0)
}

func (c *Client) ParseVirtualMachineName(name string, defaults ...pathers.VirtualMachineName) (pathers.VirtualMachineName, error) {
	r := c.Called(name, defaults)
	n, _ := r.Get(0).(pathers.VirtualMachineName)
	return n, r.Error(1)
}

func (c *Client) ParseGroupName(name string, defaults ...pathers.GroupName) pathers.GroupName {
	r := c.Called(name, defaults)
	n, _ := r.Get(0).(pathers.GroupName)
	return n
}

func (c *Client) ParseAccountName(name string, defaults ...string) string {
	r := c.Called(name, defaults)
	return r.String(0)
}

func (c *Client) ReimageVirtualMachine(name pathers.VirtualMachineName, image brain.ImageInstall) error {
	r := c.Called(name, image)
	return r.Error(0)
}

func (c *Client) ResetVirtualMachine(name pathers.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}

func (c *Client) ResizeDisc(name pathers.VirtualMachineName, id string, size int) error {
	r := c.Called(name, id, size)
	return r.Error(0)
}

func (c *Client) SetDiscIopsLimit(name pathers.VirtualMachineName, id string, size int) error {
	r := c.Called(name, id, size)
	return r.Error(0)
}

func (c *Client) RestartVirtualMachine(name pathers.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}
func (c *Client) SetVirtualMachineCDROM(name pathers.VirtualMachineName, url string) error {
	r := c.Called(name, url)
	return r.Error(0)
}
func (c *Client) SetVirtualMachineCores(name pathers.VirtualMachineName, cores int) error {
	r := c.Called(name, cores)
	return r.Error(0)
}
func (c *Client) SetVirtualMachineHardwareProfile(name pathers.VirtualMachineName, hwprofile string, locked ...bool) error {
	r := c.Called(name, hwprofile, locked)
	return r.Error(0)
}
func (c *Client) SetVirtualMachineHardwareProfileLock(name pathers.VirtualMachineName, locked bool) error {
	r := c.Called(name, locked)
	return r.Error(0)
}
func (c *Client) SetVirtualMachineMemory(name pathers.VirtualMachineName, memory int) error {
	r := c.Called(name, memory)
	return r.Error(0)
}
func (c *Client) StartVirtualMachine(name pathers.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}
func (c *Client) StopVirtualMachine(name pathers.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}
func (c *Client) ShutdownVirtualMachine(name pathers.VirtualMachineName, stayoff bool) error {
	r := c.Called(name, stayoff)
	return r.Error(0)
}

func (c *Client) UndeleteVirtualMachine(name pathers.VirtualMachineName) error {
	r := c.Called(name)
	return r.Error(0)
}

func (c *Client) CreateBackup(server pathers.VirtualMachineName, discLabelOrID string) (brain.Backup, error) {
	r := c.Called(server, discLabelOrID)
	snap, _ := r.Get(0).(brain.Backup)
	return snap, r.Error(1)
}
func (c *Client) DeleteBackup(server pathers.VirtualMachineName, discLabelOrID string, backupLabelOrID string) error {
	r := c.Called(server, discLabelOrID, backupLabelOrID)
	return r.Error(0)
}
func (c *Client) CreateBackupSchedule(server pathers.VirtualMachineName, discLabelOrID string, start string, interval int) (brain.BackupSchedule, error) {
	r := c.Called(server, discLabelOrID, start, interval)
	sched, _ := r.Get(0).(brain.BackupSchedule)
	return sched, r.Error(1)
}
func (c *Client) DeleteBackupSchedule(server pathers.VirtualMachineName, discLabelOrID string, id int) error {
	r := c.Called(server, discLabelOrID, id)
	return r.Error(0)
}
func (c *Client) GetBackups(server pathers.VirtualMachineName, discLabelOrID string) (brain.Backups, error) {
	r := c.Called(server, discLabelOrID)
	snaps, _ := r.Get(0).(brain.Backups)
	return snaps, r.Error(1)
}
func (c *Client) RestoreBackup(server pathers.VirtualMachineName, discLabelOrID string, backupLabelOrID string) (brain.Backup, error) {
	r := c.Called(server, discLabelOrID, backupLabelOrID)
	snap, _ := r.Get(0).(brain.Backup)

	return snap, r.Error(1)
}

func (c *Client) GetPrivileges(username string) (privs brain.Privileges, err error) {
	r := c.Called(username)
	privs, _ = r.Get(0).(brain.Privileges)
	return privs, r.Error(1)
}
func (c *Client) GetPrivilegesForAccount(accountName string) (privs brain.Privileges, err error) {
	r := c.Called(accountName)
	privs, _ = r.Get(0).(brain.Privileges)
	return privs, r.Error(1)
}
func (c *Client) GetPrivilegesForGroup(group pathers.GroupName) (privs brain.Privileges, err error) {
	r := c.Called(group)
	privs, _ = r.Get(0).(brain.Privileges)
	return privs, r.Error(1)
}
func (c *Client) GetPrivilegesForVirtualMachine(vm pathers.VirtualMachineName) (privs brain.Privileges, err error) {
	r := c.Called(vm)
	privs, _ = r.Get(0).(brain.Privileges)
	return privs, r.Error(1)
}
func (c *Client) GrantPrivilege(priv brain.Privilege) (err error) {
	r := c.Called(priv)
	return r.Error(0)
}
func (c *Client) RevokePrivilege(priv brain.Privilege) (err error) {
	r := c.Called(priv)
	return r.Error(0)
}
func (c *Client) GetVLANs() (brain.VLANs, error) {
	r := c.Called()
	vlans, _ := r.Get(0).(brain.VLANs)
	return vlans, r.Error(1)
}
func (c *Client) GetVLAN(num int) (brain.VLAN, error) {
	r := c.Called(num)
	vlans, _ := r.Get(0).(brain.VLAN)
	return vlans, r.Error(1)
}
func (c *Client) GetIPRanges() (brain.IPRanges, error) {
	r := c.Called()
	ipRanges, _ := r.Get(0).(brain.IPRanges)
	return ipRanges, r.Error(1)
}
func (c *Client) GetIPRange(idOrCIDR string) (brain.IPRange, error) {
	r := c.Called(idOrCIDR)
	ipRange, _ := r.Get(0).(brain.IPRange)
	return ipRange, r.Error(1)
}
func (c *Client) GetHeads() (brain.Heads, error) {
	r := c.Called()
	heads, _ := r.Get(0).(brain.Heads)
	return heads, r.Error(1)
}
func (c *Client) GetHead(idOrLabel string) (brain.Head, error) {
	r := c.Called(idOrLabel)
	head, _ := r.Get(0).(brain.Head)
	return head, r.Error(1)
}
func (c *Client) GetTails() (brain.Tails, error) {
	r := c.Called()
	tails, _ := r.Get(0).(brain.Tails)
	return tails, r.Error(1)
}
func (c *Client) GetTail(idOrLabel string) (brain.Tail, error) {
	r := c.Called(idOrLabel)
	tail, _ := r.Get(0).(brain.Tail)
	return tail, r.Error(1)
}
func (c *Client) GetStoragePools() (brain.StoragePools, error) {
	r := c.Called()
	storagePools, _ := r.Get(0).(brain.StoragePools)
	return storagePools, r.Error(1)
}
func (c *Client) GetStoragePool(idOrLabel string) (brain.StoragePool, error) {
	r := c.Called(idOrLabel)
	storagePool, _ := r.Get(0).(brain.StoragePool)
	return storagePool, r.Error(1)
}
func (c *Client) GetMigratingVMs() (brain.VirtualMachines, error) {
	r := c.Called()
	vms, _ := r.Get(0).(brain.VirtualMachines)
	return vms, r.Error(1)
}
func (c *Client) GetStoppedEligibleVMs() (brain.VirtualMachines, error) {
	r := c.Called()
	vms, _ := r.Get(0).(brain.VirtualMachines)
	return vms, r.Error(1)
}
func (c *Client) GetRecentVMs() (brain.VirtualMachines, error) {
	r := c.Called()
	vms, _ := r.Get(0).(brain.VirtualMachines)
	return vms, r.Error(1)
}
func (c *Client) MigrateDisc(disc int, newStoragePool string) error {
	r := c.Called(disc, newStoragePool)
	return r.Error(0)
}
func (c *Client) MigrateVirtualMachine(vmName pathers.VirtualMachineName, newHead string) error {
	r := c.Called(vmName, newHead)
	return r.Error(0)
}
func (c *Client) ReapVMs() error {
	r := c.Called()
	return r.Error(0)
}
func (c *Client) DeleteVLAN(id int) error {
	r := c.Called(id)
	return r.Error(0)
}
func (c *Client) AdminCreateGroup(name pathers.GroupName, vlanNum int) error {
	r := c.Called(name, vlanNum)
	return r.Error(0)
}
func (c *Client) CreateIPRange(ipRange string, vlanNum int) error {
	r := c.Called(ipRange, vlanNum)
	return r.Error(0)
}
func (c *Client) CancelDiscMigration(id int) error {
	r := c.Called(id)
	return r.Error(0)
}
func (c *Client) CancelVMMigration(id int) error {
	r := c.Called(id)
	return r.Error(0)
}
func (c *Client) EmptyStoragePool(idOrLabel string) error {
	r := c.Called(idOrLabel)
	return r.Error(0)
}
func (c *Client) EmptyHead(idOrLabel string) error {
	r := c.Called(idOrLabel)
	return r.Error(0)
}
func (c *Client) ReifyDisc(id int) error {
	r := c.Called(id)
	return r.Error(0)
}
func (c *Client) ApproveVM(name pathers.VirtualMachineName, powerOn bool) error {
	r := c.Called(name, powerOn)
	return r.Error(0)
}
func (c *Client) RejectVM(name pathers.VirtualMachineName, reason string) error {
	r := c.Called(name, reason)
	return r.Error(0)
}
func (c *Client) RegradeDisc(disc int, newGrade string) error {
	r := c.Called(disc, newGrade)
	return r.Error(0)
}
func (c *Client) UpdateVMMigration(name pathers.VirtualMachineName, speed *int64, downtime *int) error {
	r := c.Called(name, speed, downtime)
	return r.Error(0)
}
func (c *Client) CreateUser(username string, privilege string) error {
	r := c.Called(username, privilege)
	return r.Error(0)
}
func (c *Client) UpdateHead(idOrLabel string, options lib.UpdateHead) error {
	r := c.Called(idOrLabel, options)
	return r.Error(0)
}
func (c *Client) UpdateTail(idOrLabel string, options lib.UpdateTail) error {
	r := c.Called(idOrLabel, options)
	return r.Error(0)
}
func (c *Client) UpdateStoragePool(idOrLabel string, options brain.StoragePool) error {
	r := c.Called(idOrLabel, options)
	return r.Error(0)
}

func (c *Client) GetMigratingDiscs() (brain.Discs, error) {
	r := c.Called()
	discs := r.Get(0).(brain.Discs)
	return discs, r.Error(1)
}

func (c *Client) EnsureAccountName(name *pathers.AccountName) error {
	*name = "blank-account-name"
	return nil
}

func (c *Client) EnsureGroupName(name *pathers.GroupName) error {
	if name.Account == "" {
		name.Account = "blank-account-name"
	}
	if name.Group == "" {
		name.Group = "blank-group-name"
	}
	return nil
}

func (c *Client) EnsureVirtualMachineName(name *pathers.VirtualMachineName) error {
	if name.Account == "" {
		name.Account = "blank-account-name"
	}
	if name.Group == "" {
		name.Group = "blank-group-name"
	}
	if name.VirtualMachine == "" {
		name.VirtualMachine = "blank-vm-name"
	}
	return nil
}
