package common

import "errors"

var (
	ErrIsNotXml    = errors.New("template必须是xml类型的文件")
	ErrFileSize    = errors.New("file size不能超过1MB")
	ErrBuildFailed = errors.New("build failed")
)
