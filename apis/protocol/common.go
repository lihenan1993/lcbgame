package protocol

type CommonError struct {
	Code string `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}
