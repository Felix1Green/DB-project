ALTER SYSTEM SET checkpoint_completion_target = '0.9';
ALTER SYSTEM SET wal_buffers = '6912kB';
ALTER SYSTEM SET default_statistics_target = '100';
ALTER SYSTEM SET effective_io_concurrency = '200';
ALTER SYSTEM SET max_worker_processes = '4';
ALTER SYSTEM SET max_parallel_workers_per_gather = '2';
ALTER SYSTEM SET max_parallel_workers = '4';
ALTER SYSTEM SET max_parallel_maintenance_workers = '2';
ALTER SYSTEM SET random_page_cost = '0.1';
ALTER SYSTEM SET seq_page_cost = '0.1';

CREATE EXTENSION IF NOT EXISTS citext;

create unlogged table users
(
    id serial not null primary key,
    nickname citext not null unique,
    fullname varchar(256) not null,
    about text,
    email citext not null unique
);

CREATE INDEX index_get_users_info on users(nickname, fullname, about, email);

create unlogged table forum
(
    id serial not null primary key,
    title varchar(128) not null,
    user_id citext not null,
    slug citext not null unique,
    threads int not null default 0,
    posts int not null default 0
);

CREATE INDEX index_forum_user_fk on forum(user_id);
CREATE INDEX index_forum_info on forum(slug, title, user_id);

create unlogged table thread
(
    id serial not null primary key ,
    title varchar(128) not null,
    author citext not null,
    message text,
    forum citext not null,
    votes_counter int default 0,
    slug citext unique,
    created timestamp with time zone default now()
);
CREATE INDEX index_thread_forum_fk on thread(forum);
CREATE INDEX index_thread_info on thread(forum, created);


create unlogged table post
(
    id serial primary key ,
    parent int not null ,
    author citext not null,
    message text,
    isEdited boolean default false,
    forum citext not null,
    thread int,
    created timestamp with time zone default now(),
    path int [] default '{0}':: int []
);

CREATE INDEX index_post_thread_fk on post(id, thread);
CREATE INDEX index_post_thread_path on post(thread, path);
CREATE INDEX index_post_thread_first_path on post((path[1]), path);

create unlogged table vote
(
    id serial not null primary key ,
    thread_id int not null,
    user_name citext not null references users(nickname),
    rating int not null,
    UNIQUE(thread_id, user_name)
);


create unlogged table forum_users(
    id serial primary key,
    forum citext not null references forum(slug),
    user_nickname citext not null references users(nickname),
    unique(forum, user_nickname)
);


create or replace function update_forum_threads_counter()
returns trigger as $update_forum_threads_counter$
BEGIN UPDATE forum set threads = threads + 1 where slug = new.forum;
insert into forum_users(forum, user_nickname) values(new.forum, new.author) on conflict (forum, user_nickname) DO NOTHING;
return new;
end;
$update_forum_threads_counter$ language plpgsql;

create trigger update_forum_threads_counter
    before insert on thread
    for each row execute procedure update_forum_threads_counter();


create or replace function update_forum_posts_counter()
returns trigger as $update_forum_posts_counter$
    begin update forum set posts = posts + 1 where slug = new.forum;
    insert into forum_users(forum, user_nickname) values(new.forum, new.author) on conflict (forum, user_nickname) DO NOTHING ;
    return new;
end;
$update_forum_posts_counter$ language plpgsql;

create trigger update_forum_posts_counter
    before insert on post
    for each row execute procedure update_forum_posts_counter();

-- insert vote for thread trigger

create or replace function insert_thread_vote()
returns trigger as $insert_thread_vote$
    begin update thread set votes_counter = (votes_counter + new.rating) where id = new.thread_id;
    return new;
    end;
$insert_thread_vote$ language plpgsql;

create trigger insert_thread_vote
    before insert on vote
    for each row
    execute procedure insert_thread_vote();

-- update vote for thread trigger

create or replace function update_thread_vote()
returns trigger as $update_thread_vote$
    begin
        update thread set votes_counter = (select sum(rating) from vote where thread_id = new.thread_id) where id = new.thread_id;
        return new;
end;
$update_thread_vote$ language plpgsql;

create trigger update_thread_vote
    after update on vote
    for each row
    execute procedure update_thread_vote();



CREATE OR REPLACE FUNCTION update_path()
    RETURNS TRIGGER AS
$BODY$
DECLARE
    parent_path         INT[];
    first_parent_thread INT;
BEGIN
    IF (new.parent = 0) THEN
        NEW.path := ARRAY(select new.id::integer);
    ELSE
        SELECT thread, path
        FROM post
        WHERE thread = NEW.thread AND id = NEW.parent
        INTO first_parent_thread, parent_path;
        IF NOT FOUND THEN
            RAISE EXCEPTION 'Parent not exists' USING ERRCODE = '00404';
        END IF ;
        NEW.path := parent_path || NEW.id;
    END IF;
    RETURN NEW;
END;
$BODY$ LANGUAGE plpgsql;

CREATE TRIGGER path_updater
    BEFORE INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE update_path();