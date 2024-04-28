package service

import (
	"io"
	"os"
	"path"

	"github.com/johannes-kuhfuss/fileupload/config"
	"github.com/johannes-kuhfuss/fileupload/dto"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

type UploadService interface {
	Upload(dto.FileDta)
}

type DefaultUploadService struct {
	Cfg *config.AppConfig
}

func NewUploadService(cfg *config.AppConfig) DefaultUploadService {
	return DefaultUploadService{
		Cfg: cfg,
	}
}

func (s DefaultUploadService) Upload(fd dto.FileDta) (written int64, err error) {
	localFile := buildFileName(s.Cfg.Upload.Path, fd.BcDate, fd.StartTime, fd.EndTime, fd.Header.Filename)
	dst, err := os.Create(localFile)
	if err != nil {
		return 0, err
	}
	defer dst.Close()
	bw, err := io.Copy(dst, fd.File)
	if err != nil {
		return 0, err
	}
	return bw, nil
}

func buildFileName(uploadPath string, bcDate string, startTime string, endTime string, fileName string) string {
	logger.Info(bcDate + startTime + endTime)
	return path.Join(uploadPath, fileName)
}
