package main

import (
	"bufio"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

/**
  @Author:      He Bao Jing
  @Date:        5/10/2023 1:36 PM
  @Description:
*/

package main
import (    "bufio"    "context"    "fmt"    "github.com/docker/docker/api/types"
"github.com/docker/docker/client")
func main() {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)
	ir, err := cli.ContainerExecCreate(ctx, "test", types.ExecConfig{AttachStdin: true, AttachStdout: true, AttachStderr: true, Cmd: []string{"/bin/bash"}, Tty: true})
	if err != nil {
		panic(err)
	}
	hr, err := cli.ContainerExecAttach(ctx, ir.ID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		panic(err)
	}
	defer hr.Close()
	hr.Conn.Write([]byte("ls"))
	scanner := bufio.NewScanner(hr.Conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
}
