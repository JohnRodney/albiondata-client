package client

import (
	"github.com/regner/albiondata-client/log"
)

type operationGoldMarketGetInfos struct {
}

func (op operationGoldMarketGetInfos) Process(state *albionState) {
	log.Debug("Got GoldMarketGetInfo operation...")
}

type operationGoldMarketGetInfosResponse struct {
	PriceToBuyGold []int `mapstructure:"0"`
	PriceForSellingGold []int `mapstructure:"2"`
}

func (op operationGoldMarketGetInfosResponse) Process(state *albionState) {
	log.Debug("Got response to GoldMarketGetInfo operation...")
}
