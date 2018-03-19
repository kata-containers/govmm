# Virtual Machine Manager for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/intel/govmm)](https://goreportcard.com/report/github.com/intel/govmm)
[![Build Status](https://travis-ci.org/intel/govmm.svg?branch=master)](https://travis-ci.org/intel/govmm)
[![GoDoc](https://godoc.org/github.com/intel/govmm/qemu?status.svg)](https://godoc.org/github.com/intel/govmm/qemu)
[![Coverage Status](https://coveralls.io/repos/github/intel/govmm/badge.svg?branch=master)](https://coveralls.io/github/intel/govmm?branch=master)

Virtual Machine Manager for Go (govmm) is a suite of packages that
provide Go APIs for creating and managing virtual machines.  There's
currently support for two hypervisors; qemu/kvm and crosvm.

The github.com/intel/govmm/qemu package provides APIs for launching
qemu instances and for managing those instances via QMP, once launched.
VM instances can be stopped, have devices attached to them and monitored
for events via the qemu APIs.

The github.com/intel/govmm/crosvm package provides basic support for
launching and stopping crosvm instances.

The govmm packages have no external dependencies apart from the Go
standard library and so are nice and easy to vendor inside other
projects.
