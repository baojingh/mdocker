package main

/**
  @Author   : bob
  @Datetime : 2023-05-09 下午 11:08
  @File     : main.go
  @Desc     :
*/

import (
	"bufio"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	client "github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"time"
)

var (
	ctx context.Context
	cli *client.Client
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat:           "2006-01-02 15:04:05",
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
		DisableLevelTruncation:    true,
	})
	ctx = context.Background()
	//cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	cli, _ = client.NewClientWithOpts(
		client.WithHost("tcp://121.5.73.196:2375"),
		client.WithAPIVersionNegotiation())
}

func main() {

	log.Info("hello")

	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Error(err)
	}

	for _, container := range containers {
		log.Info(container.ID, container.Names, container.Status, container.State)
		//containerLogs(cli, ctx, container.ID)
		containerExec(cli, ctx, container.ID)
		//time.Sleep(time.Second * 9999)
	}

}

func containerLogs(cli *client.Client, ctx context.Context, containerID string) error {
	options := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true}
	out, err := cli.ContainerLogs(ctx, containerID, options)
	if err != nil {
		log.Errorf("fail to get the container %s logs, %s", containerID, err)
		return err
	}
	io.Copy(os.Stdout, out)
	return nil
}

func containerExec(cli *client.Client, ctx context.Context, containerID string) error {
	// 在指定容器中执行/bin/bash命令
	ir, err := cli.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"/bin/bash"},
		Tty:          true,
	})
	if err != nil {
		panic(err)
	}

	// 附加到上面创建的/bin/bash进程中
	hr, err := cli.ContainerExecAttach(ctx, ir.ID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		panic(err)
	}
	// 关闭I/O
	defer hr.Close()
	// 输入
	hr.Conn.Write([]byte("ls\r"))
	// 输出
	scanner := bufio.NewScanner(hr.Conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return nil
}

func containerExec2(cli *client.Client, ctx context.Context, containerID string) error {
	execOpts := types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"date"},
	}
	resp, err := cli.ContainerExecCreate(ctx, containerID, execOpts)
	if err != nil {
		log.Error(err)
	}
	attach, err := cli.ContainerExecAttach(ctx, resp.ID, types.ExecStartCheck{})
	if err != nil {
		log.Error(err)
	}
	defer attach.Close()
	data, _ := ioutil.ReadAll(attach.Reader)
	log.Info("正在解压...\n", string(data))
	return nil
}

func containerExec1(cli *client.Client, ctx context.Context, containerID string) error {
	execOpts := types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"date"},
	}

	resp, err := cli.ContainerExecCreate(context.Background(), containerID, execOpts)
	if err != nil {
		log.Error(err)
		return err
	}

	respTwo, err := cli.ContainerExecAttach(context.Background(), resp.ID,
		types.ExecStartCheck{Detach: true, Tty: true},
	)
	if err != nil {
		log.Error(err)

		return err
	}
	defer respTwo.Close()

	err = cli.ContainerExecStart(context.Background(), resp.ID, types.ExecStartCheck{Detach: true, Tty: true})
	if err != nil {
		log.Error(err)

		return err
	}

	running := true
	for running {
		respThree, err := cli.ContainerExecInspect(context.Background(), resp.ID)
		if err != nil {
			log.Error(err)

			panic(err)
		}

		if !respThree.Running {
			running = false
		}

		time.Sleep(250 * time.Millisecond)
	}

	return nil
}
