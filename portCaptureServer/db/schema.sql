DROP DATABASE port_capture;
CREATE DATABASE port_capture;

\connect port_capture

CREATE TABLE public.ports(
    id SERIAL NOT NULL PRIMARY KEY,
	primary_unloc CHAR(5) NOT NULL,
	name CHARACTER VARYING(100) NOT NULL,
	code CHAR(5),
	city CHARACTER VARYING(100) NOT NULL,
	country CHARACTER VARYING(100) NOT NULL,
    coordinates POINT NOT NULL,
	province CHARACTER VARYING(100) NOT NULL,
	timezone CHARACTER VARYING(100) NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
	deleted_at TIMESTAMP DEFAULT NULL
);
CREATE INDEX name_city ON public.ports(name);
CREATE INDEX name_city ON public.ports(primary_unloc);
CREATE UNIQUE INDEX no_duplicate_code ON public.ports(primary_unloc,deleted_at)
   WHERE deleted_at IS null;

CREATE TABLE public.alias(
    port_id INTEGER REFERENCES public.ports NOT NULL,
	name CHARACTER VARYING(100) NOT NULL
);
CREATE UNIQUE INDEX no_duplicate_alias ON public.alias(port_id,name);

CREATE TABLE public.regions(
    port_id INTEGER REFERENCES public.ports NOT NULL,
	name CHARACTER VARYING(100) NOT NULL
);
CREATE UNIQUE INDEX no_duplicate_regions ON public.regions(port_id,name);

CREATE TABLE public.unlocs(
    port_id INTEGER REFERENCES public.ports NOT NULL,
	name CHARACTER VARYING(100) NOT NULL
);
CREATE UNIQUE INDEX no_duplicate_unlocs ON public.unlocs(port_id,name);
