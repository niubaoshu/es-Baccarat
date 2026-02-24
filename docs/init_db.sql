-- EZ Baccarat Development Environment Database Setup
-- Run these commands as the 'postgres' superuser: `sudo -u postgres psql`

-- 1. Create the dedicated gRPC / Go backend user
CREATE ROLE baccarat_admin WITH LOGIN PASSWORD 'EZBaccarat2026';

-- 2. Create the exclusive database for the multiplayer lobby and transactions
CREATE DATABASE es_baccarat;

-- 3. Hand over ownership of the database
ALTER DATABASE es_baccarat OWNER TO baccarat_admin;

-- Credentials Reference:
-- Host: localhost
-- Port: 5432
-- User: baccarat_admin
-- Pass: EZBaccarat2026
-- DB:   es_baccarat
