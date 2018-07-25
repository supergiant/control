package aws

type InstanceConfig struct {
	Name           string
	Type           string
	ClusterName    string
	ClusterRole    string
	ID             string
	Region         string
	ImageID        string
	KeyName        string
	HasPublicAddr  bool
	SecurityGroups []*string
	SubnetID       string
	IAMRole        string
	VolumeType     string
	VolumeSize     int64
	Tags           map[string]string
	UsedData       string
}

type EC2TypeInfo struct {
	SKU           string             `json:"sku"`
	ProductFamily string             `json:"productFamily"`
	Attributes    *EC2TypeAttributes `json:"attributes"`
}

type EC2TypeAttributes struct {
	ServiceCode       string `json:"servicecode"`
	ServiceName       string `json:"servicename"`
	InstanceFamily    string `json:"instanceFamily"`
	InstanceType      string `json:"instanceType"`
	OperatingSystem   string `json:"operatingSystem"`
	Location          string `json:"location"`
	VCPU              string `json:"vcpu"`
	Memory            string `json:"memory"`
	ECU               string `json:"ecu"`
	PhysicalProcessor string `json:"physicalProcessor"`
	Storage           string `json:"storage"`
	CapacityStatus    string `json:"capacitystatus"`
}