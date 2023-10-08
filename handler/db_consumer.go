package handler

import (
	"encoding/json"
	"mdocker/handler"

	"github.com/docker/docker/api/types"
)

const (
	org = "mdocker"
	bucket = "mdocker-bucket"
)

func DbConsumer(statsChan chan types.StatsJSON) {
	cli, ctx := GetInfluxdbClient()
	writeAPI := cli.WriteAPIBlocking(org, bucket)
	for value := 0; value < 5; value++ {
		tags := map[string]string{
			"tagname1": "tagvalue1",
		}
		fields := map[string]interface{}{
			"field1": value,
		}
		point := write.NewPoint("measurement1", tags, fields, time.Now())
		time.Sleep(1 * time.Second) // separate points by 1 second
	
		if err := writeAPI.WritePoint(ctx, point); err != nil {
			log.Fatal(err)
		}
	}

	for val := range statsChan {
		statsJSONBytes, _ := json.MarshalIndent(val, "", "  ")
		log.Info(string(statsJSONBytes))
	}
}



writeAPI := client.WriteAPIBlocking(org, bucket)
for value := 0; value < 5; value++ {
	tags := map[string]string{
		"tagname1": "tagvalue1",
	}
	fields := map[string]interface{}{
		"field1": value,
	}
	point := write.NewPoint("measurement1", tags, fields, time.Now())
	time.Sleep(1 * time.Second) // separate points by 1 second

	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
		log.Fatal(err)
	}
}