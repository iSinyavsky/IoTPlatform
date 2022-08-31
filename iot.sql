--
-- PostgreSQL database dump
--

-- Dumped from database version 12.2 (Debian 12.2-2.pgdg100+1)
-- Dumped by pg_dump version 12.2 (Debian 12.2-2.pgdg100+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: events; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.events (
    id integer NOT NULL,
    if_event json,
    then_event json,
    "createdAt" timestamp without time zone DEFAULT now()
);


ALTER TABLE public.events OWNER TO postgres;

--
-- Name: events_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.events ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.events_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: migrate_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.migrate_history (
    id integer NOT NULL,
    version character varying,
    comment character varying,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.migrate_history OWNER TO postgres;

--
-- Name: migrate_history_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.migrate_history_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.migrate_history_id_seq OWNER TO postgres;

--
-- Name: migrate_history_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.migrate_history_id_seq OWNED BY public.migrate_history.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    email text,
    password text,
    token text,
    name text
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.users ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: users_variables; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users_variables (
    id bigint NOT NULL,
    userid integer,
    varid bigint
);


ALTER TABLE public.users_variables OWNER TO postgres;

--
-- Name: users_variables_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.users_variables ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.users_variables_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: values; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."values" (
    id integer NOT NULL,
    varid bigint,
    value text,
    "createdAt" timestamp without time zone DEFAULT now()
);


ALTER TABLE public."values" OWNER TO postgres;

--
-- Name: values_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public."values" ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.values_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: variables; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.variables (
    id bigint NOT NULL,
    name text,
    label text,
    integration_service text,
    integration_id text,
    integration_type text
);


ALTER TABLE public.variables OWNER TO postgres;

--
-- Name: variables_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.variables ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.variables_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: migrate_history id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.migrate_history ALTER COLUMN id SET DEFAULT nextval('public.migrate_history_id_seq'::regclass);


--
-- Data for Name: events; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.events (id, if_event, then_event, "createdAt") FROM stdin;
14	{"1":{"value":"20","operator":"\\u0026ge;","type":1}}	{"34":{"value":"1"}}	2021-01-26 06:32:41.875395
15	{"1":{"value":"20","operator":"\\u0026lt;","type":1}}	{"34":{"value":"0"}}	2021-01-26 06:33:17.311931
\.


--
-- Data for Name: migrate_history; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.migrate_history (id, version, comment, created_at) FROM stdin;
1	1.0.0	initial_schema	2020-10-03 12:15:02.025392
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, email, password, token, name) FROM stdin;
1	23	7372ae5f551a12dfc285df0f69214b51fc6a8dd6cb649121646313a111f603fd9e357cc9	535fa30d7e25dd8a49f1536779734ec8286108d115da5045d77f3b4185d8f790	\N
2	232	70ba33708cbfb103f1a8e34afef333ba7dc021022b2d9aaa583aabb8058d8d67	835d5e8314340ab852a2f979ab4cd53e994dbe38366afb6eed84fe4957b980c8	\N
3	test@admin.ru	5994471abb01112afcc18159f6cc74b4f511b99806da59b3caf5a9c173cacfc5	4a7901ddc2298d07d7433e31cc9aab5d4e4701ea4e80e6a6436024a5ccd1884c	\N
4	test	9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08	9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08	\N
5	fdvdfv	768091cbbd1c586d807d941e7a9de067763f61b405586193056944a5a82b85ca	a2091cd94e6199b79e6eea1cfa76eaa1c7d1d7f5529b77b4698aa2b97107edbd	\N
6	kulikov-k@list.ru	fe3f893f705ce5d280ad0d6ccab042f2865fa719858acd5116564b2f54a3f1c8	970a1ab9d5cb2d49c72514445c634b2848e4c6049941c1736b126ac0aefd65ac	\N
\.


--
-- Data for Name: users_variables; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users_variables (id, userid, varid) FROM stdin;
1	4	1
2	0	9
3	0	10
4	4	13
5	4	14
6	4	15
7	4	16
8	0	17
9	0	18
10	0	19
11	4	20
12	4	21
13	4	22
15	0	24
16	0	25
17	0	26
18	4	27
19	4	28
20	4	29
25	4	34
26	0	35
27	0	36
28	0	37
32	0	41
33	0	42
34	0	43
35	4	44
36	4	45
37	4	46
\.


--
-- Data for Name: values; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."values" (id, varid, value, "createdAt") FROM stdin;
168	1	22	2021-01-25 16:13:50.882017
170	15	off	2021-01-25 16:13:50.904709
171	1	24	2021-01-25 16:33:36.257077
173	15	off	2021-01-25 16:33:36.279899
174	1	25	2021-01-25 16:36:49.026313
176	15	off	2021-01-25 16:36:49.039706
177	1	100	2021-01-25 16:37:11.675356
179	15	on	2021-01-25 16:37:11.68813
180	16	10	2021-01-25 16:37:29.519746
181	1	15	2021-01-25 16:38:03.67384
183	15	on	2021-01-25 16:38:03.686008
184	1	25	2021-01-25 16:38:18.015053
186	15	off	2021-01-25 16:38:18.027425
187	23	empty	2021-01-25 19:47:32.080388
188	1	18	2021-01-25 19:48:20.325538
190	23	on	2021-01-25 19:48:20.3557
191	23	off	2021-01-25 20:24:08.932842
192	1	24	2021-01-25 20:24:17.822377
194	23	on	2021-01-25 20:24:17.839846
195	1	25	2021-01-26 06:18:47.29799
197	1	20	2021-01-26 06:21:58.522373
199	1	21	2021-01-26 06:22:52.118646
201	1	22	2021-01-26 06:33:45.694666
203	1	24	2021-01-26 06:34:07.914819
205	1	20	2021-01-26 06:34:34.766336
207	1	23	2021-01-26 06:36:40.222007
209	34	1	2021-01-26 06:36:40.262387
210	1	18	2021-01-26 06:37:08.57887
212	34	0	2021-01-26 06:37:08.605149
213	1	22	2021-01-26 06:37:46.42605
215	34	1	2021-01-26 06:37:46.443182
216	1	12	2021-01-26 06:38:30.291227
218	34	0	2021-01-26 06:38:30.311617
219	1	24	2021-01-26 06:38:38.374203
221	34	1	2021-01-26 06:38:38.397281
165	1	20	2021-01-25 16:13:33.790251
167	15	off	2021-01-25 16:13:33.833489
\.


--
-- Data for Name: variables; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.variables (id, name, label, integration_service, integration_id, integration_type) FROM stdin;
4	ccftd	cc	\N	\N	\N
2	fmt		\N	\N	\N
6	asd		\N	\N	\N
5	asdsdasd		\N	\N	\N
7	hgfj	ffff	\N	\N	\N
9	hgfj00	ffff00	\N	\N	\N
10	test10	10test	\N	\N	\N
1	Датчик температуры	aa	\N	\N	\N
34	Свет	8484e0144adafda38b3da3e26c3d2902	\N	\N	\N
44	Лампа	559765e9-0c84-467a-b2fa-dbebe78ad2c6	yandex	559765e9-0c84-467a-b2fa-dbebe78ad2c6	devices.types.light
45	Ночник	68c5c21f-7dee-4786-86c7-1377c8526824	yandex	68c5c21f-7dee-4786-86c7-1377c8526824	devices.types.light
46	Ярик	5964c88c-4bb4-442c-bdd9-53490ba74e85	yandex	5964c88c-4bb4-442c-bdd9-53490ba74e85	devices.types.vacuum_cleaner
\.


--
-- Name: events_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.events_id_seq', 15, true);


--
-- Name: migrate_history_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.migrate_history_id_seq', 1, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 6, true);


--
-- Name: users_variables_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_variables_id_seq', 37, true);


--
-- Name: values_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.values_id_seq', 221, true);


--
-- Name: variables_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.variables_id_seq', 46, true);


--
-- Name: events events_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_pkey PRIMARY KEY (id);


--
-- Name: migrate_history migrate_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.migrate_history
    ADD CONSTRAINT migrate_history_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users_variables users_variables_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_variables
    ADD CONSTRAINT users_variables_pkey PRIMARY KEY (id);


--
-- Name: values values_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."values"
    ADD CONSTRAINT values_pkey PRIMARY KEY (id);


--
-- Name: variables variables_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.variables
    ADD CONSTRAINT variables_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

