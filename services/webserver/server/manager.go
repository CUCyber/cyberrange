package server

import (
	"context"
	"github.com/cucyber/cyberrange/pkg/proto"
	"github.com/cucyber/cyberrange/services/webserver/db"
	"google.golang.org/grpc"
	"time"
)

var conn *grpc.ClientConn
var client proto.ManagerClient

func CloseManager() {
	conn.Close()
}

func ListMachines() ([]*proto.Machine, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := client.List(ctx, &proto.Empty{})
	if err != nil {
		return nil, err
	}

	return response.Machines, nil
}

func CreateMachine(machine *db.Machine) error {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	req := &proto.Machine{Name: machine.Name}
	_, err := client.Create(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func StartMachine(machine *db.Machine) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &proto.Machine{Name: machine.Name}
	_, err := client.Start(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func StopMachine(machine *db.Machine) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &proto.Machine{Name: machine.Name}
	_, err := client.Stop(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func SnapshotMachine(machine *db.Machine) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := &proto.Machine{Name: machine.Name}
	_, err := client.Snapshot(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func RestartMachine(machine *db.Machine) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := &proto.Machine{Name: machine.Name}
	_, err := client.Restart(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func RevertMachine(machine *db.Machine) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := &proto.Machine{Name: machine.Name}
	_, err := client.Revert(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func UpdateMachines() {
	for {
		machines, err := ListMachines()
		if err != nil {
			panic(err.Error())
		}

		for _, v := range machines {
			db.SetMachineIp(&db.Machine{
				Name:      v.Name,
				IpAddress: v.Ip,
			})
		}

		time.Sleep(5 * time.Second)
	}
}

func InitializeManager() {
	var err error
	conn, err = grpc.Dial("10.0.144.12:8080", grpc.WithInsecure())
	if err != nil {
		panic(err.Error())
	}
	client = proto.NewManagerClient(conn)
}
