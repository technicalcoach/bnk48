package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

type Envelope struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	Soapenv string   `xml:"xmlns:soapenv,attr"`
	Awd     string   `xml:"xmlns:awd,attr"`
	Header  string   `xml:"soapenv:Header"`
	Body    Body     `xml:"soapenv:Body"`
}
type Body struct {
	AreYouThere string `xml:"awd:areYouThere"`
}

func main() {

	url := "https://wcc.sc.egov.usda.gov/awdbWebService/services"

	request := Envelope{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Awd:     "http://www.wcc.nrcs.usda.gov/ns/awdbWebService",
	}
	b, _ := xml.Marshal(&request)

	payload := bytes.NewReader(b)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Postman-Token", "d81729ce-835f-4eac-b282-e0e09dd89218")

	bd, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(bd))

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}
