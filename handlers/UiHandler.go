package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/fileupload/config"
	"github.com/johannes-kuhfuss/fileupload/dto"
)

type UiHandler struct {
	Cfg *config.AppConfig
}

func NewUiHandler(cfg *config.AppConfig) UiHandler {
	return UiHandler{
		Cfg: cfg,
	}
}

func (uh *UiHandler) StatusPage(c *gin.Context) {
	configData := dto.GetConfig(uh.Cfg)
	c.HTML(http.StatusOK, "status.page.tmpl", gin.H{
		"title":      "Status",
		"configdata": configData,
	})
}

func (uh *UiHandler) AboutPage(c *gin.Context) {
	c.HTML(http.StatusOK, "about.page.tmpl", gin.H{
		"title": "About",
		"data":  nil,
	})
}

func (uh *UiHandler) UploadPage(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.page.tmpl", gin.H{
		"title": "Upload",
		"data":  nil,
	})
}
