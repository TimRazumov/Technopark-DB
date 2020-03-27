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

create table user_forum
(
    user_id     integer,
    forum_slug  citext,
    primary key (forum_slug, user_id)
);

create index on user_forum (forum_slug);
