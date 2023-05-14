package container

/**
  @Author:      He Bao Jing
  @Date:        5/12/2023 1:45 PM
  @Description:
*/

import (
	"context"
	client "github.com/docker/docker/client"
	log "mdocker/logger"
	"sync"
)

var (
	ctx  context.Context
	cli  *client.Client
	once sync.Once
)

// var dockerHost string = "tcp://192.168.1.130:2375"
// var dockerHost string = "tcp://127.0.0.1:2375"
var dockerHost string = "tcp://11.0.1.128:2375"

func GetDockerClient() (*client.Client, context.Context, error) {
	var err error
	once.Do(func() {
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
