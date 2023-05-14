package container

import (
	"github.com/docker/docker/api/types"
	log "mdocker/logger"
)

/**
  @Author   : bob
  @Datetime : 2023-05-14 下午 3:06
  @File     : contianer_inspect.go
  @Desc     :
*/

func ContainerInspect(containerId string) (types.ContainerJSON, error) {
	cli, ctx, err := GetDockerClient()
	if err != nil {
		log.Log.Error(err)
	}
	inspect, err := cli.ContainerInspect(ctx, containerId)
	if err != nil {
		log.Log.Error(err)
		return types.ContainerJSON{}, err
	}
	log.Log.Infof(inspect.Name)
	log.Log.Infof(inspect.ID)
	log.Log.Infof(inspect.MountLabel)
	return inspect, err
}
