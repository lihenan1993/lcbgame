package middleware

import (
	"mania/logger"
	"mania/tcpx"
)

func Logger(c *tcpx.Context) {
	log := logger.NewLog()
	c.SetCtxPerConn("logger", &log)
	c.Next()
}
