package dao

import (
	"github.com/distributed_lock/internal/model"
)

func (d *Dao) GetExclusiveLock(resource_name, owner string) (error) {
	lock := model.Lock{ResourceName: resource_name, Owner: owner}
	return lock.ForUpdateLock(d.engine)
}
