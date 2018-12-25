cd database
gem install bundle
bundle install
rake db:drop
rake db:create
rake db:migrate

# make sure structure.sql is updated
cp db/structure.sql db/structure.sql.old
rake db:structure:dump
diff db/structure.sql db/structure.sql.old

# make sure reverting migrations is working
rake db:migrate VERSION=0
