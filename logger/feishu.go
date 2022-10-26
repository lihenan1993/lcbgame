package logger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type PanicPayload struct {
	Happen       time.Time
	Stack        string
	RemoteSocket string
	LocalSocket  string
	Recover      string
	Protocol     int
	UID          int
	BuildTime    string
	CommitID     string
	LogLevel     string
	TestUser     bool
	ServerName   string
	Uuid         string
}

type PanicMsg struct {
	MsgType string   `json:"msg_type"`
	Content *Content `json:"content"`
}

type Content struct {
	Text string `json:"text"`
}

func Warning(payload *PanicPayload) {
	buf, _ := json.MarshalIndent(payload, "  ", "  ")
	go SendFeishu(buf)
}

func SendFeishu(text []byte) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("SendFeishu recover")
		}
	}()

	content := &Content{Text: string(text)}
	msg := &PanicMsg{
		MsgType: "text",
		Content: content,
	}
	jsonBuf, err := json.MarshalIndent(msg, "  ", "  ")
	if err != nil {
		return
	}
	resp, err := http.Post(feishu_url,
		"application/json",
		strings.NewReader(string(jsonBuf)))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
}
