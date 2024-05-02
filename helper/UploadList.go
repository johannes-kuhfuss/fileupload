package helper

import (
	"time"

	"github.com/johannes-kuhfuss/fileupload/config"
	"github.com/johannes-kuhfuss/fileupload/domain"
)

func AddToUploadList(cfg *config.AppConfig, fileName string, bcdate string, startime string, endtime string, status string, uploader string, size string) {
	t := time.Now()
	ul := domain.Upload{
		UploadDate: t.Format("2006-01-02 15:04:05"),
		FileName:   fileName,
		BcDate:     bcdate,
		StartTime:  startime,
		EndTime:    endtime,
		Status:     status,
		Uploader:   uploader,
		Size:       size,
	}
	cfg.RunTime.UploadList = append(cfg.RunTime.UploadList, ul)
}
