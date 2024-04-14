APP_NAME?=my-app

# сборка отдельного приложения
clean:
	rm -f ${APP_NAME}

build: clean
	go build -mod vendor -o ${APP_NAME} ./cmd/service/service.go 

run: build
	./${APP_NAME}

.PHONY: shutdown
shutdown:
	docker-compose down -v

.PHONY: test
test:
	go test -count=1 -timeout 10s ./...

apply-migrations:
	docker build -t migrator ./db
	sleep 20
	docker run --network host migrator  \
	-path=/migrations/ \
	-database "postgresql://postgres:postgres@localhost:7557/urls?sslmode=disable" up

build-docker:
	STORAGE=$(STORAGE) docker-compose up --build
	 	
	