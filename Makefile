test-db:
	cd database \
		&& rbenv exec rake db:drop \
		&& rbenv exec rake db:create \
		&& rbenv exec rake db:migrate \
		&& rbenv exec rake db:migrate VERSION=0
