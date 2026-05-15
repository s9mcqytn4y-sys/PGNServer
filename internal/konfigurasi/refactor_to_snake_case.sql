-- PGNServer Database Refactor to snake_case (FIXED)
-- Author: Senior Go Architect & QC Specialist

BEGIN;

-- 1. Rename Tables
ALTER TABLE IF EXISTS "MATERIAL" RENAME TO material;
ALTER TABLE IF EXISTS daily_inspections RENAME TO inspeksi_harian;
ALTER TABLE IF EXISTS inspection_logs RENAME TO log_inspeksi;
ALTER TABLE IF EXISTS material_defect_ledger RENAME TO buku_besar_defect_material;
ALTER TABLE IF EXISTS material_potential_defects RENAME TO pemetaan_defect_material;
ALTER TABLE IF EXISTS defect_categories RENAME TO kategori_defect;
ALTER TABLE IF EXISTS defect_master RENAME TO master_defect;
ALTER TABLE IF EXISTS production_lines RENAME TO lini_produksi;
ALTER TABLE IF EXISTS suppliers RENAME TO supplier;
ALTER TABLE IF EXISTS customers RENAME TO customer;

-- 2. Refactor products
ALTER TABLE products RENAME COLUMN "ID" TO id;
ALTER TABLE products RENAME COLUMN "PART NO." TO nomor_part;
ALTER TABLE products RENAME COLUMN "UNIQ NO" TO nomor_unik;
ALTER TABLE products RENAME COLUMN "PART NAME" TO nama_part;
ALTER TABLE products RENAME COLUMN "MODEL" TO model;
ALTER TABLE products RENAME COLUMN "CUSTOMER_ID" TO customer_id;
ALTER TABLE products RENAME COLUMN "LINE_ID" TO line_id;
ALTER TABLE products RENAME COLUMN "ASSY_NAME" TO nama_assy;
ALTER TABLE products RENAME COLUMN "PICTURE_LOCATION" TO lokasi_foto;

-- 3. Refactor material
ALTER TABLE material RENAME COLUMN "ID" TO id;
ALTER TABLE material RENAME COLUMN "SUPP_ID" TO supplier_id;
ALTER TABLE material RENAME COLUMN "UNIQ" TO nomor_unik;
ALTER TABLE material RENAME COLUMN "PART NAME" TO nama_part;
ALTER TABLE material RENAME COLUMN "TEBAL (MM)" TO tebal_mm;
ALTER TABLE material RENAME COLUMN "LEBAR (CM)" TO lebar_cm;
ALTER TABLE material RENAME COLUMN "PANJANG (CM)" TO panjang_cm;
ALTER TABLE material RENAME COLUMN "BERAT (GSM)" TO berat_gsm;
ALTER TABLE material RENAME COLUMN "MASSA (KG)" TO massa_kg;
ALTER TABLE material RENAME COLUMN "SATUAN" TO satuan;
ALTER TABLE material RENAME COLUMN "PICTURE_LOCATION" TO lokasi_foto;

-- 4. Refactor bill_of_materials
ALTER TABLE bill_of_materials RENAME COLUMN "PRODUCT_ID" TO produk_id;
ALTER TABLE bill_of_materials RENAME COLUMN "MATERIAL_ID" TO material_id;
ALTER TABLE bill_of_materials RENAME COLUMN "USAGE_QTY" TO jumlah_pemakaian;

-- 5. Refactor inspeksi_harian
ALTER TABLE inspeksi_harian RENAME COLUMN "ID" TO id;
ALTER TABLE inspeksi_harian RENAME COLUMN "INSPECTION_DATE" TO tanggal_inspeksi;
ALTER TABLE inspeksi_harian RENAME COLUMN "LEADER_ID" TO leader_id;
ALTER TABLE inspeksi_harian RENAME COLUMN "LINE_ID" TO line_id;
ALTER TABLE inspeksi_harian RENAME COLUMN "PRODUCT_ID" TO produk_id;
ALTER TABLE inspeksi_harian RENAME COLUMN "TOTAL_PRODUCED" TO total_produksi;
ALTER TABLE inspeksi_harian RENAME COLUMN "TOTAL_OK" TO total_ok;
ALTER TABLE inspeksi_harian RENAME COLUMN "TOTAL_NG" TO total_ng;
ALTER TABLE inspeksi_harian RENAME COLUMN "STAT_YEAR" TO stat_tahun;
ALTER TABLE inspeksi_harian RENAME COLUMN "STAT_MONTH" TO stat_bulan;
ALTER TABLE inspeksi_harian RENAME COLUMN "STAT_WEEK" TO stat_minggu;

-- 6. Refactor kategori_defect & master_defect
ALTER TABLE kategori_defect RENAME COLUMN "ID" TO id;
ALTER TABLE kategori_defect RENAME COLUMN "CATEGORY_NAME" TO nama_kategori;

ALTER TABLE master_defect RENAME COLUMN "ID" TO id;
ALTER TABLE master_defect RENAME COLUMN "CATEGORY_ID" TO kategori_id;
ALTER TABLE master_defect RENAME COLUMN "NG_NAME" TO nama_ng;

-- 7. Refactor log_inspeksi
ALTER TABLE log_inspeksi RENAME COLUMN "ID" TO id;
ALTER TABLE log_inspeksi RENAME COLUMN "INSPECTION_ID" TO inspeksi_id;
ALTER TABLE log_inspeksi RENAME COLUMN "DEFECT_ID" TO defect_id;
ALTER TABLE log_inspeksi RENAME COLUMN "TIME_WINDOW" TO jendela_waktu;
ALTER TABLE log_inspeksi RENAME COLUMN "OCCURRENCE_TIME" TO waktu_kejadian;
ALTER TABLE log_inspeksi RENAME COLUMN "QTY_NG" TO jumlah_ng;

-- 8. Refactor buku_besar_defect_material
ALTER TABLE buku_besar_defect_material RENAME COLUMN "ID" TO id;
ALTER TABLE buku_besar_defect_material RENAME COLUMN "INSPECTION_LOG_ID" TO log_inspeksi_id;
ALTER TABLE buku_besar_defect_material RENAME COLUMN "MATERIAL_ID" TO material_id;
ALTER TABLE buku_besar_defect_material RENAME COLUMN "DEFECT_ID" TO defect_id;
ALTER TABLE buku_besar_defect_material RENAME COLUMN "QTY_NG" TO jumlah_ng;
ALTER TABLE buku_besar_defect_material RENAME COLUMN logged_at TO dicatat_pada;

-- 9. Refactor pemetaan_defect_material
ALTER TABLE pemetaan_defect_material RENAME COLUMN "MATERIAL_ID" TO material_id;
ALTER TABLE pemetaan_defect_material RENAME COLUMN "DEFECT_ID" TO defect_id;

-- 10. Refactor other master tables
ALTER TABLE lini_produksi RENAME COLUMN "ID" TO id;
ALTER TABLE lini_produksi RENAME COLUMN "LINE_NAME" TO nama_lini;

ALTER TABLE supplier RENAME COLUMN "ID" TO id;
ALTER TABLE supplier RENAME COLUMN "SUPP" TO nama_supplier;

ALTER TABLE customer RENAME COLUMN "ID" TO id;
ALTER TABLE customer RENAME COLUMN "CUSTOMER_NAME" TO nama_customer;

-- Drop and Recreate Views and Triggers with new names
DROP VIEW IF EXISTS analytics_control_chart;
DROP VIEW IF EXISTS analytics_fishbone_logic;
DROP VIEW IF EXISTS analytics_histogram_time;
DROP VIEW IF EXISTS analytics_pareto_data;
DROP VIEW IF EXISTS analytics_stratification_trend;
DROP VIEW IF EXISTS view_form_input_qc;

-- Recreate Trigger Function
CREATE OR REPLACE FUNCTION public.trg_auto_trace_material_ng() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    v_category_name TEXT;
    v_product_id TEXT;
    v_material_id TEXT;
BEGIN
    SELECT dc.nama_kategori INTO v_category_name
    FROM master_defect dm
    JOIN kategori_defect dc ON dm.kategori_id = dc.id
    WHERE dm.id = NEW.defect_id;

    SELECT produk_id INTO v_product_id
    FROM inspeksi_harian
    WHERE id = NEW.inspeksi_id;

    IF v_category_name = 'MATERIAL' THEN
        FOR v_material_id IN
            SELECT bom.material_id
            FROM bill_of_materials bom
            JOIN pemetaan_defect_material mpd ON bom.material_id = mpd.material_id
            WHERE bom.produk_id = v_product_id AND mpd.defect_id = NEW.defect_id
        LOOP
            INSERT INTO buku_besar_defect_material (log_inspeksi_id, material_id, defect_id, jumlah_ng)
            VALUES (NEW.id, v_material_id, NEW.defect_id, NEW.jumlah_ng);
        END LOOP;
    END IF;
    RETURN NEW;
END;
$$;

-- Recreate Trigger on log_inspeksi (after rename)
DROP TRIGGER IF EXISTS after_inspection_log_insert ON log_inspeksi;
CREATE TRIGGER after_inspection_log_insert 
AFTER INSERT ON log_inspeksi 
FOR EACH ROW EXECUTE FUNCTION trg_auto_trace_material_ng();

-- Recreate Views
CREATE VIEW analytics_pareto_data AS
 WITH ng_summary AS (
         SELECT dm.nama_ng,
            sum(il.jumlah_ng) AS total_qty
           FROM (log_inspeksi il
             JOIN master_defect dm ON ((il.defect_id = dm.id)))
          GROUP BY dm.nama_ng
        ), pareto_calc AS (
         SELECT ng_summary.nama_ng,
            ng_summary.total_qty,
            sum(ng_summary.total_qty) OVER (ORDER BY ng_summary.total_qty DESC) AS cumulative_qty,
            sum(ng_summary.total_qty) OVER () AS grand_total
           FROM ng_summary
        )
 SELECT nama_ng,
    total_qty,
    round((((total_qty)::numeric / NULLIF(grand_total, (0)::numeric)) * (100)::numeric), 2) AS percentage,
    round(((cumulative_qty / NULLIF(grand_total, (0)::numeric)) * (100)::numeric), 2) AS cumulative_percentage
   FROM pareto_calc
  ORDER BY total_qty DESC;

CREATE VIEW analytics_stratification_trend AS
 SELECT di.stat_tahun,
    di.stat_bulan,
    di.stat_minggu,
    pl.nama_lini,
    sum(di.total_produksi) AS prod_total,
    sum(di.total_ng) AS ng_total,
    round((((sum(di.total_ng))::numeric / (NULLIF(sum(di.total_produksi), 0))::numeric) * (100)::numeric), 2) AS reject_rate_pct
   FROM (inspeksi_harian di
     JOIN lini_produksi pl ON ((di.line_id = pl.id)))
  GROUP BY di.stat_tahun, di.stat_bulan, di.stat_minggu, pl.nama_lini
  ORDER BY di.stat_tahun DESC, di.stat_bulan DESC, di.stat_minggu DESC;

COMMIT;
