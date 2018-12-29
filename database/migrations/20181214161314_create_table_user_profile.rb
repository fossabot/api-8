class CreateTableUserProfile < ActiveRecord::Migration[5.2]
  def up
    execute <<~SQL
      create table user_profile (
        user_id integer not null references users(id),
        activation_token varchar(100) not null,
        activation_token_expires_at timestamptz not null default now() + interval '15 minutes',
        name varchar(500) not null,
        email varchar(500),
        address varchar(2000),
        phone varchar(50)
      );
      create unique index unique_user_id_on_users on user_profile(user_id);
      create unique index unique_email_on_users on user_profile(email);
      create unique index unique_phone_on_users on user_profile(phone);
      create unique index unique_activation_token_on_users on user_profile(activation_token);
    SQL
  end

  def down
    execute <<~SQL
      drop table user_profile;
    SQL
  end
end
