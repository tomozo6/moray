package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func SSMPortForwarding(c *ssm.Client, bastionID *string, dbHost *string, dbPort *int32, localPort *int32, region *string, profile *string) {
	sessionInput := &ssm.StartSessionInput{
		DocumentName: aws.String("AWS-StartPortForwardingSessionToRemoteHost"),
		Parameters: map[string][]string{
			"host":            {*dbHost},
			"portNumber":      {strconv.Itoa(int(*dbPort))},
			"localPortNumber": {strconv.Itoa(int(*localPort))},
		},
		Target: bastionID,
	}

	sessionOutput, err := c.StartSession(context.TODO(), sessionInput)
	if err != nil {
		panic(err.Error())
	}

	encodedSessionInput, err := json.Marshal(sessionInput)
	// if err != nil {
	// return nil, err
	// }
	encodedSessionOutput, err := json.Marshal(sessionOutput)
	// if err != nil {
	// return nil, err
	// }

	// session-manager-pluginを実行するコマンドを構築
	smpCmd := "session-manager-plugin"
	if runtime.GOOS == "windows" {
		smpCmd += ".exe"
	}
	cmd := exec.Command(smpCmd, string(encodedSessionOutput), *region, "StartSession", *profile, string(encodedSessionInput))

	// コマンドの実行には標準入出力を使う
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// コマンドを実行
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error running session-manager-plugin:", err)
		return
	}

	// コマンドの実行が終了するまで待機
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error waiting for session-manager-plugin to finish:", err)
		return
	}

	// 念のため、セッションを終了
	c.TerminateSession(context.TODO(), &ssm.TerminateSessionInput{
		SessionId: sessionOutput.SessionId,
	})
}
