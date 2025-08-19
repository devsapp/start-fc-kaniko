package main

import (
	"encoding/json"
	"fmt"
	"github.com/alibabacloud-go/tea/tea"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	gr "github.com/awesome-fc/golang-runtime"
)

const WorkSpaceDir = "/kaniko/workspace"

type AliyunCredentials struct {
	AccessKeyId     *string `json:"accessKeyId"`
	AccessKeySecret *string `json:"accessKeySecret"`
	SecurityToken   *string `json:"securityToken,omitempty"`
}

type BuilderEvent struct {
	Url         string            `json:"url"`
	Type        string            `json:"type"`
	Image       string            `json:"image"`
	InstanceID  *string           `json:"instanceId"`
	Dockerfile  string            `json:"dockerfile"`
	Credentials AliyunCredentials `json:"credentials"`
}

func handler(ctx *gr.FCContext, event []byte) ([]byte, error) {
	fcLogger := gr.GetLogger().WithField("requestId", ctx.RequestID)
	fcLogger.Info(string(event))
	var builderEvent BuilderEvent
	err := json.Unmarshal(event, &builderEvent)
	if err != nil {
		fcLogger.Error("json.Unmarshal error:", err)
		panic(err)
	}
	fcLogger.Info("start init workspace ...")
	err = CloneToWorkSpace(&builderEvent, fcLogger)
	if err != nil {
		panic(err)
	}
	fcLogger.Info("start set acr auth ...")
	builderEvent.Credentials = AliyunCredentials{
		AccessKeyId:     tea.String(ctx.Credentials.AccessKeyID),
		AccessKeySecret: tea.String(ctx.Credentials.AccessKeySecret),
		SecurityToken:   tea.String(ctx.Credentials.SecurityToken),
	}
	err = SetAcrAuth(&builderEvent, fcLogger)
	if err != nil {
		panic(err)
	}
	filename := path.Join(WorkSpaceDir, builderEvent.Dockerfile)
	contextDir := filepath.Dir(filename)
	cmdStr := fmt.Sprintf("executor --force=true --cache=false --use-new-run=true --dockerfile %s --context %s --destination %s",
		filename, contextDir, builderEvent.Image)
	fcLogger.Info("start build&push ...")
	fcLogger.Infof(`cmd str = "%s"`, cmdStr)
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	return event, nil
}

func main() {
	gr.Start(handler, nil)
}
