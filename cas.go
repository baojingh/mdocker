package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	containerID := "<Your Container ID>"

	for {
		stats, err := cli.ContainerStats(ctx, containerID, false)
		if err != nil {
			panic(err)
		}
		ts.Usage
		defer stats.Body.Close()

		var v types.StatsJSON
		if err := json.NewDecoder(stats.Body).Decode(&v); err != nil {
			panic(err)
		}

		cpuUsage := v.CPUStats.CPUUsage.TotalUsage
		memUsage := v.MemoryStats.Usage

		fmt.Printf("CPU Usage: %d\n", cpuUsage)
		fmt.Printf("Memory Usage: %d\n", memUsage)

		time.Sleep(time.Second * 5)
	}
}


package main

import (
"context"
"fmt"

"github.com/docker/docker/api/types"
"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	containerID := "<Your Container ID>"
	cmd := []string{"sh", "-c", "ls"}

	execConfig := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}

	execID, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		panic(err)
	}
	defer resp.Close()

	fmt.Println("Output:")

	buf := make([]byte, 4096)
	for {
		bytesRead, err := resp.Reader.Read(buf)
		if err != nil {
			break
		}

		fmt.Print(string(buf[:bytesRead]))
	}
}





