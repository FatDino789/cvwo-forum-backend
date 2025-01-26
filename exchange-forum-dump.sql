--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5 (Debian 14.5-2.pgdg110+2)
-- Dumped by pg_dump version 14.15 (Homebrew)

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
-- Name: posts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.posts (
    id text NOT NULL,
    user_id text,
    title character varying(255) NOT NULL,
    content text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    likes_count integer DEFAULT 0,
    views_count integer DEFAULT 0,
    comments jsonb DEFAULT '[]'::jsonb,
    tags text[] DEFAULT '{}'::text[]
);


ALTER TABLE public.posts OWNER TO postgres;

--
-- Name: tags; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tags (
    id text NOT NULL,
    text character varying(50) NOT NULL,
    color character varying(7) NOT NULL,
    searches integer DEFAULT 0
);


ALTER TABLE public.tags OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id text NOT NULL,
    username character varying(50) NOT NULL,
    email character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    icon_index integer NOT NULL,
    color_index integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    likes text[] DEFAULT '{}'::text[]
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Data for Name: posts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.posts (id, user_id, title, content, created_at, updated_at, likes_count, views_count, comments, tags) FROM stdin;
2af546fb-0aef-4adb-8a47-1e3f0670a525	ff716f73-aeb9-49cc-bbb1-43465aecae8c	Cultural Shock in Seoul: A Semester's Journey	From subway etiquette to late-night study cafes, sharing my biggest culture shock moments during my semester at Yonsei University. Bonus: How I navigated the language barrier using basic Korean phrases.	2025-01-26 14:33:56.672379+00	2025-01-26 14:41:50.04316+00	1	11	[{"id": "002b0399-c386-458e-9fb5-b5d4ebbdf54e", "content": "Wow sounds like loads of fun! Glad you were able to get your first choice destination!", "user_id": "123e4567-e89b-12d3-a456-426614174000", "username": "testuser123", "created_at": "2025-01-26T14:37:30.767Z", "icon_index": 45, "color_index": 2}, {"id": "ddab534a-a98c-471f-b526-d97c5ffd658e", "content": "Congrates on completing your exchange!", "user_id": "7bddb903-6a44-4abe-9e13-1a3c0dd7a180", "username": "ronanlim", "created_at": "2025-01-26T14:39:15.962Z", "icon_index": 50, "color_index": 1}, {"id": "a7072467-448a-4b2b-9dc8-1882037e39d9", "content": "@testuser123 HAHAH I was lucky enough!", "user_id": "ff716f73-aeb9-49cc-bbb1-43465aecae8c", "username": "antonius", "created_at": "2025-01-26T14:40:35.433Z", "icon_index": 54, "color_index": 3}, {"id": "9819c0fb-4230-462f-b7d7-2060dc8eb9e2", "content": "@ronanlim Thank you see you soon!", "user_id": "ff716f73-aeb9-49cc-bbb1-43465aecae8c", "username": "antonius", "created_at": "2025-01-26T14:40:41.382Z", "icon_index": 54, "color_index": 3}]	{7079bed1-2fe8-42f7-891f-368b95258db6,d3359cb1-2e1c-49fd-8c13-71e645ad90e1,4a544809-862f-4821-9ec3-a8409474e8f3}
fd81b65c-77fa-42e6-b0d7-2cead714445e	123e4567-e89b-12d3-a456-426614174000	Volunteering Experience in Cambodia	Just completed my summer volunteering program teaching English in rural Cambodia. The kids' enthusiasm and resilience were truly inspiring. Highlight was organizing a mini English language fair for the community.	2025-01-26 14:30:45.009216+00	2025-01-26 14:41:36.408552+00	2	8	[{"id": "7f5e0799-4d33-4b69-9935-d7df6af6f316", "content": "Wow! Sounds like a great time! Just a pity that I could not come with you on this trip. See you back in Singapore!", "user_id": "ff716f73-aeb9-49cc-bbb1-43465aecae8c", "username": "antonius", "created_at": "2025-01-26T14:34:38.996Z", "icon_index": 54, "color_index": 3}, {"id": "758378f7-07e3-44ce-bc78-944ad22cdc91", "content": "@antonius Don't forget about your trip to Korea. I am sure you had loads of fun too!", "user_id": "7bddb903-6a44-4abe-9e13-1a3c0dd7a180", "username": "ronanlim", "created_at": "2025-01-26T14:39:43.515Z", "icon_index": 50, "color_index": 1}]	{4a544809-862f-4821-9ec3-a8409474e8f3,1b79f27a-1a44-4004-b84e-f3f38f953881,3fd95908-27ee-4434-90c2-1aa1b0e308ab}
11c4b0b2-6ce4-4450-85ef-49ed1d4ae3d8	7bddb903-6a44-4abe-9e13-1a3c0dd7a180	NUS SEP Application 2024	Compiled my experience applying for Student Exchange at NUS. Key deadlines, required documents, and interview preparation tips included. Remember to start your application early as spots fill up quickly!	2025-01-26 14:32:15.740052+00	2025-01-26 14:41:51.426259+00	2	12	[{"id": "a2bf9c3f-76f2-46db-b743-06d358c61b4f", "content": "THANK YOU SO MUCH! This was a life saver 🤩", "user_id": "ff716f73-aeb9-49cc-bbb1-43465aecae8c", "username": "antonius", "created_at": "2025-01-26T14:35:50.5Z", "icon_index": 54, "color_index": 3}, {"id": "8b59d9bd-88e2-4b58-bf3c-b0441fdc91db", "content": "Was looking for something like this for a long time. Thankuuu ❤️", "user_id": "123e4567-e89b-12d3-a456-426614174000", "username": "testuser123", "created_at": "2025-01-26T14:36:57.1Z", "icon_index": 45, "color_index": 2}]	{4e6c7d6c-6969-4683-90a5-f84ccdcf884d,d3359cb1-2e1c-49fd-8c13-71e645ad90e1}
\.


--
-- Data for Name: tags; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tags (id, text, color, searches) FROM stdin;
3fd95908-27ee-4434-90c2-1aa1b0e308ab	VOLUNTEERING	#F4E5ED	1
1b79f27a-1a44-4004-b84e-f3f38f953881	CAMBODIA	#E7F3F1	1
7079bed1-2fe8-42f7-891f-368b95258db6	SEOUL	#F2E8EC	1
d3359cb1-2e1c-49fd-8c13-71e645ad90e1	EXCHANGE	#EDE7F6	2
4a544809-862f-4821-9ec3-a8409474e8f3	NUS	#E5EEF9	2
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, username, email, password_hash, icon_index, color_index, created_at, likes) FROM stdin;
123e4567-e89b-12d3-a456-426614174001	johndoe	john@example.com	$2a$10$mN6CaIxk7mU0QM3B2Q490euGHJS5Dx0AOTjG7v82f9dQL/Gm.gCEC	12	4	2025-01-26 14:29:33.08874+00	{123e4567-e89b-12d3-a456-426614174003}
123e4567-e89b-12d3-a456-426614174002	janesmith	jane@example.com	$2a$10$mN6CaIxk7mU0QM3B2Q490euGHJS5Dx0AOTjG7v82f9dQL/Gm.gCEC	33	1	2025-01-26 14:29:33.08874+00	{}
9a4cd37d-e59c-408e-8b6c-676890f9017b	vincentgoh	testing3@gmail.com	$2a$10$Qxp.73q.AQLBmaRlEzccUuPsEeAQiJNBIslDX0wHqrEjjiQrORowq	31	11	2025-01-26 14:31:06.168368+00	{}
00420a70-0a2d-487d-9224-f76a5ef3f9e3	mofiras	testing5@gmail.com	$2a$10$pTD67fiyKDfq8NClQc0aH.CSTY7s.YL.PIBVsG9KFLkorAnpwDuzq	27	7	2025-01-26 14:32:55.313084+00	{}
123e4567-e89b-12d3-a456-426614174000	testuser123	testing@gmail.com	$2a$10$mN6CaIxk7mU0QM3B2Q490euGHJS5Dx0AOTjG7v82f9dQL/Gm.gCEC	45	2	2025-01-26 14:29:33.08874+00	{123e4567-e89b-12d3-a456-426614174004,2af546fb-0aef-4adb-8a47-1e3f0670a525,11c4b0b2-6ce4-4450-85ef-49ed1d4ae3d8}
7bddb903-6a44-4abe-9e13-1a3c0dd7a180	ronanlim	testing4@gmail.com	$2a$10$VCM9TzSojbtrrt0xscxYcuY8KOAye5OjrhcCdVIzf/49cyf3zw1iG	50	1	2025-01-26 14:31:39.195985+00	{fd81b65c-77fa-42e6-b0d7-2cead714445e}
ff716f73-aeb9-49cc-bbb1-43465aecae8c	antonius	testing6@gmail.com	$2a$10$RkLkSjiAU5hR0jBKGdFid.g6Jij77OtsWFQexRRRdmT8mhatmbzeC	54	3	2025-01-26 14:33:16.16218+00	{fd81b65c-77fa-42e6-b0d7-2cead714445e,11c4b0b2-6ce4-4450-85ef-49ed1d4ae3d8}
\.


--
-- Name: posts posts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);


--
-- Name: tags tags_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (id);


--
-- Name: tags tags_text_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_text_key UNIQUE (text);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


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
-- Name: posts posts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

