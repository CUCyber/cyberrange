package main

import (
	"context"
	"errors"
	"github.com/cucyber/cyberrange/pkg/proto"
	"github.com/cucyber/cyberrange/services/manager/git"
	"github.com/cucyber/cyberrange/services/manager/ovirt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/exec"
)

type server struct{}

var (
	ErrNoMachinePlaybook = errors.New("cyberrange: no machine playbook")
)

func (s *server) Create(ctx context.Context, req *proto.Machine) (*proto.Response, error) {
	MachineName := req.GetName()
	PlaybookPath := git.PlaybooksPath + "/machines/" + MachineName

	err := git.PullPlaybooks()
	if err != nil {
		return &proto.Response{Result: false}, err
	}

	if _, err = os.Stat(PlaybookPath); os.IsNotExist(err) {
		return &proto.Response{Result: false}, ErrNoMachinePlaybook
	}

	prov_cmd := exec.Command("ansible-playbook",
		PlaybookPath+"/provision/provision_vm.yml")
	err = prov_cmd.Start()
	if err != nil {
		return &proto.Response{Result: false}, err
	}
	err = prov_cmd.Wait()
	if err != nil {
		return &proto.Response{Result: false}, err
	}

	ip, err := ovirt.WaitForIPByName(MachineName)
	if err != nil {
		return &proto.Response{Result: false}, err
	}

	setup_cmd := exec.Command("ansible-playbook",
		"-i", ip+`,`,
		`--private-key`, `/home/cucyber/.ssh/id_rsa`,
		PlaybookPath+`/setup/playbook.yaml`,
	)
	err = setup_cmd.Start()
	if err != nil {
		return &proto.Response{Result: false}, err
	}
	err = setup_cmd.Wait()
	if err != nil {
		return &proto.Response{Result: false}, err
	}

	return &proto.Response{Result: true}, nil
}

func (s *server) Start(ctx context.Context, req *proto.Machine) (*proto.Response, error) {
	MachineName := req.GetName()
	err := ovirt.StartMachine(MachineName)
	if err != nil {
		return &proto.Response{Result: false}, err
	}
	return &proto.Response{Result: true}, nil
}

func (s *server) Stop(ctx context.Context, req *proto.Machine) (*proto.Response, error) {
	MachineName := req.GetName()
	err := ovirt.StopMachine(MachineName)
	if err != nil {
		return &proto.Response{Result: false}, err
	}
	return &proto.Response{Result: true}, nil
}

func (s *server) Restart(ctx context.Context, req *proto.Machine) (*proto.Response, error) {
	MachineName := req.GetName()
	err := ovirt.RestartMachine(MachineName)
	if err != nil {
		return &proto.Response{Result: false}, err
	}
	return &proto.Response{Result: true}, nil
}

func (s *server) Snapshot(ctx context.Context, req *proto.Machine) (*proto.Response, error) {
	MachineName := req.GetName()
	err := ovirt.SnapshotMachine(MachineName)
	if err != nil {
		return &proto.Response{Result: false}, err
	}
	return &proto.Response{Result: true}, nil
}

func (s *server) Revert(ctx context.Context, req *proto.Machine) (*proto.Response, error) {
	MachineName := req.GetName()
	err := ovirt.RevertMachine(MachineName)
	if err != nil {
		return &proto.Response{Result: false}, err
	}
	return &proto.Response{Result: true}, nil
}

func (s *server) List(ctx context.Context, req *proto.Empty) (*proto.MachineList, error) {
	return ovirt.ListMachines()
}

func serve() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err.Error())
	}

	srv := grpc.NewServer()
	proto.RegisterManagerServer(srv, &server{})
	reflection.Register(srv)

	if err = srv.Serve(listener); err != nil {
		panic(err.Error())
	}
}

func main() {
	ovirt.InitializeOVirt()
	defer ovirt.CloseOVirt()

	git.InitializeGit()
	git.PullPlaybooks()

	serve()
}
