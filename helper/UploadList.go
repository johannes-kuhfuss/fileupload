package helper

import (
	"math"
	"strconv"
	"time"

	"github.com/johannes-kuhfuss/fileupload/config"
	"github.com/johannes-kuhfuss/fileupload/domain"
	"github.com/johannes-kuhfuss/fileupload/dto"
)

func AddToUploadList(cfg *config.AppConfig, fd dto.FileDta, status string) {
	t := time.Now()
	ul := domain.Upload{
		UploadDate: t.Format("2006-01-02 15:04:05"),
		FileName:   fd.Header.Filename,
		BcDate:     fd.BcDate,
		StartTime:  fd.StartTime,
		EndTime:    fd.EndTime,
		Status:     status,
		Uploader:   fd.Uploader,
	}
	if fd.FileSize == 0 {
		ul.Size = ""
	} else {
		sizekb := float64(fd.FileSize) / (1 << 20)
		sizeStr := strconv.FormatInt(int64((math.Round(sizekb))), 10) + "MB"
		ul.Size = sizeStr
	}

	cfg.RunTime.UploadList = append(cfg.RunTime.UploadList, ul)
}
