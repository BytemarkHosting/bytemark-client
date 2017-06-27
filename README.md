Bytemark command-line client
============================

[![Build Status](https://travis-ci.org/BytemarkHosting/bytemark-client.svg)](https://travis-ci.org/BytemarkHosting/bytemark-client) [![Coverage Status](https://coveralls.io/repos/github/BytemarkHosting/bytemark-client/badge.svg?branch=develop)](https://coveralls.io/github/BytemarkHosting/bytemark-client?branch=develop)

Installation
------------

If you're just looking to get started you can find the most recent stable
release on our [download page](https://github.com/BytemarkHosting/bytemark-client/releases).

If you wish to build, and install, the latest stable-release from source you can do so with:

    go get github.com/BytemarkHosting/bytemark-client/cmd/bytemark


Tracking Development
====================

If you prefer to track our in-development branch you can do that via:

    go get -d github.com/BytemarkHosting/bytemark-client/cmd/bytemark
    cd $GOPATH/src/github.com/BytemarkHosting/bytemark-client/
    git checkout develop
    cd cmd/bytemark
    go build


If you have problems to report, or contributions to make, feel free to [use the issue-tracker](https://github.com/BytemarkHosting/bytemark-client/issues)

Compatibility Guarantee
=======================

In go, semantic versioning is pretty hard. We guarantee that the following types of changes will not occur between minor versions within the lib package and all packages under its path (e.g. lib/brain)

* Publicly-exported functions and methods will not be removed, renamed, nor will their prototypes change.
* Publicly-exported struct-type fields will not be removed, renamed, nor will their types change.
* Publicly-exported variables and constants will not be removed, renamed, nor will their types change.

It's suggested that you avoid using struct embedding or interface composition with multiple types if any of those types are from bytemark-client/lib or any packages inside - bytemark-client's types are wont to have extra fields and methods added.

Breaking API change
===================

The following breaking API change to the 'lib' package occurred in version 3.0

* These `lib.Client` methods have been deleted
  * `ParseVirtualMachineName`
  * `ParseGroupName`
  * `ParseAccountName`

Most `lib.Client` methods now take `lib`/`brain`/`billing` structs as values, rather than as pointers. See the `lib/interface.go` for the full list of methods available and their new type signatures.

Almost all `lib.Client` struct fields are now values or slices of values instead of pointers or slices of pointers - below is a list. Two notable exceptions are `brain.VirtualMachineSpec.ImageInstall` and `brain.VirtualMachineSpec.IPs` - which may need to be null, and so remain as pointers.

* The type of `brain.Account.Groups` has changed from `[]*Group` to `[]Group`
* The type of `brain.Backups` has changed from `[]*Group` to `[]Group`
* The type of `brain.BackupSchedules` has changed from `[]*BackupSchedule` to `[]BackupSchedule`
* The type of `brain.Group.VirtualMachines` has changed from `[]*VirtualMachine` to `[]VirtualMachine`
* The type of `brain.Head.CCAddress` has changed from `[]*net.IP` to `[]net.IP`
* The type of `brain.IP.IP` has changed from `[]*net.IP` to `[]net.IP`
* The type of `brain.IPRange.Available` has changed from `[]*math/big.Int` to `[]math/big.Int`
* The type of `brain.IPCreateRequest.IPs` has changed from `[]*net.IP` to `[]net.IP`
* The type of `brain.IPs` has changed from `[]*net.IP` to `[]net.IP`
* The type of `brain.NetworkInterface.ExtraIPs` has changed from `map[string]*net.IP` to `map[string]net.IP`
* The type of `brain.Privileges` has changed from `[]*Privilege` to `[]Privilege`
* The type of `brain.Tail.CCAddress` has changed from `*net.IP` to `net.IP`
* The type of `brain.VirtualMachine.Discs` has changed from `[]*Disc` to `[]Disc`
* The type of `brain.VirtualMachine.ManagementAddress` has changed from `*net.IP` to `net.IP`
* The type of `brain.VirtualMachine.NetworkInterfaces` has changed from `[]*NetworkInterface` to `[]NetworkInterface`
* The type of `brain.VirtualMachineSpec.VirtualMachine` has changed from `*VirtualMachine` to `VirtualMachine`
* The type of `brain.VLAN.IPRanges` has changed from `[]*IPRange` to `[]IPRange`


If you require the old API for longer you can use `gopkg.in/BytemarkHosting/bytemark-client.v2/lib` to refer to the package prior to this change.
