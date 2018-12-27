test: test-db test-api

test-db:
	cd database \
		&& rbenv exec rake db:drop \
		&& rbenv exec rake db:create \
		&& rbenv exec rake db:migrate \
		&& rbenv exec rake db:structure:dump \
		&& rbenv exec rake db:migrate VERSION=0

test-api:
	cd api \
		&& go build -o out/devlover-api github.com/devlover-id/api/cmd/devlover-api \
		&& go test -v -cover ./...

run-api:
	cd api && go build -o out/devlover-api github.com/devlover-id/api/cmd/devlover-api
	./api/out/devlover-api
