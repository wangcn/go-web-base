package cli

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"mybase/cmd/bootstrap"
	"mybase/controller"
	"mybase/controller/person"
	"mybase/middleware"
	"mybase/pkg/environment"
	"mybase/util"
)

func init() {
	// 自动添加子命令
	bootstrap.RootCmd.AddCommand(apiCmd)
}

// api 子命令
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "http 服务",
	Long:  `http 服务，提供接口`,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.Get("env") == environment.PrdEnvironment {
			gin.SetMode(gin.ReleaseMode)
		}

		ln, err := net.Listen("tcp", viper.GetString("server.addr"))
		if err != nil {
			log.Panic(err)
		}

		g := gin.New()
		if viper.Get("env") != environment.PrdEnvironment {
			g.Use(gin.Logger())
		}
		// pprof.Register(g)
		g.Use(middleware.RequestID())
		g.Use(gzip.Gzip(gzip.DefaultCompression))
		g.Use(middleware.AccessLog())
		s := &http.Server{
			Handler:        g,
			ReadTimeout:    viper.GetDuration("server.readTimeout") * time.Millisecond,
			WriteTimeout:   viper.GetDuration("server.writeTimeout") * time.Millisecond,
			MaxHeaderBytes: 1 << 20,
		}
		setRoutes(g)

		go func() {
			_ = s.Serve(ln)
		}()

		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
		sleepDuration := 1 * time.Second
		for {
			sig := <-ch
			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
			switch sig {
			case syscall.SIGTERM:
				log.Println("Received SIGTERM.")
				// Stop the service gracefully.
				util.SingleSweeper().Stop(sleepDuration)
				log.Println("server gracefully shutdown", s.Shutdown(ctx))
				log.Println("done.")
				os.Exit(0)
			case syscall.SIGINT:
				log.Println("Received SIGINT.")
				// Stop the service gracefully.
				util.SingleSweeper().Stop(sleepDuration)
				log.Println("server gracefully shutdown", s.Shutdown(ctx))
				log.Println("done.")
				os.Exit(0)
			case syscall.SIGUSR2:
				log.Println("Received SIGUSR2.")
				// Stop the service gracefully.
				util.SingleSweeper().Stop(sleepDuration)
				log.Println("server gracefully shutdown", s.Shutdown(ctx))
				log.Println("done.")
				os.Exit(0)
			}
		}
	},
}

func setRoutes(router *gin.Engine) {
	router.GET("/ping", util.Handle(controller.Ping))

	v1 := router.Group("v1")
	{
		v1.GET("/user/profile", util.Handle(person.ProfileInfo))
	}
}
