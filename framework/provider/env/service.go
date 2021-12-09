package env

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/arms/framework/contract"
)

var _ contract.Env = new(EnvService)

type EnvService struct {
	folder string
	maps   map[string]string
}

func (env *EnvService) AppEnv() string {
	return env.Get("APP_ENV")
}

func (env *EnvService) IsExist(key string) bool {
	_, has := env.maps[key]

	return has
}

func (env *EnvService) Get(key string) string {
	return env.maps[key]
}

func (env *EnvService) All() map[string]string {
	return env.maps
}

func NewEnvService(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("NewArmsEnv param error")
	}

	folder := params[0].(string)

	envService := &EnvService{
		folder: folder,
		maps:   make(map[string]string),
	}
	envService.maps["APP_ENV"] = contract.EnvDevelopment

	file := filepath.Join(folder, ".env")

	//获取env文件变量
	readEnvFile(file, envService.maps)
	//获取环境变量
	os := getOsEnv()
	//使用环境变量覆盖env文件变量
	overwrite(os, envService.maps)

	return envService, nil
}

func readEnvFile(file string, maps map[string]string) {
	//打开文件
	fi, err := os.Open(file)
	if err == nil {
		defer fi.Close()

		//读取文件
		br := bufio.NewReader(fi)
		for {
			//按照行读取
			line, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}

			//按照等号解析
			s := bytes.SplitN(line, []byte{'='}, 2)
			//如果不符合规范，则过滤
			if len(s) < 2 {
				continue
			}
			//保存map
			key := string(s[0])
			value := string(s[1])
			maps[key] = value
		}
	}

}

func getOsEnv() map[string]string {
	maps := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) < 2 {
			continue
		}
		key := pair[0]
		value := pair[1]
		maps[key] = value
	}

	return maps
}

func overwrite(os, us map[string]string) {
	for k, v := range os {
		if _, has := us[k]; has {
			us[k] = v
		}
	}

}
