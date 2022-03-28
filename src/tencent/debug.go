package tencent

import "os"

// DebugEnabled returns true if PACKER_LOG is defined as an environment variable and that it is not empty or equals "0"
func DebugEnabled() bool {
	debugEnabled := os.Getenv("PACKER_LOG")
	return debugEnabled != "" && debugEnabled != "0"
}
