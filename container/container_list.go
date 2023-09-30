package container

import (
	"github.com/docker/docker/api/types"
)

func ContainerList() ([]types.Container, error) {
	cli, ctx, err := GetDockerClient()
	if err != nil {
		log.Error(err)
		return []types.Container{}, err
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	return containers, nil
}
