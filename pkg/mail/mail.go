package mail

import (
	"apc_compare/server/setting"
	"errors"
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"os"
	"path"
)

func SendMail() error {
	ole.CoInitialize(0)

	unknown, err := oleutil.CreateObject("Outlook.Application")
	if err != nil {
		return err
	}
	outlook, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	newMail := oleutil.MustCallMethod(outlook, "CreateItem", "0").ToIDispatch()

	if setting.MailSetting.Subject != "" {
		oleutil.MustPutProperty(newMail, "Subject", setting.MailSetting.Subject)
	}

	if setting.MailSetting.Send != "" {
		oleutil.MustPutProperty(newMail, "To", setting.MailSetting.Send)
	}

	if setting.MailSetting.CC != "" {
		oleutil.MustPutProperty(newMail, "CC", setting.MailSetting.CC)
	}

	if setting.MailSetting.Content != "" {
		html := fmt.Sprintf("<p>%s</p>", setting.MailSetting.Content)
		oleutil.MustPutProperty(newMail, "HTMLBody", html)
	}

	newAtt := oleutil.MustGetProperty(newMail, "Attachments").ToIDispatch()

	rootPath, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath := path.Join(rootPath, setting.OutRecordSetting.CsvPath)
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return err
	}

	oleutil.MustCallMethod(newAtt, "Add", filePath)

	oleutil.MustCallMethod(newMail, "Send").ToIDispatch()

	return nil
}
