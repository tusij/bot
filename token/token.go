package token

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

const (
	TypeBot = "Bot"
)

// Token 用于调用接口的 token 结构
type Token struct {
	AppID       uint64
	AccessToken string
	Type        string
}

// New 创建一个新的 Token
func New(tokenType string) *Token {
	return &Token{
		Type: tokenType,
	}
}

// BotToken 机器人身份的 token
func BotToken(appID uint64, accessToken string) *Token {
	return &Token{
		AppID:       appID,
		AccessToken: accessToken,
		Type:        TypeBot,
	}
}

// GetString 获取授权头字符串
func (t *Token) GetString() string {
	return fmt.Sprintf("%v.%s", t.AppID, t.AccessToken)
}

// LoadFromConfig 从配置中读取 appid 和 token
func (t *Token) LoadFromConfig(file string) error {
	var conf struct {
		AppID uint64 `yaml:"appid"`
		Token string `yaml:"token"`
	}
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("read token from file failed, err: %v", err)
		return err
	}
	if err = yaml.Unmarshal(content, &conf); err != nil {
		log.Fatalf("parse config failed, err: %v", err)
		return err
	}
	t.AppID = conf.AppID
	t.AccessToken = conf.Token
	return nil
}
