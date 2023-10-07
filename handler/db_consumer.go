package handler

import (
	"encoding/json"

	"github.com/docker/docker/api/types"
)

func DbConsumer(statsChan chan types.StatsJSON) {
	statsValue := <-statsChan
	statsJSONBytes, _ := json.MarshalIndent(statsValue, "", "  ")
	log.Info(string(statsJSONBytes))
}
