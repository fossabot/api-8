run-migration:
	cd database \
		&& rbenv exec rake db:drop \
		&& rbenv exec rake db:create \
		&& rbenv exec rake db:migrate

run-api:
	cd api && go build -o out/devlover-api github.com/devlover-id/api/cmd/devlover-api
	./api/out/devlover-api

test-db: run-migration
	cd database && rbenv exec rake db:migrate VERSION=0

test-api: run-migration
	cd api \
		&& go build -o out/devlover-api github.com/devlover-id/api/cmd/devlover-api \
		&& go test -v -cover ./...
