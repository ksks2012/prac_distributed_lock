package model

import (
	"log"
	"time"
	"errors"

	"gorm.io/gorm"
)

type Resource struct {
	// *Model
	ResourceName  string `json:"resource_name"`
	Share   string `json:"share"`
	Version int  `json:"version"`
}

func (resource Resource) TableName() string {
	return "lock_resource"
}

func (resource Resource) ForUpdateVersion(db *gorm.DB) (error) {
	succ := false
	for !succ {
		err := db.Where("resource_name = ?", resource.ResourceName).First(&resource).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Resource %s : not foud", resource.ResourceName)
			time.Sleep(time.Duration(1)*time.Second)
			continue
		} else if err != nil {
			log.Printf("Select Resource %s error: %s", resource.ResourceName, err)
			return nil
		}
		// NOTE: update version
		resource.Version += 1
		err = db.Model(&resource).Where("resource_name = ?", resource.ResourceName).Update("version", resource.Version).Error
		if err != nil {
			log.Printf("update Resource %s error: %s", resource.ResourceName, err)
			return err
		}
		log.Printf("update Resource %s success", resource.ResourceName)
		succ = true
	}

	return nil
}
