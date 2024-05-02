package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/johannes-kuhfuss/fileupload/config"
	"github.com/johannes-kuhfuss/fileupload/dto"
	"github.com/johannes-kuhfuss/fileupload/helper"
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
	session := sessions.Default(c)
	fd.Uploader = session.Get("uploadUser").(string)
	fd.FileId = uuid.New()

	logger.Info(fmt.Sprintf("Upload request %v received from %v", fd.FileId.String(), fd.Uploader))

	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		msg := "error getting form"
		logger.Error(msg, err)
		apiErr := api_error.NewInternalServerError(msg, err)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	fd.File, fd.Header, err = c.Request.FormFile("file")
	if err != nil {
		msg := "cannot read remote file"
		logger.Error(msg, err)
		apiErr := api_error.NewInternalServerError(msg, err)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	defer fd.File.Close()

	if !misc.SliceContainsString(uh.Cfg.Upload.AllowedExtensions, filepath.Ext(fd.Header.Filename)) {
		msg := fmt.Sprintf("Cannot upload file with extension %v", filepath.Ext(fd.Header.Filename))
		helper.AddToUploadList(uh.Cfg, fd, msg)
		logger.Warn(msg)
		apiErr := api_error.NewBadRequestError(msg)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}

	logger.Info(fmt.Sprintf("Upload request %v, File: %v", fd.FileId.String(), fd.Header.Filename))

	fd.BcDate = c.PostForm("bcdate")
	fd.StartTime = c.PostForm("starttime")
	fd.EndTime = c.PostForm("endtime")
	err = uh.Cfg.RunTime.Sani.Sanitize(&fd)
	if err != nil {
		msg := "Date and/or time information missing"
		helper.AddToUploadList(uh.Cfg, fd, msg)
		logger.Warn(msg)
		apiErr := api_error.NewBadRequestError(msg)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}

	logger.Info(fmt.Sprintf("request %v metadata. Brodcast Date: %v, Start Time: %v, End Time: %v", fd.FileId.String(), fd.BcDate, fd.StartTime, fd.EndTime))

	fd.FileSize, err = uh.Svc.Upload(fd)
	if err != nil {
		msg := "Could not complete the upload"
		helper.AddToUploadList(uh.Cfg, fd, msg)
		logger.Error(msg, err)
		apiErr := api_error.NewInternalServerError(msg, err)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	helper.AddToUploadList(uh.Cfg, fd, "Successfully completed")
	logger.Info(fmt.Sprintf("Upload request %v (file: %v) sucessfully completed.", fd.FileId.String(), fd.Header.Filename))

	ret := dto.FileRet{
		FileName:     fd.Header.Filename,
		BytesWritten: fd.FileSize,
	}
	c.JSON(http.StatusCreated, ret)
}
