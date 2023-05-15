package container

import (
	"github.com/docker/docker/api/types"
	log "mdocker/logger"
)

func ContainerExec(containerId string) (types.HijackedResponse, error) {
	cli, ctx, err := GetDockerClient()
	if err != nil {
		log.Log.Error(err)
		return types.HijackedResponse{}, err
	}

	execOpts := types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"/bin/bash"},
		Tty:          true,
	}
	resp, err := cli.ContainerExecCreate(ctx, containerId, execOpts)
	if err != nil {
		log.Log.Error("fail to exec create container, ", err)
		return types.HijackedResponse{}, err
	}
	log.Log.Info("ContainerExecCreate is done")

	checkOptions := types.ExecStartCheck{
		Detach: false,
		Tty:    true,
	}
	attach, err := cli.ContainerExecAttach(ctx, resp.ID, checkOptions)
	if err != nil {
		log.Log.Error("fail to exec attach container, ", err)
		return types.HijackedResponse{}, err
	}
	log.Log.Info("ContainerExecAttach is done")
	return attach, nil
}
