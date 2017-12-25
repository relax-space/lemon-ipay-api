package core

import "sync"

type EnvParamDto struct {
	AppEnv      string `json:"app_env"`
	BmappingUrl string `json:"bmapping_url"`
	ConnEnv     string `json:"conn_env"`
	HostUrl     string `json:"host_url"`
}

var Env *EnvParamDto
var onceEnv sync.Once

func InitEnv(env *EnvParamDto) {
	onceEnv.Do(func() {
		Env = env
	})
}
