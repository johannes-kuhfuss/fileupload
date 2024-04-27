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
	logger.Info("Upload request received")
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		msg := "error getting form"
		logger.Error(msg, err)
		apiErr := api_error.NewInternalServerError(msg, err)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		msg := "cannot read remote file"
		logger.Error(msg, err)
		apiErr := api_error.NewInternalServerError(msg, err)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	defer file.Close()

	if !misc.SliceContainsString(uh.Cfg.Upload.AllowedExtensions, filepath.Ext(handler.Filename)) {
		msg := fmt.Sprintf("Cannot upload file with extension %v", filepath.Ext(handler.Filename))
		logger.Warn(fmt.Sprintf("User tried to upload file with name %v. Extension not allowed.", handler.Filename))
		apiErr := api_error.NewBadRequestError(msg)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}

	bcdate := c.PostForm("bcdate")
	starttime := c.PostForm("starttime")
	endtime := c.PostForm("endtime")
	// sanitize data
	// parse data
	// package into struct
	logger.Info(fmt.Sprintf("bcdate: %v, starttime: %v, endtime: %v", bcdate, starttime, endtime))

	localFile := path.Join(uh.Cfg.Upload.Path, handler.Filename)
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
