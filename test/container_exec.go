package test

import (
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
