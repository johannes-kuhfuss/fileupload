package helper

import (
	"github.com/johannes-kuhfuss/fileupload/config"
	"github.com/johannes-kuhfuss/fileupload/domain"
)

func AddToUploadList(cfg *config.AppConfig, fileName string, bcdate string, startime string, endtime string, status string, uploader string, size string) {
	ul := domain.Upload{
		FileName:  fileName,
		BcDate:    bcdate,
		StartTime: startime,
		EndTime:   endtime,
		Status:    status,
		Uploader:  uploader,
		Size:      size,
	}
	cfg.RunTime.UploadList = append(cfg.RunTime.UploadList, ul)
}
