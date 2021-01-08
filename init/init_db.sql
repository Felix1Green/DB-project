CREATE EXTENSION IF NOT EXISTS citext;

create table users
(
    id serial not null primary key,
    nickname citext not null unique,
    fullname varchar(256) not null,
    about text,
    email citext not null unique
);

create table forum
(
    id serial not null primary key,
    title varchar(128) not null,
    user_id citext references users(nickname),
    slug citext not null unique
);

create table thread
(
    id serial not null primary key ,
    title varchar(128) not null,
    author citext references users(nickname),
    message text,
    forum citext references forum(slug),
    votes_counter int default 0,
    slug citext unique,
    created timestamp with time zone default now()
);

create table post
(
    id serial not null primary key ,
    parent int not null ,
    author citext references users(nickname),
    message text,
    isEdited boolean default false,
    forum citext references forum(slug),
    thread int references thread(id),
    created timestamp with time zone default now()
);

create table vote
(
    id serial not null primary key ,
    thread_id int not null,
    user_name citext not null references users(nickname),
    rating int not null,
    UNIQUE(thread_id, user_name)
);