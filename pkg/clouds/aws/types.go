package aws

type InstanceConfig struct {
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
}
