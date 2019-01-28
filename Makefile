# Build the binaries.
build:
	mkdir -p build
	cd ./cmd/user; CGO_ENABLED=0 go build -o user-svc
	mv ./cmd/user/user-svc ./build
	cd ./cmd/task; CGO_ENABLED=0 go build -o task-svc
	mv ./cmd/task/task-svc ./build

# Create the docker images.
docker: build
	docker build -t stratumn/elk-go-user -f Dockerfile.user .
	docker build -t stratumn/elk-go-task -f Dockerfile.task .
