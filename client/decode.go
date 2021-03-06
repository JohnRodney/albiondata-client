package client

import (
	"encoding/hex"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/regner/albiondata-client/lib"
	"github.com/regner/albiondata-client/log"
	"fmt"
	"net/http"
	"strconv"
	"io/ioutil"
)

func decodeRequest(params map[string]interface{}) (operation operation, err error) {
	if _, ok := params["253"]; !ok {
		return nil, nil
	}

	code := params["253"].(int16)

	switch code {
	case 10:
		operation = &operationGetGameServerByCluster{}
	case 67:
		operation = &operationAuctionGetOffers{}
	case 166:
		operation = &operationGetClusterMapInfo{}
	case 217:
		operation = &operationGoldMarketGetAverageInfo{}
	case 232:
		operation = &operationRealEstateGetAuctionData{}
	case 233:
		operation = &operationRealEstateBidOnAuction{}
	default:
		return nil, nil
	}

	err = decodeParams(params, operation)

	return operation, err
}

func decodeResponse(params map[string]interface{}) (operation operation, err error) {
	if _, ok := params["253"]; !ok {
		return nil, nil
	}

	code := params["253"].(int16)

	switch code {
	case 2:
		operation = &operationJoinResponse{}
	case 69:
		operation = &operationAuctionGetOffersResponse{}
	case 70:
		operation = &operationAuctionGetRequestsResponse{}
	case 149:
		operation = &operationReadMail{}
	case 217:
		operation = &operationGoldMarketGetInfosResponse{}
		
		mapstructure.Decode(params, &operation)
		cost := strconv.FormatInt(params["0"].(int64),10)		
		sell := strconv.FormatInt(params["2"].(int64),10)
		cost = cost[:len(cost)-4]
		sell = sell[:len(sell)-4]			
		log.Debug("before url")

		url := "https://limitless-lowlands-38074.herokuapp.com/ingest?buyPrice=" + cost + "&sellPrice=" + sell
		log.Debug("after url")
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		robots, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", robots)
	case 168:
		operation = &operationGetClusterMapInfoResponse{}
	case 219:
		operation = &operationGoldMarketGetAverageInfoResponse{}
	case 235:
		operation = &operationRealEstateGetAuctionDataResponse{}
	case 236:
		operation = &operationRealEstateBidOnAuctionResponse{}
	default:
		return nil, nil
	}

	err = decodeParams(params, operation)

	return operation, err
}

func decodeEvent(params map[string]interface{}) (event operation, err error) {
	if _, ok := params["252"]; !ok {
		return nil, nil
	}

	eventType := params["252"].(int16)

	switch eventType {
	case 77:
		event = &eventPlayerOnlineStatus{}
	case 118:
		event = &eventSkillData{}
	default:
		return nil, nil
	}

	err = decodeParams(params, event)

	return event, err
}

func decodeParams(params interface{}, operation operation) error {
	convertGameObjects := func(from reflect.Type, to reflect.Type, v interface{}) (interface{}, error) {
		if from == reflect.TypeOf([]int8{}) && to == reflect.TypeOf(lib.CharacterID("")) {
			log.Debug("Parsing character ID from mixed-endian UUID")

			return decodeCharacterID(v.([]int8)), nil
		}

		return v, nil
	}

	config := mapstructure.DecoderConfig{
		DecodeHook: convertGameObjects,
		Result:     operation,
	}

	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}

	err = decoder.Decode(params)

	return err
}

func decodeCharacterID(array []int8) lib.CharacterID {
	/* So this is a UUID, which is stored in a 'mixed-endian' format.
	The first three components are stored in little-endian, the rest in big-endian.
	See https://en.wikipedia.org/wiki/Universally_unique_identifier#Encoding.
	By default, our int array is read as big-endian, so we need to swap the first
	three components of the UUID
	*/
	b := make([]byte, len(array))

	// First, convert to byte
	for k, v := range array {
		b[k] = byte(v)
	}

	// swap first component
	b[0], b[1], b[2], b[3] = b[3], b[2], b[1], b[0]

	// swap second component
	b[4], b[5] = b[5], b[4]

	// swap third component
	b[6], b[7] = b[7], b[6]

	// format it UUID-style
	var buf [36]byte
	hex.Encode(buf[:], b[:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], b[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], b[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], b[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], b[10:])

	return lib.CharacterID(buf[:])
}
