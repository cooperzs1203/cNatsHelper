package cNatsHelper

import "net"

var (
	cNConfig *cNatsConfig
)
//todo:初步改动为提取指定路径的，后续加强为接受RPC消息重新加载配置文件
func GlobalConfig() *cNatsConfig {
	if cNConfig == nil {
		cNConfig = &cNatsConfig{
			natsServeHost:"127.0.0.1",
			natsServePort:"5500",
		}
	}

	return cNConfig
}

type cNatsConfig struct {
	natsServePort string
	natsServeHost string
}

func (cnf *cNatsConfig) GetNatsServeAddr() string {
	addr := net.JoinHostPort(cnf.natsServeHost , cnf.natsServePort)
	return "nats://" + addr
}