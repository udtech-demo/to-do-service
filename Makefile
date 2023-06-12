dependents:
	docker-compose up -d
	go mod tidy

test:
	ginkgo ./...

run: dependents test
	go run cmd/main.go

run-compose:
	docker-compose up -d

stop-compose:
	docker-compose down -v

clean: stop-compose
	rm -rf ./rsa_keys
	rm -rf ./conf/config.yml
	rm -rf ./data
	rm -rf ./logs