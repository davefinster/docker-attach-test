ROOT := $(shell pwd)
DOCKERCERTPATH := /path/to/certs

DOCKERRUN := docker run -it --rm --name="docker-attach-test" \
	-v ${ROOT}:/go/src/test \
	-w /go/src/test \
	--privileged \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v ${DOCKERCERTPATH}:/etc/dockercerts \
	test_work_image

cmd: work_image
	$(DOCKERRUN)

work_image:
	docker rmi -f test_work_image > /dev/null 2>&1 || true
	docker build -t test_work_image .
	docker inspect -f "{{ .ID }}" test_work_image > work_image

