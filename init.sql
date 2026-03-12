CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ============================================================
-- ENUM TYPES
-- ============================================================
-- 'None' = genuïnament sense RAM/disc (p.ex. equip al taller sense memoria)
-- NULL   = no modificat respecte al model base (laptop_model / desktop_model)
-- Ambdós valors s'utilitzen en contextos diferents; vegeu comentaris de desktop i laptop.
CREATE TYPE ram_type_enum       AS ENUM ('DDR3', 'DDR4', 'DDR5', 'None');
CREATE TYPE storage_type_enum   AS ENUM ('HDD', 'SSD', 'NVMe', 'None');
CREATE TYPE audit_event_enum    AS ENUM ('created', 'updated', 'deleted');

-- Torn per als portàtils d'alumnes: matí o tarda
CREATE TYPE shift_enum AS ENUM ('morning', 'afternoon');


-- ============================================================
-- ROLS D'USUARI
-- Taula en comptes d'ENUM per permetre afegir rols sense migració
-- d'esquema. Es fa servir role_id TEXT com a PK descriptiva.
-- ============================================================
CREATE TABLE role (
    role_id     TEXT PRIMARY KEY,
    description TEXT NOT NULL
);


-- ============================================================
-- USUARIS DE L'APLICACIÓ
-- ============================================================
CREATE TABLE app_user (
    app_user_id   BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username      TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL CHECK (length(password_hash) >= 60),
    role_id       TEXT NOT NULL DEFAULT 'readonly' REFERENCES role (role_id) ON UPDATE CASCADE
);


-- ============================================================
-- CENTRES I AULES
-- ============================================================
CREATE TABLE center (
    center_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name      TEXT NOT NULL UNIQUE
);

CREATE TABLE room (
    room_id   BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    center_id BIGINT NOT NULL REFERENCES center (center_id) ON DELETE CASCADE,
    name      TEXT NOT NULL,

    UNIQUE (center_id, name)
);

CREATE INDEX ON room (center_id);


-- ============================================================
-- HARDWARE: CPU I OS
-- ============================================================
CREATE TABLE cpu (
    cpu_id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    model_name      TEXT NOT NULL UNIQUE,
    benchmark_score INTEGER
);

CREATE TABLE os (
    os_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name  TEXT NOT NULL UNIQUE
);


-- ============================================================
-- USUARIS SIMPLES D'EQUIP
-- Per a sobretaules i portàtils d'ús individual no-estudiantil
-- (professors, personal, etc.)
-- ============================================================
CREATE TABLE equipment_user (
    equipment_user_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name              TEXT NOT NULL UNIQUE
);


-- ============================================================
-- MARQUES (compartida entre portàtils i sobretaules)
-- ============================================================
CREATE TABLE brand (
    brand_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name     TEXT NOT NULL UNIQUE   -- 'HP', 'Lenovo', 'Acer', 'Dell', ...
);


-- ============================================================
-- MODELS DE PORTÀTIL
-- Representem la configuració de fàbrica per evitar duplicació.
-- Les unitats concretes (taula laptop) emmagatzemen les possibles
-- ampliacions respecte a aquest model base.
-- ============================================================
CREATE TABLE laptop_model (
    laptop_model_id   BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    brand_id          BIGINT NOT NULL REFERENCES brand (brand_id),
    model_name        TEXT NOT NULL,
    cpu_id            BIGINT REFERENCES cpu (cpu_id),
    base_ram_gb       INTEGER NOT NULL DEFAULT 0 CHECK (base_ram_gb >= 0),
    base_ram_type     ram_type_enum NOT NULL DEFAULT 'None',
    base_storage_gb   INTEGER NOT NULL DEFAULT 0 CHECK (base_storage_gb >= 0),
    base_storage_type storage_type_enum NOT NULL DEFAULT 'None',
    -- SO natiu del model (p. ex. ChromeOS, Windows); pot haver canviat en la unitat concreta
    base_os_id        BIGINT REFERENCES os (os_id),

    UNIQUE (brand_id, model_name)
);

CREATE INDEX ON laptop_model (brand_id);
CREATE INDEX ON laptop_model (cpu_id);


-- ============================================================
-- MODELS DE SOBRETAULA (infreqüent, per models estandarditzats)
-- Igual que laptop_model: la unitat concreta (desktop) referencia
-- el model i només emmagatzema les diferències. Si desktop_model_id
-- és NULL, els camps de specs de desktop són l'especificació real.
-- ============================================================
CREATE TABLE desktop_model (
    desktop_model_id  BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    brand_id          BIGINT NOT NULL REFERENCES brand (brand_id),
    model_name        TEXT NOT NULL,
    cpu_id            BIGINT REFERENCES cpu (cpu_id),
    base_ram_gb       INTEGER NOT NULL DEFAULT 0 CHECK (base_ram_gb >= 0),
    base_ram_type     ram_type_enum NOT NULL DEFAULT 'None',
    base_storage_gb   INTEGER NOT NULL DEFAULT 0 CHECK (base_storage_gb >= 0),
    base_storage_type storage_type_enum NOT NULL DEFAULT 'None',
    base_os_id        BIGINT REFERENCES os (os_id),

    UNIQUE (brand_id, model_name)
);

CREATE INDEX ON desktop_model (brand_id);
CREATE INDEX ON desktop_model (cpu_id);


-- ============================================================
-- EQUIP BASE
-- Camps comuns a tots els equips (sobretaules i portàtils).
-- ============================================================
CREATE TABLE computer (
    computer_id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    hostname               TEXT NOT NULL UNIQUE,
    room_id                BIGINT REFERENCES room (room_id),
    observations           TEXT DEFAULT NULL,
    created_by_app_user_id BIGINT NOT NULL REFERENCES app_user (app_user_id),
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON computer (room_id);
CREATE INDEX ON computer (created_by_app_user_id);


-- ============================================================
-- SOBRETAULA
-- Herència de taula concreta (Class Table Inheritance).
--
-- Si desktop_model_id és NOT NULL, els camps cpu_id / ram_* /
-- storage_* / os_id amb valor NULL s'interpreten com "igual que
-- el model base". Si desktop_model_id és NULL, els camps són
-- l'especificació real de la unitat.
-- ============================================================
CREATE TABLE desktop (
    computer_id       BIGINT PRIMARY KEY REFERENCES computer (computer_id) ON DELETE CASCADE,
    desktop_model_id  BIGINT DEFAULT NULL REFERENCES desktop_model (desktop_model_id),
    cpu_id            BIGINT DEFAULT NULL REFERENCES cpu (cpu_id),
    ram_gb            INTEGER DEFAULT NULL CHECK (ram_gb IS NULL OR ram_gb >= 0),
    ram_type          ram_type_enum DEFAULT NULL,
    storage_gb        INTEGER DEFAULT NULL CHECK (storage_gb IS NULL OR storage_gb >= 0),
    storage_type      storage_type_enum DEFAULT NULL,
    os_id             BIGINT DEFAULT NULL REFERENCES os (os_id),
    equipment_user_id BIGINT REFERENCES equipment_user (equipment_user_id),

    -- Targeta de xarxa wifi addicional: MAC s'apunta quan n'hi ha
    has_wifi_card     BOOLEAN NOT NULL DEFAULT FALSE,
    mac_address       TEXT UNIQUE DEFAULT NULL,

    CONSTRAINT desktop_wifi_mac_check CHECK (
        has_wifi_card = TRUE OR mac_address IS NULL
    )
);

CREATE INDEX ON desktop (desktop_model_id);
CREATE INDEX ON desktop (cpu_id);
CREATE INDEX ON desktop (equipment_user_id);


-- ============================================================
-- PORTÀTIL
-- Herència de taula concreta.
-- Els camps de specs (ram_gb, etc.) emmagatzemen l'estat actual
-- de la unitat; NULL = no modificat respecte al model base.
--
-- Assignació d'usuari (mútuament excloents):
--   · equipment_user_id NOT NULL  → portàtil de professor, personal, etc.
--                                   (no estudiantil: cap fila a laptop_student_assignment)
--   · equipment_user_id NULL      → portàtil estudiantil: l'alumne/es s'assigna
--                                   via laptop_student_assignment (morning / afternoon)
-- ============================================================
CREATE TABLE laptop (
    computer_id       BIGINT PRIMARY KEY REFERENCES computer (computer_id) ON DELETE CASCADE,
    laptop_model_id   BIGINT NOT NULL REFERENCES laptop_model (laptop_model_id),

    -- NULL = igual que laptop_model.base_*
    ram_gb            INTEGER DEFAULT NULL CHECK (ram_gb IS NULL OR ram_gb >= 0),
    ram_type          ram_type_enum DEFAULT NULL,
    storage_gb        INTEGER DEFAULT NULL CHECK (storage_gb IS NULL OR storage_gb >= 0),
    storage_type      storage_type_enum DEFAULT NULL,
    mac_address       TEXT UNIQUE DEFAULT NULL,

    -- NULL = igual que laptop_model.base_os_id
    os_id             BIGINT DEFAULT NULL REFERENCES os (os_id),
    equipment_user_id BIGINT DEFAULT NULL REFERENCES equipment_user (equipment_user_id)
);

CREATE INDEX ON laptop (laptop_model_id);
CREATE INDEX ON laptop (equipment_user_id);


-- ============================================================
-- ESTRUCTURA EDUCATIVA: CICLES I CLASSES
-- Un cicle (SMX, ASIX, GEAD, FPB-Informatica...) pot tenir
-- diversos cursos (1r, 2n) i, dins d'un curs, diverses classes
-- (1A, 1B...). Cada classe té un tutor (app_user) que pot
-- gestionar les assignacions de portàtils del seu grup.
-- ============================================================
CREATE TABLE cycle (
    cycle_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name     TEXT NOT NULL UNIQUE   -- 'SMX', 'ASIX', 'GEAD', ...
);

CREATE TABLE school_class (
    class_id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    cycle_id          BIGINT NOT NULL REFERENCES cycle (cycle_id),
    course            SMALLINT NOT NULL CHECK (course >= 1),  -- any: 1, 2, ...
    -- 'A', 'B'... per classes paral·leles. Cicles amb un sol grup usen 'A'.
    -- NOT NULL per evitar el bug de UNIQUE amb NULL (NULL != NULL en SQL).
    class_label       TEXT NOT NULL DEFAULT 'A',
    -- El torn és propietat de la classe: ASIX 1r → tarda, GEAD 1rA → matí, etc.
    shift             shift_enum NOT NULL,
    tutor_app_user_id BIGINT REFERENCES app_user (app_user_id),

    UNIQUE (cycle_id, course, class_label)
);

CREATE INDEX ON school_class (cycle_id);
CREATE INDEX ON school_class (tutor_app_user_id);


-- ============================================================
-- ALUMNES
-- Cada alumne pertany a una classe concreta, que determina el seu
-- cicle, curs, grup i torn actual. La classe històrica en cada
-- assignació de portàtil es captura a laptop_student_assignment.
-- ============================================================
CREATE TABLE student (
    student_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    full_name  TEXT NOT NULL,
    -- Classe actual de l'alumne (per consultes ràpides del curs en curs)
    class_id   BIGINT NOT NULL REFERENCES school_class (class_id)
);

CREATE INDEX ON student (class_id);


-- ============================================================
-- ASSIGNACIONS DE PORTÀTIL A ALUMNES
--
-- Cada portàtil estudiantil pot tenir fins a dues assignacions per
-- any acadèmic: una per torn morning i una per torn afternoon.
-- El torn no s'emmagatzema aquí perquè ja és un atribut de la
-- classe de l'alumne (student → school_class.shift).
--
-- class_id: captura la classe de l'alumne en el moment de
--   l'assignació per preservar l'històric (student.class_id pot
--   canviar cada curs, aquí queda la classe real de cada any).
--
-- Garanties a nivell de DB:
--   · UNIQUE (computer_id, student_id, academic_year): un alumne
--     no pot tenir el mateix portàtil assignat dues vegades al mateix any.
--   · UNIQUE (student_id, academic_year): un alumne només pot tenir
--     un portàtil assignat per any acadèmic.
--
-- Garantia a nivell d'aplicació:
--   · No pot haver-hi dos alumnes del mateix torn (school_class.shift)
--     assignats al mateix portàtil en el mateix any.
--
-- academic_year: format '2025-2026'
-- El tutor de la classe pot modificar les assignacions del seu grup.
-- ============================================================
CREATE TABLE laptop_student_assignment (
    assignment_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    -- Referencia laptop, no computer, per garantir que només s'assignen portàtils
    computer_id   BIGINT NOT NULL REFERENCES laptop (computer_id) ON DELETE CASCADE,
    student_id    BIGINT NOT NULL REFERENCES student (student_id),
    -- Snapshot de la classe de l'alumne en aquest any acadèmic
    class_id      BIGINT NOT NULL REFERENCES school_class (class_id),
    academic_year TEXT NOT NULL CHECK (academic_year ~ '^\d{4}-\d{4}$'),

    -- Un alumne només pot tenir un portàtil per any acadèmic
    UNIQUE (student_id, academic_year)
);

CREATE INDEX ON laptop_student_assignment (computer_id);
CREATE INDEX ON laptop_student_assignment (student_id);
CREATE INDEX ON laptop_student_assignment (class_id);
CREATE INDEX ON laptop_student_assignment (academic_year);


-- ============================================================
-- AUDITORIA GENERAL
-- Cobreix qualsevol taula del model: computer, desktop, laptop,
-- laptop_student_assignment, student, etc.
--
-- table_name : nom de la taula afectada (vegeu CHECK més avall)
-- record_id  : PK de la fila afectada (sempre BIGINT en aquest esquema)
-- event_type : 'created', 'updated' o 'deleted'
-- old_values : estat anterior (NULL en 'created')
-- new_values : estat posterior (NULL en 'deleted')
-- changed_by_username : snapshot del username en el moment de l'esdeveniment
--   (immutable encara que l'usuari canviï de nom o sigui eliminat)
-- ============================================================
CREATE TABLE audit_log (
    audit_id               BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    table_name             TEXT NOT NULL,
    record_id              BIGINT NOT NULL,
    event_type             audit_event_enum NOT NULL,
    old_values             JSONB,
    new_values             JSONB,
    changed_by_app_user_id BIGINT NOT NULL REFERENCES app_user (app_user_id),
    changed_by_username    TEXT NOT NULL,
    changed_at             TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT audit_log_table_name_check CHECK (
        table_name IN (
            'computer', 'desktop', 'laptop',
            'desktop_model', 'laptop_model', 'brand',
            'laptop_student_assignment', 'student', 'school_class', 'cycle',
            'equipment_user', 'room', 'center',
            'cpu', 'os', 'app_user', 'role'
        )
    ),
    CONSTRAINT audit_log_values_check CHECK (
        CASE event_type
            WHEN 'created' THEN new_values IS NOT NULL AND old_values IS NULL
            WHEN 'deleted' THEN old_values IS NOT NULL AND new_values IS NULL
            WHEN 'updated' THEN old_values IS NOT NULL AND new_values IS NOT NULL
        END
    )
);

-- Índex compost que cobreix també les consultes per (table_name, record_id) sense changed_at
CREATE INDEX ON audit_log (table_name, record_id, changed_at);
CREATE INDEX ON audit_log (changed_by_app_user_id);
CREATE INDEX ON audit_log (changed_at);
CREATE INDEX ON audit_log (event_type);


-- ============================================================
-- DADES INICIALS
-- ============================================================
INSERT INTO role (role_id, description) VALUES
('readonly', 'Només lectura de totes les dades'),
('editor',   'Pot crear i modificar equips, models, sales i usuaris d''equip'),
('admin',    'Accés total i gestió d''usuaris de l''aplicació'),
('tutor',    'Pot modificar alumnes i assignacions de portàtils de les seves classes');

INSERT INTO app_user (username, password_hash, role_id)
VALUES ('sergi', '$2a$10$43Hib5uNSgnM9NCwZ0VulOr8JU0TGl1hPn1G.cGD5Q/kvvPHHYxNC', 'admin');

INSERT INTO cpu (model_name, benchmark_score) VALUES
('Intel Core i3-4130 3.40GHz',  3311),
('Intel Core i3-6100 3.70GHz',  4151),
('Intel Core i3-2100 3.10GHz',  1845),
('Intel Celeron G550 2.60GHz',  1214),
('Intel Core i5-750 2.67GHz',   2512),
('Intel Core i3-3240 3.40GHz',  2304),
('Intel Celeron G1610 2.60GHz', 1508),
('Intel Core i5-4440',          4742),
('Intel Pentium G3260 3.30GHz', 2068),
('Intel Core i3-2105 3.10GHz',  1883),
('Intel Core i5-4590 3.30GHz',  5354),
('Intel Core i5-6500 3.20GHz',  5644),
('Intel Core i5-2400 3.10GHz',  3855),
('Intel Core i3-7100 3.90GHz',  4336),
('AMD Ryzen 5 Pro 5650U',      14719);

INSERT INTO os (name) VALUES
('Windows 10'),
('Windows 11'),
('ChromeOs Flex'),
('ChromeOS');

INSERT INTO cycle (name) VALUES
('SMX'),
('ASIX'),
('GEAD'),
('FPB-Informatica'),
('DAW'),
('DAM'),
('FPB-Administracio');

-- Usuari administrador per defecte
INSERT INTO role (role_id, description) VALUES ('admin', 'Administrador') ON CONFLICT DO NOTHING;
INSERT INTO app_user (username, password_hash, role_id)
VALUES ('admin', crypt('admin', gen_salt('bf')), 'admin')
ON CONFLICT (username) DO NOTHING;