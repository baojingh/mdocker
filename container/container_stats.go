package container

import (
	"context"
	"io"
)

func ContainerStats(containerId string) (context.Context, io.ReadCloser, error) {
	cli, ctx, err := GetDockerClient()
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}
	stats, err := cli.ContainerStats(ctx, containerId, true)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}
	return ctx, stats.Body, nil
}
