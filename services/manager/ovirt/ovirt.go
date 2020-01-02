package ovirt

import (
	"errors"
	"github.com/cucyber/cyberrange/pkg/proto"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
	"github.com/spf13/viper"
	"strings"
	"time"
)

var ovirt *ovirtsdk4.Connection
var prefix string

var (
	ErrMachineNotFound = errors.New("cyberrange: machine not found")
	ErrGetAttribute    = errors.New("cyberrange: failed to get object attribute")
	ErrNoSnapshots     = errors.New("cyberrange: the VM has no saved snapshots")
	ErrRevertFailed    = errors.New("cyberrange: revert to snapshot failed")
)

func CloseOVirt() {
	ovirt.Close()
}

func IsIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}

func ParseOVirtOptions() (map[string]string, error) {
	ovirtConfig := viper.New()
	ovirtConfig.SetConfigName("config")
	ovirtConfig.SetConfigType("yaml")
	ovirtConfig.AddConfigPath("./ovirt/")

	if err := ovirtConfig.ReadInConfig(); err != nil {
		return nil, err
	}

	return ovirtConfig.GetStringMapString("ovirt"), nil
}

func WaitForState(vmId string, target ovirtsdk4.VmStatus) error {
	vmsService := ovirt.SystemService().VmsService()
	vmService := vmsService.VmService(vmId)

	for {
		vmResp, err := vmService.Get().Send()
		if err != nil {
			return err
		}

		vm, ok := vmResp.Vm()
		if !ok {
			continue
		}

		if status, ok := vm.Status(); ok {
			if status == target {
				break
			}
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}

func GetVMServiceByName(MachineName string) (*ovirtsdk4.VmService, error) {
	vmsService := ovirt.SystemService().VmsService()

	vmsResp, err := vmsService.List().Search("name=" + MachineName).Send()
	if err != nil {
		return nil, err
	}

	vmsSlice, ok := vmsResp.Vms()
	if !ok {
		return nil, ErrGetAttribute
	}

	if len(vmsSlice.Slice()) == 0 {
		return nil, ErrMachineNotFound
	}

	vm := vmsSlice.Slice()[0]

	vmId, ok := vm.Id()
	if !ok {
		return nil, ErrGetAttribute
	}

	return vmsService.VmService(vmId), nil
}

func StartMachine(MachineName string) error {
	vmService, err := GetVMServiceByName(MachineName)
	if err != nil {
		return err
	}

	_, err = vmService.Start().Send()
	if err != nil {
		return err
	}

	return nil
}

func StopMachine(MachineName string) error {
	vmService, err := GetVMServiceByName(MachineName)
	if err != nil {
		return err
	}

	_, err = vmService.Stop().Send()
	if err != nil {
		return err
	}

	return nil
}

func RestartMachine(MachineName string) error {
	vmService, err := GetVMServiceByName(MachineName)
	if err != nil {
		return err
	}

	_, err = vmService.Reboot().Send()
	if err != nil {
		return err
	}

	return nil
}

func SnapshotMachine(MachineName string) error {
	MachineName = prefix + MachineName

	vmsService := ovirt.SystemService().VmsService()

	vmsResp, err := vmsService.List().Search("name=" + MachineName).Send()
	if err != nil {
		return err
	}

	vmsSlice, ok := vmsResp.Vms()
	if !ok {
		return ErrGetAttribute
	}

	if len(vmsSlice.Slice()) == 0 {
		return ErrMachineNotFound
	}

	vm := vmsSlice.Slice()[0]

	vmId, ok := vm.Id()
	if !ok {
		return ErrGetAttribute
	}

	snapshotsService := vmsService.VmService(vmId).SnapshotsService()

	build, err := ovirtsdk4.NewSnapshotBuilder().Description("Initial").Build()
	if err != nil {
		return err
	}

	_, err = snapshotsService.Add().Snapshot(build).Send()
	if err != nil {
		return err
	}

	return nil
}

func RevertMachine(MachineName string) error {
	MachineName = prefix + MachineName

	vmsService := ovirt.SystemService().VmsService()

	vmsResp, err := vmsService.List().Search("name=" + MachineName).Send()
	if err != nil {
		return err
	}

	vmsSlice, ok := vmsResp.Vms()
	if !ok {
		return ErrGetAttribute
	}

	if len(vmsSlice.Slice()) == 0 {
		return ErrMachineNotFound
	}

	vm := vmsSlice.Slice()[0]

	vmId, ok := vm.Id()
	if !ok {
		return ErrGetAttribute
	}

	vmService := vmsService.VmService(vmId)

	snapsService := vmService.SnapshotsService()

	snapResponse, err := snapsService.List().Send()
	if err != nil {
		return err
	}

	snapshots, ok := snapResponse.Snapshots()
	if !ok {
		return ErrGetAttribute
	}

	if len(snapshots.Slice()) <= 1 {
		return ErrNoSnapshots
	}

	for _, snapshot := range snapshots.Slice() {
		snapId, ok := snapshot.Id()
		if !ok {
			continue
		}

		snapName, ok := snapshot.Description()
		if !ok || snapName != "Initial" {
			continue
		}

		snapService := snapsService.SnapshotService(snapId)

		_, err = vmService.Stop().Send()
		if err != nil {
			return err
		}

		err = WaitForState(vmId, ovirtsdk4.VMSTATUS_DOWN)
		if err != nil {
			return err
		}

		_, err := snapService.Restore().Async(false).RestoreMemory(false).Send()
		if err != nil {
			return err
		}

		err = WaitForState(vmId, ovirtsdk4.VMSTATUS_DOWN)
		if err != nil {
			return err
		}

		_, err = vmService.Start().Send()
		if err != nil {
			return err
		}

		err = WaitForState(vmId, ovirtsdk4.VMSTATUS_UP)
		if err != nil {
			return err
		}

		return nil
	}

	return ErrRevertFailed
}

func ListMachines() (*proto.MachineList, error) {
	var machines []*proto.Machine

	vmsService := ovirt.SystemService().VmsService()

	vmsResponse, err := vmsService.List().Send()
	if err != nil {
		return &proto.MachineList{Machines: nil}, err
	}

	vms, ok := vmsResponse.Vms()
	if !ok {
		return &proto.MachineList{Machines: nil}, ErrGetAttribute
	}

	for _, vm := range vms.Slice() {
		vmName, ok := vm.Name()
		if !ok || !strings.HasPrefix(vmName, prefix) {
			continue
		}

		vmId, ok := vm.Id()
		if !ok {
			continue
		}

		vmService := vmsService.VmService(vmId)

		devicesResponse, err := vmService.ReportedDevicesService().List().Send()
		if err != nil {
			return &proto.MachineList{Machines: nil}, err
		}

		devices, ok := devicesResponse.ReportedDevice()
		if !ok {
			return &proto.MachineList{Machines: nil}, ErrGetAttribute
		}

		var ipv4Address string

		for _, device := range devices.Slice() {
			ips, ok := device.Ips()
			if !ok {
				continue
			}

			for _, ip := range ips.Slice() {
				address, ok := ip.Address()
				if !ok {
					continue
				}

				if IsIPv4(address) {
					ipv4Address = address
				}
			}
		}

		vmName = strings.TrimPrefix(vmName, prefix)
		machines = append(machines, &proto.Machine{Name: vmName, Ip: ipv4Address})
	}

	return &proto.MachineList{Machines: machines}, nil
}

func Connect() (*ovirtsdk4.Connection, error) {
	ovirtoptions, err := ParseOVirtOptions()
	if err != nil {
		panic(err.Error())
	}
	prefix = ovirtoptions["prefix"]

	ovirt, err := ovirtsdk4.NewConnectionBuilder().
		URL(ovirtoptions["apiurl"]).
		Username(ovirtoptions["user"]).
		Password(ovirtoptions["pass"]).
		Insecure(true).
		Compress(true).
		Timeout(time.Second * 120).
		Build()
	if err != nil {
		return nil, err
	}

	return ovirt, nil
}

func InitializeOVirt() {
	conn, err := Connect()
	if err != nil {
		panic(err.Error())
	}
	ovirt = conn
}
