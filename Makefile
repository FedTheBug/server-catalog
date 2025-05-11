build-run:
	docker-compose up -d
	@echo "Preparing the environment..."
	sleep 5
	go mod vendor -v
	go build -v .
	./server-catalog migration down
	./server-catalog migration up
	@echo "Killing any existing server process on port 8080..."
	-lsof -ti:8080 | xargs kill -9 2>/dev/null || true
	sleep 2 &&  echo "Uploading catalog file..." && curl --location 'http://127.0.0.1:8080/api/v1/upload' \
    --header 'app-key: PPTjT3ApHD' \
    --form 'file=@"./servers_filters_assignment.xlsx"' && echo "Data loaded successfully...." &
	@echo "Starting server in background..."
	./server-catalog serve

clean:
	@echo "Cleaning up..."
	docker-compose down
	rm -f server-catalog
	@echo "Cleanup complete"