package main

import (
	"encoding/xml"
	"net/url"
)

// DeviceInfo is ...
type DeviceInfo struct {
	XMLName     xml.Name    `xml:"root"`
	SpecVersion SpecVersion `xml:"specVersion"`
	Device      Device      `xml:"device"`
}

// SpecVersion is...
type SpecVersion struct {
	XMLName xml.Name `xml:"specVersion"`
	Major   string   `xml:"major"`
	Minor   string   `xml:"minor"`
}

// Device is...
type Device struct {
	XMLName                xml.Name               `xml:"device"`
	DeviceType             string                 `xml:"deviceType"`
	FriendlyName           string                 `xml:"friendlyName"`
	Manufacturer           string                 `xml:"manufacturer"`
	ManufacturerURL        url.URL                `xml:"manufacturerURL"`
	ModelDescription       string                 `xml:"modelDescription"`
	ModelName              string                 `xml:"modelName"`
	UDN                    string                 `xml:"UDN"`
	ServiceList            ServiceList            `xml:"serviceList"`
	ScalarWebAPIDeviceInfo ScalarWebAPIDeviceInfo `xml:"X_ScalarWebAPI_DeviceInfo"`
}

// ServiceList is ...
type ServiceList struct {
	XMLName  xml.Name             `xml:"serviceList"`
	Services []ServiceListService `xml:"service"`
}

// ServiceListService is ...
type ServiceListService struct {
	XMLName     xml.Name `xml:"service"`
	ServiceType string   `xml:"serviceType"`
	ServiceID   string   `xml:"serviceId"`
	SCPDURL     url.URL  `xml:"SCPDURL"`
	ControlURL  url.URL  `xml:"controlURL"`
	EventSubURL url.URL  `xml:"eventSubURL"`
}

// ScalarWebAPIDeviceInfo is ...
type ScalarWebAPIDeviceInfo struct {
	XMLName     xml.Name                  `xml:"X_ScalarWebAPI_DeviceInfo"`
	Version     string                    `xml:"X_ScalarWebAPI_Version"`
	ServiceList []ScalarWebAPIServiceList `xml:"X_ScalarWebAPI_ServiceList"`
}

// ScalarWebAPIServiceList is...
type ScalarWebAPIServiceList struct {
	XMLName  xml.Name                         `xml:"X_ScalarWebAPI_ServiceList"`
	Services []ScalarWebAPIServiceListService `xml:"X_ScalarWebAPI_Service"`
}

// ScalarWebAPIServiceListService is...
type ScalarWebAPIServiceListService struct {
	XMLName       xml.Name `xml:"X_ScalarWebAPI_Service"`
	ServiceType   string   `xml:"X_ScalarWebAPI_ServiceType"`
	ActionListURL url.URL  `xml:"X_ScalarWebAPI_ActionList_URL"`
}
