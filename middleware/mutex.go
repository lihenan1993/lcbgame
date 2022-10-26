package middleware

import (
	"fmt"
	"mania/tcpx"
	"sync"
)

func MutexConnection(c *tcpx.Context) {
	mut, ok := c.GetCtxPerConn("mutex")
	if !ok {
		fmt.Println("mutex is nil")
	}
	mutex := mut.(*sync.Mutex)
	mutex.Lock()
	c.Next()
	defer mutex.Unlock()
}
