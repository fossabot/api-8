SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: ar_internal_metadata; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ar_internal_metadata (
    key character varying NOT NULL,
    value character varying,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);


--
-- Name: user_profile; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_profile (
    user_id integer NOT NULL,
    activation_token character varying(100) NOT NULL,
    activation_token_expires_at timestamp with time zone DEFAULT (now() + '00:15:00'::interval) NOT NULL,
    name character varying(500) NOT NULL,
    address character varying(2000) NOT NULL,
    email character varying(500) NOT NULL,
    phone character varying(50) NOT NULL
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(100) NOT NULL,
    encrypted_password character varying(200),
    github_username character varying(500),
    active boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


--
-- Name: COLUMN users.encrypted_password; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.encrypted_password IS 'can be nil if user created the account using social login';


--
-- Name: COLUMN users.active; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.active IS 'true if user has clicked activation url';


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: ar_internal_metadata ar_internal_metadata_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ar_internal_metadata
    ADD CONSTRAINT ar_internal_metadata_pkey PRIMARY KEY (key);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: unique_activation_token_on_users; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_activation_token_on_users ON public.user_profile USING btree (activation_token);


--
-- Name: unique_email_on_users; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_email_on_users ON public.user_profile USING btree (email);


--
-- Name: unique_github_username_on_users; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_github_username_on_users ON public.users USING btree (github_username);


--
-- Name: unique_phone_on_users; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_phone_on_users ON public.user_profile USING btree (phone);


--
-- Name: unique_user_id_on_users; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_user_id_on_users ON public.user_profile USING btree (user_id);


--
-- Name: unique_username_on_users; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_username_on_users ON public.users USING btree (username);


--
-- Name: user_profile user_profile_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_profile
    ADD CONSTRAINT user_profile_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

SET search_path TO "$user", public;

INSERT INTO "schema_migrations" (version) VALUES
('20181214155421'),
('20181214161314');


