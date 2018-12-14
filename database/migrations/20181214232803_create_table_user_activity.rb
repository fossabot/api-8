class CreateTableUserActivity < ActiveRecord::Migration[5.2]
  def up
    execute <<~SQL
      create table user_activity (
        id serial4 primary key,
        user_id integer not null references users(id),
        description varchar(1000) not null,
        created_at timestamptz not null default now()
      );
      create index created_at_on_user_activity on user_activity(created_at desc);
    SQL
  end

  def down
    execute <<~SQL
      drop table user_activity;
    SQL
  end
end
