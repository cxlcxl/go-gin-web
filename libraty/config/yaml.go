package config

import (
	"gin-web/app/vars"
	ci "gin-web/libraty/config/interface"
	"gin-web/libraty/container"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"time"
)

// CreateYamlFactory 参数设置为可变参数的文件名，这样参数就可以不需要传递，如果传递了多个，我们只取第一个参数作为配置文件名
func CreateYamlFactory(fileName ...string) ci.YamlConfigInterface {
	yamlConfig := viper.New()
	// 配置文件所在目录
	yamlConfig.AddConfigPath(vars.BasePath + "/config")
	// 需要读取的文件名,默认为：web
	if len(fileName) == 0 {
		yamlConfig.SetConfigName("web")
	} else {
		yamlConfig.SetConfigName(fileName[0])
	}
	//设置配置文件类型(后缀)为 yml
	yamlConfig.SetConfigType("yaml")

	if err := yamlConfig.ReadInConfig(); err != nil {
		log.Fatal(err.Error())
	}

	return &YmlConfig{
		yamlConfig,
	}
}

type YmlConfig struct {
	viper *viper.Viper
}

// 由于 vipver 包本身对于文件的变化事件有一个bug，相关事件会被回调两次
// 常年未彻底解决，相关的 issue 清单：https://github.com/spf13/viper/issues?q=OnConfigChange
// 设置一个内部全局变量，记录配置文件变化时的时间点，如果两次回调事件事件差小于1秒，我们认为是第二次回调事件，而不是人工修改配置文件
// 这样就避免了 vipver 包的这个bug
var lastChangeTime time.Time

func init() {
	lastChangeTime = time.Now()
}

// ConfigFileChangeListen 监听文件变化
func (y *YmlConfig) ConfigFileChangeListen() {
	y.viper.OnConfigChange(func(changeEvent fsnotify.Event) {
		if time.Now().Sub(lastChangeTime).Seconds() >= 1 {
			if changeEvent.Op.String() == "WRITE" {
				y.clearCache()
				lastChangeTime = time.Now()
			}
		}
	})
	y.viper.WatchConfig()
}

// 清空已经窜换的配置项信息
func (y *YmlConfig) clearCache() {
	container.CreateContainersFactory().FuzzyDelete(vars.ConfigKeyPrefix)
}

// 判断相关键是否已经缓存
func (y *YmlConfig) keyIsCache(key string) (value interface{}, exists bool) {
	if value, exists = container.CreateContainersFactory().KeyIsExists(vars.ConfigKeyPrefix + key); exists {
		return value, true
	} else {
		return nil, false
	}
}

// 对键值进行缓存
func (y *YmlConfig) cache(key string, value interface{}) bool {
	return container.CreateContainersFactory().Set(vars.ConfigKeyPrefix+key, value)
}

// GetString 字符串类型
func (y *YmlConfig) GetString(key string) string {
	if value, ok := y.keyIsCache(key); ok {
		return value.(string)
	} else {
		value := y.viper.GetString(key)
		y.cache(key, value)
		return value
	}
}

// Get 一个原始值
func (y *YmlConfig) Get(key string) interface{} {
	if value, exists := y.keyIsCache(key); exists {
		return value
	} else {
		value = y.viper.Get(key)
		y.cache(key, value)
		return value
	}
}

// GetBool ..
func (y *YmlConfig) GetBool(keyName string) bool {
	if value, exists := y.keyIsCache(keyName); exists {
		return value.(bool)
	} else {
		value := y.viper.GetBool(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetInt ..
func (y *YmlConfig) GetInt(keyName string) int {
	if value, exists := y.keyIsCache(keyName); exists {
		return value.(int)
	} else {
		value := y.viper.GetInt(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetInt32 ..
func (y *YmlConfig) GetInt32(keyName string) int32 {
	if value, exists := y.keyIsCache(keyName); exists {
		return value.(int32)
	} else {
		value := y.viper.GetInt32(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetInt64 ..
func (y *YmlConfig) GetInt64(keyName string) int64 {
	if value, exists := y.keyIsCache(keyName); exists {
		return value.(int64)
	} else {
		value := y.viper.GetInt64(keyName)
		y.cache(keyName, value)
		return value
	}
}

// GetFloat64 ..
func (y *YmlConfig) GetFloat64(keyName string) float64 {
	if value, exists := y.keyIsCache(keyName); exists {
		return value.(float64)
	} else {
		value := y.viper.GetFloat64(keyName)
		y.cache(keyName, value)
		return value
	}
}
