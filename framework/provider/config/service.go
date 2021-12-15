package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"github.com/zzm996-zzm/arms/framework"
	"github.com/zzm996-zzm/arms/framework/contract"
	"gopkg.in/yaml.v2"
)

type ConfigService struct {
	c        framework.Container    //容器
	folder   string                 //文件夹
	keyDelim string                 //路径分隔符，默认是.
	envMaps  map[string]string      //env
	confMaps map[string]interface{} //配置文件结构，key为文件名
	confRaws map[string][]byte      // 配置文件的原始信息
	lock     *sync.RWMutex
}

var _ contract.Config = new(ConfigService)

func NewConfigService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	configEnvFolder := params[1].(string)
	envMaps   := params[2].(map[string]string)

	//检查文件夹是否存在
	if _, err := os.Stat(configEnvFolder); os.IsNotExist(err) {
		return nil, errors.New("folder " + configEnvFolder + " not exist: " + err.Error())
	}

	//实例化
	conf := &ConfigService{
		c:        container,
		folder:   configEnvFolder,
		envMaps:  envMaps,
		confMaps: map[string]interface{}{},
		confRaws: map[string][]byte{},
		keyDelim: ".",
		lock:     &sync.RWMutex{},
	}

	//读取每个文件
	files, err := ioutil.ReadDir(configEnvFolder)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		fileName := file.Name()
		err := conf.loadConfigFile(configEnvFolder, fileName)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	//监控文件夹文件
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = watch.Add(configEnvFolder)
	if err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		for {
			select {
			case ev := <-watch.Events:
				//FIXME:事件过多

				//判断事件发生的类型
				// Create 创建
				// Write 写入
				// Remove 删除
				path, _ := filepath.Abs(ev.Name)
				index := strings.LastIndex(path, string(os.PathSeparator))
				folder := path[:index]
				fileName := path[index+1:]

				if ev.Op&fsnotify.Create == fsnotify.Create {
					log.Println("创建文件 : ", ev.Name)
					conf.loadConfigFile(folder, fileName)
				}
				if ev.Op&fsnotify.Write == fsnotify.Write {
					log.Println("写入文件 : ", ev.Name)
					conf.loadConfigFile(folder, fileName)
				}
				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("删除文件 : ", ev.Name)
					conf.removeConfigFile(folder, fileName)
				}
			case err := <-watch.Errors:
				log.Println("error :", err)
				return
			}
		}
	}()

	return conf, nil
}

func (conf *ConfigService) removeConfigFile(folder string, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()
	s := strings.Split(file, ".")
	// 只有yaml或者yml后缀才执行
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]
		// 删除内存中对应的key
		delete(conf.confRaws, name)
		delete(conf.confMaps, name)
	}
	return nil
}

func (conf *ConfigService) loadConfigFile(folder string, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()

	s := strings.Split(file, ".")
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]
		bf, err := ioutil.ReadFile(filepath.Join(folder, file))
		if err != nil {
			return err
		}

		//直接针对文本做环境变量替换
		bf = replace(bf, conf.envMaps)

		//解析对应的文件
		c := map[string]interface{}{}
		if err := yaml.Unmarshal(bf, &c); err != nil {
			return err
		}

		conf.confMaps[name] = c
		conf.confRaws[name] = bf

		if name == "app" && conf.c.IsBind(contract.AppKey) {
			if p, ok := c["path"]; ok {
				appService := conf.c.MustMake(contract.AppKey).(contract.ArmsApp)
				appService.LoadAppConfig(cast.ToStringMapString(p))
			}
		}
	}

	return nil

}

func (conf *ConfigService) find(key string) interface{} {
	return searchMap(conf.confMaps, strings.Split(key, conf.keyDelim))
}

func replace(content []byte, maps map[string]string) []byte {
	if maps == nil {
		return content
	}

	for key, val := range maps {
		reKey := "env(" + key + ")"
		content = bytes.ReplaceAll(content, []byte(reKey), []byte(val))
	}
	return content
}

func searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	next, ok := source[path[0]]
	if ok {
		if len(path) == 1 {
			return next
		}

		switch next.(type) {
		case map[interface{}]interface{}:
			// 如果是interface的map，使用cast进行下value转换
			return searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			// 如果是map[string]，直接循环调用
			return searchMap(next.(map[string]interface{}), path[1:])
		default:
			// 否则的话，返回nil
			return nil
		}
	}

	return nil
}

func (conf *ConfigService) IsExist(key string) bool {
	if key == "" {
		return false
	}
	return conf.find(key) != nil
}
func (conf *ConfigService) Get(key string) interface{} {
	return conf.find(key)
}
func (conf *ConfigService) GetBool(key string) bool {
	return cast.ToBool(conf.find(key))
}
func (conf *ConfigService) GetInt(key string) int {
	return cast.ToInt(conf.find(key))
}
func (conf *ConfigService) GetFloat64(key string) float64 {
	return cast.ToFloat64(conf.find(key))
}
func (conf *ConfigService) GetTime(key string) time.Time {
	return cast.ToTime(conf.find(key))
}
func (conf *ConfigService) GetString(key string) string {
	return cast.ToString(conf.find(key))
}
func (conf *ConfigService) GetIntSlice(key string) []int {
	return cast.ToIntSlice(conf.find(key))
}
func (conf *ConfigService) GetStringSlice(key string) []string {
	return cast.ToStringSlice(conf.find(key))
}
func (conf *ConfigService) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(conf.find(key))
}
func (conf *ConfigService) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(conf.find(key))
}

func (conf *ConfigService) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(conf.find(key))
}

// Load 加载配置到某个对象
func (conf *ConfigService) Load(key string, val interface{}) error {
	return mapstructure.Decode(conf.find(key), val)
}
