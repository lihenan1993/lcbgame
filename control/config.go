package control

import (
	"encoding/json"
	"github.com/pkg/errors"
	"mania/model"
	"os"
)

var c model.Config

func LoadConfigFromFile(path string) error {
	filePtr, err := os.Open(path)
	if err != nil {
		err = errors.Wrap(err, "read config.json failed")
		return err
	}
	defer filePtr.Close()

	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&c)
	if err != nil {
		err = errors.Wrap(err, "reflect failed")
		return err
	}

	return err
}

func GetRedisSource() *model.Redis {
	return c.Redis
}
func GetMongodbSource() *model.Mongodb {
	return c.Mongodb
}

func GetFeishuUrl() string {
	return c.FeishuUrl
}
