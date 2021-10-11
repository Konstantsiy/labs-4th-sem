-- select 'create database postgres'
-- where not exists(select from pg_database where datname = 'postgres');
-- \gexec
--
-- \c postgres

create schema if not exists cinema;

-- do $$
--     begin
--         if not exists(select 1 from pg_extension where extname = 'uuid-ossp') then
--             create extension "uuid-ossp" with schema cinema;
--         end if;
--     end
-- $$;

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

set intervalstyle = 'postgres';

create or replace function random_sequence10() returns varchar(10) as
$$
    declare
chars text[] := '{0,1,2,3,4,5,6,7,8,9,A,B,C,D,E,F,G,H,I,J,K,L,M,N}';
        result varchar(10) := '';
        i integer := 0;
begin
for i in 1..10 loop
                result := result || chars[1+random()*(array_length(chars, 1)-1)];
end loop;
return result;
end;
$$ language plpgsql;

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

create table if not exists cinema.films_genres (
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
    passport_number varchar(10) not null unique default random_sequence10(),
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
    sector hall_sector not null
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
    worker_id uuid not null,
    price numeric(8, 2) not null check ( price > 0 and price < 60 ),
    foreign key (session_id) references cinema.sessions(id) on delete cascade,
    foreign key (worker_id) references cinema.workers(id)
    );

create table if not exists cinema.tickets_places (
    ticket_id uuid references cinema.tickets(id) on delete cascade,
    place_id uuid references cinema.places(id) on delete cascade,
    constraint ticket_place primary key (ticket_id, place_id)
    );

-- --------------------------------------------------------------------------------
-- --------------------------------------------------------------------------------
insert into cinema.genres (title)
values ('horror'),
       ('fantasy'),
       ('action'),
       ('drama'),
       ('comedy'),
       ('thriller');
select * from cinema.genres;

insert into cinema.films (title, duration, rental_start_date, rental_end_date)
values ('365 Days', interval '118 minutes', date '2021-10-15', date '2021-10-30'),
       ('The Courier', interval '114 minutes', date '2021-10-16', date '2021-10-27'),
       ('Tenet', interval '150 minutes', date '2021-10-11', date '2021-10-19'),
       ('News of the World', interval '120 minutes', date '2021-10-16', date '2021-10-30'),
       ('Extraction', interval '116 minutes', date '2021-11-01', date '2021-11-15'),
       ('The Night House', interval '107 minutes', date '2021-11-03', date '2021-11-17');
select * from cinema.films;

insert into cinema.films_genres (film_id, genre_id)
values ('73e176be-2607-11ec-8491-61ceeb4939bd', 'a9e58242-2606-11ec-8491-61ceeb4939bd'),
       ('73e176be-2607-11ec-8491-61ceeb4939bd', 'a9e58247-2606-11ec-8491-61ceeb4939bd'),
       ('73e176bf-2607-11ec-8491-61ceeb4939bd', 'a9e58243-2606-11ec-8491-61ceeb4939bd'),
       ('73e176c0-2607-11ec-8491-61ceeb4939bd', 'a9e58246-2606-11ec-8491-61ceeb4939bd'),
       ('73e176c1-2607-11ec-8491-61ceeb4939bd', 'a9e58245-2606-11ec-8491-61ceeb4939bd'),
       ('73e176c2-2607-11ec-8491-61ceeb4939bd', 'a9e58244-2606-11ec-8491-61ceeb4939bd'),
       ('73e176c2-2607-11ec-8491-61ceeb4939bd', 'a9e58242-2606-11ec-8491-61ceeb4939bd'),
       ('73e176c3-2607-11ec-8491-61ceeb4939bd', 'a9e58244-2606-11ec-8491-61ceeb4939bd');
select * from cinema.films_genres;

insert into cinema.positions (title)
values ('cashier'),
       ('cleaner'),
       ('technical operator');
select * from cinema.positions;

insert into cinema.workers (position_id, name, surname, passport_number)
values ('d2f7c9ae-2608-11ec-8491-61ceeb4939bd', 'John', 'Wick', random_sequence10()),
       ('d2f7c9af-2608-11ec-8491-61ceeb4939bd', 'Aldo', 'Apache', random_sequence10()),
       ('d2f7c9b0-2608-11ec-8491-61ceeb4939bd', 'Tony', 'Stark', random_sequence10()),
       ('d2f7c9b0-2608-11ec-8491-61ceeb4939bd', 'James', 'Burton', random_sequence10());
select * from cinema.workers;

insert into cinema.workers (position_id, name, surname, passport_number)
values ('d2f7c9ae-2608-11ec-8491-61ceeb4939bd', 'Emily', 'Black', random_sequence10()),
       ('d2f7c9ae-2608-11ec-8491-61ceeb4939bd', 'Jojo', 'Rabbit', random_sequence10());
select * from cinema.workers;

insert into cinema.places (row_number, place_number)
values (1, 1), (1, 2), (1, 3),
       (2, 1), (2, 2), (2, 3),
       (3, 1), (3, 2), (3, 3);
select * from cinema.places;

-- drop table cinema.tickets cascade;

-- alter table cinema.films add country varchar(40) not null ;

insert into cinema.halls (number, type)
values (1, '2D'), (2, '3D'), (3, 'IMAX');
select * from cinema.halls;

-- insert into cinema.halls (number, type)
-- values (1, '2D'),
--        (2, '3D'),
--        (3, 'IMAX');
-- select * from cinema.halls;

insert into cinema.halls_workers (hall_id, worker_id, sector)
values ('1e004ba6-2609-11ec-8491-61ceeb4939bd', 'f8f4c059-2608-11ec-8491-61ceeb4939bd', 'near the screen'),
       ('1e004ba6-2609-11ec-8491-61ceeb4939bd', '7b99b226-260d-11ec-8491-61ceeb4939bd', 'center'),
       ('1e004ba6-2609-11ec-8491-61ceeb4939bd', '7b99b227-260d-11ec-8491-61ceeb4939bd', 'balcony'),
       ('1e004ba7-2609-11ec-8491-61ceeb4939bd', 'f8f4c05a-2608-11ec-8491-61ceeb4939bd', 'near the screen'),
       ('1e004ba7-2609-11ec-8491-61ceeb4939bd', 'f8f4c059-2608-11ec-8491-61ceeb4939bd', 'center'),
       ('1e004ba7-2609-11ec-8491-61ceeb4939bd', '7b99b226-260d-11ec-8491-61ceeb4939bd', 'balcony'),
       ('1e004ba8-2609-11ec-8491-61ceeb4939bd', '7b99b227-260d-11ec-8491-61ceeb4939bd', 'near the screen'),
       ('1e004ba8-2609-11ec-8491-61ceeb4939bd', 'f8f4c05b-2608-11ec-8491-61ceeb4939bd', 'center'),
       ('1e004ba8-2609-11ec-8491-61ceeb4939bd', 'f8f4c05b-2608-11ec-8491-61ceeb4939bd', 'balcony');
select * from cinema.halls_workers;

insert into cinema.halls_places (hall_id, place_id)
values ('1e004ba6-2609-11ec-8491-61ceeb4939bd', '0d5fa422-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba6-2609-11ec-8491-61ceeb4939bd', '0d5fa423-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba6-2609-11ec-8491-61ceeb4939bd', '0d5fa424-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba6-2609-11ec-8491-61ceeb4939bd', '0d5fa425-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba6-2609-11ec-8491-61ceeb4939bd', '0d5fa426-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba6-2609-11ec-8491-61ceeb4939bd', '0d5fa427-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba6-2609-11ec-8491-61ceeb4939bd', '0d5fa428-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba6-2609-11ec-8491-61ceeb4939bd', '0d5fa429-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba6-2609-11ec-8491-61ceeb4939bd', '0d5fa42a-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba7-2609-11ec-8491-61ceeb4939bd', '0d5fa422-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba7-2609-11ec-8491-61ceeb4939bd', '0d5fa423-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba7-2609-11ec-8491-61ceeb4939bd', '0d5fa424-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba7-2609-11ec-8491-61ceeb4939bd', '0d5fa425-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba7-2609-11ec-8491-61ceeb4939bd', '0d5fa426-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba7-2609-11ec-8491-61ceeb4939bd', '0d5fa427-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba7-2609-11ec-8491-61ceeb4939bd', '0d5fa428-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba7-2609-11ec-8491-61ceeb4939bd', '0d5fa429-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba7-2609-11ec-8491-61ceeb4939bd', '0d5fa42a-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba8-2609-11ec-8491-61ceeb4939bd', '0d5fa422-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba8-2609-11ec-8491-61ceeb4939bd', '0d5fa423-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba8-2609-11ec-8491-61ceeb4939bd', '0d5fa424-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba8-2609-11ec-8491-61ceeb4939bd', '0d5fa425-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba8-2609-11ec-8491-61ceeb4939bd', '0d5fa426-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba8-2609-11ec-8491-61ceeb4939bd', '0d5fa427-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba8-2609-11ec-8491-61ceeb4939bd', '0d5fa428-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba8-2609-11ec-8491-61ceeb4939bd', '0d5fa429-2609-11ec-8491-61ceeb4939bd'),
       ('1e004ba8-2609-11ec-8491-61ceeb4939bd', '0d5fa42a-2609-11ec-8491-61ceeb4939bd');
select * from cinema.halls_places;

insert into cinema.sessions (film_id, hall_id, date, time)
values ('73e176be-2607-11ec-8491-61ceeb4939bd', '1e004ba6-2609-11ec-8491-61ceeb4939bd', date '2021-10-17', time '14:30:00'),
       ('73e176be-2607-11ec-8491-61ceeb4939bd', '1e004ba7-2609-11ec-8491-61ceeb4939bd', date '2021-10-20', time '16:30:00'),
       ('73e176be-2607-11ec-8491-61ceeb4939bd', '1e004ba8-2609-11ec-8491-61ceeb4939bd', date '2021-10-25', time '19:00:00'),

       ('73e176bf-2607-11ec-8491-61ceeb4939bd', '1e004ba6-2609-11ec-8491-61ceeb4939bd', date '2021-10-19', time '11:30:00'),
       ('73e176bf-2607-11ec-8491-61ceeb4939bd', '1e004ba7-2609-11ec-8491-61ceeb4939bd', date '2021-10-18', time '18:30:00'),
       ('73e176bf-2607-11ec-8491-61ceeb4939bd', '1e004ba8-2609-11ec-8491-61ceeb4939bd', date '2021-10-26', time '21:00:00'),

       ('73e176c0-2607-11ec-8491-61ceeb4939bd', '1e004ba6-2609-11ec-8491-61ceeb4939bd', date '2021-10-12', time '11:30:00'),
       ('73e176c0-2607-11ec-8491-61ceeb4939bd', '1e004ba7-2609-11ec-8491-61ceeb4939bd', date '2021-10-15', time '13:30:00'),
       ('73e176c0-2607-11ec-8491-61ceeb4939bd', '1e004ba8-2609-11ec-8491-61ceeb4939bd', date '2021-10-18', time '20:00:00'),

       ('73e176c1-2607-11ec-8491-61ceeb4939bd', '1e004ba6-2609-11ec-8491-61ceeb4939bd', date '2021-10-18', time '12:30:00'),
       ('73e176c1-2607-11ec-8491-61ceeb4939bd', '1e004ba7-2609-11ec-8491-61ceeb4939bd', date '2021-10-22', time '15:30:00'),
       ('73e176c1-2607-11ec-8491-61ceeb4939bd', '1e004ba8-2609-11ec-8491-61ceeb4939bd', date '2021-10-28', time '21:15:00'),

       ('73e176c2-2607-11ec-8491-61ceeb4939bd', '1e004ba6-2609-11ec-8491-61ceeb4939bd', date '2021-11-03', time '13:30:00'),
       ('73e176c2-2607-11ec-8491-61ceeb4939bd', '1e004ba7-2609-11ec-8491-61ceeb4939bd', date '2021-11-10', time '16:30:00'),
       ('73e176c2-2607-11ec-8491-61ceeb4939bd', '1e004ba8-2609-11ec-8491-61ceeb4939bd', date '2021-11-12', time '21:40:00'),

       ('73e176c3-2607-11ec-8491-61ceeb4939bd', '1e004ba6-2609-11ec-8491-61ceeb4939bd', date '2021-11-05', time '14:00:00'),
       ('73e176c3-2607-11ec-8491-61ceeb4939bd', '1e004ba7-2609-11ec-8491-61ceeb4939bd', date '2021-11-11', time '17:30:00'),
       ('73e176c3-2607-11ec-8491-61ceeb4939bd', '1e004ba8-2609-11ec-8491-61ceeb4939bd', date '2021-11-16', time '22:00:00');
select * from cinema.sessions;

insert into cinema.tickets (session_id, worker_id, price)
values ('463a9ba4-261c-11ec-8491-61ceeb4939bd', 'f8f4c058-2608-11ec-8491-61ceeb4939bd', 30.50),
       ('463a9ba4-261c-11ec-8491-61ceeb4939bd', '9907212e-261f-11ec-8491-61ceeb4939bd', 20.25),
       ('463a9ba4-261c-11ec-8491-61ceeb4939bd', '9907212f-261f-11ec-8491-61ceeb4939bd', 48),

       ('463a9ba5-261c-11ec-8491-61ceeb4939bd', 'f8f4c058-2608-11ec-8491-61ceeb4939bd', 21.50),
       ('463a9ba5-261c-11ec-8491-61ceeb4939bd', '9907212e-261f-11ec-8491-61ceeb4939bd', 25.25),
       ('463a9ba5-261c-11ec-8491-61ceeb4939bd', '9907212f-261f-11ec-8491-61ceeb4939bd', 30),

       ('463a9ba6-261c-11ec-8491-61ceeb4939bd', 'f8f4c058-2608-11ec-8491-61ceeb4939bd', 12.50),
       ('463a9ba6-261c-11ec-8491-61ceeb4939bd', '9907212e-261f-11ec-8491-61ceeb4939bd', 21.50),
       ('463a9ba6-261c-11ec-8491-61ceeb4939bd', '9907212f-261f-11ec-8491-61ceeb4939bd', 35),

       ('463a9ba7-261c-11ec-8491-61ceeb4939bd', 'f8f4c058-2608-11ec-8491-61ceeb4939bd', 17),
       ('463a9ba7-261c-11ec-8491-61ceeb4939bd', '9907212e-261f-11ec-8491-61ceeb4939bd', 25.50),
       ('463a9ba7-261c-11ec-8491-61ceeb4939bd', '9907212f-261f-11ec-8491-61ceeb4939bd', 50),

       ('463a9ba8-261c-11ec-8491-61ceeb4939bd', 'f8f4c058-2608-11ec-8491-61ceeb4939bd', 30),
       ('463a9ba8-261c-11ec-8491-61ceeb4939bd', '9907212e-261f-11ec-8491-61ceeb4939bd', 20.25),
       ('463a9ba8-261c-11ec-8491-61ceeb4939bd', '9907212f-261f-11ec-8491-61ceeb4939bd', 48),

       ('463a9ba9-261c-11ec-8491-61ceeb4939bd', 'f8f4c058-2608-11ec-8491-61ceeb4939bd', 35.50),
       ('463a9ba9-261c-11ec-8491-61ceeb4939bd', '9907212e-261f-11ec-8491-61ceeb4939bd', 25.25),
       ('463a9ba9-261c-11ec-8491-61ceeb4939bd', '9907212f-261f-11ec-8491-61ceeb4939bd', 48),

       ('463a9baa-261c-11ec-8491-61ceeb4939bd', 'f8f4c058-2608-11ec-8491-61ceeb4939bd', 30.50),
       ('463a9baa-261c-11ec-8491-61ceeb4939bd', '9907212e-261f-11ec-8491-61ceeb4939bd', 20.25),
       ('463a9baa-261c-11ec-8491-61ceeb4939bd', '9907212f-261f-11ec-8491-61ceeb4939bd', 48),

       ('463a9bab-261c-11ec-8491-61ceeb4939bd', 'f8f4c058-2608-11ec-8491-61ceeb4939bd', 20),
       ('463a9bab-261c-11ec-8491-61ceeb4939bd', '9907212e-261f-11ec-8491-61ceeb4939bd', 20),
       ('463a9bab-261c-11ec-8491-61ceeb4939bd', '9907212f-261f-11ec-8491-61ceeb4939bd', 35),

       ('463a9bac-261c-11ec-8491-61ceeb4939bd', 'f8f4c058-2608-11ec-8491-61ceeb4939bd', 15),
       ('463a9bac-261c-11ec-8491-61ceeb4939bd', '9907212e-261f-11ec-8491-61ceeb4939bd', 30),
       ('463a9bac-261c-11ec-8491-61ceeb4939bd', '9907212f-261f-11ec-8491-61ceeb4939bd', 50),

       ('463a9bad-261c-11ec-8491-61ceeb4939bd', 'f8f4c058-2608-11ec-8491-61ceeb4939bd', 35),
       ('463a9bad-261c-11ec-8491-61ceeb4939bd', '9907212e-261f-11ec-8491-61ceeb4939bd', 15.50),
       ('463a9bad-261c-11ec-8491-61ceeb4939bd', '9907212f-261f-11ec-8491-61ceeb4939bd', 35),

       ('463a9bae-261c-11ec-8491-61ceeb4939bd', 'f8f4c058-2608-11ec-8491-61ceeb4939bd', 35.50),
       ('463a9bae-261c-11ec-8491-61ceeb4939bd', '9907212e-261f-11ec-8491-61ceeb4939bd', 15),
       ('463a9bae-261c-11ec-8491-61ceeb4939bd', '9907212f-261f-11ec-8491-61ceeb4939bd', 40.50);
select * from cinema.tickets;










