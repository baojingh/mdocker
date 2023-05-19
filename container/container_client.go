package container

/**
  @Author:      He Bao Jing
  @Date:        5/12/2023 1:45 PM
  @Description:
*/

import (
	"context"
	client "github.com/docker/docker/client"
	"mdocker/config"
	log "mdocker/logger"
	"strings"
	"sync"
)

var (
	ctx  context.Context
	cli  *client.Client
	once sync.Once
)

func composeDockerHost() string {
	host := config.MDocker.Docker.Host
	port := config.MDocker.Docker.Port
	var builder strings.Builder
	builder.WriteString("tcp://")
	builder.WriteString(host)
	builder.WriteString(":")
	builder.WriteString(port)
	res := builder.String()
	log.Log.Infof("get docker host info: %s", res)
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
	log.Log.Info("create docker client success")
	return cli, ctx, nil
}
