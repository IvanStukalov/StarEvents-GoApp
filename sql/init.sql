DROP TABLE IF EXISTS "events" CASCADE;

DROP TABLE IF EXISTS "stars" CASCADE;

DROP TABLE IF EXISTS "star_events" CASCADE;

DROP TABLE IF EXISTS "users" CASCADE;

CREATE TABLE stars(
    star_id serial NOT NULL CONSTRAINT star_pk PRIMARY KEY,
    name varchar(30) NOT NULL UNIQUE,
    description varchar(200),
    distance real,
    age real,
    magnitude real,
    image varchar(100),
    is_active boolean
);

ALTER TABLE stars OWNER TO postgres;

CREATE TABLE "users"(
    user_id serial NOT NULL CONSTRAINT user_pk PRIMARY KEY,
    login varchar(50),
    password varchar(200),
    is_moderator boolean
);

ALTER TABLE "users" OWNER TO postgres;

CREATE TABLE events(
    event_id serial NOT NULL CONSTRAINT event_pk PRIMARY KEY,
    name varchar(50),
    status varchar(20),
    creation_date timestamp,
    formation_date timestamp,
    completion_date timestamp,
    moderator_id integer,
    creator_id integer CONSTRAINT creator_id_fk REFERENCES "users"(user_id),
    scanned_percent integer
);

ALTER TABLE events OWNER TO postgres;

CREATE TABLE star_events(
    star_id integer CONSTRAINT star_event_star_star_id_fk REFERENCES stars(star_id) ON DELETE CASCADE,
    event_id integer CONSTRAINT star_event_event_event_id_fk REFERENCES events(event_id) ON DELETE CASCADE,
    PRIMARY KEY (star_id, event_id)
);

ALTER TABLE star_events OWNER TO postgres;

