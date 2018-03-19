/*
// Copyright (c) 2018 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
*/

package crosvm

import (
	"errors"
	"fmt"
	"net"
)

// SMP is the multi processors configuration structure.
type SMP struct {
	// CPUs is the number of VCPUs made available to qemu.
	CPUs uint
}

// Memory is the guest memory configuration structure.
type Memory struct {
	// Size is the amount of memory made available to the guest in MiB
	Size uint
}

// Kernel is the guest kernel configuration structure.
type Kernel struct {
	// Path is the guest kernel path on the host filesystem.
	Path string

	// Params is the kernel parameters string.
	Params string
}

// NetDevice represents a guest networking device
type NetDevice struct {
	// VHost enables virtio device emulation from the host kernel instead of from crosvm.
	VHost bool

	// MACAddress is the networking device interface MAC address.
	MACAddress string

	// IPv4 address to assign to the host TAP interface
	HostIP net.IP

	// Netmask for the subnet of the virtual machine
	NetMask net.IPMask
}

// DiskType indicates that the type of a disk image
type DiskType string

const (
	// FlatFile indicates that an image is a flat file
	FlatFile DiskType = "flatfile"

	// QCOW indicates that an image is a QCOW file
	QCOW = "qcow"
)

// Disk represents a host disk image to be made accessible in the guest
type Disk struct {
	// Type provides the type of the disk image
	Type DiskType

	// Writeable indicates whether the image is writable or not
	Writeable bool

	// Path provides the path on the host to the disk image
	Path string
}

// Security contains the security settings for the guest
type Security struct {
	// All devices run in a single non-sanboxed process
	DisableSandbox bool

	// Path on host to seccomp policy files
	SeccompPolicyDir string
}

// Config is the crosvm configuration structure used to configure crosvm
// virtual machine instances.
type Config struct {
	// Path is the crosvm binary path.
	Path string

	// RootFS is a path to a read-only root image
	RootFS string

	// Socket is a path to domain socket used to control the guest
	Socket string

	// Disks a slice of disks to attach to the guest
	Disks []Disk

	// Kernel is the guest kernel configuration.
	Kernel Kernel

	// Memory is the guest memory configuration.
	Memory Memory

	// SMP is the quest multi processors configuration.
	SMP SMP

	// Net specifies the networking configuration of the guest
	Net NetDevice

	// Sec specifies the security settings for the guest
	Sec Security
}

func (s SMP) appendArgs(args []string) ([]string, error) {
	if s.CPUs > 0 {
		args = append(args, "-c", fmt.Sprintf("%d", s.CPUs))
	}
	return args, nil
}

func (m Memory) appendArgs(args []string) ([]string, error) {
	if m.Size > 0 {
		args = append(args, "-m", fmt.Sprintf("%d", m.Size))
	}
	return args, nil
}

// Must be the last function called when appending arguments
func (k Kernel) appendArgs(args []string) ([]string, error) {
	if k.Path == "" {
		return nil, errors.New("Kernel path must be specified")
	}

	if k.Params != "" {
		args = append(args, "-p", k.Params)
	}

	return append(args, k.Path), nil
}

func (n NetDevice) appendArgs(args []string) ([]string, error) {
	if n.MACAddress == "" && len(n.HostIP) == 0 && len(n.NetMask) == 0 {
		return args, nil
	}

	if n.MACAddress == "" || len(n.HostIP) == 0 || len(n.NetMask) == 0 {
		return nil, errors.New("MACAddress, HostIP and NetMask must be defined")
	}

	args = append(args, "--mac", n.MACAddress)
	args = append(args, "--host_ip", n.HostIP.String())
	args = append(args, "--netmask", net.IP(n.NetMask).String())

	if n.VHost {
		args = append(args, "--vhost-net")
	}

	return args, nil
}

func (d Disk) appendArgs(args []string) ([]string, error) {
	var opt string

	if d.Path == "" {
		return nil, errors.New("Path not specified for disk")
	}

	if d.Type == FlatFile {
		opt = "disk"
	} else if d.Type == QCOW {
		opt = "qcow"
	} else {
		return nil, fmt.Errorf("Unknown disk type %s", d.Type)
	}

	if d.Writeable {
		opt = "rw" + opt
	}

	return append(args, "--"+opt, d.Path), nil
}

func (c Security) appendArgs(args []string) ([]string, error) {
	if c.DisableSandbox {
		return append(args, "--disable-sandbox"), nil
	}

	return append(args, "-u", "--seccomp-policy-dir", c.SeccompPolicyDir), nil
}

func (c Config) appendArgs(args []string) ([]string, error) {
	if c.RootFS == "" && len(c.Disks) == 0 {
		return nil, errors.New("No rootfs or disks specified")
	}

	if c.Socket == "" {
		return nil, errors.New("Socket path must be defined")
	}
	args = append(args, "-s", c.Socket)

	if c.RootFS != "" {
		args = append(args, "-r", c.RootFS)
	}

	var err error
	for _, d := range c.Disks {
		args, err = d.appendArgs(args)
		if err != nil {
			return nil, err
		}
	}

	for _, fn := range []func([]string) ([]string, error){
		c.Net.appendArgs,
		c.Memory.appendArgs,
		c.SMP.appendArgs,
		c.Sec.appendArgs,
		c.Kernel.appendArgs, // order is important, kernel must be last
	} {
		args, err = fn(args)
		if err != nil {
			return nil, err
		}
	}

	return args, nil
}
