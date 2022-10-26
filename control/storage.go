package control

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"mania/constant"
	"mania/logger"
	"time"
)

type Storage struct {
	MongoDB   *mongo.Client
	RedisConn *redis.Pool
}

var Store Storage
var expireTime = 180

const Online = 1
const Offline = 0

type serverInfo struct {
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	ServerKey string `json:"serverKey"`
	Node      string `json:"node"`
	Time      string `json:"time"`
}

func (s *Storage) Connect(configSrvName string) error {
	defer func(start time.Time) {
		fmt.Println("Connect:", time.Now().Sub(start).Milliseconds())
	}(time.Now())

	mongoPool := GetMongodbSource()
	mongodb, err := ConnectMongoDB(mongoPool, configSrvName)
	if err != nil {
		return err
	}

	redisPool := GetRedisSource()
	var cache *redis.Pool
	if len(redisPool.SentinelPath) > 0 {
		//cache, err = InitRedisSentinelConnPool(redisPool)
	} else {
		cache, err = ConnectRedis(redisPool)
	}
	if err != nil {
		return err
	}

	s.MongoDB = mongodb
	s.RedisConn = cache

	return nil
}

func (s *Storage) RegisterTimedTasks(d time.Duration, node string) {
	defer func() {
		if x := recover(); x != nil {
			logger.Error("RegisterTimedTasks", "x", x)
		}
	}()

	t := time.NewTimer(d)

	for {
		select {
		case <-t.C:
			//conn := s.RedisConn.Get()
			//s.setBlackList(conn)
			//s.setWhiteList(conn)
			s.concurrentUser(node)
			s.setUserOnlineTTL(node)
			s.setServerRunningTTL(node)
			//conn.Close()
			t.Reset(d)
		}
	}
}

func (s *Storage) CloseDB() {
	s.DelServerInfo(constant.SERVER_NAME)
	_ = s.DelUserOnline(constant.SERVER_NAME)

	err := s.RedisConn.Close()
	if err != nil {
		logger.Error("Redis关闭异常", "err", err.Error())
	}
	err = s.MongoDB.Disconnect(context.TODO())
	if err != nil {
		logger.Error("MongoDB关闭异常", "err", err.Error())
	}

	logger.Error("release DB success")
}

func (s *Storage) setServerRunningTTL(node string) {
	conn := s.RedisConn.Get()
	defer conn.Close()

	key := fmt.Sprintf("SERVERRUNNING::%s", node)
	_, err := conn.Do("SET", key, 1, "EX", expireTime)
	if err != nil {
		logger.Error("setServerRunningTTL", "key", key, "err", err.Error())
		return
	}
}

func (s *Storage) delServerRunningTTL(node string) {
	conn := s.RedisConn.Get()
	defer conn.Close()

	key := fmt.Sprintf("SERVERRUNNING::%s", node)
	_, err := conn.Do("DEL", key)
	if err != nil {
		logger.Error("delServerRunningTTL", "key", key, "err", err.Error())
		return
	}
}

func (s *Storage) RegisterServerInfo(node string, configSrvName string, port int, ip string) error {
	conn := s.RedisConn.Get()
	defer conn.Close()

	info := serverInfo{
		IP:        ip,
		Port:      port,
		ServerKey: configSrvName,
		Node:      node,
		Time:      time.Now().String(),
	}

	i, err := json.Marshal(info)
	if err != nil {
		return err
	}

	_, err = conn.Do("HSET", "SERVERRUNNING", node, string(i))
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DelServerInfo(node string) {
	conn := s.RedisConn.Get()
	defer conn.Close()

	_, err := conn.Do("HDEL", "SERVERRUNNING", node)
	if err != nil {
		logger.Error("DelServerInfo", "err", err.Error())
		return
	}

	s.delServerRunningTTL(node)

	return
}

func (s *Storage) DelUserOnline(node string) (err error) {
	conn := s.RedisConn.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s_USERONLINE_BIT", node)
	_, err = conn.Do("DEL", key)
	if err != nil {
		logger.Error("DelUserOnline", "err", err.Error())
		return
	}

	return
}

func (s *Storage) setUserOnlineTTL(node string) {
	conn := s.RedisConn.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s_USERONLINE_BIT", node)
	_, err := conn.Do("EXPIRE", key, expireTime)
	if err != nil {
		logger.Error("setUserOnlineTTL", "key", key, "err", err.Error())
		return
	}
}

func (s *Storage) UserOnline(uid int, node string) error {
	conn := s.RedisConn.Get()
	defer conn.Close()
	key := fmt.Sprintf("%s_USERONLINE_BIT", node)
	_, err := conn.Do("SETBIT", key, uid, Online)
	if err != nil {
		return err
	}
	logger.Info("", "key", "UserOnline", "uid", uid, "node", node, "redis_key", key)
	return nil
}
func (s *Storage) UserOffline(uid int, node string) error {
	conn := s.RedisConn.Get()
	defer conn.Close()
	key := fmt.Sprintf("%s_USERONLINE_BIT", node)
	_, err := conn.Do("SETBIT", key, uid, Offline)
	if err != nil {
		return err
	}
	logger.Info("", "key", "UserOffline", "uid", uid, "node", node, "redis_key", key)
	return nil
}
func (s *Storage) IsUserOnline(uid int) (bool, string, error) {
	conn := s.RedisConn.Get()
	defer conn.Close()
	nodes, err := redis.StringMap(conn.Do("HGETALL", "SERVERRUNNING"))
	if err != nil {
		return false, "", err
	}
	for k := range nodes {
		key := fmt.Sprintf("%s_USERONLINE_BIT", k)
		if ok, _ := redis.Bool(conn.Do("GETBIT", key, uid)); ok {
			return true, k, nil
		}
	}
	return false, "", nil
}

func (s *Storage) concurrentUser(node string) (num int) {
	conn := s.RedisConn.Get()
	defer conn.Close()
	key := fmt.Sprintf("%s_USERONLINE_BIT", node)
	_, err := redis.Int(conn.Do("BITCOUNT", key))
	if err != nil {
		logger.Error("concurrentUser", "err", err.Error())
		return
	}
	return
}

func (s *Storage) CreateRankID(key string) (int, error) {
	conn := s.RedisConn.Get()
	defer conn.Close()
	gid, err := redis.Int(conn.Do("INCR", key))
	if err != nil {
		return 0, err
	}
	return gid, nil
}
func (s *Storage) CreateUserID() (int, error) {
	conn := s.RedisConn.Get()
	defer conn.Close()
	uid, err := redis.Int(conn.Do("INCR", "USERID"))
	if err != nil {
		return 0, err
	}
	return uid, nil
}
func (s *Storage) CreateBillID() (int, error) {
	conn := s.RedisConn.Get()
	defer conn.Close()
	uid, err := redis.Int(conn.Do("INCR", "BILLID"))
	if err != nil {
		return 0, err
	}
	return uid, nil
}
func (s *Storage) CreateInboxID() (int, error) {
	conn := s.RedisConn.Get()
	defer conn.Close()
	uid, err := redis.Int(conn.Do("DECR", "INBOXID"))
	if err != nil {
		return 0, err
	}
	return uid, nil
}
func (s *Storage) GetUser(filter bson.D) (singleResult *mongo.SingleResult, err error) {
	collection := s.MongoDB.Database(constant.DB).Collection("user")

	res := collection.FindOne(context.TODO(), filter)

	return res, res.Err()
}

func (s *Storage) CreateUser(document interface{}) error {
	collection := s.MongoDB.Database(constant.DB).Collection("user")

	_, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetUserToken(filter bson.D) (resp *mongo.SingleResult, err error) {
	collection := s.MongoDB.Database(constant.DB).Collection("user_token")

	res := collection.FindOne(context.TODO(), filter)

	return res, res.Err()
}

func (s *Storage) CreateUserToken(document interface{}) error {
	collection := s.MongoDB.Database(constant.DB).Collection("user_token")

	_, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateUserToken(filter bson.D, update bson.D) error {
	collection := s.MongoDB.Database(constant.DB).Collection("user_token")
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
func (s *Storage) UpdateUser(filter bson.D, update bson.D) error {
	collection := s.MongoDB.Database(constant.DB).Collection("user")
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
