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
	// ApiUrl url.
	ApiUrl = "http://192.168.122.1:8080/sony"

	// AppName is const for app name.
	AppName = "Sony app"

	// AppVersion is const defined app version.
	AppVersion = "0.0.1"
)

type Options struct {
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`
	Version bool   `short:"V" long:"version" description:"Show program version"`
}

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
	}

	fmt.Println("Successfully Opened sony.xml")
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	var deviceInfo DeviceInfo
	xml.Unmarshal(byteValue, &deviceInfo)

	spew.Dump(deviceInfo)

	// checkAPIVersion()
	// checkAvailableAPI()

	fmt.Println("")
}

func showAppVersion() {
	fmt.Printf("%s version: %s\n\n", AppName, AppVersion)
}

func discoveryCamera() error {
	responses, err := ssdp.Search("urn:schemas-sony-com:service:ScalarWebAPI:1", 2*time.Second)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if len(responses) == 0 {
		return errors.New("Can not found devices. Please connect to camery WiFi.")
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

	err := makeRequest("/camera", &p, &r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Camera API name:", r.GetResult()[0])
	fmt.Println("Camera API version:", r.GetResult()[1])
}

// checkAvailableAPI is ...
func checkAvailableAPI() {
	p := SonyRequest{1, "getAvailableApiList", "1.0", []int{}}
	r := SonyArrayOfArrayResponse{}

	err := makeRequest("/camera", &p, &r)
	if err != nil {
		log.Fatal(err)
	}

	var s []string
	for _, val := range r.GetResult() {
		s = append(s, val.(string))
	}

	fmt.Print("Supported: " + strings.Join(s, ", "))
}
