package handler

import (
	"context"
	"encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

const (
	org    = "mdocker"
	bucket = "mdocker-bucket"
)

func DbConsumer(statsChan chan types.StatsJSON) {

	cli, ctx := GetInfluxdbClient()
	writeAPI := cli.WriteAPIBlocking(org, bucket)
	for {
		select {
		case val, ok := <-statsChan:
			if !ok {
				log.Warn("Consumer stop consumming container stat metrics")
				return
			} else {
				log.Infof("cpu: %v, mem: %v", val.CPUStats.CPUUsage.TotalUsage,
					val.MemoryStats.Usage)

				statsJSONBytes, _ := json.MarshalIndent(val, "", "  ")
				writeData2DB(ctx, writeAPI, statsJSONBytes)
			}
		default:
		}
	}
}

func writeData2DB(ctx context.Context, writeAPI api.WriteAPIBlocking, statsJSONBytes []byte) {
	// statsStr := string(statsJSONBytes)
	// tags := map[string]string{
	// 	"tagname1": "tagvalue1",
	// }
	// fields := map[string]interface{}{
	// 	"field1": statsStr,
	// }
	// point := write.NewPoint("measurement1", tags, fields, time.Now())
	// if err := writeAPI.WritePoint(ctx, point); err != nil {
	// 	log.Fatal(err)
	// }
	// log.Info("Data is save success")
	// DbDataView()
}

// func DbConsumer(statsChan chan types.StatsJSON) {
// 	cli, ctx := GetInfluxdbClient()
// 	writeAPI := cli.WriteAPIBlocking(org, bucket)
// 	for val := range statsChan {
// 		statsJSONBytes, _ := json.MarshalIndent(val, "", "  ")
// 		statsStr := string(statsJSONBytes)
// 		// log.Info(statsStr)
// 		tags := map[string]string{
// 			"tagname1": "tagvalue1",
// 		}
// 		fields := map[string]interface{}{
// 			"field1": statsStr,
// 		}
// 		point := write.NewPoint("measurement1", tags, fields, time.Now())
// 		if err := writeAPI.WritePoint(ctx, point); err != nil {
// 			log.Fatal(err)
// 		}
// 		log.Info("Data is save success")
// 		// DbDataView()
// 	}
// }

func DbDataView() {
	client, ctx := GetInfluxdbClient()
	// Get query client
	queryAPI := client.QueryAPI(org)
	// Get parser flux query result
	result, err := queryAPI.Query(ctx, `from(bucket:"mdocker-bucket")|> range(start: -10h) |> filter(fn: (r) => r._measurement == "measurement1")`)
	if err == nil {
		// Use Next() to iterate over query result lines
		for result.Next() {
			// Observe when there is new grouping key producing new table
			if result.TableChanged() {
				log.Infof("table: %s\n", result.TableMetadata().String())
			}
			// read result
			log.Infof("row: %s\n", result.Record().String())
		}
		if result.Err() != nil {
			log.Infof("Query error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}
	// Ensures background processes finishes
	client.Close()
}
