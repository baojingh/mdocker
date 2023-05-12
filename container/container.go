package container

/**
  @Author:      He Bao Jing
  @Date:        5/12/2023 1:45 PM
  @Description:
*/

import (
	"context"
	"github.com/docker/docker/api/types"
	client "github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"io"
	"sync"
)

var (
	ctx  context.Context
	cli  *client.Client
	once sync.Once
)

// var dockerHost string = "tcp://192.168.1.130:2375"
var dockerHost string = "tcp://127.0.0.1:2375"

func init() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat:           "2006-01-02 15:04:05",
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
		DisableLevelTruncation:    true,
	})
}

func GetDockerClient() (*client.Client, error) {
	var err error
	once.Do(func() {
		cli, err = client.NewClientWithOpts(
			client.WithHost(dockerHost),
			client.WithAPIVersionNegotiation())
		ctx = context.Background()
	})
	if err != nil {
		return nil, err
	}
	log.Info("create docker client success")
	return cli, nil
}

func conList(containerId string) []types.Container {
	cli, err := GetDockerClient()
	defer cli.Close()
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Error(err)
	}
	log.Infof("get container list, size: %d", len(containers))
	return containers
}

func ConLogs(containerId string) (io.Reader, error) {
	cli, err := GetDockerClient()
	if err != nil {
		log.Error(err)
	}
	defer cli.Close()
	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: true,
		Details:    true,
	}
	resp, err := cli.ContainerLogs(ctx, containerId, options)
	if err != nil {
		log.Errorf("fail to get the container %s logs, %s", containerId, err)
		return nil, err
	}
	return resp, nil
}
