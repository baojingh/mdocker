package test

/**
 @Author   : bob
 @Datetime : 2023-05-09 下午 11:08
 @File     : main.go
 @Desc     :
*/

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	client "github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var (
	ctx context.Context
	cli *client.Client
	once         sync.Once

)

var dockerHost string = "tcp://192.168.1.130:2375"
//var dockerHost string = "tcp://121.5.73.196:2375"




func init() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat:           "2006-01-02 15:04:05",
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
		DisableLevelTruncation:    true,
	})
}

func GetDockerClient() (*client.Client, error) {
	var err error
	once.Do(func() {
		cli, err = client.NewClientWithOpts(
			client.WithHost(dockerHost),
			client.WithAPIVersionNegotiation())
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}


func main() {

	ctx = context.Background()

	log.Info("hello")
	cli, err := GetDockerClient()
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Error(err)
	}

	for _, container := range containers {
		log.Info(container.ID, container.Names, container.Status, container.State)
		containerMonitor(cli, ctx, container.ID)
		//containerLogs(cli, ctx, container.ID)
		//containerExec(cli, ctx, container.ID)
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

type myStruct struct {
	Id        string `json:"id"`
	Read      string `json:"read"`
	Preread   string `json:"preread"`
	CpuStats  cpu    `json:"cpu_stats"`
	num_procs string `json:"num_procs"`
}

type cpu struct {
	Usage cpuUsage `json:"cpu_usage"`
}

type cpuUsage struct {
	Total float64 `json:"total_usage"`
}

func containerMonitor(cli *client.Client, ctx context.Context, containerID string) {
	// https://stackoverflow.com/questions/47154036/decode-json-from-stream-of-data-docker-go-sdk
	// https://github.com/infracloudio/cstats/blob/master/src/getStats.go
	stats, e := cli.ContainerStats(ctx, containerID, true)
	if e != nil {
		fmt.Errorf("%s", e.Error())
	}
	decoder := json.NewDecoder(stats.Body)
	var containerStats myStruct
	for {
		select {
		case <-ctx.Done():
			stats.Body.Close()
			fmt.Println("Stop logging")
			return
		default:
			if err := decoder.Decode(&containerStats); err == io.EOF {
				return
			} else if err != nil {
				log.Error(err)
			}
			fmt.Println(containerStats.CpuStats.Usage.Total)
		}
	}

}

func containerMonitor1(cli *client.Client, ctx context.Context, containerID string) {
	containerStats, err := cli.ContainerStats(ctx, containerID, false)
	if err != nil {
		log.Error(err)
	}
	log.Info(containerStats.Body)

	// https://www.cnblogs.com/xwxz/p/13637730.html
	// https://stackoverflow.com/questions/47154036/decode-json-from-stream-of-data-docker-go-sdk
	//ContainerStats的返回的结构如下 注意这个Body的类型是io.ReadCloser 好奇怪的类型 下面我们给他转成json
	type ContainerStats struct {
		Body   io.ReadCloser `json:"body"`
		OSType string        `json:"ostype"`
	}

	fmt.Println(containerStats)
	fmt.Println("containerStats.Body的内容是: ", containerStats.Body)
	buf := new(bytes.Buffer)
	//io.ReadCloser 转换成 Buffer 然后转换成json字符串
	buf.ReadFrom(containerStats.Body)
	newStr := buf.String()
	fmt.Printf(newStr)

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
