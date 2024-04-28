package service

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

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
	bcd, err := time.Parse("2006-01-02", bcDate)
	if err != nil {
		logger.Error("Could not parse date.", err)
	}
	yStr := fmt.Sprintf("%02d", bcd.Year())
	mStr := fmt.Sprintf("%02d", bcd.Month())
	dStr := fmt.Sprintf("%02d", bcd.Day())
	st, err := time.Parse("15:04", startTime)
	if err != nil {
		logger.Error("Could not start time.", err)
	}
	et, err := time.Parse("15:04", endTime)
	if err != nil {
		logger.Error("Could not end time.", err)
	}
	stStr := fmt.Sprintf("%02d%02d", st.Hour(), st.Minute())
	etStr := fmt.Sprintf("%02d%02d", et.Hour(), et.Minute())
	fn := "UL__" + stStr + "-" + etStr + "__" + fileName

	os.MkdirAll(path.Join(uploadPath, yStr, mStr, dStr), os.ModePerm)
	return path.Join(uploadPath, yStr, mStr, dStr, fn)
}
