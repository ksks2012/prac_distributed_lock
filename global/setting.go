package global

import (
	"github.com/distributed_lock/pkg/setting"
	"github.com/distributed_lock/pkg/logger"
)

var (
	AppSetting *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
	Logger          *logger.Logger
)
