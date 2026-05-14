-- Migration script para criar tabelas de tarefas no banco de dados smid
-- Este script foi aplicado manualmente na VPS em 2026-05-14

CREATE DATABASE IF NOT EXISTS smid;

USE smid;

-- Tabela criada manualmente com AUTO_INCREMENT
CREATE TABLE IF NOT EXISTS tarefas (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    tarefa VARCHAR(300) NOT NULL,
    dt_tarefa DATE NOT NULL,
    hr_tarefa TIME NOT NULL,
    status CHAR(1) NOT NULL DEFAULT 'N',
    login VARCHAR(255) NOT NULL,
    lead_id INT
);
