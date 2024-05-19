package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strconv"
	"bytes"
	"encoding/json"

	consulapi "github.com/hashicorp/consul/api"
)

func Discovery(serviceName string) []*consulapi.ServiceEntry {
	config := consulapi.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Printf("consul client error: %v", err)
	}
	service, _, err := client.Health().Service(serviceName, "", false, nil)
	if err != nil {
		fmt.Printf("consul client get serviceIp error: %v", err)
	}
	return service
}

type ChannelModel struct {
	LargeScaleModel string    `json:"largescalemodel"`
	SmallScaleModel string    `json:"smallscalemodel"`
}

type Position struct {
	X float64    `json:"x"`
	Y float64    `json:"y"`
	Z float64    `json:"z"`
}

type WirelessNode struct {
	Frequency  float64    `json:"frequency"`
	BitRate    float64    `json:"bitrate"`
	Modulation string     `json:"modulation"`
	BandWidth  float64    `json:"bandwidth"`
	M          float64    `json:"m"`
	PowerInDbm float64    `json:"powerindbm"`
}

type ReqParams struct {
	LinkId 	      int64			`json:"linkid"`
	TxNode 		  WirelessNode	`json:"txnode"`
	RxNode		  WirelessNode	`json:"rxnode"`
	TxPosition 	  Position		`json:"txposition"`
	RxPosition    Position		`json:"rxposition"`
	Model 		  ChannelModel	`json:"model"`
}

func main() {
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response:", string(body))
}

