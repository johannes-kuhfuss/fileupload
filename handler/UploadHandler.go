package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/fileupload/config"
	"github.com/johannes-kuhfuss/fileupload/dto"
	"github.com/johannes-kuhfuss/fileupload/service"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"github.com/johannes-kuhfuss/services_utils/misc"
)

type UploadHandler struct {
	Svc service.DefaultUploadService
	Cfg *config.AppConfig
}

func NewUploadHandler(cfg *config.AppConfig, svc service.DefaultUploadService) UploadHandler {
	return UploadHandler{
		Cfg: cfg,
		Svc: svc,
	}
}

func (uh UploadHandler) Receive(c *gin.Context) {
	var (
		fd dto.FileDta
	)
	logger.Info("Upload request received")
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		msg := "error getting form"
		logger.Error(msg, err)
		apiErr := api_error.NewInternalServerError(msg, err)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		msg := "cannot read remote file"
		logger.Error(msg, err)
		apiErr := api_error.NewInternalServerError(msg, err)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	defer file.Close()

	if !misc.SliceContainsString(uh.Cfg.Upload.AllowedExtensions, filepath.Ext(header.Filename)) {
		msg := fmt.Sprintf("Cannot upload file with extension %v", filepath.Ext(header.Filename))
		logger.Warn(fmt.Sprintf("User tried to upload file with name %v. Extension not allowed.", header.Filename))
		apiErr := api_error.NewBadRequestError(msg)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}

	fd.File = file
	fd.Header = header
	fd.BcDate = c.PostForm("bcdate")
	fd.StartTime = c.PostForm("starttime")
	fd.EndTime = c.PostForm("endtime")

	err = uh.Cfg.RunTime.Sani.Sanitize(&fd)

	if err != nil {
		msg := "Date and/or time not present"
		logger.Warn(msg)
		apiErr := api_error.NewBadRequestError(msg)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}

	logger.Info(fmt.Sprintf("file: %v, bcdate: %v, starttime: %v, endtime: %v", fd.Header.Filename, fd.BcDate, fd.StartTime, fd.EndTime))

	uploadUser := c.MustGet(gin.AuthUserKey).(string)

	bw, err := uh.Svc.Upload(fd, uploadUser)
	if err != nil {
		msg := "cannot create local file"
		logger.Error(msg, err)
		apiErr := api_error.NewInternalServerError(msg, err)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	ret := dto.FileRet{
		FileName:     fd.Header.Filename,
		BytesWritten: bw,
	}

	c.JSON(http.StatusCreated, ret)
}
