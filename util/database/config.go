package database

// Config ----------------------------------------
//
//	配置项（TOML 支持, Viper 支持）
//
// ----------------------------------------
type Config struct {
	Driver      string `toml:"driver" mapstructure:"driver"`
	Host        string `toml:"host" mapstructure:"host"`
	Port        int    `toml:"port" mapstructure:"port"`
	Username    string `toml:"username" mapstructure:"username"`
	Password    string `toml:"password" mapstructure:"password"`
	Database    string `toml:"database" mapstructure:"database"`
	Charset     string `toml:"charset" mapstructure:"charset"`
	MaxConn     int    `toml:"max_conn" mapstructure:"max_conn"`
	MaxIdleConn int    `toml:"max_idle_conn" mapstructure:"max_idle_conn"`
	Ping        bool   `toml:"ping" mapstructure:"ping"`
	Debug       bool   `toml:"debug" mapstructure:"debug"`
}
