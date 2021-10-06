include: .envrc

.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: run
run:
	go run ./cmd/api -db-dsn ${TODOAPP_DB_DSN} -smtp-username ${smtp-username} -smtp-password ${smtp-password}

.PHONY: db/migration/up
db/migration/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migration -database ${TODOAPP_DB_DSN} up

.PHONY: audit
audit:
	@echo 'Tidying and verify module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

current_time = $(shell date --iso-8601=seconds)
git_description = $(shell git describe --always --dirty --tags --long)
# linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'
linker_flags = '-s -X main.version=${git_description}'

.PHONY: build/api
build/api:
	echo 'Building cmd/api...'
	go build -ldflags=${linker_flags} -o=./bin/api.bin ./cmd/api
#	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/api.bin ./cmd/api
#	GOOS=windows GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/api.exe ./cmd/api

production_host_ip = ''

.PHONE: production/connect
production/connect:
	ssh todoapp@${production_host_ip}

.PHONY: production/deploy/api
production/deploy/api:
	rsync -rP --delete ./bin/api.bin ./migration todoapp@${production_host_ip}:~
	ssh -t todoapp@${production_host_ip} 'migrate -path ~/migration -database $$TODOAPP_DB_DSN up'

.PHONY: production/configure/api.service
production/configure/api.service:
	rsync -P ./remote/production/api.service todoapp@${production_host_ip}:~
	ssh -t todoapp@${production_host_ip} '\
		sudo mv ~/api.service /etc/systemd/system/ \
		&& sudo systemctl enable api \
		&& sudo systemctl restart api \
	'

.PHONY: production/configure/caddyfile
production/configure/caddyfile:
	rsync -P ./remote/production/Caddyfile greenlight@${production_host_ip}:~
	ssh -t greenlight@${production_host_ip} '\
		sudo mv ~/Caddyfile /etc/caddy/ \
		&& sudo systemctl reload caddy \
	'