package control

var SrvConfig ServerConfig

type ServerConfig struct {
	Name       string `json:"name"`
	InnerIP    string `json:"ip"`
	InstanceID string `json:"instance_id"`
	Port       int    `json:"port"`
	BuildTime  string
	CommitID   string
	LogLevel   string
}
