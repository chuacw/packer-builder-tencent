package tencent

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	//  "crypto/sha1"
	//  "encoding/hex"
	"encoding/base64"
	"encoding/json"
	"net/url"

	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/template/interpolate"
)

type (
	// Response returned by a call to CloudAPICall
	CloudAPICallResponse struct {
		Response map[string]interface{} `json:"Response"`
	}

	// CVMPlacement represents the structure required for creating a VM.
	CVMPlacement struct {
		Zone      string   `mapstructure:"Zone"`
		ProjectId int64    `mapstructure:"ProjectId"`
		HostIds   []string `mapstructure:"HostIds"`
	}

	// Error example
	// "Code": "AuthFailure.SecretIdNotFound",
	// "Message": "The SecretId is not found, please ensure that your SecretId is correct."
	CVMError struct {
		Code    string
		Message string
	}

	// CVMError example
	// "Error": {
	// 	"Code": "AuthFailure.SecretIdNotFound",
	// 	"Message": "The SecretId is not found, please ensure that your SecretId is correct."
	// },
	// "RequestId": "0f2bdb80-74dc-418b-b9e0-16e8c98288ed"
	CVMErrorResponse struct {
		Error     CVMError
		RequestId string
	}

	// CVMCreateInstanceResult represents the structure for the result returned from creating a VM
	CVMCreateInstanceResponse struct {
		InstanceIdSet []string
		RequestId     string
		Error         CVMError
	}

	// JSONResponse struct {
	// 	Response struct {
	// 		InstanceIdSet []string
	// 		RequestId     string
	// 		Error         struct {
	// 			Code    string
	// 			Message string
	// 		}
	// 	}
	// }

	CVMVirtualPrivateCloud struct {
		VpcId              string
		SubnetId           string
		AsVpcGateway       bool
		PrivateIpAddresses []string
	}

	CVMSystemDisk struct {
		DiskType string
		DiskId   string
		DiskSize int64
	}

	CVMDataDisk struct {
		DiskType string
		DiskId   string
		DiskSize int64
	}

	CVMInternetAccessible struct {
		InternetChargeType      string
		InternetMaxBandwidthOut int64
		PublicIpAssigned        bool
	}

	CVMInstanceSet struct {
		Placement           CVMPlacement
		InstanceId          string
		InstanceType        string
		CPU                 int64
		Memory              int64
		InstanceName        string
		InstanceChargeType  string
		SystemDisk          CVMSystemDisk
		DataDisks           []CVMDataDisk
		PrivateIpAddresses  []string
		PublicIpAddresses   []string
		InternetAccessible  CVMInternetAccessible
		VirtualPrivateCloud CVMVirtualPrivateCloud
		ImageId             string
		RenewFlag           string
		CreatedTime         string
		ExpiredTime         string
	}

	LoginSettings struct {
		Password       string
		KeyIds         []string
		KeepImageLogin string
	}

	RunSecurityServiceEnabled struct {
		Enabled bool
	}

	RunMonitorServiceEnabled struct {
		Enabled bool
	}

	EnhancedService struct {
		SecurityService RunSecurityServiceEnabled
		MonitorService  RunMonitorServiceEnabled
	}

	InstanceChargePrepaid struct {
		Period    int
		RenewFlag string
	}

	CVMInstanceInfo struct {
		InstanceId string
		Region     string
	}

	CVMRunInstancesResponse struct {
		InstanceIdSet []string
		RequestId     string
	}

	CVMDescribeProjectResponse struct {
		projectName string
		projectId   string
		createTime  string
		creatorUin  string
		projectInfo string
	}
)

var (
	CloudAPIDebug       bool
	CloudProviderPrefix = "cvm.tencentcloudapi.com/"
)

// The request should look like as follows.
// see https://cloud.tencent.com/document/api/213/9384#example-2
// https://cvm.api.qcloud.com/v2/index.php?Action=RunInstances
// &Version=2017-03-12
// &Placement.Zone=ap-guangzhou-2
// &InstanceChargeType=PREPAID
// &InstanceChargePrepaid.Period=1
// &InstanceChargePrepaid.RenewFlag=NOTIFY_AND_AUTO_RENEW
// &ImageId=img-pmqg1cw7
// &InstanceType=S1.SMALL1
// &SystemDisk.DiskType=LOCAL_BASIC
// &SystemDisk.DiskSize=50
// &DataDisks.0.DiskType=LOCAL_BASIC
// &DataDisks.0.DiskSize=100
// &InternetAccessible.InternetChargeType=TRAFFIC_POSTPAID_BY_HOUR
// &InternetAccessible.InternetMaxBandwidthOut=10
// &InternetAccessible.PublicIpAssigned=TRUE
// &InstanceName=QCLOUD-TEST
// &LoginSettings.Password=Qcloud@TestApi123++
// &EnhancedService.SecurityService.Enabled=TRUE
// &EnhancedService.MonitorService.Enabled=TRUE
// &InstanceCount=1
// &<Common request parameters>
// All parameters must be sorted: https://cloud.tencent.com/document/api/213/11652#2.1.-sort-parameters
func CreateVM(c *Config) (CVMError, CVMInstanceInfo) {
	configInfo := c.CreateVMmap()
	if c.PackerDebug || CloudAPIDebug {
		log.Printf("CreateVM configInfo: %+v", configInfo)
	}
	response := CloudAPICall("RunInstances", configInfo, nil)
	var (
		runInstancesResponse CVMRunInstancesResponse
		cvmErrorResponse     CVMErrorResponse
		instanceInfo         CVMInstanceInfo
	)
	if c.PackerDebug || CloudAPIDebug {
		log.Printf("CloudAPICall response content: %s", string(response))
	}
	err := DecodeResponse(response, &runInstancesResponse)
	error := err != nil
	instanceid := ""
	if error {
		DecodeResponse(response, &cvmErrorResponse)
		if c.PackerDebug || CloudAPIDebug {
			log.Printf("The error is: %+v", cvmErrorResponse.Error)
		}
	} else {
		log.Printf("CloudAPICall successful: %+v", runInstancesResponse)
		instanceid = runInstancesResponse.InstanceIdSet[0]
		instanceInfo = CVMInstanceInfo{instanceid, configInfo["Region"].(string)}
	}
	// log.Printf("%s", string(response))
	return cvmErrorResponse.Error, instanceInfo
}

func DecodeResponse(data []byte, target interface{}) error {
	var decodedResponse CloudAPICallResponse
	json.Unmarshal(data, &decodedResponse)
	err := config.Decode(target, &config.DecodeOpts{
		Interpolate: true,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, decodedResponse.Response)
	if CloudAPIDebug {
		log.Printf("%+v", err)
	}
	return err
}

// CloudAPICall sample response for error
// Error 1
// {
// 	"Response": {
// 		"Error": {
// 			"Code": "AuthFailure.SecretIdNotFound",
// 			"Message": "The SecretId is not found, please ensure that your SecretId is correct."
// 		},
// 		"RequestId": "9e79c786-c4a8-48a6-9419-8825e2466ae0"
// 	}
// }
// Error 2
// {
// 	"Response": {
// 		"Error": {
// 			"Code": "AuthFailure.SignatureFailure",
// 			"Message": "The provided credentials could not be validated. Please check your signature is correct."
// 		},
// 		"RequestId": "4dd40e21-4c7d-4166-a3ca-0da802ad2852"
// 	}
// }
// sample response for successful call
// {
// 	"Response": {
// 		"TotalCount": 17,
// 		"RegionSet": [
// 			{
// 				"Region": "ap-bangkok",
// 				"RegionName": "亚太地区(曼谷)",
// 				"RegionState": "AVAILABLE"
// 			},
// 			{
// 				"Region": "ap-beijing",
// 				"RegionName": "华北地区(北京)",
// 				"RegionState": "AVAILABLE"
// 			},
// 			{
// 				"Region": "ap-chengdu",
// 				"RegionName": "西南地区(成都)",
// 				"RegionState": "AVAILABLE"
// 			},
// 			{
// 				"Region": "ap-chongqing",
// 				"RegionName": "西南地区(重庆)",
// 				"RegionState": "AVAILABLE"
// 			},
// 			{
// 				"Region": "ap-guangzhou",
// 				"RegionName": "华南地区(广州)",
// 				"RegionState": "AVAILABLE"
// 			},
// 			{
// 				"Region": "ap-guangzhou-open",
// 				"RegionName": "华南地区(广州Open)",
// 				"RegionState": "AVAILABLE"
// 			},
// 			{
// 				"Region": "ap-hongkong",
// 				"RegionName": "东南亚地区(香港)",
// 				"RegionState": "AVAILABLE"
// 			},
// 			{
// 				"Region": "ap-mumbai",
// 				"RegionName": "亚太地区(孟买)",
// 				"RegionState": "AVAILABLE"
// 			},
// 			{
// 				"Region": "ap-seoul",
// 				"RegionName": "东南亚地区(首尔)",
// 				"RegionState": "AVAILABLE"
// 			},
// 			{
// 				"Region": "ap-shanghai",
// 				"RegionName": "华东地区(上海)",
// 				"RegionState": "AVAILABLE"
// 			},
// 			{
// 				"Region": "ap-shanghai-fsi",
// 				"RegionName": "华东地区(上海金融)",
// 				"RegionState": "AVAILABLE"
// 			},
// 			{
// 				"Region": "ap-shenzhen-fsi",
// 				"RegionName": "华南地区(深圳金融)",
// 				"RegionState": "AVAILABLE"
// 			},
// 			{
// 				"Region": "ap-singapore",
// 				"RegionName": "东南亚地区(新加坡)",
// 				"RegionState": "AVAILABLE"
// 			},
// 			{
// 				"Region": "eu-frankfurt",
// 				"RegionName": "欧洲地区(德国)",
// 				"RegionState": "AVAILABLE"
// 			},
// 			{
// 				"Region": "na-ashburn",
// 				"RegionName": "美国东部(弗吉尼亚)",
// 				"RegionState": "AVAILABLE"
// 			},
// 			{
// 				"Region": "na-siliconvalley",
// 				"RegionName": "美国西部(硅谷)",
// 				"RegionState": "AVAILABLE"
// 			},
// 			{
// 				"Region": "na-toronto",
// 				"RegionName": "北美地区(多伦多)",
// 				"RegionState": "AVAILABLE"
// 			}
// 		],
// 		"RequestId": "769cdd57-f4d6-415c-a8e2-5776c2d6d939"
// 	}
// }
// To call the CloudAPICall function, provide it the name
// of an action, eg, RunInstances, DescribeZones, DescribeRegions
// a map[string]interface{} containing the required parameters
// and another parameter containing
func CloudAPICall(action string, config map[string]interface{},
	extraParams2 map[string]string) []byte {
	var (
		c    *Config
		err1 error
	)
	c, _, err1 = NewSimpleConfig(config)
	if err1 != nil {
		log.Printf("%+v: ", err1)
	}

	extraParams := c.Keys()
	// merges the keys and values into extraParams
	// if extraParams2 is nil, this is skipped
	for k, v := range extraParams2 {
		extraParams[k] = v
	}
	secretID := c.SecretID
	secretKey := c.SecretKey
	signaturestring := SignatureString(action, secretID, extraParams)
	url := CloudProviderPrefix
	if CloudAPIDebug {
		log.Printf("url: %s\n", url)
		log.Printf("URI: %s\n", signaturestring)
	}

	signature := SignatureGet(url, signaturestring, secretKey)
	response, err2 := RequestGet(url, signaturestring, signature)
	if err2 != nil && CloudAPIDebug {
		log.Printf("Error in CloudAPICall: %s", err2)
	}
	if CloudAPIDebug {
		log.Printf("HTTP response: %s", string(response))
	}
	return response
}

func DescribeProject() {
	CloudAPICall("DescribeProject", configInfo, nil)
}

func SignatureStringNonceTimestamp(action, nonce, timestamp, secretId string, extraParams map[string]string) string {
	var sortparams = []string{}

	// Common request parameters
	params := make(map[string]string)
	params["Action"] = action
	sortparams = append(sortparams, "Action")
	params["Nonce"] = nonce
	sortparams = append(sortparams, "Nonce")
	params["Timestamp"] = timestamp
	sortparams = append(sortparams, "Timestamp")
	params["SecretId"] = secretId
	sortparams = append(sortparams, "SecretId")

	// params["SignatureMethod"] = "HmacSHA256"
	// sortparams = append(sortparams, "SignatureMethod")

	for k, v := range extraParams {
		params[k] = v
		sortparams = append(sortparams, k)
	}

	sort.Strings(sortparams)

	requestParamString := ""
	var paramstr = []string{}
	for _, requestKey := range sortparams {
		if params[requestKey] != "" {
			paramstr = append(paramstr, requestKey+"="+params[requestKey])
		}
	}

	requestParamString += strings.Join(paramstr, "&")
	return requestParamString

}

// Sample usage
//   log.Println(tencent.SignatureString("action", "SecretId", make(map[string]string)))
//   log.Println(tencent.SignatureString("action", "SecretId", nil))
func SignatureString(action, secretId string, extraParams map[string]string) string {
	nonce := GenerateNonce()
	timestamp := CurrentTimeStamp()
	return SignatureStringNonceTimestamp(action, nonce, timestamp, secretId, extraParams)
	// var sortparams = []string{}

	// // Common request parameters
	// params := make(map[string]string)
	// params["Action"] = action
	// sortparams = append(sortparams, "Action")
	// params["Nonce"] = GenerateNonce()
	// sortparams = append(sortparams, "Nonce")
	// params["Timestamp"] = CurrentTimeStamp()
	// sortparams = append(sortparams, "Timestamp")
	// params["SecretId"] = secretId
	// sortparams = append(sortparams, "SecretId")

	// // params["SignatureMethod"] = "HmacSHA256"
	// // sortparams = append(sortparams, "SignatureMethod")

	// for paramKey, paramValue := range extraParams {
	// 	params[paramKey] = paramValue
	// 	sortparams = append(sortparams, paramKey)
	// }

	// sort.Strings(sortparams)

	// requestParamString := ""
	// var paramstr = []string{}
	// for _, requestKey := range sortparams {
	// 	if params[requestKey] != "" {
	// 		paramstr = append(paramstr, requestKey+"="+params[requestKey])
	// 	}
	// }

	// requestParamString += strings.Join(paramstr, "&")
	// return requestParamString

}

// Signature generates the Base64 encoded HMAC key for the given parameters
func SignatureGet(requestURL, requestParamString, secretKey string) string {
	// this POST string below, requires the request to be sent using http.Post!
	signstr := "GET" + requestURL + "?" + requestParamString
	// signature := Hmac256ToBase64(secretKey, signstr, true)
	signature := Hmac1ToBase64(secretKey, signstr, true) // if post, last param is false
	return signature
}

// Signature generates the Base64 encoded HMAC key for the given parameters
func SignaturePost(requestURL, requestParamString, secretKey string) string {
	// this POST string below, requires the request to be sent using http.Post!
	signstr := "POST" + requestURL + "?" + requestParamString
	// signature := Hmac256ToBase64(secretKey, signstr, true)
	signature := Hmac1ToBase64(secretKey, signstr, false) // if post, last param is false
	return signature
}

// Signature generates the Base64 encoded HMAC key for the given parameters
func Signature(requestURL, requestParamString, secretKey string) string {
	// this POST string below, requires the request to be sent using http.Post!
	signstr := "POST" + requestURL + "?" + requestParamString
	// signature := Hmac256ToBase64(secretKey, signstr, true)
	signature := Hmac1ToBase64(secretKey, signstr, true) // if post, last param is false
	return signature
}

// Request makes a call to the HTTP endpoint
func RequestGet(requestURL, requestParamString, signature string) ([]byte, error) {
	defer func() {
		recover()
	}()
	url := "https://" + requestURL + "?" + requestParamString + "&Signature=" + signature
	if CloudAPIDebug {
		log.Printf("url: %s", url)
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	res, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return res, err
}

// Request makes a call to the HTTP endpoint
func RequestPost(requestURL, requestParamString, signature string) ([]byte, error) {
	defer func() {
		recover()
	}()
	data := requestParamString + "&Signature=" + signature
	// data := url.Values{}

	resp, err := http.Post("https://"+requestURL, "application/x-www-form-urlencoded",
		strings.NewReader(data))
	if CloudAPIDebug {
		log.Printf("data: %s", data)
	}
	if err != nil {
		return nil, err
	}
	res, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return res, err
}

// GenerateNonce generates a unique key based on the current time
func GenerateNonce() string {
	rand.Seed(time.Now().UnixNano())
	time := rand.Intn(10000) + 10000
	return strconv.Itoa(time)
}

// CurrentTimeStamp returns the current time in Unix seconds as a string
func CurrentTimeStamp() string {
	ts := time.Now().Unix()
	s := strconv.FormatInt(ts, 10)
	return s
}

// Hmac256ToBase64 Base64 encodes the given parameters as a string
func Hmac1ToBase64(key string, str string, IsUrl bool) string {
	s := hmac.New(sha1.New, []byte(key))
	s.Write([]byte(str))
	return EncodingBase64(s.Sum(nil), IsUrl)
}

// Hmac256ToBase64 Base64 encodes the given parameters as a string
func Hmac256ToBase64(key string, str string, IsUrl bool) string {
	s := hmac.New(sha256.New, []byte(key))
	s.Write([]byte(str))
	return EncodingBase64(s.Sum(nil), IsUrl)
}

func EncodingBase64(b []byte, IsURL bool) string {
	if IsURL {
		return url.QueryEscape(base64.StdEncoding.EncodeToString(b))
	}
	return base64.StdEncoding.EncodeToString(b)
}

// IntToString converts an integer to a string
func IntToString(num int) string {
	return strconv.Itoa(num)
}

// Int64ToString converts an Int64 integer to a string
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}
