package main

/**
  @Author   : bob
  @Datetime : 2023-05-09 下午 11:08
  @File     : main.go
  @Desc     :
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	containerId := "ecd673b5f6e0"

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	readerArr := []io.ReadCloser{}
	stats, err := cli.ContainerStats(ctx, containerId, true)
	readerArr = append(readerArr, stats.Body)
	for {
		for _, rr := range readerArr {
			decoder := json.NewDecoder(rr)
			var statsValue types.StatsJSON
			if err := decoder.Decode(&statsValue); err == io.EOF {
				fmt.Println("Producer receive EOF flag and exit")
				return
			} else if err != nil {
				fmt.Println("Something Error occured to producer", err)
				return
			} else {
				// fmt.Printf("cpu: %v, mem: %v",
				// 	statsValue.CPUStats.CPUUsage.TotalUsage,
				// 	statsValue.MemoryStats.Usage)

				statsJSONBytes, _ := json.MarshalIndent(statsValue, "", "  ")
				statsStr := string(statsJSONBytes)
				fmt.Println(statsStr)
			}
		}
		time.Sleep(time.Second * 1)
		fmt.Println()
		fmt.Println()
	}

}
