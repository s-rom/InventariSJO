CREATE TYPE computer_type_enum AS ENUM ('Desktop', 'Laptop');
CREATE TYPE ram_type_enum      AS ENUM ('DDR3', 'DDR4', 'DDR5', 'None');
CREATE TYPE storage_type_enum  AS ENUM ('HDD', 'SSD', 'NVMe', 'None');
CREATE TYPE audit_event_enum   AS ENUM ('updated', 'deleted');
-- CREATE TYPE net_interface_type_enum AS ENUM ('eth', 'wifi', 'both');

CREATE TABLE app_user (
    app_user_id   BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username      TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL CHECK (length(password_hash) >= 60),

    -- Flags de permisos
    -- is_meta: pot crear usuaris i cambiar els permisos d'altres usuaris
    can_create    BOOLEAN NOT NULL DEFAULT FALSE, 
    can_update    BOOLEAN NOT NULL DEFAULT FALSE,
    can_delete    BOOLEAN NOT NULL DEFAULT FALSE,
    is_meta       BOOLEAN NOT NULL DEFAULT FALSE
);


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


CREATE TABLE cpu (
    cpu_id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    model_name      TEXT UNIQUE,
    benchmark_score INTEGER
);


CREATE TABLE os (
    os_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name  TEXT NOT NULL UNIQUE
);


CREATE TABLE equipment_user (
    equipment_user_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name              TEXT NOT NULL UNIQUE
);


CREATE TABLE computer (
    computer_id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    hostname               TEXT NOT NULL UNIQUE,
    cpu_id                 BIGINT REFERENCES cpu (cpu_id),
    ram_gb                 INTEGER NOT NULL CHECK (ram_gb >= 0), -- pot ser 0 si es al taller
    ram_type               ram_type_enum NOT NULL DEFAULT 'None', -- ddr2, 3, 4 o none
    storage_gb             INTEGER NOT NULL CHECK (storage_gb >= 0), -- pot ser 0 si es al taller
    storage_type           storage_type_enum NOT NULL DEFAULT 'None', -- ssd, hdd, nvme o none
    computer_type          computer_type_enum NOT NULL DEFAULT 'Desktop', -- desktop o laptop
    observations           TEXT DEFAULT NULL,admin
 
    -- Opcional en principi, des de l'app s'hauría d'assignar un usuari i 
    -- una sala "default" millor
    equipment_user_id      BIGINT REFERENCES equipment_user (equipment_user_id),
    room_id                BIGINT REFERENCES room (room_id),
    --------------------------------------------------------------------------

    -- net_interface_type     net_interface_type_enum DEFAULT 'eth',
    mac_address            TEXT UNIQUE DEFAULT NULL,

    -- model_id/name     
    -- Es podría tenir un taula Model per ficar els models dels portàtils o be fer una herència
    -- laptop,desktop --> computer
    -- On: 
    -- computer(id, hostname, ram_gb, os_id, observations, ..)
    -- laptop(laptop_id, computer_id, model_name)
    -- dekstop(desk_id, computer_id, cpu_id, ram_type, storage_type, ...)

    created_by_app_user_id BIGINT NOT NULL REFERENCES app_user (app_user_id),
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON computer (cpu_id);
CREATE INDEX ON computer (equipment_user_id);
CREATE INDEX ON computer (room_id);
CREATE INDEX ON computer (created_by_app_user_id);


CREATE TABLE computer_os (
    computer_id BIGINT NOT NULL REFERENCES computer (computer_id) ON DELETE CASCADE,
    os_id       BIGINT NOT NULL REFERENCES os (os_id) ON DELETE CASCADE,

    PRIMARY KEY (computer_id, os_id)
);

CREATE INDEX ON computer_os (os_id);

CREATE TABLE computer_audit (
    audit_id               BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    event_type             audit_event_enum NOT NULL, -- updated o deleted
    computer_id            BIGINT NOT NULL,
    old_values             JSONB NOT NULL,
    new_values             JSONB,       
    changed_by_app_user_id BIGINT NOT NULL,
    changed_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON computer_audit (computer_id, changed_at);
CREATE INDEX ON computer_audit (changed_by_app_user_id);
CREATE INDEX ON computer_audit (changed_at);
CREATE INDEX ON computer_audit (event_type);


INSERT INTO app_user (username, password_hash, can_create, can_update, can_delete, is_meta)
VALUES ('sergi', '$2a$10$43Hib5uNSgnM9NCwZ0VulOr8JU0TGl1hPn1G.cGD5Q/kvvPHHYxNC', true, true, true, true);

INSERT INTO cpu (model_name, benchmark_score) VALUES
('Intel Core i3-4130 3.40GHz',                3311),
('Intel Core i3-6100 3.70GHz',                4151),
('Intel Core i3-2100 3.10GHz',                1845),
('Intel Celeron G550 2.60GHz',                1214),
('Intel Core i5-750 2.67GHz',                 2512),
('Intel Core i3-3240 3.40GHz',                2304),
('Intel Celeron G1610 2.60GHz',               1508),
('Intel Core i5-4440',                        4742),
('Intel Pentium G3260 3.30GHz',               2068),
('Intel Core i3-2105 3.10GHz',                1883),
('Intel Core i5-4590 3.30GHz',                5354),
('Intel Core i5-6500 3.20GHz',                5644),
('Intel Core i5-2400 3.10GHz',                3855),
('Intel Core i3-7100 3.90GHz',                4336),
('AMD Ryzen 5 Pro 5650U',                    14719);