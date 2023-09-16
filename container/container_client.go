package container

/**
  @Author:      He Bao Jing
  @Date:        5/12/2023 1:45 PM
  @Description:
*/

import (
	"context"
	entity "mdocker/entity"
	logger "mdocker/logger"
	"strings"
	"sync"

	client "github.com/docker/docker/client"
)

var log = logger.New()

var (
	ctx  context.Context
	cli  *client.Client
	once sync.Once
)

func composeDockerHost() string {
	host := entity.MDocker.Docker.Host
	port := entity.MDocker.Docker.Port
	var builder strings.Builder
	builder.WriteString("tcp://")
	builder.WriteString(host)
	builder.WriteString(":")
	builder.WriteString(port)
	res := builder.String()
	log.Infof("get docker host info: %s", res)
	return res
}

func GetDockerClient() (*client.Client, context.Context, error) {
	var err error
	once.Do(func() {
		dockerHost := composeDockerHost()
		cli, err = client.NewClientWithOpts(
			client.WithHost(dockerHost),
			client.WithAPIVersionNegotiation())
		ctx = context.Background()
	})
	if err != nil {
		return nil, nil, err
	}
	log.Info("create docker client success")
	return cli, ctx, nil
}
