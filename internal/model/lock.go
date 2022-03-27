package model

import (
	"log"
	"time"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Lock struct {
	// *Model
	ResourceName  string `json:"resource_name"`
	Owner string  `json:"owner"`
}

func (lock Lock) TableName() string {
	return "lock_exclusive_lock"
}

func (lock Lock) ForUpdateLock(db *gorm.DB) (error) {
	err := db.Clauses(clause.Locking{
		Strength: "UPDATE",
	}).Where("resource_name = ?", lock.ResourceName).First(&lock).Error
	// SELECT * FROM `lock_exclusive_lock` FOR UPDATE

	// NOTE: lock not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Owner %s %s: lock not foud", lock.Owner, err)
		db.Omit().Create(&lock)

	} else if err != nil {
		log.Printf("Resource %s %s: locked", lock.ResourceName, err)
		return nil
	}

	// NOTE: Do something
	log.Printf("%s", lock.Owner)
	time.Sleep(time.Duration(5)*time.Second)

	db.Where("resource_name = ?", lock.ResourceName).Where("owner = ?", lock.Owner).Delete(&lock)
	// DELETE * FROM `lock_exclusive_lock` WHERE resource_name = ? AND owner = ?

	return nil
}
