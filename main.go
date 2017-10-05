package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	ssdp "github.com/dimkalinux/go-ssdp"
	flags "github.com/jessevdk/go-flags"
	"gopkg.in/jmcvetta/napping.v3"
)

const (
	// APIURL url.
	APIURL = "http://192.168.122.1:8080/sony"

	// AppName is const for app name.
	AppName = "Sony app"

	// AppVersion is const defined app version.
	AppVersion = "0.0.2"
)

// Options is...
type Options struct {
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`
	Version bool   `short:"V" long:"version" description:"Show program version"`
}

// APIEndpoint  is...
type APIEndpoint struct {
	Action string
	URL    string
}

var endpoints []APIEndpoint
var options Options
var parser = flags.NewParser(&options, flags.Default)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func main() {
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	if true == options.Version {
		showAppVersion()
		os.Exit(0)
	}

	// if err := discoveryCamera(); err != nil {
	// 	log.Fatal(err)
	// }

	xmlFile, err := os.Open("sony.xml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	var deviceInfo DeviceInfo
	xml.Unmarshal(byteValue, &deviceInfo)

	fmt.Printf("Device name: %s (%s)\n", deviceInfo.Device.FriendlyName, deviceInfo.Device.Manufacturer)

	for _, s := range deviceInfo.Device.ScalarWebAPIDeviceInfo.ServiceList.Services {
		fmt.Printf("Action: %s, url: %s\n", s.ServiceType, s.ActionListURL)
		endpoints = append(endpoints, APIEndpoint{s.ServiceType, s.ActionListURL + "/" + s.ServiceType})
	}

	spew.Dump(endpoints)
	// checkAPIVersion()
	// checkAvailableAPI()

	fmt.Println("")
}

func showAppVersion() {
	fmt.Printf("%s version: %s\n\n", AppName, AppVersion)
}

func getAPIUrl(api string) string {
	for _, s := range endpoints {
		if s.Action == api {
			return s.URL
		}
	}

	return ""
}

func discoveryCamera() error {
	responses, err := ssdp.Search("urn:schemas-sony-com:service:ScalarWebAPI:1", 2*time.Second)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if len(responses) == 0 {
		return errors.New("can not found devices. Please connect to camery WiFi")
	}

	location := responses[0].Location.String()
	session := napping.Session{Log: false}
	resp, _ := session.Get(location, nil, nil, nil)

	log.Print(resp.RawText())

	return nil
}

// checkAPIVersion is ...
func checkAPIVersion() {
	p := SonyRequest{1, "getApplicationInfo", "1.0", []int{}}
	r := SonyArrayResponse{}

	err := makeRequest(getAPIUrl("camera"), &p, &r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Camera API name: ", r.GetResult()[0])
	fmt.Println("Camera API version: ", r.GetResult()[1])
}

// checkAvailableAPI is ...
func checkAvailableAPI() {
	p := SonyRequest{1, "getAvailableApiList", "1.0", []int{}}
	r := SonyArrayOfArrayResponse{}

	err := makeRequest(getAPIUrl("camera"), &p, &r)
	if err != nil {
		log.Fatal(err)
	}

	var s []string
	for _, val := range r.GetResult() {
		s = append(s, val.(string))
	}

	fmt.Print("Supported: " + strings.Join(s, ", "))
}
