package service

import (
	"log"
)

func (srv *Service) GetLock(owner string) (error) {
	log.Printf("GetLock")
	return srv.dao.GetExclusiveLock("mem", owner)
}
