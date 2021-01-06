create table users
(
    id serial not null primary key,
    nickname varchar(64) not null unique,
    fullname varchar(128) not null,
    about text,
    email varchar(64) not null unique
);

create table forum
(
    id serial not null primary key,
    title varchar(64) not null,
    user_id varchar(64) references users(nickname),
    slug varchar(64) not null unique
);

create table thread
(
    id serial not null primary key ,
    title varchar(64) not null,
    author varchar(64) references users(nickname),
    message text,
    forum varchar(64) references forum(slug),
    votes_counter int default 0,
    created timestamp default now()
);

create table post
(
    id serial not null primary key ,
    parent int not null ,
    author varchar(64) references users(nickname),
    message text,
    isEdited boolean default false,
    forum varchar(64) references forum(slug),
    thread int references thread(id),
    created timestamp default now()
);

create table vote
(
    id serial not null primary key ,
    thread_id int not null,
    user_name varchar(64) not null references users(nickname),
    rating int not null,
    UNIQUE(thread_id, user_name)
);