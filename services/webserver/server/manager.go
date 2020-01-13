package server

import (
	"context"
	"time"

	"github.com/cucyber/cyberrange/pkg/proto"
	"github.com/cucyber/cyberrange/services/webserver/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
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
		c.logger.Printf("ERROR: ListMachines: %+v", err)
		return nil, err
	}

	return response.Machines, nil
}

func CheckCreateMachine(machine *db.Machine) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &proto.Machine{Name: machine.Name}
	_, err := client.CheckCreate(ctx, req)
	if err != nil {
		c.logger.Printf("ERROR: CheckCreateMachine: %+v", err)
		return err
	}

	return nil
}

func CreateMachine(machine *db.Machine) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	req := &proto.Machine{Name: machine.Name}
	_, err := client.Create(ctx, req)
	if err != nil {
		c.logger.Printf("ERROR: CreateMachine: %+v", err)
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
		c.logger.Printf("ERROR: StartMachine: %+v", err)
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
		c.logger.Printf("ERROR: StopMachine: %+v", err)
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
		c.logger.Printf("ERROR: SnapshotMachine: %+v", err)
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
		c.logger.Printf("ERROR: RestartMachine: %+v", err)
		return err
	}

	return nil
}

func RevertMachine(machine *db.Machine) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	req := &proto.Machine{Name: machine.Name}
	_, err := client.Revert(ctx, req)
	if err != nil {
		c.logger.Printf("ERROR: RevertMachine: %+v", err)
		return err
	}

	return nil
}

func UpdateMachines() {
	for {
		time.Sleep(5 * time.Second)

		machines, err := ListMachines()
		if err != nil {
			continue
		}

		for i := range machines {
			db.SetMachineIp(&db.Machine{
				Name:      machines[i].Name,
				IpAddress: machines[i].Ip,
			})
			db.SetMachineStatus(&db.Machine{
				Name:   machines[i].Name,
				Status: machines[i].Status,
			})
		}
	}
}

func ConnectManager(address string) *grpc.ClientConn {
	var err error
	var managerConn *grpc.ClientConn

	for {
		managerConn, err = grpc.Dial(address,
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithTimeout(10*time.Second),
		)

		if err != nil && conn == nil {
			/* There was an error on initial connection */
			panic(err.Error())
		} else if err != nil {
			/* There was a connection break, keep retrying */
			continue
		}

		/* We've recovered the manager connection */
		break
	}

	return managerConn
}

func MonitorManager() {
	for {
		conn.WaitForStateChange(context.Background(), conn.GetState())

		for {
			state := conn.GetState()
			if state == connectivity.TransientFailure || state == connectivity.Shutdown {
				c.logger.Printf("ERROR: Manager Connection Break!")
				CloseManager()
				InitializeManager()
			} else {
				break
			}
		}
	}
}

func InitializeManager() {
	conn = ConnectManager("manager.vm.cucyber.net:8080")
	client = proto.NewManagerClient(conn)
}
