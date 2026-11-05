// Package filestore 提供导出文件存储抽象。
//
// 当前实现使用本地磁盘，后续可以替换为 S3/OSS 等对象存储实现。
package filestore

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Store 文件存储接口。
type Store interface {
	// FilePath 返回相对路径对应的本地绝对路径，并做目录穿越校验。
	FilePath(relative string) (string, error)
	// URL 返回可下载 URL。
	URL(relative string) string
	// Open 返回可读取的本地文件路径。
	Open(relative string) (string, error)
}

// LocalStore 本地磁盘存储。
type LocalStore struct {
	root      string
	urlPrefix string
}

// Default 全局存储实例，由 bootstrap 初始化。
var Default Store

// Init 初始化本地存储。
func Init(root, urlPrefix string) {
	if root == "" {
		root = "data/export"
	}
	if urlPrefix == "" {
		urlPrefix = "/manage/export/download"
	}
	_ = os.MkdirAll(root, 0o755)
	Default = &LocalStore{
		root:      root,
		urlPrefix: strings.TrimRight(urlPrefix, "/"),
	}
}

func store() Store {
	if Default == nil {
		Init("data/export", "/manage/export/download")
	}
	return Default
}

// FilePath 返回相对路径对应的本地绝对路径。
func FilePath(relative string) (string, error) {
	return store().FilePath(relative)
}

// URL 返回可下载 URL。
func URL(relative string) string {
	return store().URL(relative)
}

// Open 返回可读取的本地文件路径。
func Open(relative string) (string, error) {
	return store().Open(relative)
}

func (s *LocalStore) FilePath(relative string) (string, error) {
	relative = filepath.Clean(strings.TrimPrefix(relative, "/"))
	if relative == "." || relative == "" || filepath.IsAbs(relative) || strings.HasPrefix(relative, "..") {
		return "", errors.New("invalid file path")
	}

	root, err := filepath.Abs(s.root)
	if err != nil {
		return "", err
	}
	full, err := filepath.Abs(filepath.Join(root, relative))
	if err != nil {
		return "", err
	}
	if full != root && !strings.HasPrefix(full, root+string(os.PathSeparator)) {
		return "", errors.New("file path escapes root")
	}
	return full, nil
}

func (s *LocalStore) URL(relative string) string {
	return s.urlPrefix + "/" + filepath.ToSlash(filepath.Clean(strings.TrimPrefix(relative, "/")))
}

func (s *LocalStore) Open(relative string) (string, error) {
	path, err := s.FilePath(relative)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", errors.New("is a directory")
	}
	return path, nil
}
