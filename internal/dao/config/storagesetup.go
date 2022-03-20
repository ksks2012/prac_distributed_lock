package config

import (
	"errors"

	interfaces "github.com/distributed_lock/interfaces"
	"github.com/distributed_lock/pkg/setting"
)

// StorageSetup contains storage type and storage instance
type StorageSetup struct {
	Type     string
	Instance interfaces.StorageEngine
}

func (s *StorageSetup) NewDBEngine(databaseSetting *setting.DatabaseSettingS) (err error) {
	switch databaseSetting.DBType {
	case "pxc":
		s.Instance, err = setupMySQLRoundRobinStorageEngine(databaseSetting)
	case "mysql", "mariadb":
		s.Instance, err = setupMySQLStorageEngine(databaseSetting)
	default:
		err = errors.New("unknown storage engine type: " + databaseSetting.DBType)
	}
	if nil != err {
		s.Instance = nil
	} else {
		s.Type = databaseSetting.DBType
	}
	return err
}
