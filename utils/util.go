package utils

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"net/http"
	"path"
	"runtime"
	"time"
)

func LoadConfig(file string, config interface{}) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(content, config); err != nil {
		return err
	}
	return nil
}

func getConfigPath(name string) string {
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		fmt.Println(fmt.Sprintf("%s/%s", path.Dir(filename), name))
		return fmt.Sprintf("%s/%s", path.Dir(filename), name)
	}
	return ""
}

func CreateTransport(localAddr net.Addr, idleConns int) *http.Transport {
	dialer := &net.Dialer{
		Timeout:   60 * time.Second,
		KeepAlive: 60 * time.Second,
	}
	if localAddr != nil {
		dialer.LocalAddr = localAddr
	}
	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          idleConns,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   idleConns,
		MaxConnsPerHost:       idleConns,
	}
}
