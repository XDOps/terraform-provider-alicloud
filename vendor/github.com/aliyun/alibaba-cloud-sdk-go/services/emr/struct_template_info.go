package emr

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// TemplateInfo is a nested struct in emr response
type TemplateInfo struct {
	UserDefinedEmrEcsRole  string                                       `json:"UserDefinedEmrEcsRole" xml:"UserDefinedEmrEcsRole"`
	InitCustomHiveMetaDb   bool                                         `json:"InitCustomHiveMetaDb" xml:"InitCustomHiveMetaDb"`
	DataDiskKMSKeyId       string                                       `json:"DataDiskKMSKeyId" xml:"DataDiskKMSKeyId"`
	TemplateName           string                                       `json:"TemplateName" xml:"TemplateName"`
	SecurityGroupId        string                                       `json:"SecurityGroupId" xml:"SecurityGroupId"`
	MachineType            string                                       `json:"MachineType" xml:"MachineType"`
	Configurations         string                                       `json:"Configurations" xml:"Configurations"`
	MetaStoreConf          string                                       `json:"MetaStoreConf" xml:"MetaStoreConf"`
	IsOpenPublicIp         bool                                         `json:"IsOpenPublicIp" xml:"IsOpenPublicIp"`
	CreateSource           string                                       `json:"CreateSource" xml:"CreateSource"`
	UseCustomHiveMetaDb    bool                                         `json:"UseCustomHiveMetaDb" xml:"UseCustomHiveMetaDb"`
	EasEnable              bool                                         `json:"EasEnable" xml:"EasEnable"`
	IoOptimized            bool                                         `json:"IoOptimized" xml:"IoOptimized"`
	UserId                 string                                       `json:"UserId" xml:"UserId"`
	Id                     string                                       `json:"Id" xml:"Id"`
	EmrVer                 string                                       `json:"EmrVer" xml:"EmrVer"`
	VpcId                  string                                       `json:"VpcId" xml:"VpcId"`
	SecurityGroupName      string                                       `json:"SecurityGroupName" xml:"SecurityGroupName"`
	AllowNotebook          bool                                         `json:"AllowNotebook" xml:"AllowNotebook"`
	DataDiskEncrypted      bool                                         `json:"DataDiskEncrypted" xml:"DataDiskEncrypted"`
	NetType                string                                       `json:"NetType" xml:"NetType"`
	ClusterType            string                                       `json:"ClusterType" xml:"ClusterType"`
	MasterNodeTotal        int                                          `json:"MasterNodeTotal" xml:"MasterNodeTotal"`
	VSwitchId              string                                       `json:"VSwitchId" xml:"VSwitchId"`
	DepositType            string                                       `json:"DepositType" xml:"DepositType"`
	UseLocalMetaDb         bool                                         `json:"UseLocalMetaDb" xml:"UseLocalMetaDb"`
	GmtCreate              int64                                        `json:"GmtCreate" xml:"GmtCreate"`
	ZoneId                 string                                       `json:"ZoneId" xml:"ZoneId"`
	SshEnable              bool                                         `json:"SshEnable" xml:"SshEnable"`
	KeyPairName            string                                       `json:"KeyPairName" xml:"KeyPairName"`
	InstanceGeneration     string                                       `json:"InstanceGeneration" xml:"InstanceGeneration"`
	GmtModified            int64                                        `json:"GmtModified" xml:"GmtModified"`
	LogPath                string                                       `json:"LogPath" xml:"LogPath"`
	MetaStoreType          string                                       `json:"MetaStoreType" xml:"MetaStoreType"`
	HighAvailabilityEnable bool                                         `json:"HighAvailabilityEnable" xml:"HighAvailabilityEnable"`
	LogEnable              bool                                         `json:"LogEnable" xml:"LogEnable"`
	MasterPwd              string                                       `json:"MasterPwd" xml:"MasterPwd"`
	SoftwareInfoList       SoftwareInfoListInDescribeClusterTemplate    `json:"SoftwareInfoList" xml:"SoftwareInfoList"`
	Tags                   TagsInDescribeClusterTemplate                `json:"Tags" xml:"Tags"`
	HostGroupList          HostGroupListInDescribeClusterTemplate       `json:"HostGroupList" xml:"HostGroupList"`
	BootstrapActionList    BootstrapActionListInDescribeClusterTemplate `json:"BootstrapActionList" xml:"BootstrapActionList"`
	ConfigList             ConfigListInDescribeClusterTemplate          `json:"ConfigList" xml:"ConfigList"`
}