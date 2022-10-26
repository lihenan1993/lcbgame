package apis

import (
	"mania/apis/protocol"
	"mania/tcpx"
	"time"
)

func HeartBeat(c *tcpx.Context) {
	_, ok := c.GetCtxPerConn("uid")
	if !ok {
		return
	}

	c.RecvHeartBeat()

	resp := &protocol.RespHeartBeat{ServerTime: time.Now().Unix()}
	_ = c.JSON(200, resp)
}
