package client

import (
	"github.com/JohnRodney/albiondata-client/log"
)

type operationGoldMarketGetAverageInfo struct {
}

func (op operationGoldMarketGetAverageInfo) Process(state *albionState) {
	log.Debug("Got GoldMarketGetAverageInfo operation...")
}

type operationGoldMarketGetAverageInfoResponse struct {
	GoldPrices []int   `mapstructure:"0"`
	TimeStamps []int64 `mapstructure:"1"`
}

func (op operationGoldMarketGetAverageInfoResponse) Process(state *albionState) {
	log.Debug("Got response to GoldMarketGetAverageInfo operation...")
}
