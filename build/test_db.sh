cd database
gem install bundle
bundle install
rake db:drop
rake db:create
rake db:migrate
rake db:migrate VERSION=0

# run migration again for testing
rake db:migrate
