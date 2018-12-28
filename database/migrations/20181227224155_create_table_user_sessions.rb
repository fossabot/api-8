class CreateTableUserSessions < ActiveRecord::Migration[5.2]
  def up
    execute <<~SQL
      create table user_sessions (
        user_id integer not null references users(id),
        token varchar(300) not null,
        user_agent varchar(5000),
        ip inet,
        created_at timestamptz not null default now(),
        updated_at timestamptz,
        deleted_at timestamptz
      );
      create unique index unique_token_on_user_sessions on user_sessions(token) where deleted_at is not null;
    SQL
  end

  def down
    execute <<~SQL
      drop table user_sessions;
    SQL
  end
end
