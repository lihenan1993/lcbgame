package service

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"mania/apis/protocol"
	"mania/constant"
	"mania/control"
	"mania/logger"
	"mania/model"
	"mania/tcpx"
)

func CreateUser(t model.Token) (*model.User, *model.Token) {
	var u model.User
	uid, err := control.Store.CreateUserID()
	if err != nil {
		logger.Error("create user id failed", "err", err)
		return nil, nil
	}

	u.UID = uid
	t.UID = uid
	token := CreateUserToken(t)
	if token == nil {
		return nil, nil
	}

	err = control.Store.CreateUser(&u)
	if err != nil {
		logger.Error("create user failed", "err", err)
		return nil, nil
	}
	return &u, token
}
func GetUser(filter bson.D) (u *model.User, err error) {
	raw, err := control.Store.GetUser(filter)
	if err != nil {
		logger.Error("get user failed", "err", err.Error())
		return nil, err
	}

	err = raw.Decode(&u)
	if err != nil {
		logger.Error("get user decode failed", "err", err.Error())
		return nil, err
	}
	return u, err
}
func UpdateUser(uid int) (err error) {
	usr, err := model.Get(uid)
	if err != nil {
		return fmt.Errorf("user is nil")
	}

	filter := bson.D{{Key: "uid", Value: uid}}
	update := bson.D{
		{Key: "$set", Value: usr},
	}
	return control.Store.UpdateUser(filter, update)
}

func UpdateUserByObject(usr *model.User) (err error) {
	filter := bson.D{{Key: "uid", Value: usr.UID}}
	update := bson.D{
		{Key: "$set", Value: usr},
	}
	return control.Store.UpdateUser(filter, update)
}

func CreateUserToken(token model.Token) *model.Token {
	// 生成证书
	token.GenerateCertificate()

	err := control.Store.CreateUserToken(&token)
	if err != nil {
		logger.Error("create user token failed", "err", err)
		return nil
	}
	return &token
}

func GetUserToken(filter bson.D) (t *model.Token, err error) {
	raw, err := control.Store.GetUserToken(filter)
	if err != nil {
		return
	}
	var token model.Token
	err = raw.Decode(&token)
	if err != nil {
		logger.Error("get user token failed", "err", err)
		return
	}
	return &token, nil
}
func UpdateUserToken(token *model.Token) (err error) {
	filter := bson.D{{Key: "uid", Value: token.UID}}
	update := bson.D{
		{Key: "$set", Value: token},
	}
	return control.Store.UpdateUserToken(filter, update)
}

func UserStatus(uid int) (int, string, error) {
	if _, err := model.Get(uid); err == nil {
		return constant.LOCAL_ONLINE, "", nil
	}

	online, serverName, err := control.Store.IsUserOnline(uid)
	if err != nil {
		return 0, "", err
	}
	if online {
		if serverName == constant.SERVER_NAME {
			return constant.LOCAL_ONLINE, "", nil
		}

		return constant.OTHER_ONLINE, serverName, nil
	} else {
		return constant.OFFLINE, serverName, nil
	}
}

func SendForcedUserOffline(uid int) {
	//rmq := control.RedisStreamMQClient{
	//	ConnPool:     control.Store.RedisConn,
	//	StreamKey:    "stream::kick::user",
	//	StreamLen:    100,
	//	GroupName:    "",
	//	ConsumerName: "",
	//}
	//
	//strMsgId, err := rmq.PutMsg("", strconv.Itoa(uid), "1")
	//if err != nil {
	//	return
	//}

	//logger.Info("SendForcedUserOffline", "key", "kick", "msgId", strMsgId, "uid", uid)
}

func Login(loginToken string, platform int, deviceID string) (uid int, credential string, u *model.User, err error, isFirstLogin bool) {
	var filter bson.D
	switch platform {
	case 0:
		filter = bson.D{{Key: "duid", Value: loginToken}}
	case 1:
		filter = bson.D{{Key: "appleid", Value: loginToken}}
	case 2:
		filter = bson.D{{Key: "facebookid", Value: loginToken}}
	case 3:
		filter = bson.D{{Key: "googleid", Value: loginToken}}
	case 4:
		filter = bson.D{{Key: "credential", Value: loginToken}}
	}
	token, err := GetUserToken(filter)
	if err != nil && err != mongo.ErrNoDocuments {
		return
	}
	var user *model.User
	if platform == 3 && (token == nil || err == mongo.ErrNoDocuments) {
		err = fmt.Errorf("error credential")
		return
	}

	// 没找到对应token 创建帐号
	if token == nil || err == mongo.ErrNoDocuments {
		// 用duid查找
		filter = bson.D{{Key: "duid", Value: deviceID}}
		getUserToken, err := GetUserToken(filter)

		if err != nil && err != mongo.ErrNoDocuments {
			return uid, credential, u, err, false
		}

		// duid找到了 更新帐号
		if getUserToken != nil {
			if platform == 1 {
				if getUserToken.AppleID == "" {
					getUserToken.AppleID = loginToken
					token = getUserToken
					filter = bson.D{{Key: "uid", Value: token.UID}}
					user, err = GetUser(filter)
					if err != nil {
						return uid, credential, u, err, false
					}
					err := UpdateUserToken(token)
					if err != nil {
						return uid, credential, u, err, false
					}
				}

			} else if platform == 2 {
				if getUserToken.FacebookID == "" {
					getUserToken.FacebookID = loginToken
					token = getUserToken
					filter = bson.D{{Key: "uid", Value: token.UID}}
					user, err = GetUser(filter)
					if err != nil {
						return uid, credential, u, err, false
					}
					err := UpdateUserToken(token)
					if err != nil {
						return uid, credential, u, err, false
					}
				}
			} else if platform == 3 {
				if getUserToken.GoogleID == "" {
					getUserToken.GoogleID = loginToken
					token = getUserToken
					filter = bson.D{{Key: "uid", Value: token.UID}}
					user, err = GetUser(filter)
					if err != nil {
						return uid, credential, u, err, false
					}
					err := UpdateUserToken(token)
					if err != nil {
						return uid, credential, u, err, false
					}
				}
			}
		}

		if user == nil {
			// 新帐号
			newToken := model.Token{
				Duid:       deviceID,
				UID:        0,
				FacebookID: "",
				AppleID:    "",
				Credential: "",
				GoogleID:   "",
			}

			switch platform {
			case 1:
				newToken.AppleID = loginToken
			case 2:
				newToken.FacebookID = loginToken
			case 3:
				newToken.GoogleID = loginToken
			}

			user, token = CreateUser(newToken)
		}
	} else {
		filter = bson.D{{Key: "uid", Value: token.UID}}
		user, err = GetUser(filter)
		if err != nil {
			return uid, credential, u, err, false
		}
	}

	if user != nil && token != nil {
		return user.UID, token.Credential, user, nil, isFirstLogin
	}

	return 0, "", nil, fmt.Errorf("user or token is nil"), false
}

// UserOut
// sign = 0 ,不返回信息
// sign = 1 ,重连
// sign = 2 ,重新登陆
// sign = 3 ,踢人下线
func UserOut(c *tcpx.Context, sign int, closeConn bool) {
	var resp protocol.RespUserOut

	uid := 0
	rawUid, ok := c.GetCtxPerConn("uid")
	if !ok {
		logger.Warn("uid is wrong", "clientIP", c.Conn.RemoteAddr().String())

		return
	}
	uid = rawUid.(int)

	if sign == 3 {
		resp.Sign = sign
		err := c.JSON(constant.USER_OUT, &resp)
		if err != nil {
			logger.Warn("USER_OUT", "uid", uid, "err", err.Error(), "sign", sign)
		}
	}

	if closeConn {
		_ = c.CloseConn()
		return
	} else {
		defer func() {
			if sign == 1 || sign == 2 {
				resp.Sign = sign
				err := c.JSON(constant.USER_OUT, &resp)
				if err != nil {
					logger.Warn("USER_OUT", "uid", uid, "err", err.Error(), "sign", sign)
				}
			}
		}()
	}

	resp.Code = constant.SUCCESS
}
