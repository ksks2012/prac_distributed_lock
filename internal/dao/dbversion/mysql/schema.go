package mysql

import (
	bsstorages "github.com/distributed_lock/internal/dao/dbversion"
)

func (m *schemaManager) UpgradeSchema(currentRevs *schemaRevision) (schemaChanged bool, err error) {
	status := bsstorages.SchemaUpgradeStatus{
		Changed:   false,
		LastError: nil,
	}
	status.RunUpgrade("exclusive-locks", m.UpgradeSchemaExclusiveLocks, currentRevs.ExclusiveLocks)
	status.RunUpgrade("resource", m.UpgradeSchemaResource, currentRevs.Resource)
	return status.Changed, status.LastError
}
