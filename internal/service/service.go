package service

import (
	"context"

	// otgorm "github.com/eddycjy/opentracing-gorm"

	"github.com/distributed_lock/global"
	"github.com/distributed_lock/internal/dao"
)

type Service struct {
	ctx context.Context
	dao *dao.Dao
}

// FIXME: use context in gorm v2
func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = dao.New(global.DBEngine)
	return svc
}
