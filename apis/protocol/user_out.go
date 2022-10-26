package protocol

type RespUserOut struct {
	CommonError
	Sign int `json:"sign"`
}
