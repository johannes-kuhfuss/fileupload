package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
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
	Service service.DefaultUploadService
	Cfg     *config.AppConfig
}

func NewUploadHandler(cfg *config.AppConfig, svc service.DefaultUploadService) UploadHandler {
	return UploadHandler{
		Cfg:     cfg,
		Service: svc,
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
		msg := "Date or time not correct"
		logger.Warn(msg)
		apiErr := api_error.NewBadRequestError(msg)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}

	logger.Info(fmt.Sprintf("bcdate: %v, starttime: %v, endtime: %v", fd.BcDate, fd.StartTime, fd.EndTime))

	localFile := path.Join(uh.Cfg.Upload.Path, header.Filename)
	dst, err := os.Create(localFile)
	if err != nil {
		msg := "cannot create local file"
		logger.Error(msg, err)
		apiErr := api_error.NewInternalServerError(msg, err)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		msg := "cannot copy to local file"
		logger.Error(msg, err)
		apiErr := api_error.NewInternalServerError(msg, err)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}

	c.JSON(http.StatusCreated, nil)
}
