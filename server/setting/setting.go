package setting

import (
	"github.com/go-ini/ini"
	"log"
)

type App struct {
	User        string
	Password    string
	FabSite     string
	ToolID      []string
	ProcessType string
	DataType    string
	Ziped       string
}

type OutRecord struct {
	CsvPath string
}

type Mail struct {
	Subject string
	Content string
	Send    string
	CC      string
}

var cfg *ini.File
var AppSetting = &App{}
var OutRecordSetting = &OutRecord{}
var MailSetting = &Mail{}

func Setup() {
	var err error
	cfg, err = ini.Load("conf/config.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/config.ini': %v", err)
	}
	mapTo("app", AppSetting)
	mapTo("outRecord", OutRecordSetting)
	mapTo("mail", MailSetting)

}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
