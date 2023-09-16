// package container

// import (
// 	"io"
// 	logger "mdocker/logger"
// )

// func ContainerStats(containerId string) (io.ReadCloser, error) {
// 	cli, ctx, err := GetDockerClient()
// 	if err != nil {
// 		log.Log.Error(err)
// 		return nil, err
// 	}
// 	stats, err := cli.ContainerStats(ctx, containerId, true)
// 	if err != nil {
// 		log.Log.Error(err)
// 		return nil, err
// 	}
// 	return stats.Body, nil
// }
