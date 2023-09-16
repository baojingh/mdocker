package container

import (
	"io"

	"github.com/docker/docker/api/types"
)

/**
  @Author   : bob
  @Datetime : 2023-05-14 下午 12:04
  @File     : ca.go
  @Desc     :
*/

func ContainerLogs(containerId string) (io.ReadCloser, error) {
	cli, ctx, err := GetDockerClient()
	if err != nil {
		// log.Error(err)
		return nil, err
	}
	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
	}
	reader, err := cli.ContainerLogs(ctx, containerId, options)
	if err != nil {
		log.Error("fail to get the container, ", err)
		return nil, err
	}
	return reader, nil
}
