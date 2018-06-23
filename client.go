package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pallat/tis620"
)

func main() {

	url := "http://192.168.100.3:8080/thai"

	req, _ := http.NewRequest("GET", url, nil)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	// tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}
	res, _ := client.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(tis620.ToUTF8(string(body)))

}
