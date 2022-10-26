package protocol

type ReqHeartBeat struct {
}

type RespHeartBeat struct {
	ServerTime int64 `json:"server_time"`
}
