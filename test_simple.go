package main

import (
	"os"
	"fmt"
	"bytes"
	"path/filepath"
	"github.com/fsouza/go-dockerclient"
)

func main() {
	argsWithoutProg := os.Args[1:]
	endpoint := "unix:///var/run/docker.sock"
	tlsCertPath := "/etc/dockercerts"
	if (len(argsWithoutProg) > 0) {
		endpoint = argsWithoutProg[0]
	}
	var client *docker.Client
	if (endpoint == "unix:///var/run/docker.sock") {
		dockerClient, err := docker.NewClient(endpoint)
		if (err != nil) {
			panic(err)
		}
		client = dockerClient;
	} else {
		dockerClient, err := docker.NewVersionedTLSClient(
			endpoint,
			filepath.Join(tlsCertPath, "cert.pem"),
			filepath.Join(tlsCertPath, "key.pem"),
			filepath.Join(tlsCertPath, "ca.pem"),
			"1.24",
		)
		if (err != nil) {
			panic(err)
		}
		client = dockerClient;
	}
	fmt.Println("Pulling Image");
	pullError := client.PullImage(docker.PullImageOptions{
	    Repository: "library/ubuntu",
	    Tag: "16.10",
	}, docker.AuthConfiguration{
		ServerAddress: "docker.io",
	})
	if (pullError != nil) {
		fmt.Println("Error Pulling Image");
		panic(pullError)
	}
	fmt.Println("Image Pulled - Building Container");
	container, err := client.CreateContainer(docker.CreateContainerOptions{
	    Config: &docker.Config{
	        Image:     "ubuntu:16.10",
	        OpenStdin: true,
	        Tty:		false,
	        StdinOnce: true,
	    },
	    HostConfig: &docker.HostConfig{
			RestartPolicy: docker.NeverRestart(),
			LogConfig: docker.LogConfig{
				Type: "json-file",
			},
		},
	})
	if (err != nil) {
		fmt.Println("Error Creating Container");
		panic(err)
	}
	containerId := container.ID
	fmt.Println("Built Container with ID", containerId);
	fmt.Println("Starting Container");
	err = client.StartContainer(containerId, nil)
	if (err != nil) {
		fmt.Println("Error Starting Container");
		panic(err)
	}
	fmt.Println("Container started - Attaching");
	err = client.AttachToContainer(docker.AttachToContainerOptions{
	    Container:   containerId,
	    InputStream: bytes.NewBufferString("echo \"Working\";"),
	    Stdin:       true,
	    Stdout:		 true,
	    Stderr:		 true,
	    Stream:      true,
	    OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		RawTerminal: false,
	})
	if (err != nil) {
		fmt.Println("Error Attaching to Container");
		panic(err)
	}
	fmt.Println("Attachment complete - Waiting for Container to Exit");
	exitCode, err := client.WaitContainer(containerId)
	if (err != nil) {
		fmt.Println("Error Waiting for Container to Exit");
		panic(err)
	}
	fmt.Println("Exit Code (which should be 0)", exitCode);
	fmt.Println("Removing container");
	err = client.RemoveContainer(docker.RemoveContainerOptions{
		ID: containerId,
		RemoveVolumes: true,
		Force: true,
	})
	if (err != nil) {
		fmt.Println("Error Removing Container", err);
	}
	return
}