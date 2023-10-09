package handler

import (
	"encoding/json"
	"io"
	"mdocker/container"

	"github.com/docker/docker/api/types"
)

/**
  @Author:      He Bao Jing
  @Date:        5/12/2023 10:34 AM
  @Description:
*/

func StatsProducer(containerId string, statsChan chan types.StatsJSON, shutdownChan chan int) {
	_, reader, err := container.ContainerStats(containerId)
	if err != nil {
		log.Error("fail to get the container stats reader, ", err)
		return
	}

	decoder := json.NewDecoder(reader)
	var statsValue types.StatsJSON

	for {
		// 可以同时监听多个通道的数据流动，通过 case 分支来处理具体的通道操作。
		// 当有多个通道同时可用时，select 会随机选择一个可用的通道进行操作。
		// 当所有的通道都阻塞时，select 可以执行默认的 default 分支，实现非阻塞的操作。
		select {
		case <-shutdownChan:
			reader.Close()
			// producer should close the channel nor consumer
			close(statsChan)
			log.Warn("Produer stop producing container stat metrics")
			return
		default:
			if err := decoder.Decode(&statsValue); err == io.EOF {
				log.Warn("Receive EOF flag and exit")
				return
			} else if err != nil {
				log.Error("Something Error occured", err)
				return
			} else {
				statsChan <- statsValue
			}
		}
	}
}

type ContainerSimple struct {
	Id   string
	Name string
}

func ContainerList() []ContainerSimple {
	containers, _ := container.ContainerList()
	arr := make([]ContainerSimple, 0)

	for _, element := range containers {
		name := ""
		if len(element.Names) >= 1 {
			name = element.Names[0]
		}

		c := ContainerSimple{
			Id:   element.ID,
			Name: name,
		}
		arr = append(arr, c)
	}
	return arr
}
