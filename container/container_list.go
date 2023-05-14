package container

import (
	"github.com/docker/docker/api/types"
	log "mdocker/logger"
)

func ContainerList() ([]types.Container, error) {
	cli, ctx, err := GetDockerClient()
	if err != nil {
		log.Log.Error(err)
		return nil, err
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Log.Error(err)
		return nil, err
	}
	log.Log.Infof("get container list, size: %d", len(containers))
	return containers, nil
}
