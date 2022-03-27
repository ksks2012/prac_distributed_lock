package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"strconv"
	"golang.org/x/sys/unix"
	"gopkg.in/natefinch/lumberjack.v2"


	"github.com/distributed_lock/pkg/setting"
	"github.com/distributed_lock/global"
	"github.com/distributed_lock/internal/service"
	"github.com/distributed_lock/internal/model"
	"github.com/distributed_lock/internal/dao/config"
	"github.com/distributed_lock/pkg/logger"
)

var (
	runMode string
	cfg     string
	owner 	string
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

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}

	err = updateDB()
	if err != nil {
		log.Fatalf("init.updateDB err: %v", err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	global.Logger.Infof(ctx, "db_lock")
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, os.Interrupt, unix.SIGTERM)

	srv := service.New(ctx)
	// for loop
	for i := 1; i <= 10; i++ {
		srv.GetLock(strconv.Itoa(i))
	}
	cancel()
}

func setupFlag() error {
	flag.StringVar(&runMode, "mode", "", "running level (info, debug)")
	flag.StringVar(&cfg, "config", "etc/", "assgin the path of config file")
	flag.StringVar(&owner, "owner", "A", "set resource owner")
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

func updateDB() error {
	var err error
	updateDBSetup := &config.StorageSetup{}
	err = updateDBSetup.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	if err = updateDBSetup.Instance.Open(); nil != err {
		log.Fatalf("open storage connection failed: %v", err)
		return err
	}

	return nil
}

func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}

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
