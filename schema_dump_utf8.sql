--
-- PostgreSQL database dump
--

\restrict aOh4bAmuhZxTR1AjFq1BKYgFGIerLlVTAv99rQcWUeuaP4SxlMdzCdzEQAAkBPg

-- Dumped from database version 17.9
-- Dumped by pg_dump version 17.9

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: admin
--

-- *not* creating schema, since initdb creates it


ALTER SCHEMA public OWNER TO admin;

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: admin
--

COMMENT ON SCHEMA public IS '';


--
-- Name: trg_auto_trace_material_ng(); Type: FUNCTION; Schema: public; Owner: admin
--

CREATE FUNCTION public.trg_auto_trace_material_ng() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    v_category_name TEXT;
    v_product_id TEXT;
    v_material_id TEXT;
BEGIN
    SELECT dc."CATEGORY_NAME" INTO v_category_name
    FROM "defect_master" dm
    JOIN "defect_categories" dc ON dm."CATEGORY_ID" = dc."ID"
    WHERE dm."ID" = NEW."DEFECT_ID";

    SELECT "PRODUCT_ID" INTO v_product_id
    FROM "daily_inspections"
    WHERE "ID" = NEW."INSPECTION_ID";

    IF v_category_name = 'MATERIAL' THEN
        FOR v_material_id IN
            SELECT bom."MATERIAL_ID"
            FROM "bill_of_materials" bom
            JOIN "material_potential_defects" mpd ON bom."MATERIAL_ID" = mpd."MATERIAL_ID"
            WHERE bom."PRODUCT_ID" = v_product_id AND mpd."DEFECT_ID" = NEW."DEFECT_ID"
        LOOP
            INSERT INTO "material_defect_ledger" ("INSPECTION_LOG_ID", "MATERIAL_ID", "DEFECT_ID", "QTY_NG")
            VALUES (NEW."ID", v_material_id, NEW."DEFECT_ID", NEW."QTY_NG");
        END LOOP;
    END IF;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.trg_auto_trace_material_ng() OWNER TO admin;

--
-- Name: seq_mat_id; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.seq_mat_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.seq_mat_id OWNER TO admin;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: MATERIAL; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public."MATERIAL" (
    "ID" text DEFAULT ('MAT-'::text || lpad((nextval('public.seq_mat_id'::regclass))::text, 5, '0'::text)) NOT NULL,
    "SUPP_ID" text,
    "UNIQ" text,
    "PART NAME" text NOT NULL,
    "TEBAL (MM)" numeric(12,3),
    "LEBAR (CM)" numeric(15,3),
    "PANJANG (CM)" numeric(15,3),
    "BERAT (GSM)" numeric(12,3),
    "MASSA (KG)" numeric(12,3),
    "SATUAN" text,
    "PICTURE_LOCATION" text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public."MATERIAL" OWNER TO admin;

--
-- Name: seq_insp_id; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.seq_insp_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.seq_insp_id OWNER TO admin;

--
-- Name: daily_inspections; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.daily_inspections (
    "ID" text DEFAULT ((('INSP-'::text || to_char((CURRENT_DATE)::timestamp with time zone, 'YYYYMMDD'::text)) || '-'::text) || lpad((nextval('public.seq_insp_id'::regclass))::text, 4, '0'::text)) NOT NULL,
    "INSPECTION_DATE" date DEFAULT CURRENT_DATE NOT NULL,
    "STAT_YEAR" integer GENERATED ALWAYS AS (EXTRACT(year FROM "INSPECTION_DATE")) STORED,
    "STAT_MONTH" integer GENERATED ALWAYS AS (EXTRACT(month FROM "INSPECTION_DATE")) STORED,
    "STAT_WEEK" integer GENERATED ALWAYS AS (EXTRACT(week FROM "INSPECTION_DATE")) STORED,
    "LEADER_ID" bigint,
    "LINE_ID" text,
    "PRODUCT_ID" text,
    "TOTAL_PRODUCED" integer DEFAULT 0 NOT NULL,
    "TOTAL_OK" integer DEFAULT 0 NOT NULL,
    "TOTAL_NG" integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_totals CHECK ((("TOTAL_OK" + "TOTAL_NG") = "TOTAL_PRODUCED"))
);


ALTER TABLE public.daily_inspections OWNER TO admin;

--
-- Name: analytics_control_chart; Type: VIEW; Schema: public; Owner: admin
--

CREATE VIEW public.analytics_control_chart AS
 WITH daily_rate AS (
         SELECT daily_inspections."INSPECTION_DATE",
            sum(daily_inspections."TOTAL_PRODUCED") AS n,
            sum(daily_inspections."TOTAL_NG") AS np,
            ((sum(daily_inspections."TOTAL_NG"))::numeric / (NULLIF(sum(daily_inspections."TOTAL_PRODUCED"), 0))::numeric) AS p_bar_daily
           FROM public.daily_inspections
          GROUP BY daily_inspections."INSPECTION_DATE"
        ), global_stats AS (
         SELECT avg(daily_rate.p_bar_daily) AS cl,
            stddev(daily_rate.p_bar_daily) AS sigma
           FROM daily_rate
        )
 SELECT d."INSPECTION_DATE",
    round(d.p_bar_daily, 4) AS daily_fraction_defective,
    round(g.cl, 4) AS center_line,
    round(GREATEST((0)::numeric, (g.cl + ((3)::numeric * g.sigma))), 4) AS ucl,
    round(LEAST((1)::numeric, (g.cl - ((3)::numeric * g.sigma))), 4) AS lcl
   FROM daily_rate d,
    global_stats g
  ORDER BY d."INSPECTION_DATE";


ALTER VIEW public.analytics_control_chart OWNER TO admin;

--
-- Name: seq_def_cat_id; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.seq_def_cat_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.seq_def_cat_id OWNER TO admin;

--
-- Name: defect_categories; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.defect_categories (
    "ID" text DEFAULT ('CAT-'::text || lpad((nextval('public.seq_def_cat_id'::regclass))::text, 2, '0'::text)) NOT NULL,
    "CATEGORY_NAME" text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.defect_categories OWNER TO admin;

--
-- Name: seq_defect_id; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.seq_defect_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.seq_defect_id OWNER TO admin;

--
-- Name: defect_master; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.defect_master (
    "ID" text DEFAULT ('DEF-'::text || lpad((nextval('public.seq_defect_id'::regclass))::text, 4, '0'::text)) NOT NULL,
    "CATEGORY_ID" text,
    "NG_NAME" text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.defect_master OWNER TO admin;

--
-- Name: inspection_logs; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.inspection_logs (
    "ID" bigint NOT NULL,
    "INSPECTION_ID" text,
    "DEFECT_ID" text,
    "TIME_WINDOW" character varying(50) NOT NULL,
    "OCCURRENCE_TIME" time without time zone NOT NULL,
    "QTY_NG" integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "inspection_logs_QTY_NG_check" CHECK (("QTY_NG" > 0))
);


ALTER TABLE public.inspection_logs OWNER TO admin;

--
-- Name: analytics_fishbone_logic; Type: VIEW; Schema: public; Owner: admin
--

CREATE VIEW public.analytics_fishbone_logic AS
 SELECT dc."CATEGORY_NAME" AS "Main_Cause",
    dm."NG_NAME" AS "Sub_Cause",
    count(il."ID") AS occurrence_count
   FROM ((public.defect_master dm
     JOIN public.defect_categories dc ON ((dm."CATEGORY_ID" = dc."ID")))
     LEFT JOIN public.inspection_logs il ON ((dm."ID" = il."DEFECT_ID")))
  GROUP BY dc."CATEGORY_NAME", dm."NG_NAME";


ALTER VIEW public.analytics_fishbone_logic OWNER TO admin;

--
-- Name: analytics_histogram_time; Type: VIEW; Schema: public; Owner: admin
--

CREATE VIEW public.analytics_histogram_time AS
 SELECT "TIME_WINDOW",
    count(*) AS frequency,
    sum("QTY_NG") AS total_qty_ng
   FROM public.inspection_logs
  GROUP BY "TIME_WINDOW"
  ORDER BY "TIME_WINDOW";


ALTER VIEW public.analytics_histogram_time OWNER TO admin;

--
-- Name: analytics_pareto_data; Type: VIEW; Schema: public; Owner: admin
--

CREATE VIEW public.analytics_pareto_data AS
 WITH ng_summary AS (
         SELECT dm."NG_NAME",
            sum(il."QTY_NG") AS total_qty
           FROM (public.inspection_logs il
             JOIN public.defect_master dm ON ((il."DEFECT_ID" = dm."ID")))
          GROUP BY dm."NG_NAME"
        ), pareto_calc AS (
         SELECT ng_summary."NG_NAME",
            ng_summary.total_qty,
            sum(ng_summary.total_qty) OVER (ORDER BY ng_summary.total_qty DESC) AS cumulative_qty,
            sum(ng_summary.total_qty) OVER () AS grand_total
           FROM ng_summary
        )
 SELECT "NG_NAME",
    total_qty,
    round((((total_qty)::numeric / NULLIF(grand_total, (0)::numeric)) * (100)::numeric), 2) AS percentage,
    round(((cumulative_qty / NULLIF(grand_total, (0)::numeric)) * (100)::numeric), 2) AS cumulative_percentage
   FROM pareto_calc
  ORDER BY total_qty DESC;


ALTER VIEW public.analytics_pareto_data OWNER TO admin;

--
-- Name: seq_line_id; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.seq_line_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.seq_line_id OWNER TO admin;

--
-- Name: production_lines; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.production_lines (
    "ID" text DEFAULT ('LINE-'::text || lpad((nextval('public.seq_line_id'::regclass))::text, 3, '0'::text)) NOT NULL,
    "LINE_NAME" text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.production_lines OWNER TO admin;

--
-- Name: analytics_stratification_trend; Type: VIEW; Schema: public; Owner: admin
--

CREATE VIEW public.analytics_stratification_trend AS
 SELECT di."STAT_YEAR",
    di."STAT_MONTH",
    di."STAT_WEEK",
    pl."LINE_NAME",
    sum(di."TOTAL_PRODUCED") AS prod_total,
    sum(di."TOTAL_NG") AS ng_total,
    round((((sum(di."TOTAL_NG"))::numeric / (NULLIF(sum(di."TOTAL_PRODUCED"), 0))::numeric) * (100)::numeric), 2) AS reject_rate_pct
   FROM (public.daily_inspections di
     JOIN public.production_lines pl ON ((di."LINE_ID" = pl."ID")))
  GROUP BY di."STAT_YEAR", di."STAT_MONTH", di."STAT_WEEK", pl."LINE_NAME"
  ORDER BY di."STAT_YEAR" DESC, di."STAT_MONTH" DESC, di."STAT_WEEK" DESC;


ALTER VIEW public.analytics_stratification_trend OWNER TO admin;

--
-- Name: bill_of_materials; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.bill_of_materials (
    "PRODUCT_ID" text NOT NULL,
    "MATERIAL_ID" text NOT NULL,
    "USAGE_QTY" numeric(15,4) DEFAULT 1.0000 NOT NULL
);


ALTER TABLE public.bill_of_materials OWNER TO admin;

--
-- Name: seq_cust_id; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.seq_cust_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.seq_cust_id OWNER TO admin;

--
-- Name: customers; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.customers (
    "ID" text DEFAULT ('CUST-'::text || lpad((nextval('public.seq_cust_id'::regclass))::text, 3, '0'::text)) NOT NULL,
    "CUSTOMER_NAME" text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.customers OWNER TO admin;

--
-- Name: inspection_logs_ID_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public."inspection_logs_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."inspection_logs_ID_seq" OWNER TO admin;

--
-- Name: inspection_logs_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public."inspection_logs_ID_seq" OWNED BY public.inspection_logs."ID";


--
-- Name: material_defect_ledger; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.material_defect_ledger (
    "ID" bigint NOT NULL,
    "INSPECTION_LOG_ID" bigint,
    "MATERIAL_ID" text,
    "DEFECT_ID" text,
    "QTY_NG" integer NOT NULL,
    logged_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.material_defect_ledger OWNER TO admin;

--
-- Name: material_defect_ledger_ID_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public."material_defect_ledger_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."material_defect_ledger_ID_seq" OWNER TO admin;

--
-- Name: material_defect_ledger_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public."material_defect_ledger_ID_seq" OWNED BY public.material_defect_ledger."ID";


--
-- Name: material_potential_defects; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.material_potential_defects (
    "MATERIAL_ID" text NOT NULL,
    "DEFECT_ID" text NOT NULL
);


ALTER TABLE public.material_potential_defects OWNER TO admin;

--
-- Name: seq_prod_id; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.seq_prod_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.seq_prod_id OWNER TO admin;

--
-- Name: products; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.products (
    "ID" text DEFAULT ('PROD-'::text || lpad((nextval('public.seq_prod_id'::regclass))::text, 5, '0'::text)) NOT NULL,
    "PART NO." text NOT NULL,
    "UNIQ NO" text,
    "PART NAME" text NOT NULL,
    "MODEL" text,
    "CUSTOMER_ID" text,
    "LINE_ID" text,
    "ASSY_NAME" text,
    "PICTURE_LOCATION" text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.products OWNER TO admin;

--
-- Name: seq_supp_id; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.seq_supp_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.seq_supp_id OWNER TO admin;

--
-- Name: suppliers; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.suppliers (
    "ID" text DEFAULT ('SUPP-'::text || lpad((nextval('public.seq_supp_id'::regclass))::text, 5, '0'::text)) NOT NULL,
    "SUPP" text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.suppliers OWNER TO admin;

--
-- Name: users; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    nip character varying(20) NOT NULL,
    password character varying(255) NOT NULL,
    name character varying(100) NOT NULL,
    role character varying(50) DEFAULT 'Leader QC'::character varying,
    remember_token character varying(100),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.users OWNER TO admin;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO admin;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: view_form_input_qc; Type: VIEW; Schema: public; Owner: admin
--

CREATE VIEW public.view_form_input_qc AS
 SELECT p."PART NAME",
    p."UNIQ NO",
    dm."NG_NAME",
    dc."CATEGORY_NAME"
   FROM ((((public.products p
     LEFT JOIN public.bill_of_materials bom ON ((p."ID" = bom."PRODUCT_ID")))
     LEFT JOIN public.material_potential_defects mpd ON ((bom."MATERIAL_ID" = mpd."MATERIAL_ID")))
     JOIN public.defect_master dm ON (((dm."ID" = mpd."DEFECT_ID") OR (dm."CATEGORY_ID" = ( SELECT defect_categories."ID"
           FROM public.defect_categories
          WHERE (defect_categories."CATEGORY_NAME" = 'PROCESS'::text))))))
     JOIN public.defect_categories dc ON ((dm."CATEGORY_ID" = dc."ID")))
  ORDER BY p."PART NAME", dc."CATEGORY_NAME";


ALTER VIEW public.view_form_input_qc OWNER TO admin;

--
-- Name: inspection_logs ID; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.inspection_logs ALTER COLUMN "ID" SET DEFAULT nextval('public."inspection_logs_ID_seq"'::regclass);


--
-- Name: material_defect_ledger ID; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.material_defect_ledger ALTER COLUMN "ID" SET DEFAULT nextval('public."material_defect_ledger_ID_seq"'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: MATERIAL MATERIAL_UNIQ_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public."MATERIAL"
    ADD CONSTRAINT "MATERIAL_UNIQ_key" UNIQUE ("UNIQ");


--
-- Name: MATERIAL MATERIAL_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public."MATERIAL"
    ADD CONSTRAINT "MATERIAL_pkey" PRIMARY KEY ("ID");


--
-- Name: bill_of_materials bill_of_materials_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.bill_of_materials
    ADD CONSTRAINT bill_of_materials_pkey PRIMARY KEY ("PRODUCT_ID", "MATERIAL_ID");


--
-- Name: customers customers_CUSTOMER_NAME_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT "customers_CUSTOMER_NAME_key" UNIQUE ("CUSTOMER_NAME");


--
-- Name: customers customers_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_pkey PRIMARY KEY ("ID");


--
-- Name: daily_inspections daily_inspections_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.daily_inspections
    ADD CONSTRAINT daily_inspections_pkey PRIMARY KEY ("ID");


--
-- Name: defect_categories defect_categories_CATEGORY_NAME_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.defect_categories
    ADD CONSTRAINT "defect_categories_CATEGORY_NAME_key" UNIQUE ("CATEGORY_NAME");


--
-- Name: defect_categories defect_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.defect_categories
    ADD CONSTRAINT defect_categories_pkey PRIMARY KEY ("ID");


--
-- Name: defect_master defect_master_NG_NAME_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.defect_master
    ADD CONSTRAINT "defect_master_NG_NAME_key" UNIQUE ("NG_NAME");


--
-- Name: defect_master defect_master_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.defect_master
    ADD CONSTRAINT defect_master_pkey PRIMARY KEY ("ID");


--
-- Name: inspection_logs inspection_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.inspection_logs
    ADD CONSTRAINT inspection_logs_pkey PRIMARY KEY ("ID");


--
-- Name: material_defect_ledger material_defect_ledger_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.material_defect_ledger
    ADD CONSTRAINT material_defect_ledger_pkey PRIMARY KEY ("ID");


--
-- Name: material_potential_defects material_potential_defects_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.material_potential_defects
    ADD CONSTRAINT material_potential_defects_pkey PRIMARY KEY ("MATERIAL_ID", "DEFECT_ID");


--
-- Name: production_lines production_lines_LINE_NAME_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.production_lines
    ADD CONSTRAINT "production_lines_LINE_NAME_key" UNIQUE ("LINE_NAME");


--
-- Name: production_lines production_lines_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.production_lines
    ADD CONSTRAINT production_lines_pkey PRIMARY KEY ("ID");


--
-- Name: products products_UNIQ NO_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT "products_UNIQ NO_key" UNIQUE ("UNIQ NO");


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY ("ID");


--
-- Name: suppliers suppliers_SUPP_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT "suppliers_SUPP_key" UNIQUE ("SUPP");


--
-- Name: suppliers suppliers_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_pkey PRIMARY KEY ("ID");


--
-- Name: users users_nip_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_nip_key UNIQUE (nip);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_bi_date; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_bi_date ON public.daily_inspections USING btree ("INSPECTION_DATE");


--
-- Name: idx_bi_week; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_bi_week ON public.daily_inspections USING btree ("STAT_WEEK");


--
-- Name: idx_bi_year_month; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_bi_year_month ON public.daily_inspections USING btree ("STAT_YEAR", "STAT_MONTH");


--
-- Name: idx_log_time; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_log_time ON public.inspection_logs USING btree ("OCCURRENCE_TIME");


--
-- Name: inspection_logs after_inspection_log_insert; Type: TRIGGER; Schema: public; Owner: admin
--

CREATE TRIGGER after_inspection_log_insert AFTER INSERT ON public.inspection_logs FOR EACH ROW EXECUTE FUNCTION public.trg_auto_trace_material_ng();


--
-- Name: MATERIAL MATERIAL_SUPP_ID_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public."MATERIAL"
    ADD CONSTRAINT "MATERIAL_SUPP_ID_fkey" FOREIGN KEY ("SUPP_ID") REFERENCES public.suppliers("ID") ON DELETE RESTRICT;


--
-- Name: bill_of_materials bill_of_materials_MATERIAL_ID_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.bill_of_materials
    ADD CONSTRAINT "bill_of_materials_MATERIAL_ID_fkey" FOREIGN KEY ("MATERIAL_ID") REFERENCES public."MATERIAL"("ID") ON DELETE RESTRICT;


--
-- Name: bill_of_materials bill_of_materials_PRODUCT_ID_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.bill_of_materials
    ADD CONSTRAINT "bill_of_materials_PRODUCT_ID_fkey" FOREIGN KEY ("PRODUCT_ID") REFERENCES public.products("ID") ON DELETE CASCADE;


--
-- Name: daily_inspections daily_inspections_LEADER_ID_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.daily_inspections
    ADD CONSTRAINT "daily_inspections_LEADER_ID_fkey" FOREIGN KEY ("LEADER_ID") REFERENCES public.users(id);


--
-- Name: daily_inspections daily_inspections_LINE_ID_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.daily_inspections
    ADD CONSTRAINT "daily_inspections_LINE_ID_fkey" FOREIGN KEY ("LINE_ID") REFERENCES public.production_lines("ID");


--
-- Name: daily_inspections daily_inspections_PRODUCT_ID_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.daily_inspections
    ADD CONSTRAINT "daily_inspections_PRODUCT_ID_fkey" FOREIGN KEY ("PRODUCT_ID") REFERENCES public.products("ID");


--
-- Name: defect_master defect_master_CATEGORY_ID_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.defect_master
    ADD CONSTRAINT "defect_master_CATEGORY_ID_fkey" FOREIGN KEY ("CATEGORY_ID") REFERENCES public.defect_categories("ID") ON DELETE RESTRICT;


--
-- Name: inspection_logs inspection_logs_DEFECT_ID_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.inspection_logs
    ADD CONSTRAINT "inspection_logs_DEFECT_ID_fkey" FOREIGN KEY ("DEFECT_ID") REFERENCES public.defect_master("ID");


--
-- Name: inspection_logs inspection_logs_INSPECTION_ID_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.inspection_logs
    ADD CONSTRAINT "inspection_logs_INSPECTION_ID_fkey" FOREIGN KEY ("INSPECTION_ID") REFERENCES public.daily_inspections("ID") ON DELETE CASCADE;


--
-- Name: material_defect_ledger material_defect_ledger_DEFECT_ID_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.material_defect_ledger
    ADD CONSTRAINT "material_defect_ledger_DEFECT_ID_fkey" FOREIGN KEY ("DEFECT_ID") REFERENCES public.defect_master("ID");


--
-- Name: material_defect_ledger material_defect_ledger_INSPECTION_LOG_ID_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.material_defect_ledger
    ADD CONSTRAINT "material_defect_ledger_INSPECTION_LOG_ID_fkey" FOREIGN KEY ("INSPECTION_LOG_ID") REFERENCES public.inspection_logs("ID") ON DELETE CASCADE;


--
-- Name: material_defect_ledger material_defect_ledger_MATERIAL_ID_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.material_defect_ledger
    ADD CONSTRAINT "material_defect_ledger_MATERIAL_ID_fkey" FOREIGN KEY ("MATERIAL_ID") REFERENCES public."MATERIAL"("ID") ON DELETE CASCADE;


--
-- Name: material_potential_defects material_potential_defects_DEFECT_ID_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.material_potential_defects
    ADD CONSTRAINT "material_potential_defects_DEFECT_ID_fkey" FOREIGN KEY ("DEFECT_ID") REFERENCES public.defect_master("ID") ON DELETE CASCADE;


--
-- Name: material_potential_defects material_potential_defects_MATERIAL_ID_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.material_potential_defects
    ADD CONSTRAINT "material_potential_defects_MATERIAL_ID_fkey" FOREIGN KEY ("MATERIAL_ID") REFERENCES public."MATERIAL"("ID") ON DELETE CASCADE;


--
-- Name: products products_CUSTOMER_ID_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT "products_CUSTOMER_ID_fkey" FOREIGN KEY ("CUSTOMER_ID") REFERENCES public.customers("ID") ON DELETE RESTRICT;


--
-- Name: products products_LINE_ID_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT "products_LINE_ID_fkey" FOREIGN KEY ("LINE_ID") REFERENCES public.production_lines("ID") ON DELETE RESTRICT;


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: admin
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

\unrestrict aOh4bAmuhZxTR1AjFq1BKYgFGIerLlVTAv99rQcWUeuaP4SxlMdzCdzEQAAkBPg

