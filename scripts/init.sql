create extension if not exists citext;

drop table if exists users cascade;
drop table if exists forums cascade;
drop table if exists user_forum cascade;
drop table if exists threads cascade;
drop table if exists posts cascade;
drop table if exists votes cascade;

create unlogged table users
(
    nickname citext    not null primary key,
    fullname text      not null,
    email    citext    not null unique,
    about    text
);

create unlogged table forums
(
    slug    citext    not null primary key,
    title   text      not null,
    usr     citext    not null references users (nickname),
    posts   integer   not null default 0,
    threads integer   not null default 0
);

create unlogged table user_forum
(
    nickname  citext not null references users (nickname),
    slug      citext not null references forums (slug),
    unique    (nickname, slug)
);

create unlogged table threads
(
    id         serial      not null primary key,
    title      text        not null,
    author     citext      not null references users (nickname),
    forum      citext      not null references forums (slug),
    message    text,
    votes      integer     default 0,
    slug       citext      default null unique,   
    created    timestamptz default current_timestamp
);

create unlogged table posts
(
    id         serial         not null primary key,
    parent     integer        not null default 0,
    author     citext         not null references users (nickname),
    message    text,
    is_edited  bool           not null default false,
    forum      citext         not null references forums (slug),
    thread     integer        not null references threads (id),
    created    timestamptz    default current_timestamp,
    path       integer[]      default array[]::integer[]
);

create unlogged table votes
(
    nickname citext   not null references users (nickname),
    voice    smallint check (voice in (-1, 1)),
    thread   integer  not null references threads (id),
    unique (nickname, thread)
);
