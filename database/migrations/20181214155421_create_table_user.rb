class CreateTableUser < ActiveRecord::Migration[5.2]
  def up
    execute <<~SQL
      create table users (
        id serial4 primary key,
        username varchar(100) not null,
        encrypted_password varchar(200) not null,
        email varchar(500) not null,
        created_at timestamptz not null default now(),
        updated_at timestamptz,
        deleted_at timestamptz
      );
      create unique index unique_username_on_users on users(username);
      create unique index unique_email_on_users on users(email);
    SQL
  end

  def down
    execute <<~SQL
      drop table users;
    SQL
  end
end
