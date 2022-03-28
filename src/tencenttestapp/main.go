package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"tencent"

	"github.com/hashicorp/packer/template/interpolate"
)

const (
	HttpGet  = iota
	HttpPost = iota
)

const (
	secretid         = "AKIDELI5jCYbXIsERVLGsJfrsmd1KONeumdO"
	secretkey        = "53sZ5RAoLiuwsgacUVvFWar32eKB5tb9"
	tencentapiprefix = "cvm.tencentcloudapi.com/"
)

type (
	JSONResponse struct {
		Response struct {
			InstanceIdSet []string
			Error         struct {
				Code    string
				Message string
			}
			RequestId string
		}
	}
)

func testhttp(requestType int, requestUrl string, signaturestring string, signature string, returnError bool) {
	const shttp = "http://"
	var jsonResponse JSONResponse
	response := new(http.Response)
	errorstr := ""
	var err error
	if returnError {
		errorstr = "&error=1"
	}
	switch requestType {
	case HttpPost:
		BodyContents := signaturestring + "&Signature=" + signature + errorstr
		log.Println("Body Contents: " + BodyContents)
		Url := shttp + requestUrl
		log.Println("Final URL: " + Url)
		response, err = http.Post(Url, "application/x-www-form-urlencoded",
			strings.NewReader(BodyContents))
	default:
		Url := shttp + requestUrl + signaturestring + errorstr
		log.Println("Final URL: " + Url)
		response, err = http.Get(Url)
	}

	if err == nil {
		defer response.Body.Close()
		contents, _ := ioutil.ReadAll(response.Body)
		_ = json.Unmarshal(contents, &jsonResponse)
	}
}

func describeInstanceTypeConfigs() {
	VMConfig := map[string]interface{}{
		"SecretId":  secretid,
		"SecretKey": secretkey,
	}

	c, warns, errs := tencent.NewSimpleConfig(VMConfig)
	if warns != nil {
	}
	if errs != nil {
	}
	result := &tencent.CVMCreateInstanceResult{}
	log.Printf("%+v", result)

	extraParams := c.Keys()
	secretID := c.SecretID
	secretKey := c.SecretKey
	signaturestring := tencent.SignatureString("DescribeInstanceTypeConfigs", secretID, extraParams)
	url := tencentapiprefix // "cvm.api.qcloud.com/v2/index.php"

	log.Printf("url: %s\n", url)
	log.Printf("signature string: %s\n", signaturestring)

	signature := tencent.SignatureGet(url, signaturestring, secretKey)
	response, err := tencent.RequestGet(url, signaturestring, signature)
	log.Printf("response: %s", string(response))
	log.Printf("err: %s", err)
}

type CloudAPICallResponse struct {
	Response interface{} `json:"Response"`
}

type (
	CVMZoneSet struct {
		Zone      string
		ZoneName  string
		ZoneId    string
		ZoneState string
	}
	CVMZoneResponse struct {
		TotalCount int
		ZoneSet    []CVMZoneSet
		RequestId  string
	}
)

func describeRegions() {
	config := map[string]interface{}{
		"SecretId":  secretid,
		"SecretKey": secretkey,
	}
	response := CloudAPICall("DescribeRegions", config)
	log.Printf("%+v", response)
	var f CloudAPICallResponse
	json.Unmarshal(response, &f)
}

func describeZones(region string) CVMZoneResponse {
	config := map[string]interface{}{
		"SecretId":  secretid,
		"SecretKey": secretkey,
		"Region":    region,
	}
	response := CloudAPICall("DescribeZones", config)
	log.Printf("%+v", response)
	var data CloudAPICallResponse
	json.Unmarshal(response, &data)
	result := new(CVMZoneResponse)
	err := config.Decode(result, &config.DecodeOpts{
		Interpolate: true,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, data)
	return result
}

func testCreateVM() {

	config := map[string]interface{}{
		"ImageId":   "img-3wnd9xpl",
		"Placement": map[string]interface{}{"Zone": "ap-singapore-1"},
		"Region":    "ap-singapore",
		"SecretId":  secretid,
		"SecretKey": secretkey,

		"KeyName":      "notusefulhere",
		"ssh_username": "ubuntu",
	}
	response := CloudAPICall("RunInstances", config)
	log.Printf("%+v", response)
}

func main() {

	if 1 == 2 {
		testhttp(HttpGet, "blah", "blah", "blah", true)
	}
	describeRegions()
	describeZones("ap-singapore")
	// describeInstanceTypeConfigs()
	testCreateVM()

	// signaturestring := tencent.SignatureString("RunInstances", "SecretId", make(map[string]string))
	// signature := "signature"
	// requestUrl := "localhost:81/index.php?"
	// log.Printf("Signature: %s\n", signaturestring)

	// log.Println("HTTP GET")
	// testhttp(HttpGet, requestUrl, signaturestring, signature, false)
	// log.Println("")

	// log.Println("HTTP POST")
	// testhttp(HttpPost, requestUrl, signaturestring, signature, false)
	// log.Println("")

	// log.Println("HTTP GET ERROR")
	// testhttp(HttpGet, requestUrl, signaturestring, signature, true)
	// log.Println("")

	// log.Println("HTTP POST ERROR")
	// testhttp(HttpPost, requestUrl, signaturestring, signature, true)
	// log.Println("")

	// log.Println("HTTP DONE")

}
