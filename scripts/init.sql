create extension if not exists citext;

drop table if exists users cascade;
drop table if exists forums cascade;
drop table if exists user_forum cascade;
drop table if exists threads cascade;
drop table if exists posts cascade;
drop table if exists votes cascade;

create table users
(
    nickname citext    not null primary key collate "C",
    fullname text      not null,
    email    citext    not null unique collate "C",
    about    text
);
create index if not exists idx_users_nickname on users (nickname);
create index if not exists idx_users_email on users (nickname);

create table forums
(
    slug    citext    not null primary key collate "C",
    title   text      not null,
    usr     citext    not null references users (nickname) collate "C",
    posts   integer   not null default 0,
    threads integer   not null default 0
);
create index if not exists idx_forum_user on forums (usr);

create table user_forum
(
    nickname  citext not null references users (nickname) collate "C",
    slug      citext not null references forums (slug) collate "C",
    unique    (nickname, slug)
);
create index if not exists idx_users_forum_nickname_slug on user_forum (nickname, slug);

create table threads
(
    id         serial      not null primary key,
    title      text        not null,
    author     citext      not null references users (nickname) collate "C",
    forum      citext      not null references forums (slug) collate "C",
    message    text,
    votes      integer     default 0,
    slug       citext      default null unique collate "C",
    created    timestamptz default current_timestamp
);
create index if not exists idx_threads_slug on threads (slug);
create index if not exists idx_threads_forum on threads (forum);
create index if not exists idx_threads_author on threads (author);

create table posts
(
    id         serial         not null primary key,
    parent     integer        not null default 0,
    author     citext         not null references users (nickname) collate "C",
    message    text,
    is_edited  bool           not null default false,
    forum      citext         not null references forums (slug) collate "C",
    thread     integer        not null references threads (id),
    created    timestamptz    default current_timestamp,
    path       integer[]      default array[]::integer[]
);
create index if not exists idx_posts_id on posts (id);
create index if not exists idx_posts_path on posts (path);
create index if not exists idx_posts_path_1 on posts ((path [1]));
create index if not exists idx_posts_thread on posts (thread);
create index if not exists idx_posts_author on posts (author);
create index if not exists idx_posts_forum on posts (forum);

create table votes
(
    nickname citext   not null references users (nickname) collate "C",
    voice    smallint check (voice in (-1, 1)),
    thread   integer  not null references threads (id),
    unique (nickname, thread)
);
create index if not exists idx_votes_nickname_thread on votes (nickname, thread);

create or replace function insert_user_forum() returns trigger as
$$
begin
    insert into user_forum (nickname, slug) values (new.author, new.forum)
    on conflict do nothing;
    return new;
end;
$$ language plpgsql;

create trigger insert_forum_user_trigger after insert on threads
for each row execute procedure insert_user_forum();

create or replace function thread_counter() returns trigger as
$$
begin
   update forums set threads = threads + 1 where slug = new.forum;
   return new;
end;
$$ language plpgsql;

create trigger thread_counter_trigger after insert on threads
for each row execute procedure thread_counter();
