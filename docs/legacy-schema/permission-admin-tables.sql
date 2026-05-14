-- Migração para criar tabelas faltantes do SPEC_ADMIN
-- Este script cria as tabelas system_groups, system_roles, system_programs e system_units
-- no banco de dados permission para o SMID 10

USE permission;

-- Criar tabela system_groups
CREATE TABLE IF NOT EXISTS system_groups (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(256),
    frontpage_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY(frontpage_id) REFERENCES system_program(id)
);

-- Criar tabela system_roles
CREATE TABLE IF NOT EXISTS system_roles (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(256),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Criar tabela system_programs
CREATE TABLE IF NOT EXISTS system_programs (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(256),
    controller VARCHAR(256),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Criar tabela system_units
CREATE TABLE IF NOT EXISTS system_units (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(256),
    parent_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
