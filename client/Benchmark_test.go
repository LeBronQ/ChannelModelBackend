package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strconv"
	"bytes"
	"encoding/json"
	"testing"

)

//API test
func Benchmark(b *testing.B) {
	b.StopTimer()
	se := Discovery("Default_Model")
	port := se[0].Service.Port
	address := se[0].Service.Address
	request := "http://" + address + ":" + strconv.Itoa(port) + "/model"
	wirelessNode := WirelessNode{
		Frequency:  2.4e+9,
		BitRate:    5.0e+7,
		Modulation: "BPSK",
		BandWidth:  2.0e+7,
		M:          0,
		PowerInDbm: 20,
	}
	txPos := Position{
		X:  0.0,
		Y:  0.0,
		Z:  100.0,
	}
	rxPos := Position{
		X:  1000.0,
		Y:  1000.0,
		Z:  100.0,
	}
	mod := ChannelModel{
		LargeScaleModel: "FreeSpacePathLossModel",
		SmallScaleModel: "NakagamiFadingModel",
	}
	param := ReqParams{
		LinkId:  0,
		TxNode: wirelessNode,
		RxNode: wirelessNode,
		TxPosition: txPos,
		RxPosition: rxPos,
		Model: mod,
	}
	b.StartTimer()

	for i := 0; i < b.N; i++{
		jsonData, err := json.Marshal(param)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}
		requestBody := bytes.NewBuffer(jsonData)
	req, err := http.NewRequest("POST", request, requestBody)
    if err != nil {
        fmt.Println(err)
        return
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Content-Length", fmt.Sprintf("%d", requestBody.Len()))
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.StatusCode)
		return
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	}
	//fmt.Println("Response:", string(body))
}
