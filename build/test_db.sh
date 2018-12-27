cd database
gem install bundle
bundle install
rake db:drop
rake db:create
rake db:migrate

# make sure reverting migrations is working
rake db:migrate VERSION=0

# run migrations again for api test
rake db:migrate
