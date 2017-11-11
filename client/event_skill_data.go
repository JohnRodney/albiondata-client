package client

import (
	"strconv"

	"github.com/regner/albiondata-client/lib"
	"github.com/regner/albiondata-client/log"
	
	"net/http"
	"fmt"
	"bytes"
	
)

type eventSkillData struct {
	SkillIds    []int     `mapstructure:"1"`
	Levels      []int     `mapstructure:"2"`
	Percentages []float64 `mapstructure:"3"`
	Fame        []string  `mapstructure:"4"`
}
func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func (event eventSkillData) Process(state *albionState) {
	log.Debug("Got skill data event...")

	skills := []*lib.Skill{}
	var buffer bytes.Buffer
	buffer.WriteString("{\"skills\": [")
	for k, _ := range event.SkillIds {
		skill := &lib.Skill{}
		skill.ID = event.SkillIds[k]
		skill.Level = event.Levels[k]
		skill.PercentNextLevel = event.Percentages[k]
		// for some reason, the value is enclosed in [[]]. trying to get rid of them
		fame, err := strconv.Atoi(event.Fame[k][2 : len(event.Fame[k])-2])
		if err != nil {
			log.Error("Could not parse fame value. ", err)
			continue
		}
		skill.Fame = fame

		skills = append(skills, skill)
		buffer.WriteString("{\"SID\":\""+strconv.Itoa(skill.ID)+"\",\"SLVL\":\""+strconv.Itoa(skill.Level)+"\",\"SPER\":\""+FloatToString(skill.PercentNextLevel)+"\"},")
	}
	
	buffer.Truncate(buffer.Len()-1)
	buffer.WriteString("],\"player\": \""+state.CharacterName+"\"}")
	

	if len(skills) < 1 {
		return
	}

/*
	upload := lib.SkillsUpload{
		Skills: skills,
	}
*/
    url := "https://albion-data-revival.herokuapp.com/postSkills/"
    

    var jsonStr = []byte((buffer.String()))
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
	log.Infof("Sending %d skills of %v to ingest", len(skills), state.CharacterName)

    //fmt.Println("response Status:", resp.Status)
    //fmt.Println("response Headers:", resp.Header)
    //body, _ := ioutil.ReadAll(resp.Body)
    //fmt.Println("response Body:", string(body))



}
