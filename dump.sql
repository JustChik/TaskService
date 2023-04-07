--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.24
-- Dumped by pg_dump version 14.6 (Ubuntu 14.6-0ubuntu0.22.04.1)

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

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--


--
-- Name: status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.status AS ENUM (
    'Open',
    'In Progress',
    'On Hold',
    'Done',
    'Canceled'
);


ALTER TYPE public.status OWNER TO postgres;

SET default_tablespace = '';

--
-- Name: comments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.comments (
    id uuid NOT NULL,
    text character varying(500),
    created_at timestamp without time zone DEFAULT now(),
    task_id uuid
);


ALTER TABLE public.comments OWNER TO postgres;

--
-- Name: tasks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tasks (
    id uuid NOT NULL,
    tittle character varying(50),
    description character varying(500),
    status public.status,
    user_id uuid
);


ALTER TABLE public.tasks OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    username character varying(50),
    password character varying(128),
    resetpasswordtoken character varying(256)
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Data for Name: comments; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.comments (id, text, created_at, task_id) FROM stdin;
24605f83-e093-4c75-b4a2-82793f51cab7	Проверка коментариев	2023-03-15 14:37:37.770187	67f571be-634e-48dd-a0af-204e42bea3cb
f54139df-883a-430f-87fc-cf92ca33e072	Проверка коментариев	2023-03-15 14:37:46.695204	67f571be-634e-48dd-a0af-204e42bea3cb
a0ba9977-0920-4c0d-a3ce-eab6c5a13daf	Проверка коментариев	2023-03-15 13:45:26.956	67f571be-634e-48dd-a0af-204e42bea3cb
31925b90-c343-11ed-afa1-0242ac120002	Тест через Grip	2023-02-15 13:41:26.956	67f571be-634e-48dd-a0af-204e42bea3cb
da098dc4-bc54-48e5-8dd9-f06894eed4fb	1 Проверка	2023-03-21 14:13:46.722287	215c0059-82d9-4542-8739-031433b3de9f
d1270dbc-9de4-49ed-aff9-9089ff1e5582	1 Проверка	2023-03-21 14:14:18.175153	215c0059-82d9-4542-8739-031433b3de9f
be093427-4c07-4154-acad-79fc6d5c5a8b	1 Проверка	2023-03-21 14:14:26.161496	215c0059-82d9-4542-8739-031433b3de9f
778ae05f-cb7f-4b20-9a4e-0ec90aad7523	1 Проверка	2023-03-21 14:18:17.291031	215c0059-82d9-4542-8739-031433b3de9f
dc3b147a-5684-42b9-8a0f-d8f3407e2f98	1 Проверка	2023-03-21 14:19:56.50705	215c0059-82d9-4542-8739-031433b3de9f
357c569f-1244-4b2b-8927-38c5a9d5cb19	1 Проверка	2023-03-21 14:24:43.643614	215c0059-82d9-4542-8739-031433b3de9f
b30dab1a-64fa-45b3-9d21-40061dad59d2	1 Проверка	2023-03-21 14:26:43.653916	215c0059-82d9-4542-8739-031433b3de9f
75e51c43-9931-4823-9494-5d27cf1f95a9	1 Проверка	2023-03-21 14:27:10.924274	215c0059-82d9-4542-8739-031433b3de9f
b1f8db5f-d389-4f39-9635-3760830ebd0a	1 Проверка	2023-03-21 14:27:16.425287	215c0059-82d9-4542-8739-031433b3de9f
bf179bab-4caa-4d04-b99e-b361ba0da1ca	1 Проверка	2023-03-21 14:29:02.196877	215c0059-82d9-4542-8739-031433b3de9f
6a0c3b29-c1d4-4f14-b8d1-d9e89b8de3aa	1 Проверка	2023-03-21 14:29:48.459352	215c0059-82d9-4542-8739-031433b3de9f
f7ed15c5-5734-4a8b-909b-279801081a61	1 Проверка	2023-03-21 14:31:29.119178	215c0059-82d9-4542-8739-031433b3de9f
2c07b66c-87d5-4ef1-926f-535ce2bc6329	1 Проверка	2023-03-21 14:31:29.120241	215c0059-82d9-4542-8739-031433b3de9f
8dcc7866-fa66-488c-a37b-9b8ac1184d68	1 Проверка	2023-03-21 14:34:52.356797	215c0059-82d9-4542-8739-031433b3de9f
3f72de35-ad29-4838-8351-86109815239a	1 Проверка	2023-03-21 17:43:03.499594	215c0059-82d9-4542-8739-031433b3de9f
5d131a98-c6be-467e-bcdd-5e45d7c98643		2023-03-21 17:48:18.276016	215c0059-82d9-4542-8739-031433b3de9f
34f22603-1f73-4ba5-80cd-c218c3b46b9a		2023-03-21 17:52:41.512847	215c0059-82d9-4542-8739-031433b3de9f
190fa1a4-34af-42ae-bb20-8260c2f7b413		2023-03-21 17:53:00.658021	215c0059-82d9-4542-8739-031433b3de9f
71479d07-9b07-4dfd-8db2-b31c4ea9a638		2023-03-21 17:56:34.84761	215c0059-82d9-4542-8739-031433b3de9f
04788e5b-506b-4c54-95a0-8d6aae75fdd6		2023-03-21 17:56:42.839917	215c0059-82d9-4542-8739-031433b3de9f
60f83d5a-eb04-4993-88d5-5bf4aa1a2d97		2023-03-21 17:56:45.455435	215c0059-82d9-4542-8739-031433b3de9f
796cad68-cc01-418e-8f90-c1dddf6db13a		2023-03-21 17:56:46.531914	215c0059-82d9-4542-8739-031433b3de9f
987a4e7c-1139-43e1-9865-287081a250a8		2023-03-21 17:57:00.254036	215c0059-82d9-4542-8739-031433b3de9f
d2fb8266-48a9-46d6-9ac3-098d9bbb9f58		2023-03-21 17:57:56.877855	215c0059-82d9-4542-8739-031433b3de9f
ff7482f8-99f6-4fcb-bf6a-061c0460292e		2023-03-21 17:58:07.781901	215c0059-82d9-4542-8739-031433b3de9f
\.


--
-- Data for Name: tasks; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tasks (id, tittle, description, status, user_id) FROM stdin;
a0707430-e1f6-457d-9eeb-6485616295d3	ПРоверка	ПРоверка	Open	30375020-1f2c-49a9-a5b6-3ea22757aff7
6bfe2c2c-ce34-44e3-84a9-8659ddd8e535	Тест	ghdhd	Open	2fd134a9-85ec-44cd-9606-020354f03e9f
67f571be-634e-48dd-a0af-204e42bea3cb	FNFFNNFNF	adndadadadandnadna	Open	2fd134a9-85ec-44cd-9606-020354f03e9f
215c0059-82d9-4542-8739-031433b3de9f	хе2р	adnawndnadna	Done	30375020-1f2c-49a9-a5b6-3ea22757aff7
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, username, password) FROM stdin;
2fd134a9-85ec-44cd-9606-020354f03e9f	с	32d6f813078c445492f1674057c4aed4
1ed57461-9d85-412b-9353-3e2c3fb35bff	сchak	4b5927b4a77c56fd6aeab0bb9ead0769
e9ef60e0-d9fd-4387-bd2b-5cd28bb97511	nik	fc878a3d2665a9266bbab5c1a726a294
eed0b9fe-9796-403c-b277-1f13fb5028a7	misha	4b5927b4a77c56fd6aeab0bb9ead0769
43f41980-74ef-4dc0-a829-98c3006068e0	misha2	4b5927b4a77c56fd6aeab0bb9ead0769
36858f69-18ac-47f4-821a-7ff6a1aca4d7	misha3	4b5927b4a77c56fd6aeab0bb9ead0769
cf8e2f06-add7-42e0-845a-6ee3faaad473	misha4	4b5927b4a77c56fd6aeab0bb9ead0769
651190a6-2ff6-45af-9e50-2f0fbf72d1f4	nikitnov	79bebd91e50737e4fa4aaa529a51c584
903ea7e3-eab5-47da-9c01-99139e18a2ba		79bebd91e50737e4fa4aaa529a51c584
30375020-1f2c-49a9-a5b6-3ea22757aff7	papa	79bebd91e50737e4fa4aaa529a51c584
543da3c1-2264-4da5-a9b5-9c910ceea3b6	new	79bebd91e50737e4fa4aaa529a51c584
\.


--
-- Name: comments comments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_pkey PRIMARY KEY (id);


--
-- Name: tasks tasks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: comments comments_task_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_task_id_fkey FOREIGN KEY (task_id) REFERENCES public.tasks(id) ON DELETE CASCADE;


--
-- Name: tasks tasks_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

