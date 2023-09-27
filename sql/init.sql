--sudo docker run --name postgresql -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=12345  -p 5432:5432 -d postgres

drop table if exists "event" CASCADE;
drop table if exists "star" CASCADE;
drop table if exists "star_event" CASCADE;
drop table if exists "user" CASCADE;

create table star
(
    star_id     integer not null
        constraint star_pk
            primary key,
    name        varchar(30),
    description varchar(200),
    distance    real,
    magnitude   real,
    image       varchar(30),
    is_active   varchar(20)
);

alter table star
    owner to postgres;

create table "user"
(
    user_id integer not null
        constraint user_pk
            primary key,
    name    varchar(50)
);

alter table "user"
    owner to postgres;

create table event
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
            references "user"
);

alter table event
    owner to postgres;

create table star_event
(
    star_event_id integer not null
        constraint star_event_pk
            primary key,
    star_id       integer
        constraint star_event_star_star_id_fk
            references star,
    event_id      integer
        constraint star_event_event_event_id_fk
            references event
);

alter table star_event
    owner to postgres;

