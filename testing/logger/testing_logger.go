package main

import(
	"context"
	"flag"
	"log"
	"strings"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/distributed_lock/pkg/setting"
	"github.com/distributed_lock/global"
	"github.com/distributed_lock/pkg/logger"
)

var (
	runMode string
	cfg     string
)

func init() {
	err := setupFlag()
	if err != nil {
		log.Fatalf("init.setupFlag err: %v", err)
	}

	err = setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}

}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	global.Logger.Infof(ctx, "db_lock")

	cancel()
}

func setupFlag() error {
	flag.StringVar(&runMode, "mode", "", "running level (info, debug)")
	flag.StringVar(&cfg, "config", "etc/", "assgin the path of config file")
	flag.Parse()

	return nil
}


func setupSetting() error {
	s, err := setting.NewSetting(strings.Split(cfg, ",")...)
	if err != nil {
		return err
	}
	err = s.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	// TODO: run mode

	return nil
}


func setupLogger() error {
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename: global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)

	return nil
}
