package environment

const (
	LocalEnvironment Environment = "local" // 本地开发环境（默认的运行环境）
	DevEnvironment   Environment = "dev"   // 开发环境（默认的运行环境）
	QaEnvironment    Environment = "qa"    // 测试环境
	PgEnvironment    Environment = "pg"    // 演练环境
	LptEnvironment   Environment = "lpt"   // 压力测试环境
	PreEnvironment   Environment = "pre"   // 预发布环境
	PrdEnvironment   Environment = "prd"   // 生产环境
)

// ----------------------------------------
//
//	运行环境类型定义
//
// ----------------------------------------
type Environment string

// 获取全局运行环境字符串形式
func (e Environment) String() string { return string(e) }

// 检查当前全局运行环境是否是给定的值
func (e Environment) Is(env Environment) bool { return e == env }

// 检查当前全局运行环境是否是在给定的值范围内
func (e Environment) In(items []Environment) bool {
	for i, j := 0, len(items); i < j; i++ {
		if e == items[i] {
			return true
		}
	}
	return false
}

// 检查当前全局运行环境是否是开发环境
func (e Environment) IsDev() bool { return e.Is(DevEnvironment) }

// 检查当前全局运行环境是否是测试环境
func (e Environment) IsQa() bool { return e.Is(QaEnvironment) }

// 检查当前全局运行环境是否是演练环境
func (e Environment) IsPg() bool { return e.Is(PgEnvironment) }

// 检查当前全局运行环境是否是压测环境
func (e Environment) IsLpt() bool { return e.Is(LptEnvironment) }

// 检查当前全局运行环境是否是预发环境
func (e Environment) IsPre() bool { return e.Is(PreEnvironment) }

// 检查当前全局运行环境是否是生产环境
func (e Environment) IsPrd() bool { return e.Is(PrdEnvironment) }
