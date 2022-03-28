package tencent

import (

	// "fmt"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
)

// Config contains the configuration for Builder
type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	// Fields from config file
	ClientToken           string                 `mapstructure:"ClientToken"`
	DataDisks             []CVMDataDisk          `mapstructure:"DataDisks"`
	EnhancedService       EnhancedService        `mapstructure:"EnhancedService"`
	ImageID               string                 `mapstructure:"ImageId"`
	InstanceChargePrepaid InstanceChargePrepaid  `mapstructure:"InstanceChargePrepaid"`
	InstanceChargeType    string                 `mapstructure:"InstanceChargeType"`
	InstanceCount         int                    `mapstructure:"InstanceCount"`
	InstanceName          string                 `mapstructure:"InstanceName"`
	InstanceType          string                 `mapstructure:"InstanceType"`
	InternetAccessible    CVMInternetAccessible  `mapstructure:"InternetAccessible"`
	LoginSettings         LoginSettings          `mapstructure:"LoginSettings"`
	Placement             CVMPlacement           `mapstructure:"Placement"`
	Region                string                 `mapstructure:"Region"`
	SecretID              string                 `mapstructure:"SecretId"`
	SecretKey             string                 `mapstructure:"SecretKey"`
	SSHKeyName            string                 `mapstructure:"KeyName"`
	SSHUserName           string                 `mapstructure:"ssh_username"`
	SystemDisk            CVMSystemDisk          `mapstructure:"SystemDisk"`
	tmiDescription        string                 `mapstructure:"tmi_description"`
	Version               string                 `mapstructure:"Version"`
	VirtualPrivateCloud   CVMVirtualPrivateCloud `mapstructure:"VirtualPrivateCloud"`

	Ctx interpolate.Context
}

// NewSimpleConfig parses the given raws parameter
func NewSimpleConfig(raws ...interface{}) (*Config, []string, error) {
	c := &Config{}
	warnings := []string{}

	if CloudAPIDebug {
		log.Printf("NewSimpleConfig raws1: %+v\n", raws)
	}

	err := config.Decode(c, &config.DecodeOpts{}, raws...)

	if err != nil {
		return nil, warnings, err
	}

	if c.PackerDebug || CloudAPIDebug {
		log.Printf("NewSimpleConfig raws: %+v\n", raws)
	}

	if c.Version == "" {
		c.Version = "2017-03-12"
	}

	return c, warnings, nil
}

// NewConfig parses the given raws parameter
func NewConfig(raws ...interface{}) (*Config, []string, error) {
	c := &Config{}
	warnings := []string{}

	if CloudAPIDebug {
		log.Printf("NewConfig raws1: %+v\n", raws)
	}

	err := config.Decode(c, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &c.Ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)

	if err != nil {
		return nil, warnings, err
	}

	if c.PackerDebug || CloudAPIDebug {
		log.Printf("NewConfig raws2: %+v\n", raws)
	}

	if c.Version == "" {
		c.Version = "2017-03-12"
	}

	if c.InstanceChargeType == "" {
		c.InstanceChargeType = "POSTPAID_BY_HOUR"
	}

	var errs *packer.MultiError

	if c.Placement.Zone == "" {
		errs = packer.MultiErrorAppend(errs, errors.New("Placement.Zone needs to be set"))
	}

	if c.Region == "" {
		errs = packer.MultiErrorAppend(errs, errors.New("Region needs to be set"))
	}

	if c.ImageID == "" {
		errs = packer.MultiErrorAppend(errs, errors.New("ImageId needs to be set"))
	}

	if c.SSHKeyName != "" {
		// Ensure it's a directory, and the file specified doesn't exist
		path := c.SSHKeyName
		dir := filepath.Dir(path)
		if !DirectoryExists(dir) {
			s := fmt.Sprintf("Directory specified: %s in KeyName doesn't exist", path)
			errs = packer.MultiErrorAppend(errs, errors.New(s))
		}
		file := filepath.Base(path) // returns the filename from the path
		if FileExists(file) {
			s := fmt.Sprintf("Filename specified in KeyName: %s already exists", file)
			errs = packer.MultiErrorAppend(errs, errors.New(s))
		}
	} else {
		errs = packer.MultiErrorAppend(errs, errors.New("KeyName needs to be set"))
	}

	if c.SecretID == "" {
		errs = packer.MultiErrorAppend(errs, errors.New("SecretID needs to be set"))
	}

	if c.SecretKey == "" {
		errs = packer.MultiErrorAppend(errs, errors.New("SecretKey needs to be set"))
	}

	if c.SSHUserName == "" {
		errs = packer.MultiErrorAppend(errs, errors.New("ssh_username needs to be set"))
	}

	if errs != nil && len(errs.Errors) > 0 {
		return nil, warnings, errs
	}

	return c, warnings, nil
}

func (c *Config) CreateVMmap() map[string]interface{} {
	result := make(map[string]interface{})
	result["ImageId"] = c.ImageID
	if c.Placement.Zone != "" {
		// result["Placement.Zone"] = c.Placement.Zone
		// if c.Placement.ProjectId != 0 {
		// 	result["Placement.ProjectId"] = Int64ToString(c.Placement.ProjectId)
		// }
		// for i := 0; i < len(c.Placement.HostIds); i++ {
		// 	LHostID := fmt.Sprintf("Placement.HostIds.%d", i)
		// 	result[LHostID] = c.Placement.HostIds[i]
		// }
		result["Placement"] = c.Placement
	}
	result["Region"] = c.Region
	result["SecretId"] = c.SecretID
	result["SecretKey"] = c.SecretKey
	result["KeyName"] = c.SSHKeyName
	result["ssh_username"] = c.SSHUserName
	if c.PackerDebug || CloudAPIDebug {
		log.Printf("%+v", result)
	}
	return result
}

// Keys generates a dictionary of keys and values given in the c Config.
func (c *Config) Keys() map[string]string {
	result := make(map[string]string)

	if c.ClientToken != "" {
		result["ClientToken"] = c.ClientToken
	}

	if len(c.DataDisks) > 0 {
		for i := 0; i < len(c.DataDisks); i++ {
			LDiskSize := fmt.Sprintf("DataDisks.%d.DiskSize", i)
			LDiskType := fmt.Sprintf("DataDisks.%d.DiskType", i)
			result[LDiskSize] = Int64ToString(c.DataDisks[i].DiskSize)
			result[LDiskType] = c.DataDisks[i].DiskType
		}
	}

	// Test for not empty
	if c.EnhancedService != (EnhancedService{}) {
		result["EnhancedService.MonitorService.Enabled"] = strconv.FormatBool(c.EnhancedService.MonitorService.Enabled)
		result["EnhancedService.SecurityService.Enabled"] = strconv.FormatBool(c.EnhancedService.SecurityService.Enabled)
	}

	if c.ImageID != "" {
		result["ImageId"] = c.ImageID
	}

	// Test for not empty
	if c.InstanceChargePrepaid != (InstanceChargePrepaid{}) {
		result["InstanceChargePrepaid.Period"] = IntToString(c.InstanceChargePrepaid.Period)
		result["InstanceChargePrepaid.RenewFlag"] = c.InstanceChargePrepaid.RenewFlag
	}

	// if c.InstanceChargeType != "" {
	// 	result["InstanceChargeType"] = c.InstanceChargeType
	// }

	if c.InstanceCount != 0 {
		result["InstanceCount"] = IntToString(c.InstanceCount)
	}

	if c.InstanceName != "" {
		result["InstanceName"] = c.InstanceName
	}

	if c.InstanceType != "" {
		result["InstanceType"] = c.InstanceType
	}

	if c.InternetAccessible != (CVMInternetAccessible{}) {
		result["InternetAccessible.InternetChargeType"] = c.InternetAccessible.InternetChargeType
		result["InternetAccessible.InternetMaxBandwidthOut"] = Int64ToString(c.InternetAccessible.InternetMaxBandwidthOut)
		result["InternetAccessible.PublicIpAssigned"] = strconv.FormatBool(c.InternetAccessible.PublicIpAssigned)
	}

	if c.LoginSettings.Password != "" {
		result["LoginSettings.Password"] = c.LoginSettings.Password
		result["LoginSettings.KeepImageLogin"] = c.LoginSettings.KeepImageLogin
		for i := 0; i < len(c.LoginSettings.KeyIds); i++ {
			LKeyID := fmt.Sprintf("LoginSettings.KeyIds.%d", i)
			result[LKeyID] = c.LoginSettings.KeyIds[i]
		}
	}

	if c.Placement.Zone != "" {
		result["Placement.Zone"] = c.Placement.Zone
		if c.Placement.ProjectId != 0 {
			result["Placement.ProjectId"] = Int64ToString(c.Placement.ProjectId)
		}
		for i := 0; i < len(c.Placement.HostIds); i++ {
			LHostID := fmt.Sprintf("Placement.HostIds.%d", i)
			result[LHostID] = c.Placement.HostIds[i]
		}
	}

	if c.Region != "" {
		result["Region"] = c.Region
	}

	// Test for not empty
	if c.SystemDisk != (CVMSystemDisk{}) {
		result["SystemDisk.DiskSize"] = Int64ToString(c.SystemDisk.DiskSize)
		result["SystemDisk.DiskType"] = c.SystemDisk.DiskType
	}

	if c.Version != "" {
		result["Version"] = c.Version
	}

	if c.VirtualPrivateCloud.VpcId != "" {
		result["VirtualPrivateCloud.VpcId"] = c.VirtualPrivateCloud.VpcId
		result["VirtualPrivateCloud.AsVpcGateway"] = strconv.FormatBool(c.VirtualPrivateCloud.AsVpcGateway)
		result["VirtualPrivateCloud.SubnetId"] = c.VirtualPrivateCloud.SubnetId
		for i := 0; i < len(c.VirtualPrivateCloud.PrivateIpAddresses); i++ {
			LIPName := fmt.Sprintf("VirtualPrivateCloud.PrivateIpAddresses.%d", i)
			result[LIPName] = c.VirtualPrivateCloud.PrivateIpAddresses[i]
		}
	}

	return result
}
