class CreateTablePosts < ActiveRecord::Migration[5.2]
  def up
    execute <<~SQL
      create table posts (
        id serial4 primary key,
        reply_to int4 references posts(id),
        creator_id int4 not null references users(id),
        type varchar(50) not null,
        payload jsonb not null,
        created_at timestamptz not null default now(),
        updated_at timestamptz,
        deleted_at timestamptz
      );
      create index active_posts_load_order on posts(id desc, created_at desc, updated_at desc nulls last) where deleted_at is null;
    SQL
  end

  def down
    execute <<~SQL
      drop table posts;
    SQL
  end
end
