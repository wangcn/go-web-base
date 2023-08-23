package environment

import (
	"fmt"
	"sync"
)

var (
	mutex   = new(sync.Mutex) // 当前运行时环境写锁
	current = DevEnvironment  // 当前运行时环境（默认为开发环境）

	// 所有可用的运行时环境列表
	environments = []Environment{
		LocalEnvironment,
		DevEnvironment,
		QaEnvironment,
		PgEnvironment,
		LptEnvironment,
		PreEnvironment,
		PrdEnvironment,
	}
)

// 获取当前运行环境
func GetEnv() Environment { return current }

// 获取当前运行环境（字符串形式）
func GetEnvString() string { return current.String() }

// 确定给定的运行环境是否有效
func IsValidEnv(env Environment) bool { return env.In(environments) }

// 确定给定的运行环境是否有效（字符串形式）
func IsValidEnvString(env string) bool { return IsValidEnv(Environment(env)) }

// 设置全局运行环境
func SetEnv(env Environment) error {
	mutex.Lock()
	defer mutex.Unlock()

	if !IsValidEnv(env) {
		return fmt.Errorf("the given environment %s is invalid", env)
	}
	if !current.Is(env) {
		current = env
	}
	return nil
}

// 设置全局运行环境（字符串形式）
func SetEnvString(env string) error { return SetEnv(Environment(env)) }

// 注册自定义运行环境
func RegisterEnv(env Environment) {
	mutex.Lock()
	defer mutex.Unlock()

	if !env.In(environments) {
		environments = append(environments, env)
	}
}

// 注册自定义运行环境（字符串自动解析）
func RegisterEnvString(env string) { RegisterEnv(Environment(env)) }

func Env() Environment { return GetEnv() }

func InitEnv(env Environment) error { return SetEnv(env) }
