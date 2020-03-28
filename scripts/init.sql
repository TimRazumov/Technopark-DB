create extension if not exists citext;

drop table if exists users cascade;
drop table if exists forums cascade;
drop table if exists user_forum cascade;
drop table if exists threads cascade;
drop table if exists posts cascade;
drop table if exists vote cascade;

create unlogged table users
(
    nickname citext    not null primary key,
    fullname varchar   not null,
    email    citext    not null unique,
    about    text
);

create unlogged table forums
(
    slug    citext    not null primary key,
    title   varchar   not null,
    usr     citext    not null references users (nickname),
    posts   integer   not null default 0,
    threads integer   not null default 0
);

create unlogged table user_forum
(
    user_id     integer,
    forum_slug  citext,
    primary key (forum_slug, user_id)
);

create unlogged table threads
(
    id         serial      not null primary key,
    title      varchar     not null,
    author     citext      not null references users (nickname),
    forum      citext      not null references forums (slug),
    message    text,
    votes      integer     default 0,
    slug       citext      default null unique,   
    created    timestamptz default current_timestamp
);

create unlogged table posts
(
    id         integer        not null primary key,
    parent     integer        not null default 0,
    author     citext         not null references users (nickname),
    message    text,
    is_edited  bool           not null default false,
    forum      citext         not null references forums (slug),
    thread     integer        not null references threads (id),
    created    timestamptz    default current_timestamp,
    path       integer array  default array[]::integer[]
);
