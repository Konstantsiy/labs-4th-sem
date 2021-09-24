select 'create database postgres'
    where not exists(select from pg_database where datname = 'postgres');
\gexec

\c postgres

create schema if not exists cinema;

do $$
begin
        if not exists(select 1 from pg_extension where extname = 'uuid-ossp') then
            create extension "uuid-ossp";
end if;
end
$$;

do $$
begin
        if not exists(select 1 from pg_type where typname = 'hall_sector') then
            create type hall_sector as enum ('near the screen', 'center', 'balcony');
end if;
end
$$;

do $$
begin
        if not exists(select 1 from pg_type where typname = 'hall_type') then
            create type hall_type as enum ('2D', '3D', 'IMAX');
end if;
end
$$;

create table if not exists cinema.genres (
    id uuid default uuid_generate_v1() primary key,
    title varchar(80) not null
    );

create table if not exists cinema.films (
    id uuid default uuid_generate_v1() primary key,
    title varchar(120) not null unique,
    duration interval not null check (interval '40 minutes' < duration and duration < interval '3 hours 30 minutes'),
    rental_start_date date not null check ( rental_start_date > current_date ),
    rental_end_date date not null check ( rental_end_date > rental_start_date )
    );

create table if not exists cinema.films_genres ( -- todo
    film_id uuid references cinema.films(id) on delete cascade,
    genre_id uuid references cinema.genres(id) on delete cascade,
    constraint film_genre_pkey primary key (film_id, genre_id)
    );

create table if not exists cinema.positions (
    id uuid default uuid_generate_v1() primary key,
    title varchar(120) not null unique
    );

create table if not exists cinema.workers (
    id uuid default uuid_generate_v1() primary key,
    position_id uuid not null,
    name varchar(45) not null,
    surname varchar(45) not null,
    passport_number varchar(80) not null unique,
    foreign key (position_id) references cinema.positions(id) on delete cascade
    );

create table if not exists cinema.halls (
    id uuid default uuid_generate_v1() primary key,
    number integer not null unique,
    type hall_type not null
    );

create table if not exists cinema.halls_workers (
    hall_id uuid references cinema.halls(id) on delete cascade,
    worker_id uuid references cinema.workers(id) on delete cascade,
    sector hall_sector not null,
    constraint hall_worker primary key (hall_id, worker_id)
    );

create table if not exists cinema.places (
    id uuid default uuid_generate_v1() primary key,
    row_number integer not null,
    place_number integer not null
    );

create table if not exists cinema.halls_places (
    hall_id uuid references cinema.halls(id) on delete cascade,
    place_id uuid references cinema.places(id) on delete cascade,
    constraint hall_place primary key (hall_id, place_id)
    );

create table if not exists cinema.sessions (
    id uuid default uuid_generate_v1() primary key,
    film_id uuid not null,
    hall_id uuid not null,
    date date not null check ( date > current_date ),
    time time not null check ( time < time '23:00' and time > time '10:00'),
    foreign key (film_id) references cinema.films(id) on delete cascade,
    foreign key (hall_id) references cinema.halls(id) on delete cascade
    );

create table if not exists cinema.tickets (
    id uuid default uuid_generate_v1() primary key,
    session_id uuid not null,
    price numeric(8, 2) not null check ( price > 10 and price < 60 ),
    foreign key (session_id) references cinema.sessions(id) on delete cascade
    );

create table if not exists cinema.tickets_places (
    ticket_id uuid references cinema.tickets(id) on delete cascade,
    place_id uuid references cinema.places(id) on delete cascade,
    constraint ticket_place primary key (ticket_id, place_id)
    )