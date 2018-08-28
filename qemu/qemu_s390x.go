// +build s390x

/*
// Copyright (c) 2018 Yash Jain
// Copyright (c) 2016 Intel Corporation
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

package qemu

// -pci devices don't play well with Z hence corresponding -ccw devices
// See https://wiki.qemu.org/Documentation/Platforms/S390X
const (
	// NVDIMM is the Non Volatile DIMM device driver.
	NVDIMM DeviceDriver = "nvdimm"

	// Virtio9P is the 9pfs device driver.
	Virtio9P DeviceDriver = "virtio-9p-ccw"

	// VirtioNet is the virt-io networking device driver.
	VirtioNet DeviceDriver = "virtio-net"

	// VirtioNetPCI is the virt-io pci networking device driver.
	VirtioNetPCI DeviceDriver = "virtio-net-ccw"

	// VirtioSerial is the serial device driver.
	VirtioSerial DeviceDriver = "virtio-serial-ccw"

	// VirtioBlock is the block device driver.
	VirtioBlock DeviceDriver = "virtio-blk"

	// Console is the console device driver.
	Console DeviceDriver = "virtconsole"

	// VirtioSerialPort is the serial port device driver.
	VirtioSerialPort DeviceDriver = "virtserialport"

	// VHostVSockPCI is the vhost vsock pci driver.
	VHostVSockPCI DeviceDriver = "vhost-vsock-ccw"

	// Vfio is the vfio driver
	Vfio DeviceDriver = "vfio-ccw"

	// VirtioSCSI is the virtio-scsi device
	VirtioSCSI DeviceDriver = "virtio-scsi-ccw"

	//VhostUserSCSI represents a SCSI vhostuser device type
	VhostUserSCSI = "vhost-user-scsi-pci"

	//VhostUserNet represents a net vhostuser device type
	VhostUserNet = "virtio-net-ccw"

	//VhostUserBlk represents a block vhostuser device type
	VhostUserBlk = ""

	// VhostVSOCKCCW is the VSOCK vhost device type for s390x.
	VhostVSOCKCCW = "vhost-vsock-ccw"

	// VhostVSOCK is a generic Vsock Device
	VhostVSOCK = VhostVSOCKCCW
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
		return string(VirtioNet) // -netdev type=tap -device virtio-net-pci
	case VFIO:
		return string(Vfio) // -device vfio-pci (no netdev) or vfio-ccw on Z
	case VHOSTUSER:
		return "" // -netdev type=vhost-user (no device)
	default:
		return ""
	}
}
