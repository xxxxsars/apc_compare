package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}


func getApcData() error{

	//defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body) // response body is []byte
	//fmt.Println(string(body))


	tools := []string{
		"SPG02A", "SPG03A", "SPI01A",
	}

	apcParams := map[string]string{
		"user":        "ML3DT1_AP",
		"password":    "ML3DT1$AP",
		"FABSITE":     "3D",
		"TOOLID":      "",
		"CHAMBER":     "",
		"PARAMETER":   "",
		"ziped":       "false",
		"DATATYPE":    "Runsummary",
		"PROCESSTYPE": "ARRAY",
	}

	for _, toolId := range tools {
		apcParams["TOOLID"] = toolId

		client := &http.Client{}

		req, err := http.NewRequest("GET","http://10.21.10.124/APCWEB3D/APC_WBS/APCService.asmx/GetAPC_SpecList",nil)
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")


		q := req.URL.Query()

		for key ,value := range apcParams{

			q.Add(key, value)
			//rep.URL.RawQuery = q.Encode()
		}

		resp, err := client.Do(req)
		if err != nil {

			return err
		}

		defer resp.Body.Close()


		//fmt.Println(rep.URL.String())
		//
		//fmt.Println(rep.Body)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Printf("%s", body)

		respJson := make(map[interface{}]interface{})
		json.Unmarshal(body, &respJson)
		fmt.Printf("Results: %v\n", respJson)
	}

	return nil
}