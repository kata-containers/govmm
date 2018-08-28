// +build !s390x

/*
// Copyright (c) 2018 Yash Jain
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

// Package qemu provides methods and types for launching and managing QEMU
// instances.  Instances can be launched with the LaunchQemu function and
// managed thereafter via QMPStart and the QMP object that this function
// returns.  To manage a qemu instance after it has been launched you need
// to pass the -qmp option during launch requesting the qemu instance to create
// a QMP unix domain manageent socket, e.g.,
// -qmp unix:/tmp/qmp-socket,server,nowait.  For more information see the
// example below.
package qemu

const (
	// NVDIMM is the Non Volatile DIMM device driver.
	NVDIMM DeviceDriver = "nvdimm"

	// Virtio9P is the 9pfs device driver.
	Virtio9P DeviceDriver = "virtio-9p-pci"

	// VirtioNet is the virt-io networking device driver.
	VirtioNet DeviceDriver = "virtio-net"

	// VirtioNetPCI is the virt-io pci networking device driver.
	VirtioNetPCI DeviceDriver = "virtio-net-pci"

	// VirtioSerial is the serial device driver.
	VirtioSerial DeviceDriver = "virtio-serial-pci"

	// VirtioBlock is the block device driver.
	VirtioBlock DeviceDriver = "virtio-blk"

	// Console is the console device driver.
	Console DeviceDriver = "virtconsole"

	// VirtioSerialPort is the serial port device driver.
	VirtioSerialPort DeviceDriver = "virtserialport"

	// VHostVSockPCI is the vhost vsock pci driver.
	VHostVSockPCI DeviceDriver = "vhost-vsock-pci"

	// Vfio is the vfio driver
	Vfio = "vfio-pci"

	// VirtioSCSI is the virtio-scsi device
	VirtioSCSI = "virtio-scsi-pci"

	//VhostUserSCSI represents a SCSI vhostuser device type
	VhostUserSCSI = "vhost-user-scsi-pci"

	//VhostUserNet represents a net vhostuser device type
	VhostUserNet = "virtio-net-pci"

	//VhostUserBlk represents a block vhostuser device type
	VhostUserBlk = "vhost-user-blk-pci"

	// VhostVSOCKPCI is the VSOCK vhost device type.
	VhostVSOCKPCI = "vhost-vsock-pci"

	// VhostVSOCK is a generic Vsock vhost device
	VhostVSOCK = VhostVSOCKPCI
)

// QemuDeviceParam converts to the QEMU -device parameter notation
func (n NetDeviceType) QemuDeviceParam() string {
	switch n {
	case TAP:
		return string(VirtioNet)
	case MACVTAP:
		return string(VirtioNet)
	case IPVTAP:
		return string(VirtioNet)
	case VETHTAP:
		return string(VirtioNet)
	case VFIO:
		return string(Vfio)
	case VHOSTUSER:
		return "" // -netdev type=vhost-user (no device)
	default:
		return ""
	}
}
