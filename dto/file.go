package dto

import "mime/multipart"

type FileDta struct {
	File      multipart.File
	Header    *multipart.FileHeader
	BcDate    string `san:"trim,xss,max=10"`
	StartTime string `san:"trim,xss,max=10"`
	EndTime   string `san:"trim,xss,max=10"`
}
