package apc

import (
	"apc_compare/pkg/common"
	"apc_compare/server/setting"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"
)

type ApcData struct {
	Chamber     string    `mapstructure:"CHAMBER"`
	ChangeTime  time.Time `mapstructure:"CHANGETIME"`
	LSL         float32   `mapstructure:"LL"`
	Parameter   string    `mapstructure:"PARAMETER"`
	Recipe      string    `mapstructure:"RECIPE"`
	RuleContent string    `mapstructure:"RULECONTENT"`
	RuleType    string    `mapstructure:"RULETYPE"`
	ToolID      string    `mapstructure:"TOOLID"`
	USL         float32   `mapstructure:"UL"`
}

type QueryMap struct {
	ToolID    string
	Chamber   string
	Parameter string
}

var apcMap *[]map[string]interface{}

// Setup Initialize the APC instance
func Setup() {
	var err error

	apcMap, err = getApcMap()
	if err != nil {
		log.Fatal(err)
	}
}

func GetChartID(queryData QueryMap) (string, error) {
	for _, apcRow := range *apcMap {
		compareData := QueryMap{apcRow["tool_id"].(string), apcRow["chamber"].(string), apcRow["parameter"].(string)}
		if reflect.DeepEqual(queryData, compareData) {
			chartID := apcRow["spc_chart_id"].(string)
			return chartID, nil
		}
	}
	return "", errors.New("can't find any chart id")
}

func GetApcData() ([]map[string]interface{}, error) {
	//var result []map[string]interface{}

	toolIDs, err := getToolIDs()
	if err != nil {
		return nil, err
	}

	var eg errgroup.Group

	var result []map[string]interface{}

	for _, toolId := range toolIDs {
		toolId := toolId
		eg.Go(func() error {
			if err := apcServiceResp(toolId, &result); err != nil {
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

func getApcMap() (*[]map[string]interface{}, error) {
	// Open jsonFile
	jsonFile, err := os.Open("conf/apc_mp_spc.json")

	if err != nil {
		return nil, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var jsonData []map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &jsonData)
	if err != nil {
		return nil, err
	}

	return &jsonData, nil
}

func getToolIDs() ([]string, error) {
	var ToolIDs []string

	if setting.AppSetting.ToolID != nil {
		ToolIDs = setting.AppSetting.ToolID
		return ToolIDs, nil
	}

	for _, data := range *apcMap {
		ID := data["tool_id"].(string)
		if !(common.StringInSlice(ID, ToolIDs)) {
			ToolIDs = append(ToolIDs, ID)
		}
	}

	return ToolIDs, nil
}

func apcServiceResp(toolId string, result *[]map[string]interface{}) error {

	apcParams := map[string]string{
		"user":        setting.AppSetting.User,
		"password":    setting.AppSetting.Password,
		"FABSITE":     setting.AppSetting.FabSite,
		"TOOLID":      "",
		"PROCESSTYPE": setting.AppSetting.ProcessType,
		"DATATYPE":    setting.AppSetting.DataType,
		"ziped":       setting.AppSetting.Ziped,
	}
	apcParams["TOOLID"] = toolId

	jsonString, err := json.Marshal(apcParams)
	if err != nil {
		return err
	}

	url := "http://10.21.10.124/APCWEB3D/APC_WBS/APCService.asmx/GetAPC_SpecList?"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Get APC data has an error, the error was %s ", resp.Status))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var respJson map[string]interface{}
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		return err
	}

	*result = append(*result, respJson)

	return nil
}
