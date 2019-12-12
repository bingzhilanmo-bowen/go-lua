package main

import (
	"bingzhilanmo/go-lua/config"
	"bingzhilanmo/go-lua/controller"
	"bingzhilanmo/go-lua/models"
	"bingzhilanmo/go-lua/pkg/utils"
	"bingzhilanmo/go-lua/router"
	"bingzhilanmo/go-lua/service"
	"flag"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/mcuadros/go-gin-prometheus"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)


var (
	ConfigPath string
	LogLevel string
	LogFile  string
	file *os.File
)

func init() {
	flag.StringVar(&ConfigPath, "config_path", "config.toml", "config path.")
	flag.StringVar(&LogLevel, "log_level", "debug", "config path.")
	flag.StringVar(&LogFile, "log_file", "go_lua.log", "config path.")
}

func logConfig() {

	level,err := log.ParseLevel(LogLevel)

	if err != nil {
		level = log.DebugLevel
	}

	log.SetLevel(level)
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	log.SetReportCaller(true)
	file, err = os.OpenFile(LogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	log.SetOutput(file)
}

func run() {
	log.Info("start go-lua state engine .....")

	if err := config.LoadGlobalConfig(ConfigPath); err != nil {
		log.Panicf("init load config file errr %s", err.Error())
	}

	if err := models.InitDb() ; err != nil {
		log.Panicf("init database error %s", err.Error())
	}

	if config.GetGlobalConfig().DB.OpenListener {
		models.OpenListener()
	}

	//Init Cache
	utils.InitCache()

	//process指标采集
	utils.DoResourceMonitor()

	//init vm cache
	if config.GetGlobalConfig().Base.CacheLuaVm {
		service.WarmLuaVm()
	}

	//初始化http引擎
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.Use(controller.RequestIdMiddleware())

	//路由注册
	router.RegisterRouter(r)

	if config.GetGlobalConfig().Base.OpenPprof {
		ginpprof.Wrap(r)
	}

	err := r.Run(config.GetGlobalConfig().Http.Host + ":" + config.GetGlobalConfig().Http.Port)

	if err != nil {
		log.Panicf("start server error %s", err.Error())
	}

	log.Infof("success start server ip %s, port %s", config.GetGlobalConfig().Http.Host, config.GetGlobalConfig().Http.Port)
	log.Info("complete go-lua state engine .....")
}


// @title Svc engine API
// @version 0.1
// @description Server ft.
// @termsOfService

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func main() {

	logConfig()

	go run()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <- c:
			models.CloseDb()
			log.Debug("start go-lua state engine exit .....")
			file.Close()
			return
		}
	}


}
