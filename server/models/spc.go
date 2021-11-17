package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

type SpcData struct {
	ChartID string  `gorm:"column:CHART_ID"` // column name is `CHART_ID`
	USL     float32 `gorm:"column:USL"`      // column name is `USL`
	LSL     float32 `gorm:"column:LSL"`      // column name is `LSL`
}

func GetSpcData(chartID string) (SpcData, error) {
	var spcData SpcData

	//mssql 2008 not support fetch , gorm take first function will use fetch to get data.
	err := db.Table("GRAPH").Select(" MAX(LM_TIME),CHART_ID,USL,LSL").Where("CHART_ID = ? ", chartID).Group("CHART_ID,USL,LSL").Scan(&spcData).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return SpcData{}, err
	}

	//handler empty result
	if (spcData == SpcData{}) {
		return SpcData{}, errors.New(fmt.Sprintf("Not found data with chart id ='%s'", chartID))
	}

	return spcData, nil
}
