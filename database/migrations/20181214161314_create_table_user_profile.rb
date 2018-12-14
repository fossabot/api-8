class CreateTableUserProfile < ActiveRecord::Migration[5.2]
  def up
    execute <<~SQL
      create table user_profile (
        user_id int4 not null references users(id),
        name varchar(500) not null,
        address varchar(2000) not null,
        phone varchar(50) not null,
        github_username varchar(500),
        updated_at timestamptz default now()
      );
      create unique index unique_user_id_on_users on user_profile(user_id);
      create unique index unique_phone_on_users on user_profile(phone);
      create unique index unique_github_username_on_users on user_profile(github_username);
    SQL
  end

  def down
    execute <<~SQL
      drop table user_profile;
    SQL
  end
end
