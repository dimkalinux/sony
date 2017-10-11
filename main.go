package main

import (
	"encoding/xml"
	"errors"
	"fmt"
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

var endpoints = make(map[string]string, 10)
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

	if err := discoveryCamera(endpoints); err != nil {
		panic(err)
	}

	checkAPIVersion()
	checkAvailableAPI()
	checkSupportedFunctions()
	listFiles()

	fmt.Println("")
	fmt.Println("")
}

func showAppVersion() {
	fmt.Printf("%s version: %s\n\n", AppName, AppVersion)
}

func getAPIUrl(api string) string {
	u, ok := endpoints[api]
	if ok == true {
		return u
	}

	return ""
}

func discoveryCamera(e map[string]string) error {
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
	b := []byte(resp.RawText())

	var deviceInfo DeviceInfo
	xml.Unmarshal(b, &deviceInfo)

	fmt.Printf("Discovered device:\t%s (%s)\n", deviceInfo.Device.FriendlyName, deviceInfo.Device.Manufacturer)

	for _, s := range deviceInfo.Device.ScalarWebAPIDeviceInfo.ServiceList.Services {
		e[s.ServiceType] = s.ActionListURL + "/" + s.ServiceType
	}

	return nil
}

// checkAPIVersion is ...
func checkAPIVersion() {
	p := SonyRequest{1, "getApplicationInfo", "1.0", nil}
	r := SonyArrayResponse{}

	err := makeRequest(getAPIUrl("camera"), &p, &r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Camera API name:\t", strings.TrimSpace(r.GetResult()[0].(string)))
	fmt.Println("Camera API version:\t", strings.TrimSpace(r.GetResult()[1].(string)))
}

// checkAvailableAPI is ...
func checkAvailableAPI() {
	p := SonyRequest{1, "getAvailableApiList", "1.0", nil}
	r := SonyArrayOfArrayResponse{}

	err := makeRequest(getAPIUrl("camera"), &p, &r)
	if err != nil {
		log.Fatal(err)
	}

	var s []string
	for _, val := range r.GetResult() {
		s = append(s, val.(string))
	}

	fmt.Print("Supported:\t\t" + strings.Join(s, ", "))
}

func checkSupportedFunctions() {
	p := SonyRequest{1, "getSupportedCameraFunction", "1.0", []string{}}
	r := SonyArrayResponse{}

	err := makeRequest(getAPIUrl("camera"), &p, &r)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(r.GetResult())

	// var s []string
	// for _, val := range r.GetResult() {
	// 	s = append(s, val.(string))
	// }
	//
	// fmt.Print("Supported:\t\t" + strings.Join(s, ", "))

}

func listFiles() {
	p := SonyRequest{1, "setCameraFunction", "1.0", []string{"Contents Transfer"}}
	r := SonyArrayResponse{}

	err := makeRequest(getAPIUrl("camera"), &p, &r)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(r.GetResult())
}
