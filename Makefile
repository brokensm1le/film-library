.PHONY: buildrun
buildrun:
	docker-compose build
	docker-compose up -d

.PHONY: stop
stop:
	docker-compose down

.PHONY: gen
gen:
	mockgen -source=internal/auth/repository.go \
	-destination=internal/auth/mocks/mock_repository.go
	mockgen -source=internal/auth/usecase.go \
    	-destination=internal/auth/mocks/mock_usecase.go
	mockgen -source=internal/service/repository.go \
    	-destination=internal/service/mocks/mock_repository.go
	mockgen -source=internal/service/usecase.go \
    	-destination=internal/service/mocks/mock_usecase.go

.PHONY: cover
cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: test
test:
	go test ./...