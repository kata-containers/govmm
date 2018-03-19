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
	"net"
	"reflect"
	"testing"
)

func TestConfigGood(t *testing.T) {
	cfg := Config{
		RootFS: "crosvm/rootfs.ext4",
		Socket: "/run/user/1000/crosvm.sock",
		Kernel: Kernel{
			Path:   "linux/arch/x86/boot/compressed/vmlinux.bin",
			Params: "ip=192.168.30.2::192.168.30.1:255.255.255.0:crosvm:eth0",
		},
		Memory: Memory{
			Size: 1024,
		},
		SMP: SMP{
			CPUs: 2,
		},
		Net: NetDevice{
			MACAddress: "AA:BB:CC:00:00:12",
			HostIP:     net.ParseIP("192.168.30.1"),
			NetMask:    net.IPv4Mask(255, 255, 255, 0),
			VHost:      true,
		},
		Sec: Security{
			DisableSandbox: true,
		},
	}

	args, err := cfg.appendArgs(nil)
	if err != nil {
		t.Errorf("appendArgs failed: %v", err)
	}

	expected := []string{
		"-s", "/run/user/1000/crosvm.sock",
		"-r", "crosvm/rootfs.ext4",
		"--mac", "AA:BB:CC:00:00:12",
		"--host_ip", "192.168.30.1",
		"--netmask", "255.255.255.0",
		"--vhost-net",
		"-m", "1024", "-c", "2",
		"--disable-sandbox",
		"-p", "ip=192.168.30.2::192.168.30.1:255.255.255.0:crosvm:eth0",
		"linux/arch/x86/boot/compressed/vmlinux.bin",
	}

	if !reflect.DeepEqual(expected, args) {
		t.Errorf("result of appendArgs does not match expected result")
		t.Errorf("Got %v", args)
		t.Errorf("Expected %v", expected)
	}

	// With security

	cfg.Sec.DisableSandbox = false
	cfg.Sec.SeccompPolicyDir = "secdir"

	expected = []string{
		"-s", "/run/user/1000/crosvm.sock",
		"-r", "crosvm/rootfs.ext4",
		"--mac", "AA:BB:CC:00:00:12",
		"--host_ip", "192.168.30.1",
		"--netmask", "255.255.255.0",
		"--vhost-net",
		"-m", "1024", "-c", "2",
		"-u", "--seccomp-policy-dir", "secdir",
		"-p", "ip=192.168.30.2::192.168.30.1:255.255.255.0:crosvm:eth0",
		"linux/arch/x86/boot/compressed/vmlinux.bin",
	}

	args, err = cfg.appendArgs(nil)
	if err != nil {
		t.Errorf("appendArgs failed: %v", err)
	}

	if !reflect.DeepEqual(expected, args) {
		t.Errorf("result of appendArgs does not match expected result")
		t.Errorf("Got %v", args)
		t.Errorf("Expected %v", expected)
	}

	// With Disks but without network, mem or CPU

	cfg.Net = NetDevice{}
	cfg.Memory = Memory{}
	cfg.SMP = SMP{}
	cfg.Disks = []Disk{
		{
			Type:      QCOW,
			Writeable: true,
			Path:      "image.qcow",
		},
		{
			Type:      FlatFile,
			Writeable: false,
			Path:      "image.ext4",
		},
	}

	args, err = cfg.appendArgs(nil)
	if err != nil {
		t.Errorf("appendArgs failed: %v", err)
	}

	expected = []string{
		"-s", "/run/user/1000/crosvm.sock",
		"-r", "crosvm/rootfs.ext4",
		"--rwqcow", "image.qcow",
		"--disk", "image.ext4",
		"-u", "--seccomp-policy-dir", "secdir",
		"-p", "ip=192.168.30.2::192.168.30.1:255.255.255.0:crosvm:eth0",
		"linux/arch/x86/boot/compressed/vmlinux.bin",
	}

	if !reflect.DeepEqual(expected, args) {
		t.Errorf("result of appendArgs does not match expected result")
		t.Errorf("Got %v", args)
		t.Errorf("Expected %v", expected)
	}
}

func TestConfigBad(t *testing.T) {
	cfg := Config{}

	_, err := cfg.appendArgs(nil)
	if err == nil {
		t.Errorf("appendArgs on an empty config expected to fail")
	}

	minimalCfg := Config{
		RootFS: "crosvm/rootfs.ext4",
		Socket: "/run/user/1000/crosvm.sock",
		Kernel: Kernel{
			Path: "linux/arch/x86/boot/compressed/vmlinux.bin",
		},
	}

	_, err = minimalCfg.appendArgs(nil)
	if err != nil {
		t.Errorf("appendArgs on minimal config failed: %v", err)
	}

	cfg = minimalCfg
	cfg.Disks = []Disk{
		{
			Type:      QCOW,
			Writeable: true,
		},
	}
	_, err = cfg.appendArgs(nil)
	if err == nil {
		t.Errorf("appendArgs on a config with missing disk path expected to fail")
	}

	cfg = minimalCfg
	cfg.Net = NetDevice{
		MACAddress: "AA:BB:CC:00:00:12",
		NetMask:    net.IPv4Mask(255, 255, 255, 0),
	}
	_, err = cfg.appendArgs(nil)
	if err == nil {
		t.Errorf("appendArgs on a config with missing HostIP expected to fail")
	}
}
