package checkpoint

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"httpInterceptor/config"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Restore() {
	startTime := time.Now().Unix()
	podName := config.GetPodName()
	stateManager := config.GetStateManagerUrl()

	resp, err := http.Get(stateManager + "/" + podName)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	latestSnapshot := string(body)

	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		panic(err)
	}
	defer client.Close()
	ctx := namespaces.WithNamespace(context.Background(), "default")
	registry := config.GetRegistry()
	containerSnapshotVersion := registry + ":" + latestSnapshot
	checkpoint, err := client.Pull(ctx, containerSnapshotVersion)

	application, err := client.NewContainer(ctx, podName, containerd.WithNewSnapshot("application-"+podName, checkpoint))
	if err != nil {
		panic(err)
	}

	task, err := application.NewTask(ctx, cio.NewCreator(cio.WithStdio), containerd.WithTaskCheckpoint(checkpoint))
	err = task.Start(ctx)

	endTime := time.Now().Unix()
	deltaTime := endTime - startTime
	loggingPath := config.GetLogginPath()
	f, err := os.OpenFile(loggingPath+"/restorer.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	if _, err = f.WriteString(strconv.FormatInt(deltaTime, 10)); err != nil {
		panic(err)
	}
}
