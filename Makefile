all: push

TAG = 1.0
PREFIX = k8s-elector
GODEP ?= godep

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o server main.go

docker-build: build
	docker build --pull -t $(PREFIX):$(TAG) .

docker-push: docker-build
	gcloud docker -- push $(PREFIX):$(TAG)

save-vendor:
	rm -rf Godeps/ vendor/
	$(GODEP) save ./...

clean:
	rm -f server
	rm -rf Godeps/ vendor/
