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

package crosvm_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/intel/govmm/crosvm"
)

type StderrLogger struct{}

func (l StderrLogger) V(level int32) bool {
	return true
}

func (l StderrLogger) Infof(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
	fmt.Fprintln(os.Stderr, "")
}

func (l StderrLogger) Warningf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
	fmt.Fprintln(os.Stderr, "")
}

func (l StderrLogger) Errorf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
	fmt.Fprintln(os.Stderr, "")
}

func Example() {
	socket := fmt.Sprintf("/run/user/%d/crosvm.sock", os.Getuid())
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	cfg := crosvm.Config{
		RootFS: filepath.Join(user.HomeDir, "crosvm/rootfs.ext4"),
		Socket: socket,
		Kernel: crosvm.Kernel{
			Path:   filepath.Join(user.HomeDir, "linux/arch/x86/boot/compressed/vmlinux.bin"),
			Params: "ip=192.168.30.2::192.168.30.1:255.255.255.0:crosvm:eth0",
		},
		Net: crosvm.NetDevice{
			MACAddress: "AA:BB:CC:00:00:12",
			HostIP:     net.ParseIP("192.168.30.1"),
			NetMask:    net.IPv4Mask(255, 255, 255, 0),
		},
		Sec: crosvm.Security{
			DisableSandbox: true,
		},
	}

	logger := StderrLogger{}
	err = crosvm.LaunchCrosvmAsync(context.Background(), cfg, nil, logger)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stderr, "Crosvm guest running.")

	fmt.Fprintln(os.Stderr, "Adjusting balloon size of guest by 64 pages")
	err = crosvm.ExecuteBalloon(context.Background(), "", 64, []string{socket},
		nil, logger)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to adjust balloon size of guest")
	}

	fmt.Fprintln(os.Stderr, "Wait for 10 seconds")

	<-time.After(time.Second * 10)

	err = crosvm.ExecuteStop(context.Background(), "", []string{socket}, nil,
		logger)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to stop crosvm guest")
	}
}
