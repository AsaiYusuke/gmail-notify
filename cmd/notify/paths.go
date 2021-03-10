package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
)

type paths struct {
	executableDir    string
	imageSrcFilePath string
	serviceDirs      []string
}

func (p *paths) init() {
	executeFile, err := os.Executable()
	if err != nil {
		log.Fatalf("Unable to get current directory: %v", err)
	}
	p.executableDir = filepath.Dir(executeFile)

	imageSrcFilePath := path.Join(p.executableDir, pathImageFilename)
	if _, err := os.Stat(imageSrcFilePath); err == nil {
		p.imageSrcFilePath = imageSrcFilePath
	}
}

func (p *paths) getLogFilePath() string {
	return path.Join(p.executableDir, pathLogFilename)
}

func (p *paths) getConfigJSONFilePath() string {
	return path.Join(p.executableDir, pathConfigFilename)
}

func (p *paths) getImageSrcFilePath() string {
	return p.imageSrcFilePath
}

func (p *paths) getServiceDirs() []string {
	if p.serviceDirs == nil {
		p.checkServiceDirs()
	}
	return p.serviceDirs
}

func (p *paths) checkServiceDirs() {
	dirs, err := os.ReadDir(p.executableDir)
	if err != nil {
		log.Fatalf("Unable to read directory: %v", err)
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		subDir := path.Join(p.executableDir, dir.Name())

		if !p.isServiceDir(subDir) {
			continue
		}

		p.serviceDirs = append(p.serviceDirs, subDir)
	}
}

func (p *paths) isServiceDir(subDir string) bool {
	_, err := os.Stat(p.getCredentialsFilePath(subDir))
	return err == nil || os.IsExist(err)
}

func (p *paths) getCredentialsFilePath(subDir string) string {
	return path.Join(subDir, pathCredentialsFilename)
}

func (p *paths) getTokenFilePath(subDir string) string {
	return path.Join(subDir, pathTokenFilename)
}

func (p *paths) getUnreadJSONFilePath(subDir string) string {
	return path.Join(subDir, pathUnreadFilename)
}
