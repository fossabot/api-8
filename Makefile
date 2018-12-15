test-db:
	cd database \
		&& rbenv exec rake db:drop \
		&& rbenv exec rake db:create \
		&& rbenv exec rake db:migrate \
		&& rbenv exec rake db:migrate VERSION=0

run:
	cd api \
		&& go build -o out/pinterkode gitlab.com/pinterkode/pinterkode/api/cmd/pinterkode \
		&& ./out/pinterkode
