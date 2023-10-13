drop table if exists "events" CASCADE;
drop table if exists "stars" CASCADE;
drop table if exists "star_events" CASCADE;
drop table if exists "users" CASCADE;

create table stars (
    star_id serial constraint star_pk primary key,
    name varchar(30),
    description varchar(200),
    distance real,
    age real,
    magnitude real,
    image varchar(30),
    is_active varchar(20)
);

alter table
    stars owner to postgres;

create table "users" (
    user_id integer not null constraint user_pk primary key,
    name varchar(50),
    is_moderator boolean
);

alter table
    "users" owner to postgres;

create table events (
    event_id integer not null constraint event_pk primary key,
    name varchar(50),
    status varchar(20),
    creation_date timestamp,
    formation_date timestamp,
    completion_date timestamp,
    moderator_id integer constraint event_moderator_user_id_fk references "users",
    creator_id integer constraint event_creator_user_id_fk references "users"
);

alter table
    events owner to postgres;

create table star_events (
    star_id integer constraint star_event_star_star_id_fk references stars,
    event_id integer constraint star_event_event_event_id_fk references events,
    primary key (star_id, event_id)
);

alter table
    star_events owner to postgres;