-- +migrate Up
-- Table: public.address

CREATE TABLE IF NOT EXISTS public.address
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    building_number integer,
    street character varying(45),
    city character varying(45),
    district character varying(45),
    region character varying(45),
    postal_code character varying(45),
    CONSTRAINT address_id PRIMARY KEY (id)
)

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.address
    OWNER to postgres;

INSERT INTO public.address(
    building_number, street, district, city, region, postal_code)
VALUES (1, 'polska', 'polska', 'polska', 'polska', '58000');

-- Table: public.person

CREATE TABLE IF NOT EXISTS public.person
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    name character varying(45),
    phone character varying(30),
    email character varying(45),
    birthday timestamp,
    address_id integer,
    CONSTRAINT person_id PRIMARY KEY (id),
    CONSTRAINT address_id FOREIGN KEY (address_id)
        REFERENCES public.address (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE RESTRICT
)

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.person
    OWNER to postgres;

INSERT INTO public.person(
    name, phone, email, address_id)
VALUES ('Derek', '+380435815532', 'your.funny.email@lol.tik', 1);

-- Table: public.position

CREATE TABLE IF NOT EXISTS public.position
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    name character varying(45),
    access_level integer,
    CONSTRAINT position_id PRIMARY KEY (id)
)

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.position
    OWNER to postgres;

INSERT INTO public.position(
    name, access_level)
VALUES ('manager', 5);

-- Table: public.staff

CREATE TABLE IF NOT EXISTS public.staff
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    employment_date date,
    salary double precision,
    status text,
    person_id integer,
    cafe_id integer,
    position_id integer,
    CONSTRAINT staff_id PRIMARY KEY (id),
    CONSTRAINT person_id FOREIGN KEY (person_id)
        REFERENCES public.person (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT position_id FOREIGN KEY (position_id)
        REFERENCES public.position (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE RESTRICT
)

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.staff
    OWNER to postgres;

INSERT INTO public.staff(
    employment_date, salary, status, person_id, cafe_id, position_id)
VALUES ('1996-12-02', 100, 'BUSY', 1, 1, 1);




-- +migrate Down
DROP TABLE IF EXISTS public.staff;
DROP TABLE IF EXISTS public.position;
DROP TABLE IF EXISTS public.person;
DROP TABLE IF EXISTS public.address;