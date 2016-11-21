# Docker Attach Test

## Description
This project is to demonstrate the different between the official Docker Engine implementation of docker attach and that of sdc-docker. This has been tested using the latest sdc-docker release in both JPC and an on-premise Triton deployment and a Docker for Mac install running 1.12.3 on stable channel.

## Running

* Customise the makefile (specifically DOCKERCERTPATH) to point to local (host) filesystem path containing certificates/keys for a remote Docker host
* Run ```make cmd``` which will construct the working image and then drop you into a shell

There are two invocations available:
* ```go run test_simple.go``` will cause the app to perform all Docker operations against the Docker Engine installation on your local machine
* ```go run test_simple.go tcp://us-west-1.docker.joyent.com:2376``` will cause all Docker operations to be performed against the specified remote API endpoint. The certificates located at ```DOCKERCERTPATH``` will be used.

The code will attempt to remove the created container, but given that operation is performed at the end of the code it is not currently functional when performed against sdc-docker.

## Design
The test is a golang program that creates a container with 
```go
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
```
The Dockerfile of the ```"ubuntu:16.10"``` image doesn't define an ENTRYPOINT, rather but defines CMD as bash.

The container is then started and attach is called as so:
```go
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
```

This should attach to the container, run the command in InputStream and then the container should exit with code 0. Changing ```echo``` to something that doesn't exist, like ```echoz``` should cause the container to exit with code 127. 

The expected output is
```bash
Pulling Image
Image Pulled - Building Container
Built Container with ID 743a8b6cf1abf83f4dea5f307b0a71360080e10c681cf722e96f88342ae3404f
Starting Container
Container started - Attaching
Working
Attachment complete - Waiting for Container to Exit
Exit Code (which should be 0) 0
Removing container
```

