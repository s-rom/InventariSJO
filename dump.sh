podman exec inventari_db_1 pg_dump --data-only -U admin inventaridb > inserts.sql
