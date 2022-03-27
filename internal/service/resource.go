package service

import (
	"log"
)

func (srv *Service) GetResource() (error) {
	log.Printf("GetResource")
	return srv.dao.GetResourceVersion("mem")
}
