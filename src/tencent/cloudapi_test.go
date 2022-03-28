package tencent

import (
	"math"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestInt64ToString(t *testing.T) {
	type args struct {
		num int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// args below specifies the args type, not the args field
		{"Convert 1 to string", args{1}, "1"},
		{"Convert 2 to string", args{2}, "2"},
		{"Convert 2 to string", args{math.MinInt64}, "-9223372036854775808"},
		{"Convert 2 to string", args{math.MaxInt64}, "9223372036854775807"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int64ToString(tt.args.num); got != tt.want {
				t.Errorf("Int64ToString() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestIntToString(t *testing.T) {
	type args struct {
		num int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Convert 1 to string", args{1}, "1"},
		{"Convert 2 to string", args{2}, "2"},
		{"Convert 2 to string", args{math.MinInt32}, "-2147483648"},
		{"Convert 2 to string", args{math.MaxInt32}, "2147483647"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntToString(tt.args.num); got != tt.want {
				t.Errorf("IntToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHmac1ToBase64(t *testing.T) {
	secretKey1 := "Gu5t9xGARNpq86cd98joQYCN3Cozk1qA"
	secretKey2 := "Fu5t9xGARNpq86cd98joQYCN3Cozk1qX"
	srcstr1 := "GETcvm.api.qcloud.com/v2/index.php?Action=DescribeInstances&Nonce=11886&Region=gz&SecretId=AKIDz8krbsJ5yKBZQpn74WFkmLPx3gnPhESA&Timestamp=1465185768&instanceIds.0=ins-09dx96dg&limit=20&offset=0"
	srcstr2 := "GETcvm.tencentcloudapi.com/?Action=DescribeInstances&Nonce=11886&Region=gz&SecretId=AKIDz8krbsJ5yKBZQpn74WFkmLPx3gnPhESA&Timestamp=1465185768&instanceIds.0=ins-09dx96dg&limit=20&offset=0"
	type args struct {
		key   string
		str   string
		IsUrl bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// See https://intl.cloud.tencent.com/document/product/362/4208?!preview=true&lang=en#2.4.-generating-signature-string
		{"Test Hmac1ToBase64 case 1", args{secretKey1, srcstr1, false}, "NSI3UqqD99b/UJb4tbG/xZpRW64="},
		{"Test Hmac1ToBase64 case 2", args{secretKey1, srcstr1, true}, "NSI3UqqD99b%2FUJb4tbG%2FxZpRW64%3D"},
		{"Test Hmac1ToBase64 case 3", args{secretKey2, srcstr2, false}, "ZLvslzTWZTUyzVPopuMw3fKMQkg="},
		{"Test Hmac1ToBase64 case 4", args{secretKey2, srcstr2, true}, "ZLvslzTWZTUyzVPopuMw3fKMQkg%3D"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hmac1ToBase64(tt.args.key, tt.args.str, tt.args.IsUrl); got != tt.want {
				t.Errorf("Hmac1ToBase64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSignatureGet(t *testing.T) {
	url := "cvm.api.qcloud.com/"
	requestParams := "Action=DescribeInstances&Nonce=11886&Region=gz&SecretId=AKIDz8krbsJ5yKBZQpn74WFkmLPx3gnPhESA&Timestamp=1465185768&instanceIds.0=ins-09dx96dg&limit=20&offset=0"
	secretKey := "Gu5t9xGARNpq86cd98joQYCN3Cozk1qA"
	type args struct {
		requestURL         string
		requestParamString string
		secretKey          string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"SignatureGet test 1", args{url, requestParams, secretKey}, "2wmvFvB6R7CAVEzYcjO8BKTsvj4%3D"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SignatureGet(tt.args.requestURL, tt.args.requestParamString, tt.args.secretKey); got != tt.want {
				t.Errorf("SignatureGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodingBase64(t *testing.T) {
	type args struct {
		b     []byte
		IsURL bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"EncodeBase64 test 1", args{[]byte("chuacw rocks!"), false}, "Y2h1YWN3IHJvY2tzIQ=="},
		{"EncodeBase64 test 2", args{[]byte("chuacw rocks!"), true}, "Y2h1YWN3IHJvY2tzIQ%3D%3D"},
		{"EncodeBase64 test 3", args{[]byte("chuacw isn't cool!"), false}, "Y2h1YWN3IGlzbid0IGNvb2wh"},
		{"EncodeBase64 test 4", args{[]byte("chuacw isn't cool!"), true}, "Y2h1YWN3IGlzbid0IGNvb2wh"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodingBase64(tt.args.b, tt.args.IsURL); got != tt.want {
				t.Errorf("EncodingBase64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHmac256ToBase64(t *testing.T) {
	type args struct {
		key   string
		str   string
		IsUrl bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Hmac256ToBase64 test 1", args{"chuacw", "chuacw rocks!", false}, "WNIlcsw2mci4IZ+B5CJS6vRNZ7WTRnQjM003R/bOd9A="},
		{"Hmac256ToBase64 test 2", args{"chuacw", "chuacw rocks!", true}, "WNIlcsw2mci4IZ%2BB5CJS6vRNZ7WTRnQjM003R%2FbOd9A%3D"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hmac256ToBase64(tt.args.key, tt.args.str, tt.args.IsUrl); got != tt.want {
				t.Errorf("Hmac256ToBase64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurrentTimeStamp(t *testing.T) {
	ts1 := time.Now().Unix()
	s1 := strconv.FormatInt(ts1, 10)

	result := CurrentTimeStamp()

	ts2 := time.Now().Unix()
	s2 := strconv.FormatInt(ts2, 10)

	// the CurrentTimeStamp should be >= ts1 and <= ts2
	if !(result >= s1 && result <= s2) {
		t.Errorf("CurrentTimeStamp() = %s, expecting it to be >= %s and <= %s", result, s1, s2)
	}

}

func TestGenerateNonce(t *testing.T) {
	nonce1 := GenerateNonce()
	nonce2 := GenerateNonce()
	if nonce1 == nonce2 {
		t.Errorf("GenerateNonce %s == %s", nonce1, nonce2)
	}
}

// This tests SignatureString, since it is a wrapper to SignatureStringNonceTimestamp
// SignatureStringNonceTimestamp calls GenerateNonce() and CurrentTimeStamp() then passes
// these values to SignatureStringNonceTimestamp
func TestSignatureStringNonceTimestamp(t *testing.T) {
	action1 := "act1"
	action2 := "act2"
	nonce1 := "167890"
	nonce2 := "223569"
	ts1 := "1525252824"
	ts2 := "1525252900"
	secretId1 := "AKIDz8krbsJ5yKBZQpn74WFkmLPx3gnPhESA"
	secretId2 := "AKIDz8kkksJ5yKBZQpn74WFkmLPx3gnPhEXY"
	epK1 := "str1"
	epV1 := "value1"
	epK2 := "str2"
	epV2 := "value2"
	extraParams1 := map[string]string{
		epK1: epV1,
		epK2: epV2,
	}

	expected1 := strings.Join([]string{"Action=" + action1, "Nonce=" + nonce1, "SecretId=" + secretId1, "Timestamp=" + ts1,
		epK1 + "=" + epV1, epK2 + "=" + epV2}, "&")
	expected2 := strings.Join([]string{"Action=" + action2, "Nonce=" + nonce2, "SecretId=" + secretId2, "Timestamp=" + ts2,
		epK1 + "=" + epV1, epK2 + "=" + epV2}, "&")

	type args struct {
		action      string
		nonce       string
		timestamp   string
		secretId    string
		extraParams map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Test SignatureStringNonceTimestamp test 1", args{action1, nonce1, ts1, secretId1, extraParams1}, expected1},
		{"Test SignatureStringNonceTimestamp test 2", args{action2, nonce2, ts2, secretId2, extraParams1}, expected2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SignatureStringNonceTimestamp(tt.args.action, tt.args.nonce, tt.args.timestamp, tt.args.secretId, tt.args.extraParams); got != tt.want {
				t.Errorf("SignatureStringNonceTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}
