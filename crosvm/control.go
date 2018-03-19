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
	"context"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/intel/govmm/vmmlog"
)

// LaunchCustomCrosvm launches a crosvm guest using the parameters specified in the args
// slice.  The attrs parameter can be used to control aspects of the newly created
// crosvm process, such as the user and group under which it runs.  It may be nil.  If
// missing, LaunchCustomCrosvm will search your path for the crosvm binary.    This function blocks
// until the guest either fails to start or exits.
func LaunchCustomCrosvm(ctx context.Context, path string, args []string,
	attr *syscall.SysProcAttr, logger vmmlog.Log) error {
	if logger == nil {
		logger = vmmlog.NullLogger{}
	}
	args = append([]string{"run"}, args...)
	return execCrosvm(ctx, path, args, attr, logger)
}

// LaunchCrosvm launches a crosvm guest using the parameters specified in the config
// object.  The attrs parameter can be used to control aspects of the newly created
// crosvm process, such as the user and group under which it runs.  It may be nil.
// The path to the crosvm binary is specified in the config.Path field.  If missing,
// LaunchCrosvm will search your path for the crosvm binary.  This function blocks
// until the guest either fails to start or exits.
func LaunchCrosvm(ctx context.Context, config Config, attr *syscall.SysProcAttr,
	logger vmmlog.Log) error {
	args, err := config.appendArgs(nil)
	if err != nil {
		return err
	}

	return LaunchCustomCrosvm(ctx, config.Path, args, attr, logger)
}

// LaunchCrosvmAsync launches a crosvm guest using the parameters specified in the
// config object.  The function blocks until it has detected that the crosvm process
// has successfully launched and the domain socket is open and ready to receive
// commands.  As the function may return before the guest exits, it is not possible
// to retrieve the exit status of the guest if it starts correctly.  Once this
// function returns the ctx parameter cannot be used to stop the guest.  To
// stop the guest, users will need to call ExecuteStop.
// The parameters have the same names and meaning as LaunchCrosvm.
func LaunchCrosvmAsync(ctx context.Context, config Config, attr *syscall.SysProcAttr,
	logger vmmlog.Log) error {
	if logger == nil {
		logger = vmmlog.NullLogger{}
	}

	var launchErr error
	localCtx, cancelFn := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func(ctx context.Context, errCh chan<- error) {
		errCh <- LaunchCrosvm(ctx, config, nil, logger)
	}(localCtx, errCh)

	waitCompleteCh := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		select {
		case launchErr = <-errCh:
			cancelFn()
		case <-ctx.Done():
			cancelFn()
		case <-waitCompleteCh:
		}
		wg.Done()
	}()

	logger.Infof("Waiting for guest to start")

	err := WaitForGuest(localCtx, config.Socket, logger)
	close(waitCompleteCh)
	wg.Wait()
	if launchErr != nil {
		logger.Infof("Guest failed to launch: %v", launchErr)
		return launchErr
	}

	if err != nil {
		logger.Infof("Launch cancelled: %v", err)
	}

	return err
}

func execCrosvm(ctx context.Context, path string, args []string,
	attr *syscall.SysProcAttr, logger vmmlog.Log) error {
	if path == "" {
		path = "crosvm"
	}

	logger.Infof("running %s %s", path, args)

	cmd := exec.CommandContext(ctx, path, args...)
	cmd.SysProcAttr = attr

	return cmd.Run()
}

// ExecuteStop attempts to stop one or more crosvm guests.  The path to the crosvm
// binary is specified in the config.Path field.  If missing, ExecuteStop will search
// your path for the crosvm binary.  The domain sockets of the guests to be stopped are
// provided in the sockets slice.  This function actually launches the crosvm binary
// to stop the guests.  The attrs parameter can be used to control aspects
// of this newly created crosvm process, such as the user and group under which it runs.
// It may be nil.
func ExecuteStop(ctx context.Context, path string, sockets []string,
	attr *syscall.SysProcAttr, logger vmmlog.Log) error {
	if len(sockets) == 0 {
		return errors.New("Expected at least one socket path")
	}

	if logger == nil {
		logger = vmmlog.NullLogger{}
	}

	args := append([]string{"stop"}, sockets...)
	return execCrosvm(ctx, path, args, attr, logger)
}

// ExecuteBalloon attempts to adjust the balloon size of one or more crosvm guests by the
// value specified in the numPages parameter.  The path to the crosvm binary is specified
// in the config.Path field.  If missing, ExecuteBalloon will search your path for the crosvm
// binary.  The domain sockets of the guests are provided in the sockets slice.
// This function actually launches the crosvm binary to adjust the balloon size of the guests.
// The attrs parameter can be used to control aspects of this newly created crosvm process,
// such as the user and group under which it runs. It may be nil.
func ExecuteBalloon(ctx context.Context, path string, numPages int32, sockets []string,
	attr *syscall.SysProcAttr, logger vmmlog.Log) error {
	if len(sockets) == 0 {
		return errors.New("Expected at least one socket path")
	}

	if logger == nil {
		logger = vmmlog.NullLogger{}
	}

	args := append([]string{"balloon", fmt.Sprintf("%d", numPages)}, sockets...)
	return execCrosvm(ctx, path, args, attr, logger)
}

// WaitForGuest blocks until the VM identified by the path to a domain socket in the
// path parameter has started or the ctx object is cancelled or times out.
func WaitForGuest(ctx context.Context, path string, logger vmmlog.Log) error {
	var conn net.Conn
	var err error

	if logger == nil {
		logger = vmmlog.NullLogger{}
	}

	dialer := net.Dialer{}

	for {
		conn, err = dialer.DialContext(ctx, "unixgram", path)
		if err == nil {
			break
		}

		logger.Infof("Attempt to open guest socket failed %v", err)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Millisecond * 50):
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
		}

		logger.Infof("Trying to open guest socket again")
	}

	logger.Infof("Successfully connected to guest socket")

	_ = conn.Close()
	return nil
}
