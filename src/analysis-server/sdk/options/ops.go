package options

type CreatePoolOptions struct {
	Name      string
	Type      int
	MediaType string
	Rgc       int
	ObjSize   uint64
	Status    int
}

type CreateRackOptions struct {
	Tag string
}

type CreateHostOptions struct {
	Name        string
	MgmtIp      string
	DataIp      string
	ClientIp    string
	AdminStatus int
	Rack        string
	PoolId      string
}

type CreateDiskOptions struct {
	DeviceId       string
	HostId         string
	ManageAddr     string
	StorageAddr    string
	ClientAddr     string
	TraceAddr      string
	VolumeTypeName string
	Capacity       uint64
	AdminStatus    int
}

type CreateProxyOptions struct {
	Addr string
}

type DeleteOptions struct {
	Id string
}

type ReschedulerOptions struct {
	PoolId string
}

type TransferLeaderOptions struct {
	RgId      string
	ReplicaId string
}

type MoveReplicaOptions struct {
	ReplicaId    string
	TargetDiskId string
}

type DeleteReplicaOptions struct {
	ReplicaId string
	Force     bool
}
type AddReplicaOptions struct {
	RgID   string
	DiskID string
}

type ZbsProxyOptions struct {
	Addr string
}
