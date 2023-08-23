package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"mybase/pkg/environment"
)

var config *configuration

type configuration struct {
	paths       []string
	currentPath string
	vipers      map[string]*viper.Viper
	sync.Mutex
}

func Cfg(file string) *viper.Viper {
	config.Lock()
	defer config.Unlock()

	if cfg, ok := config.vipers[file]; ok {
		return cfg
	}

	// 读取单项配置
	subConfig := viper.New()
	subConfigFile := config.getConfigFile(file, environment.Env())
	if subConfigFile != "" {
		subConfig.SetConfigFile(subConfigFile)
		if err := subConfig.ReadInConfig(); err != nil {
			Log().With(zap.Error(err)).Error("read env config file err")
			return nil
		}
	}

	config.vipers[file] = subConfig
	return subConfig
}

func GetCfgPath() string {
	return config.currentPath
}

func InitConfig() {
	config = &configuration{
		paths: []string{
			"/data/etc/cc/" + strings.ReplaceAll(ProjectName, "-", "_"),
			".",
			filepath.Dir(os.Args[0]),
			"..",
			"../..",
		},
		vipers: make(map[string]*viper.Viper),
	}

	// 设置默认配置
	defaultConfig := config.getConfigFile("app", environment.Env())
	viper.SetConfigFile(defaultConfig)
	if err := viper.ReadInConfig(); err != nil {
		Log().With(zap.Error(err)).Error("read default config file err")
	}

}

func (c *configuration) getConfigFile(file string, env environment.Environment) string {

	configFile := ""
	for _, path := range c.paths {
		tmpConfigPath := fmt.Sprintf("%s/conf/%s", path, env)
		tmpConfigFile := fmt.Sprintf("%s/%s.toml", tmpConfigPath, file)
		if _, err := os.Stat(tmpConfigFile); err == nil {
			c.currentPath = tmpConfigPath
			configFile = tmpConfigFile
			break
		}
	}
	return configFile
}
