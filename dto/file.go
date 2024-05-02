package dto

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type FileDta struct {
	File      multipart.File
	Header    *multipart.FileHeader
	BcDate    string `san:"trim,xss,max=10"`
	StartTime string `san:"trim,xss,max=10"`
	EndTime   string `san:"trim,xss,max=10"`
	Uploader  string
	FileSize  int64
	FileId    uuid.UUID
}

type FileRet struct {
	FileName     string `json:"file_name"`
	BytesWritten int64  `json:"bytes_written"`
}
