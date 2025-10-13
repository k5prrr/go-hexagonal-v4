package filemanager

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var Version = "2025-02-18"

// File
func CreateFile(path, text string) (bool, error) {
	if isDir(path) {
		return false, nil
	}
	if text == "" {
		text = " "
	}

	file, err := os.Create(path)
	if err != nil {
		return false, fmt.Errorf("не удалось открыть файл: %s", path)
	}
	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		return false, err
	}

	return fileExists(path), nil
}

func WriteFile(path, text string) (bool, error) {
	if !fileExists(path) || isDir(path) {
		return false, nil
	}
	if text == "" {
		text = " "
	}

	err := ioutil.WriteFile(path, []byte(text), 0644)
	if err != nil {
		return false, err
	}

	return true, nil
}

func CreatePathAndFile(path, text string) (bool, error) {
	dir := filepath.Dir(path)
	if !fileExists(dir) {
		if err := fm.CreatePath(dir); err != nil {
			return false, err
		}
	}

	return fm.CreateFile(path, text)
}

func ReadFile(path string) (string, error) {
	if !fileExists(path) || isDir(path) {
		return "", nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func PropertiesFile(path string) (map[string]interface{}, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"path":   path,
		"type":   getFileType(info),
		"size":   info.Size(),
		"time":   info.ModTime().Unix(),
		"perms":  info.Mode(),
		"UserId": getFileOwner(info),
		"inode":  getFileInode(info),
		"group":  getFileGroup(info),
	}, nil
}

// Directory
func CreateDirectory(path string) error {
	if fileExists(path) && isDir(path) {
		return nil
	}
	return os.Mkdir(path, 0755)
}

func CreatePath(path string) error {
	if fileExists(path) {
		return nil
	}

	dirs := strings.Split(path, string(filepath.Separator))
	currentPath := ""

	for _, dir := range dirs {
		currentPath = filepath.Join(currentPath, dir)
		if !fileExists(currentPath) {
			if err := fm.CreateDirectory(currentPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func ReadDirectory(path string) (map[string]interface{}, error) {
	if !fileExists(path) || !isDir(path) {
		return map[string]interface{}{
			"path":  path,
			"dirs":  []string{},
			"files": []string{},
			"err":   "dir not exists",
		}, nil
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	dirs := []string{}
	fileNames := []string{}

	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		} else {
			fileNames = append(fileNames, file.Name())
		}
	}

	return map[string]interface{}{
		"path":  path,
		"dirs":  dirs,
		"files": fileNames,
	}, nil
}

func ReadTreeDirectory(path string) (map[string]interface{}, error) {
	result, err := fm.ReadDirectory(path)
	if err != nil {
		return nil, err
	}

	dirs := []string{}
	files := []string{}

	for _, dir := range result["dirs"].([]string) {
		fullDirPath := filepath.Join(path, dir)
		dirs = append(dirs, fullDirPath)
		subResult, err := fm.ReadTreeDirectory(fullDirPath)
		if err != nil {
			return nil, err
		}
		dirs = append(dirs, subResult["dirs"].([]string)...)
		files = append(files, subResult["files"].([]string)...)
	}

	for _, file := range result["files"].([]string) {
		files = append(files, filepath.Join(path, file))
	}

	return map[string]interface{}{
		"path":  path,
		"dirs":  dirs,
		"files": files,
	}, nil
}

func PropertiesDirectory(path string) (map[string]interface{}, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"path":   path,
		"type":   getFileType(info),
		"size":   fm.SizeDirectory(path),
		"time":   info.ModTime().Unix(),
		"perms":  info.Mode(),
		"UserId": getFileOwner(info),
		"inode":  getFileInode(info),
		"group":  getFileGroup(info),
	}, nil
}

func SizeDirectory(path string) int64 {
	var size int64

	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size
}

// Additionally
func Delete(path string) (bool, error) {
	if !fileExists(path) {
		return true, nil
	}

	if !isDir(path) {
		if err := os.Remove(path); err != nil {
			return false, err
		}
		return !fileExists(path), nil
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return false, err
	}

	for _, file := range files {
		fullPath := filepath.Join(path, file.Name())
		if file.IsDir() {
			if _, err := fm.Delete(fullPath); err != nil {
				return false, err
			}
		} else {
			if err := os.Remove(fullPath); err != nil {
				return false, err
			}
		}
	}

	if err := os.Remove(path); err != nil {
		return false, err
	}

	return !fileExists(path), nil
}

func Copy(src, dst string) (bool, error) {
	if !fileExists(src) {
		return false, nil
	}

	if isFile(src) {
		return copyFile(src, dst)
	}

	if isDir(src) {
		if !fileExists(dst) {
			if err := fm.CreateDirectory(dst); err != nil {
				return false, err
			}
		}

		files, err := ioutil.ReadDir(src)
		if err != nil {
			return false, err
		}

		for _, file := range files {
			srcPath := filepath.Join(src, file.Name())
			dstPath := filepath.Join(dst, file.Name())

			if file.IsDir() {
				if _, err := fm.Copy(srcPath, dstPath); err != nil {
					return false, err
				}
			} else {
				if _, err := copyFile(srcPath, dstPath); err != nil {
					return false, err
				}
			}
		}

		return true, nil
	}

	return false, nil
}

func Rename(path, newName string) (bool, error) {
	newPath := filepath.Join(filepath.Dir(path), newName)
	if err := os.Rename(path, newPath); err != nil {
		return false, err
	}
	return true, nil
}

func Move(path, newPath string) (bool, error) {
	if err := os.Rename(path, newPath); err != nil {
		return false, err
	}
	return true, nil
}

// Utilities
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func getFileType(info os.FileInfo) string {
	if info.IsDir() {
		return "dir"
	}
	return "file"
}

func getFileOwner(info os.FileInfo) int {
	return int(info.Sys().(*syscall.Stat_t).Uid)
}

func getFileInode(info os.FileInfo) uint64 {
	return info.Sys().(*syscall.Stat_t).Ino
}

func getFileGroup(info os.FileInfo) int {
	return int(info.Sys().(*syscall.Stat_t).Gid)
}

func copyFile(src, dst string) (bool, error) {
	sourceFile, err := os.Open(src)
	if err != nil {
		return false, err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return false, err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return false, err
	}

	return true, nil
}
