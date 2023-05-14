package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gorilla/websocket"
	"github.com/webtty/webtty"
)

func main() {
	// 创建一个 Docker 客户端
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	// 获取容器 ID
	containerID := "CONTAINER_ID"

	// 启动容器
	err = startContainer(cli, containerID)
	if err != nil {
		panic(err)
	}

	// 为容器创建 WebTTY 终端
	handler, err := newWebTTYHandler(cli, containerID)
	if err != nil {
		panic(err)
	}

	// 启动 HTTP 服务器并处理 WebTTY 请求
	http.HandleFunc("/terminal", handler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func startContainer(cli *client.Client, containerID string) error {
	ctx := context.Background()

	// 检查容器是否存在并处于停止状态
	container, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return err
	}
	if container.State.Status == "running" {
		return fmt.Errorf("container is already running")
	}

	// 启动容器
	err = cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return nil
}

func newWebTTYHandler(cli *client.Client, containerID string) (http.HandlerFunc, error) {
	ctx := context.Background()

	// 获取容器信息
	container, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}

	// 创建WebTTY Options
	options := webtty.DefaultOptions()
	options.Command = []string{"/bin/sh"}
	options.TerminalCols = container.HostConfig.Tty.Width
	options.TerminalRows = container.HostConfig.Tty.Height

	// 创建WebTTY终端处理程序
	handler := webtty.New(options, func() (webtty.StartFunc, webtty.ResizeFunc, webtty.StopFunc) {
		execCfg := types.ExecConfig{
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Tty:          true,
			Cmd:          []string{"/bin/sh"},
		}
		resp, err := cli.ContainerExecCreate(ctx, containerID, execCfg)
		if err != nil {
			return nil, nil, nil
		}

		start := func(conn *websocket.Conn) error {
			execStartCheck := types.ExecStartCheck{Tty: true}
			hijackedResp, err := cli.ContainerExecAttach(ctx, resp.ID, execStartCheck)
			if err != nil {
				return err
			}
			defer hijackedResp.Close()

			go func() {
				err = stdcopy.StdCopy(conn, conn, hijackedResp.Reader)
				if err != nil {
					conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				}
			}()
			return nil
		}

		resize := func(cols uint16, rows uint16) error {
			return cli.ContainerExecResize(ctx, resp.ID, types.ResizeOptions{
				Height: uint(rows),
				Width:  uint(cols),
			})
		}

		stop := func() error {
			return cli.ContainerExecKill(ctx, resp.ID, "SIGKILL")
		}

		return start, resize, stop
	})

	return handler.HandleHTTP, nil
}
