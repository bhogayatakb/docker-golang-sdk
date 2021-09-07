package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// pass docker image address
	imageName := "bfirsh/reticulate-splines"

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println("resp ID ", resp.ID)

	stats, errInFetchingStats := cli.ContainerStats(ctx, resp.ID, false)

	if errInFetchingStats != nil {
		panic(errInFetchingStats)
	}

	fmt.Println("stats: => ", stats.Body)
	fmt.Printf("stats type: %T => ", stats.Body, "\n")

	file, err := os.Create("./response.json")

	defer file.Close()

	var buf bytes.Buffer
	tee := io.TeeReader(stats.Body, &buf)

	io.Copy(file, tee)

	// query := url.Values{}
	// query.Set("stream", "0")
	// if true {
	// 	query.Set("stream", "1")
	// }

	// response, error := cli.get(ctx, "/containers/"+resp.ID+"/stats", query, nil)
	// if error == nil {
	// 	fmt.Println("response: => ", response)
	// } else {
	// 	fmt.Println("error: => ", response)
	// }

}
