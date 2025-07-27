package dash

import "time"

type AdminData struct {
	Username          string              `json:"username"`
	TotalVolumes      int                 `json:"total_volumes"`
	TotalFiles        int64               `json:"total_files"`
	TotalSize         int64               `json:"total_size"`
	VolumeSizeLimitMB uint64              `json:"volume_size_limit_mb"`
	MasterNodes       []MasterNode        `json:"master_nodes"`
	VolumeServers     []VolumeServer      `json:"volume_servers"`
	FilerNodes        []FilerNode         `json:"filer_nodes"`
	MessageBrokers    []MessageBrokerNode `json:"message_brokers"`
	DataCenters       []DataCenter        `json:"datacenters"`
	LastUpdated       time.Time           `json:"last_updated"`
}

type MasterNode struct {
	Address  string `json:"address"`
	IsLeader bool   `json:"is_leader"`
}
type VolumeServer struct {
	ID            string    `json:"id"`
	Address       string    `json:"address"`
	DataCenter    string    `json:"datacenter"`
	Rack          string    `json:"rack"`
	PublicURL     string    `json:"public_url"`
	Volumes       int       `json:"volumes"`
	MaxVolumes    int       `json:"max_volumes"`
	DiskUsage     int64     `json:"disk_usage"`
	DiskCapacity  int64     `json:"disk_capacity"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
}
type MessageBrokerNode struct {
	Address     string    `json:"address"`
	DataCenter  string    `json:"datacenter"`
	Rack        string    `json:"rack"`
	LastUpdated time.Time `json:"last_updated"`
}
type DataCenter struct {
	ID    string `json:"id"`
	Racks []Rack `json:"racks"`
}
type FilerNode struct {
	Address     string    `json:"address"`
	DataCenter  string    `json:"datacenter"`
	Rack        string    `json:"rack"`
	LastUpdated time.Time `json:"last_updated"`
}
type Rack struct {
	ID    string         `json:"id"`
	Nodes []VolumeServer `json:"nodes"`
}

// Object Store Users management structures
type ObjectStoreUser struct {
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	AccessKey   string   `json:"access_key"`
	SecretKey   string   `json:"secret_key"`
	Permissions []string `json:"permissions"`
}
type ObjectStoreUsersData struct {
	Username    string            `json:"username"`
	Users       []ObjectStoreUser `json:"users"`
	TotalUsers  int               `json:"total_users"`
	LastUpdated time.Time         `json:"last_updated"`
}

// User management request structures
type CreateUserRequest struct {
	Username    string   `json:"username" binding:"required"`
	Email       string   `json:"email"`
	Actions     []string `json:"actions"`
	GenerateKey bool     `json:"generate_key"`
}

type UpdateUserRequest struct {
	Email   string   `json:"email"`
	Actions []string `json:"actions"`
}

type UpdateUserPoliciesRequest struct {
	Actions []string `json:"actions" binding:"required"`
}
