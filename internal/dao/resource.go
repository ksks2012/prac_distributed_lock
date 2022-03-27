package dao

import (
	"github.com/distributed_lock/internal/model"
)

func (d *Dao) GetResourceVersion(resource_name string) (error) {
	resource := model.Resource{ResourceName: resource_name}
	return resource.ForUpdateVersion(d.engine)
}
