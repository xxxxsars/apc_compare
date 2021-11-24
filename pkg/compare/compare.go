package compare

import (
	"apc_compare/pkg/apc"
	"apc_compare/pkg/common"
	"apc_compare/server/models"
	"apc_compare/server/setting"
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/jszwec/csvutil"
	"github.com/mitchellh/mapstructure"
	"os"
	"reflect"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

type Record struct {
	ToolID     string
	Chamber    string
	Parameter  string
	ChartID    string
	ChangeTime time.Time
	LSL        string
	Recipe     string
	USL        string
	Result     string
}

func SpaceCompare() (bool, error) {
	//jsonfile, err := os.Open("apc.json")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer jsonfile.Close()
	//byteValue, _ := ioutil.ReadAll(jsonfile)
	//var apcServiceData []map[string]interface{}
	//if err := json.Unmarshal([]byte(byteValue), &apcServiceData); err != nil {
	//	log.Fatal(err)
	//}

	hadErr := false

	apcServiceData, err := apc.GetApcData()
	if err != nil {
		return hadErr, err
	}

	//time convert function
	stringToDateTimeHook := func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t == reflect.TypeOf(time.Time{}) && f == reflect.TypeOf("") {
			return time.Parse(timeFormat, data.(string))
		}
		return data, nil
	}

	var records []Record
	var existChartID []string

	for _, toolData := range apcServiceData {
		for _, data := range (toolData["data"].(map[string]interface{})["Data"]).([]interface{}) {

			var apcData apc.ApcData
			mapData := data.(map[string]interface{})

			config := mapstructure.DecoderConfig{
				DecodeHook: stringToDateTimeHook,
				Result:     &apcData,
			}

			decoder, err := mapstructure.NewDecoder(&config)
			if err != nil {
				return hadErr, err
			}

			if err = decoder.Decode(mapData); err != nil {
				return hadErr, err
			}

			//TODO: goroutine enhance performance
			chartID, err := apc.GetChartID(apc.QueryMap{ToolID: apcData.ToolID, Chamber: apcData.Chamber, Parameter: apcData.Parameter})
			// If it can't find chart id will ignore the record (force the apc_mp_spc.json).
			if err == nil {
				spcData, err := models.GetSpcData(chartID)
				if err != nil {
					return hadErr, err
				} else {

					csvRow := Record{
						ToolID:     apcData.ToolID,
						Chamber:    apcData.Chamber,
						Parameter:  apcData.Parameter,
						ChartID:    chartID,
						ChangeTime: apcData.ChangeTime,
						LSL:        fmt.Sprintf("%.2f/%.2f", apcData.LSL, spcData.LSL),
						Recipe:     apcData.Recipe,
						USL:        fmt.Sprintf("%.2f/%.2f", apcData.USL, spcData.USL),
					}

					if apcData.LSL >= spcData.LSL && apcData.USL <= spcData.USL {
						csvRow.Result = "OK"
					} else {
						csvRow.Result = "NG"
						hadErr = true
					}

					if !common.StringInSlice(chartID, existChartID) {
						existChartID = append(existChartID, chartID)
						records = append(records, csvRow)
					} else {
						records = updateRecord(csvRow, records)
					}
				}
			}
		}

	}

	if err := recordToCsv(records, setting.OutRecordSetting.CsvPath); err != nil {
		return hadErr, err
	}

	return hadErr, err
}

func recordToCsv(records []Record, outPath string) error {
	//format csv changeTime format
	marshalTime := func(t time.Time) ([]byte, error) {
		return t.AppendFormat(nil, timeFormat), nil
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	enc := csvutil.NewEncoder(w)

	enc.Register(marshalTime)

	if err := enc.Encode(records); err != nil {
		return err
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}

	if err := os.WriteFile(outPath, buf.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func updateRecord(ckRecord Record, records []Record) []Record {

	needAppend := false

	for index, record := range records {
		if (record.ChartID == ckRecord.ChartID) && (record.Recipe == ckRecord.Recipe) && (ckRecord.ChangeTime.After(record.ChangeTime)) {
			records[index] = ckRecord
		} else {
			needAppend = true
		}
	}

	if needAppend {
		records = append(records, ckRecord)
	}

	return records
}
