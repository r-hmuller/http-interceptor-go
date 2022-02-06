package checkpoint

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/google/uuid"
	"httpInterceptor/config"
	"os"
	"strconv"
	"time"
)

func Monitor() {
	period := config.GetSnapshotPeriodicity()
	periodicity := time.Duration(period) * time.Millisecond
	for _ = range time.Tick(periodicity) {
		generateSnapshot()
	}
}

func generateSnapshot() {
	startTime := time.Now().Unix()
	// create a new client connected to the default socket path for containerd
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// create a new context with an "example" namespace
	ctx := namespaces.WithNamespace(context.Background(), "default")

	containers, err := client.Containers(ctx)
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		labels, _ := container.Labels(ctx)
		if labels["io.kubernetes.container.name"] == config.GetPodName() {
			task, err := container.Task(ctx, cio.NewAttach(cio.WithStdio))
			if err != nil {
				panic(err)
			}
			checkpoint, err := task.Checkpoint(ctx)
			if err != nil {
				panic(err)
			}

			registry := config.GetRegistry()
			containerSnapshotVersion := registry + ":" + uuid.NewString()
			err = client.Push(ctx, containerSnapshotVersion, checkpoint.Target())
			if err != nil {
				panic(err)
			}
		}
	}

	endTime := time.Now().Unix()
	deltaTime := endTime - startTime
	loggingPath := config.GetLogginPath()
	f, err := os.OpenFile(loggingPath+"/checkpoint.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
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
