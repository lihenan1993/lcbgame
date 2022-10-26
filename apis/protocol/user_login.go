package protocol

type ReqUserLogin struct {
	Package    string `json:"package"`
	Version    string `json:"version"`
	Duid       string `json:"duid"`
	Credential string `json:"credential"`
	FaceBookID string `json:"facebook_id"`
	AppleID    string `json:"apple_id"`
	GoogleID   string `json:"google_id"`
	FCMToken   string `json:"token"`
	PID        int    `json:"pid"`
	P          int    `json:"p"`
}

type RespUserLogin struct {
	Credential string `json:"credential"`
	Coins      int64  `json:"coins"`
	UserID     int    `json:"user_id"`
	CommonError
}
