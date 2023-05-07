package services

import "github.com/lenarsaitov/metrics-tpl/internal/models/services/memstorage"

type Services struct {
	MemStorageService memstorage.Service
}

func New(memStorageService memstorage.Service) *Services {
	return &Services{
		MemStorageService: memStorageService,
	}
}
