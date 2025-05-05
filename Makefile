build-run:
	docker-compose up -d
	go build -v .
	./server-catalog serve