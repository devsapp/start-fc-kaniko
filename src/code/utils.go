package main

import (
	"archive/zip"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/cavaliergopher/grab/v3"
	"github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type AcrUser struct {
	User     string `json:"user"`
	PassWord string `json:"passWord"`
}

func SetAcrAuth(builderEvent *BuilderEvent, logger *logrus.Entry) error {
	li := strings.Split(builderEvent.Image, ".")
	region := li[1]
	acrUser := AcrUser{}
	if builderEvent.InstanceID != nil { // acree
		if len(li) != 5 {
			msg := fmt.Sprintf("Invalid acr image %s", builderEvent.Image)
			logger.Infof(msg)
			panic(errors.New(msg))
		}
		resp, _err := GetAcrEETmpUser(builderEvent.Credentials, region, *builderEvent.InstanceID)
		if _err != nil {
			return _err
		}
		acrUser.User = *resp.TempUsername
		acrUser.PassWord = *resp.AuthorizationToken
	} else { // acr
		if len(li) != 4 {
			msg := fmt.Sprintf("Invalid acree image %s", builderEvent.Image)
			logger.Infof(msg)
			panic(errors.New(msg))
		}
		resp, _err := GetAcrTmpUser(builderEvent.Credentials, region)
		if _err != nil {
			return _err
		}
		acrUser.User = resp.Data.TempUserName
		acrUser.PassWord = resp.Data.AuthorizationToken
	}
	li = strings.Split(builderEvent.Image, "/")
	registry := li[0]
	authStr := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", acrUser.User, acrUser.PassWord)))
	dockerAuthStr := fmt.Sprintf(`{"auths":{"%s":{"auth":"%s"}}}`, registry, authStr)
	logger.Infof("docker auth = %s", dockerAuthStr)
	// overwrite
	err := ioutil.WriteFile("/kaniko/.docker/config.json", []byte(dockerAuthStr), 0664)
	if err != nil {
		return err
	}
	return nil
}

func CloneToWorkSpace(builderEvent *BuilderEvent, logger *logrus.Entry) error {
	if builderEvent.Type == "git" {
		li := strings.Split(builderEvent.Url, "/")
		if len(li) == 0 {
			msg := fmt.Sprintf("Invalid git url %s", builderEvent.Url)
			logger.Error(msg)
			return errors.New(msg)
		}
		_, err := git.PlainClone(WorkSpaceDir, false, &git.CloneOptions{
			URL:      builderEvent.Url,
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}
		logger.Infof("git clone  saved to %s", WorkSpaceDir)
	} else {
		resp, err := grab.Get(WorkSpaceDir, builderEvent.Url)
		if err != nil {
			logger.Errorf("download failed. Url: %s Error: %v", builderEvent.Url, err)
			return err
		}
		logger.Infof("Download saved to %s", resp.Filename)
		err = unzipSource(resp.Filename, WorkSpaceDir)
		if err != nil {
			logger.Errorf("unzip %s failed. Error: %v", path.Join(WorkSpaceDir, resp.Filename), err)
			return err
		}
	}
	return nil
}

func unzipSource(source, destination string) error {
	// 1. Open the zip file
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 2. Get the absolute destination path
	destination, err = filepath.Abs(destination)
	if err != nil {
		return err
	}

	// 3. Iterate over zip files inside the archive and unzip each of them
	for _, f := range reader.File {
		err := unzipFile(f, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

func unzipFile(f *zip.File, destination string) error {
	// 4. Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// 5. Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 6. Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 7. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}
