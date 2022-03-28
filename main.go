package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"tencent"

	"github.com/hashicorp/packer/packer"
)

const (
	HttpGet  = iota
	HttpPost = iota
)

const (
	secretid  = "AKIDELI5jCYbXIsERVLGsJfrsmd1KONeumdO"
	secretkey = "53sZ5RAoLiuwsgacUVvFWar32eKB5tb9"
	// tencentapiprefix = "cvm.tencentcloudapi.com/"
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

// func describeInstanceTypeConfigs() {
// 	VMConfig := map[string]interface{}{
// 		"SecretId":  secretid,
// 		"SecretKey": secretkey,
// 	}

// 	c, warns, errs := tencent.NewSimpleConfig(VMConfig)
// 	if warns != nil {
// 	}
// 	if errs != nil {
// 	}
// 	result := &tencent.CVMCreateInstanceResult{}
// 	log.Printf("%+v", result)

// 	extraParams := c.Keys()
// 	secretID := c.SecretID
// 	secretKey := c.SecretKey
// 	signaturestring := tencent.SignatureString("DescribeInstanceTypeConfigs", secretID, extraParams)
// 	url := tencent.CloudProviderPrefix // "cvm.api.qcloud.com/v2/index.php"

// 	log.Printf("url: %s\n", url)
// 	log.Printf("signature string: %s\n", signaturestring)

// 	signature := tencent.SignatureGet(url, signaturestring, secretKey)
// 	response, err := tencent.RequestGet(url, signaturestring, signature)
// 	log.Printf("response: %s", string(response))
// 	log.Printf("err: %s", err)
// }

type (
	CVMZoneSet struct {
		Zone      string
		ZoneName  string
		ZoneId    string
		ZoneState string
	}

	CVMZones struct {
		TotalCount int
		ZoneSet    []CVMZoneSet
		RequestId  string
	}

	CVMRegionSet struct {
		Region      string
		RegionName  string
		RegionState string
	}

	CVMRunInstancesResponse struct {
		InstanceIdSet []string
		RequestId     string
	}

	CVMInstanceInfo struct {
		InstanceId string
		Region     string
	}

	CVMInstanceStatus struct {
		InstanceId    string // ins-ggitl6mi
		InstanceState string // PENDING, RUNNING, STOPPED, REBOOTING, STARTING, STOPPING
	}

	// CVMInstanceStatusResponse is what's returned from calling DescribeInstancesStatus
	CVMInstanceStatusResponse struct {
		TotalCount        int
		InstanceStatusSet []CVMInstanceStatus
		RequestId         string
	}

	// CVMRegions example
	// "TotalCount": 2,
	// "RegionSet": [
	// 	{
	// 		"Region": "ap-bangkok",
	// 		"RegionName": "亚太地区(曼谷)",
	// 		"RegionState": "AVAILABLE"
	// 	},
	// 	{
	// 		"Region": "ap-beijing",
	// 		"RegionName": "华北地区(北京)",
	// 		"RegionState": "AVAILABLE"
	// 	},
	// ],
	// "RequestId": "1219fa76-a6c0-462a-b8a8-fa40f3f2637e"
	CVMRegions struct {
		TotalCount int
		RegionSet  []CVMRegionSet
		RequestId  string
	}
)

func describeRegions() (CVMRegions, tencent.CVMError) {
	configInfo := map[string]interface{}{
		"SecretId":  secretid,
		"SecretKey": secretkey,
	}
	response := tencent.CloudAPICall("DescribeRegions", configInfo, nil)
	var (
		regionInfo CVMRegions
		cvmError   tencent.CVMError
	)
	err := tencent.DecodeResponse(response, &regionInfo)
	if err != nil {
		tencent.DecodeResponse(response, cvmError)
	}
	return regionInfo, cvmError
}

func describeZones(region string) (CVMZones, tencent.CVMError) {
	configinfo := map[string]interface{}{
		"SecretId":  secretid,
		"SecretKey": secretkey,
	}
	extraParams := map[string]string{
		"Region": region,
	}
	response := tencent.CloudAPICall("DescribeZones", configinfo, extraParams)
	var (
		zoneInfo CVMZones
		cvmError tencent.CVMError
	)
	err := tencent.DecodeResponse(response, &zoneInfo)

	if err != nil {
		tencent.DecodeResponse(response, &cvmError)
	}
	return zoneInfo, cvmError
}

func testCreateVM1() (bool, tencent.CVMInstanceInfo) {
	// configInfo := map[string]string{
	// 	"SecretId":  "error" + secretid,
	// 	"SecretKey": secretkey,
	// }
	// extraParams := map[string]string{
	// 	"ImageID":        "img-3wnd9xpl",
	// 	"Placement.Zone": "ap-singapore-1",
	// 	"Region":         "ap-singapore",
	// 	"KeyName":        "notusefulhere",
	// 	"ssh_username":   "ubuntu",
	// 	"Version":        "2017-03-12",
	// }
	var c tencent.Config
	c.SecretID = secretid
	c.SecretKey = secretkey
	c.ImageID = "img-3wnd9xpl"
	c.Placement.Zone = "ap-singapore-1"
	c.Region = "ap-singapore"
	c.SSHKeyName = "notusehere"
	c.SSHUserName = "ubuntu"
	c.Version = "2017-03-12"
	var ui packer.Ui

	driver := tencent.NewTencentDriver(ui, &c)
	err1, _, instanceInfo := driver.CreateImage(c)

	return err1, instanceInfo
}

// {
//     "Response": {
//         "InstanceIdSet": [
//             "ins-ggitl6mi"
//         ],
//         "RequestId": "55c68d81-808e-4e1e-9590-007064a0b5e6"
//     }
// }
func testCreateVM2() (bool, CVMInstanceInfo) {
	// configInfo := map[string]string{
	// 	"SecretId":  "error" + secretid,
	// 	"SecretKey": secretkey,
	// }
	// extraParams := map[string]string{
	// 	"ImageID":        "img-3wnd9xpl",
	// 	"Placement.Zone": "ap-singapore-1",
	// 	"Region":         "ap-singapore",
	// 	"KeyName":        "notusefulhere",
	// 	"ssh_username":   "ubuntu",
	// 	"Version":        "2017-03-12",
	// }
	var c tencent.Config
	c.SecretID = secretid
	c.SecretKey = secretkey
	c.ImageID = "img-3wnd9xpl"
	c.Placement.Zone = "ap-singapore-1"
	c.Region = "ap-singapore"
	c.SSHKeyName = "notusehere"
	c.SSHUserName = "ubuntu"
	c.Version = "2017-03-12"
	configInfo := c.CreateVMmap()

	response := tencent.CloudAPICall("RunInstances", configInfo, nil)

	// response :=
	// 	`{
	// 		"Response": {
	// 			"InstanceIdSet": [
	// 				"ins-ggitl6mi"
	// 			],
	// 			"RequestId": "55c68d81-808e-4e1e-9590-007064a0b5e6"
	// 		}
	// 	}`
	var runInstancesResponse CVMRunInstancesResponse
	err2 := tencent.DecodeResponse([]byte(response), &runInstancesResponse)
	error := err2 != nil
	instanceid := ""
	if !error {
		instanceid = runInstancesResponse.InstanceIdSet[0]
	}
	instanceInfo := CVMInstanceInfo{instanceid, configInfo["Region"].(string)}
	// log.Printf("%s", string(response))
	return error, instanceInfo
}

func GetAllVMStatus() []CVMInstanceStatus {
	configInfo := map[string]interface{}{
		"SecretId":  secretid,
		"SecretKey": secretkey,
	}
	params := map[string]string{
		"Region": "ap-singapore",
	}
	response := tencent.CloudAPICall("DescribeInstancesStatus", configInfo, params)
	log.Printf("%s", string(response))
	var (
		instanceStatusResponse CVMInstanceStatusResponse
	)
	err := tencent.DecodeResponse(response, &instanceStatusResponse)
	if err == nil && instanceStatusResponse.TotalCount > 0 {
		return instanceStatusResponse.InstanceStatusSet
	}
	return []CVMInstanceStatus{}
}

// GetVMStatus will return the status of the instance in the given Region
func GetVMStatus(instanceInfo tencent.CVMInstanceInfo) CVMInstanceStatus {
	configInfo := map[string]interface{}{
		"SecretId":  secretid,
		"SecretKey": secretkey,
	}
	params := map[string]string{
		"InstanceIds.0": instanceInfo.InstanceId,
		"Region":        instanceInfo.Region,
	}
	response := tencent.CloudAPICall("DescribeInstancesStatus", configInfo, params)
	log.Printf("%s", string(response))
	var (
		instanceStatusResponse CVMInstanceStatusResponse
		emptyInstanceStatus    CVMInstanceStatus
	)
	err := tencent.DecodeResponse(response, &instanceStatusResponse)
	if err != nil && instanceStatusResponse.TotalCount > 0 {
		return instanceStatusResponse.InstanceStatusSet[0]
	}
	return emptyInstanceStatus
}

func main() {
	tencent.CloudAPIDebug = true
	if 1 == 2 {
		testhttp(HttpGet, "blah", "blah", "blah", true)
		describeRegions()
		describeZones("ap-singapore")
		// describeInstanceTypeConfigs()
		allVMsStatus := GetAllVMStatus()
		log.Printf("%+v", allVMsStatus)
	}
	err, instanceInfo := testCreateVM1()
	if !err {
		GetVMStatus(instanceInfo)
	}
	// tencent.StopVM(instanceInfo)

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
