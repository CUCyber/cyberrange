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
	ErrNoMachinePath      = errors.New("cyberrange: no machine playbook path")
	ErrNoMachineSetup     = errors.New("cyberrange: no machine playbook setup")
	ErrNoMachineProvision = errors.New("cyberrange: no machine playbook provision")
	ErrSetupFailed        = errors.New("cyberrange: machine setup failed, review ansible scripts")
	ErrProvisionFailed    = errors.New("cyberrange: machine provision failed, review ansible scripts")
)

func (s *server) CheckCreate(ctx context.Context, req *proto.Machine) (*proto.Response, error) {
	MachineName := req.GetName()
	PlaybookPath := git.PlaybooksPath + "/machines/" + MachineName
	SetupConfig := PlaybookPath + "/setup/playbook.yaml"
	ProvisionConfig := PlaybookPath + "/provision/provision_vm.yml"

	err := git.PullPlaybooks()
	if err != nil {
		return &proto.Response{Result: false}, err
	}

	if _, err = os.Stat(PlaybookPath); os.IsNotExist(err) {
		return &proto.Response{Result: false}, ErrNoMachinePath
	}

	if _, err = os.Stat(SetupConfig); os.IsNotExist(err) {
		return &proto.Response{Result: false}, ErrNoMachineSetup
	}

	if _, err = os.Stat(ProvisionConfig); os.IsNotExist(err) {
		return &proto.Response{Result: false}, ErrNoMachineProvision
	}

	return &proto.Response{Result: true}, nil
}

func (s *server) Create(ctx context.Context, req *proto.Machine) (*proto.Response, error) {
	MachineName := req.GetName()
	PlaybookPath := git.PlaybooksPath + "/machines/" + MachineName
	SetupConfig := PlaybookPath + "/setup/playbook.yaml"
	ProvisionConfig := PlaybookPath + "/provision/provision_vm.yml"

	prov_cmd := exec.Command("ansible-playbook", ProvisionConfig)

	err := prov_cmd.Start()
	if err != nil {
		return &proto.Response{Result: false}, ErrProvisionFailed
	}

	err = prov_cmd.Wait()
	if err != nil {
		return &proto.Response{Result: false}, ErrProvisionFailed
	}

	ip, err := ovirt.WaitForIPByName(MachineName)
	if err != nil {
		return &proto.Response{Result: false}, err
	}

	setup_cmd := exec.Command("ansible-playbook",
		"-i", ip+`,`,
		`--private-key`, `/home/cucyber/.ssh/id_rsa`,
		SetupConfig,
	)

	setup_cmd.Env = append(os.Environ(), "ANSIBLE_HOST_KEY_CHECKING=False")

	err = setup_cmd.Start()
	if err != nil {
		return &proto.Response{Result: false}, ErrSetupFailed
	}

	err = setup_cmd.Wait()
	if err != nil {
		return &proto.Response{Result: false}, ErrSetupFailed
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
