//package container
//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/docker/docker/api/types"
//	"github.com/docker/docker/client"
//	"time"
//)
//
//func main() {
//	cli, err := client.NewClientWithOpts(client.FromEnv)
//	if err != nil {
//		panic(err)
//	}
//
//	ctx := context.Background()
//
//	containerID := "<Your Container ID>"
//
//	for {
//		stats, err := cli.ContainerStats(ctx, containerID, false)
//		if err != nil {
//			panic(err)
//		}
//		ts.Usage
//		defer stats.Body.Close()
//
//		var v types.StatsJSON
//		if err := json.NewDecoder(stats.Body).Decode(&v); err != nil {
//			panic(err)
//		}
//
//		cpuUsage := v.CPUStats.CPUUsage.TotalUsage
//		memUsage := v.MemoryStats.Usage
//
//		fmt.Printf("CPU Usage: %d\n", cpuUsage)
//		fmt.Printf("Memory Usage: %d\n", memUsage)
//
//		time.Sleep(time.Second * 5)
//	}
//}
