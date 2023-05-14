package test

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/daemon/logger"
)

func main() {
	cli, err := GetDockerClient()
	defer cli.Close()
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("get container list, size: %d", len(containers))
	return containers
}
}
