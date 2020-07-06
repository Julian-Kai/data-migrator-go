package config

import (
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var once sync.Once
var instance *viper.Viper

func init() {
	once.Do(func() {
		// default config
		viper := viper.New()
		viper.SetConfigType("yaml")
		viper.SetConfigName("app")
		viper.AddConfigPath("./configs")
		viper.AddConfigPath("../configs")
		// bind env variable and modify key mapping(e.g. envKey "SYSTEM_PORT" in k8s yaml mapping to configKey "system.port")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			//logs.Logger.Fatal("Parse config failed", logs.MakeLogFields(logs.FieldError, err.Error()), nil)
			panic(err)
		}

		instance = viper
	})
}

func GetConfig(key string) string {
	if key != "" {
		return instance.GetString(key)
	}
	return ""
}

func GetConfigBool(key string) bool {
	if key != "" {
		return instance.GetBool(key)
	}
	return false
}

func GetInteger(key string) int {
	if key != "" {
		return instance.GetInt(key)
	}
	return 0
}
