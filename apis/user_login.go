package apis

import (
	"go.mongodb.org/mongo-driver/bson"
	"mania/apis/protocol"
	"mania/constant"
	"mania/control"
	"mania/logger"
	"mania/model"
	"mania/service"
	"mania/tcpx"
	"time"
)

func UserLogin(c *tcpx.Context) {
	openSocketTime := time.Time{}
	openSocket, ok := c.GetCtxPerConn("open_socket")
	if ok {
		ost, ok := openSocket.(time.Time)
		if ok {
			openSocketTime = ost
		}
	}

	rawLog, _ := c.GetCtxPerConn("logger")
	log := rawLog.(*logger.Logger)
	log.Append("", "key", "user_login")

	if !openSocketTime.IsZero() {
		log.Append("", "socket_delta", time.Now().Sub(openSocketTime).Milliseconds(), "socket_open", openSocketTime.Unix())
	}

	old := time.Now()
	var req protocol.ReqUserLogin
	var resp protocol.RespUserLogin

	defer func() {
		err := c.JSON(constant.USER_LOGIN, resp)
		if err != nil {
			log.SetLogLevel("error")
			log.Append("send resp failed", "err", err)
		}
		now := time.Now()
		delta := now.Sub(old).Milliseconds()
		log.Append("", "delta", delta)
		log.LogExceptInfo()
	}()
	_, err := c.BindWithMarshaller(&req, tcpx.JsonMarshaller{})
	if err != nil {
		log.SetLogLevel("error")
		log.Append("decode req failed", "err", err)
		resp.Code = constant.ERR_JSON_MARSHALLER_FAIL
		return
	}
	log.Append("", "req", req)

	if req.Duid == "" {
		log.SetLogLevel("error")
		log.Append("err req argument")
		resp.Code = constant.ERR_JSON_MARSHALLER_FAIL
		return
	}

	duid := ""
	platform := 0
	if req.Credential != "" {
		platform = 4
		duid = req.Credential
	} else if req.GoogleID != "" {
		platform = 3
		duid = req.GoogleID
	} else if req.FaceBookID != "" {
		platform = 2
		duid = req.FaceBookID
	} else if req.AppleID != "" {
		platform = 1
		duid = req.AppleID
	} else {
		platform = 0
		duid = req.Duid
	}

	uid, credential, usr, err, _ := service.Login(duid, platform, req.Duid)
	if err != nil {
		if err.Error() == "error credential" {
			log.SetLogLevel("warn")
			log.Append("user.Login failed", "err", err.Error())
			resp.Code = constant.ERR_TOKEN
			return
		}

		log.SetLogLevel("error")
		log.Append("user.Login failed", "err", err)
		resp.Code = constant.ERR_JSON_MARSHALLER_FAIL
		return
	}
	if uid == 0 || credential == "" || usr == nil {
		log.SetLogLevel("error")
		log.Append("login or create user failed", "err", err)
		resp.Code = constant.ERR_GET_TCP_USER_FAIL
		return
	}
	log.Append("", "uid", uid)

	// 检查用户数据是否在服务器内存中
	online, serverName, err := service.UserStatus(usr.UID)
	if err != nil {
		log.SetLogLevel("error")
		log.Append("get user online status failed", "err", err)
		resp.Code = constant.ERR_GET_TCP_USER_FAIL
		return
	}

	// 在本机内存中
	if online == constant.LOCAL_ONLINE {
		log.Append("Re login")
		// 获取用户数据
		oldUser, _ := model.Get(uid)
		if oldUser != nil {
			// 保存
			err = service.UpdateUserByObject(oldUser)
			if err != nil {
				log.SetLogLevel("error")
				log.Append("user.UpdateUserByObject", "err", err.Error())
			}

			// 如果是新的socket,关闭旧的socket
			oldTCPXContext := oldUser.Conn
			if oldTCPXContext != nil && oldTCPXContext.Conn != c.Conn {
				log.Append("oldTCPXContext.Conn != c.Conn")
				service.UserOut(oldTCPXContext, 3, true)
			}

			time.Sleep(time.Millisecond * 500)

			filter := bson.D{{Key: "uid", Value: uid}}
			usr, err = service.GetUser(filter)
			if err != nil {
				log.SetLogLevel("error")
				log.Append("user.GetUser error", "err", err.Error())
				resp.Code = constant.ERR_GET_TCP_USER_FAIL
				return
			}
		}
	} else if online == constant.OTHER_ONLINE {
		log.Append("another server", "serverName", serverName)
		service.SendForcedUserOffline(uid)

		start := time.Now()
		for {
			time.Sleep(time.Millisecond * 500)

			online, serverName, _ := service.UserStatus(usr.UID)
			if online != constant.OTHER_ONLINE {
				// kick成功
				log.Append("kick success", "serverName", serverName)
				break
			}

			if time.Now().Unix()-start.Unix() > 5 {
				log.SetLogLevel("error")
				log.Append("kick failed", "serverName", serverName)
				break
			}
		}

		filter := bson.D{{Key: "uid", Value: uid}}
		usr, err = service.GetUser(filter)
		if err != nil {
			log.SetLogLevel("error")
			log.Append("user.GetUser error", "err", err)
			resp.Code = constant.ERR_GET_TCP_USER_FAIL
			return
		}
	}

	_ = model.Set(usr)
	err = control.Store.UserOnline(usr.UID, constant.SERVER_NAME)
	if err != nil {
		log.SetLogLevel("error")
		log.Append("user online failed", "err", err)
		return
	}
	usr.Log(log)

	resp.Credential = credential
	resp.UserID = usr.UID

	usr.Package = req.Package
	usr.Conn = c
	c.SetCtxPerConn("uid", uid)

	log.Append("success")
	resp.Code = constant.SUCCESS
	logger.Debug("login debug", "uid", usr.UID, "req", req, "resp", resp)
}
