class CreateTableUser < ActiveRecord::Migration[5.2]
  def up
    execute <<~SQL
      create table users (
        id serial4 primary key,
        username varchar(100) not null,
        encrypted_password varchar(200),
        github_username varchar(500),
        authentication_token varchar(100) not null,
        active bool not null default false,
        created_at timestamptz not null default now(),
        updated_at timestamptz,
        deleted_at timestamptz
      );
      create unique index unique_username_on_users on users(username);
      create unique index unique_github_username_on_users on users(github_username);
      create unique index unique_authentication_token_on_users on users(authentication_token);

      comment on column users.active is 'true if user has clicked activation url';
      comment on column users.encrypted_password is 'can be nil if user created the account using social login';
    SQL
  end

  def down
    execute <<~SQL
      drop table users;
    SQL
  end
end
