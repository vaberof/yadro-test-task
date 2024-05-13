PROJECT_NAME=computer-club

DOCKER_LOCAL_IMAGE_NAME=$(PROJECT_NAME)

WORK_DIR_LINUX=./cmd/yadro-test-task

docker.run: docker.build
	docker run $(DOCKER_LOCAL_IMAGE_NAME) $$TEST_FILE

docker.build: build.linux
	docker build -t $(DOCKER_LOCAL_IMAGE_NAME) -f $(WORK_DIR_LINUX)/Dockerfile .

build.linux: build.clean
	mkdir -p $(WORK_DIR_LINUX)/build
	go build -o $(WORK_DIR_LINUX)/build/main $(WORK_DIR_LINUX)/*.go

build.local: build.clean
	mkdir -p $(WORK_DIR_LINUX)/build
	go build -o $(WORK_DIR_LINUX)/build/main $(WORK_DIR_LINUX)/*.go
	@echo "build.local: OK"

build.clean:
	rm -rf $(WORK_DIR_LINUX)/build

tests.run:
	go test -v ./internal/domain/computerclub