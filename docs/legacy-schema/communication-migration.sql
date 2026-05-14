-- Migration script para criar tabelas de comunicação no banco de dados communication
-- Este script foi aplicado manualmente na VPS em 2026-05-14

CREATE DATABASE IF NOT EXISTS communication;

USE communication;

-- Tabelas criadas manualmente com AUTO_INCREMENT
CREATE TABLE IF NOT EXISTS system_notification (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    system_user_id INT,
    system_user_to_id INT,
    subject VARCHAR(256),
    message TEXT,
    dt_message VARCHAR(20),
    action_url TEXT,
    action_label VARCHAR(256),
    icon VARCHAR(100),
    checked CHAR(1)
);

CREATE TABLE IF NOT EXISTS system_message (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    system_user_id INT,
    system_user_to_id INT,
    subject VARCHAR(256),
    message TEXT,
    dt_message VARCHAR(20),
    checked CHAR(1),
    removed CHAR(1),
    viewed CHAR(1),
    attachments TEXT
);
