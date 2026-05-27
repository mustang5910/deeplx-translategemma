// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/mustang5910/deeplx-translategemma/internal/config"
	"golang.org/x/sync/semaphore"
)

type ServiceContext struct {
	Config    config.Config
	Semaphore *semaphore.Weighted
}

func NewServiceContext(c config.Config) *ServiceContext {
	maxConcurrent := int64(c.MaxConcurrent)
	if maxConcurrent <= 0 {
		maxConcurrent = 10
	}
	return &ServiceContext{
		Config:    c,
		Semaphore: semaphore.NewWeighted(maxConcurrent),
	}
}
