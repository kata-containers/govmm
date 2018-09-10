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

import "testing"

// -pci devices don't play well with Z hence replace them with corresponding -ccw devices
// See https://wiki.qemu.org/Documentation/Platforms/S390X
var (
	deviceFSString                 = "-device virtio-9p-ccw,disable-modern=true,fsdev=workload9p,mount_tag=rootfs -fsdev local,id=workload9p,path=/var/lib/docker/devicemapper/mnt/e31ebda2,security_model=none"
	deviceNetworkString            = "-netdev tap,id=tap0,vhost=on,ifname=ceth0,downscript=no,script=no -device driver=virtio-net,netdev=tap0,mac=01:02:de:ad:be:ef,disable-modern=true"
	deviceNetworkStringMq          = "-netdev tap,id=tap0,vhost=on,fds=3:4 -device driver=virtio-net,netdev=tap0,mac=01:02:de:ad:be:ef,disable-modern=true,mq=on,vectors=6"
	deviceNetworkPCIString         = "-netdev tap,id=tap0,vhost=on,ifname=ceth0,downscript=no,script=no -device driver=virtio-net,netdev=tap0,mac=01:02:de:ad:be:ef,bus=/pci-bus/pcie.0,addr=ff,disable-modern=true"
	deviceNetworkPCIStringMq       = "-netdev tap,id=tap0,vhost=on,fds=3:4 -device driver=virtio-net,netdev=tap0,mac=01:02:de:ad:be:ef,bus=/pci-bus/pcie.0,addr=ff,disable-modern=true,mq=on,vectors=6"
	deviceSerialString             = "-device virtio-serial-ccw,disable-modern=true,id=serial0"
	deviceVhostUserNetString       = "-chardev socket,id=char1,path=/tmp/nonexistentsocket.socket -netdev type=vhost-user,id=net1,chardev=char1,vhostforce -device virtio-net-ccw,netdev=net1,mac=00:11:22:33:44:55"
	deviceVSOCKString              = "-device vhost-vsock-ccw,disable-modern=true,id=vhost-vsock-pci0,guest-cid=4"
	deviceVFIOString               = "-device vfio-ccw,host=02:10.0"
	deviceSCSIControllerStr        = "-device virtio-scsi-ccw,id=foo"
	deviceSCSIControllerBusAddrStr = "-device virtio-scsi-ccw,id=foo,bus=pci.0,addr=00:04.0,disable-modern=true,iothread=iothread1"
	deviceVhostUserSCSIString      = "-chardev socket,id=char1,path=/tmp/nonexistentsocket.socket -device vhost-user-scsi-pci,id=scsi1,chardev=char1"
	deviceVhostUserBlkString       = "-chardev socket,id=char2,path=/tmp/nonexistentsocket.socket -device ,logical_block_size=4096,size=512M,chardev=char2"
)

func TestAppendDeviceVhostUser(t *testing.T) {

	vhostuserSCSIDevice := VhostUserDevice{
		SocketPath:    "/tmp/nonexistentsocket.socket",
		CharDevID:     "char1",
		TypeDevID:     "scsi1",
		Address:       "",
		VhostUserType: VhostUserSCSI,
	}
	testAppend(vhostuserSCSIDevice, deviceVhostUserSCSIString, t)

	vhostuserNetDevice := VhostUserDevice{
		SocketPath:    "/tmp/nonexistentsocket.socket",
		CharDevID:     "char1",
		TypeDevID:     "net1",
		Address:       "00:11:22:33:44:55",
		VhostUserType: VhostUserNet,
	}
	testAppend(vhostuserNetDevice, deviceVhostUserNetString, t)
}
