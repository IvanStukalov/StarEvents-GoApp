--sudo docker run --name postgresql -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=12345  -p 5432:5432 -d postgres

drop table if exists "events" CASCADE;
drop table if exists "stars" CASCADE;
drop table if exists "star_events" CASCADE;
drop table if exists "users" CASCADE;

create table stars
(
    id     integer not null
        constraint star_pk
            primary key,
    name        varchar(30),
    description varchar(200),
    distance    real,
    age    real,
    magnitude   real,
    image       varchar(30),
    is_active   varchar(20)
);

alter table stars
    owner to postgres;

create table "users"
(
    user_id integer not null
        constraint user_pk
            primary key,
    name    varchar(50)
);

alter table "users"
    owner to postgres;

create table events
(
    event_id        integer not null
        constraint event_pk
            primary key,
    name            varchar(50),
    status          varchar(20),
    creation_date   timestamp,
    formation_date  timestamp,
    completion_date integer,
    moderator_id    integer
        constraint event_user_user_id_fk
            references "users"
);

alter table events
    owner to postgres;

create table star_events
(
    star_event_id integer not null
        constraint star_event_pk
            primary key,
    star_id       integer
        constraint star_event_star_star_id_fk
            references stars,
    event_id      integer
        constraint star_event_event_event_id_fk
            references events
);

alter table star_events
    owner to postgres;

