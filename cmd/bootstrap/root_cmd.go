package bootstrap

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"mybase/pkg/environment"
	"mybase/util"
)

var RootCmd = &cobra.Command{
	Use:  util.ProjectName,
	Long: util.ProjectName,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		// 取得环境
		env, err := cmd.Flags().GetString("env")
		if err != nil {
			return errors.Wrap(err, "root cmd get env error")
		}

		// 初始化随机数
		rand.Seed(time.Now().UnixNano())

		// 初始化环境、日志、配置
		if err := environment.InitEnv(environment.Environment(env)); err != nil {
			return err
		}
		util.InitConfig()
		util.InitLog()

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var env string

func Run() {
	RootCmd.PersistentFlags().StringVar(&env, "env", "dev", "[local|dev|qa|pre|prd]")

	if err := RootCmd.Execute(); err != nil {
		msg := fmt.Sprintf("init root cmd error=%s", err.Error())
		panic(msg)
	}
}
