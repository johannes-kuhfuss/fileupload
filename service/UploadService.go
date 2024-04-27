package service

import "github.com/johannes-kuhfuss/fileupload/config"

type UploadService interface {
	Upload()
}

type DefaultUploadService struct {
	Cfg *config.AppConfig
}

func NewUploadService(cfg *config.AppConfig) DefaultUploadService {
	return DefaultUploadService{
		Cfg: cfg,
	}
}

func (s DefaultUploadService) Upload() {

}
