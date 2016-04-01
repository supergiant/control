build:
	GOOS=linux GOARCH=amd64 go build -o supergiant-api .
	docker login -e $(DOCKER_EMAIL) -u $(DOCKER_USER) -p $(DOCKER_PASS)
	docker build -f Dockerfile -t $(REPO):$(TAG) .
	docker push $(REPO)
