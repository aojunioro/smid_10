-- Migration: 001_create_admin_tables.sql
-- Description: Cria tabelas do domínio Admin no banco permission
-- Banco: permission

-- Tabela de usuários
CREATE TABLE IF NOT EXISTS system_users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    login VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    system_unit_id BIGINT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    frontpage_id BIGINT,
    equipe_id BIGINT,
    phone VARCHAR(20),
    address TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_login (login),
    INDEX idx_email (email),
    INDEX idx_active (active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabela de grupos
CREATE TABLE IF NOT EXISTS system_groups (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    frontpage_id BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabela de papéis
CREATE TABLE IF NOT EXISTS system_roles (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabela de programas
CREATE TABLE IF NOT EXISTS system_programs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    controller VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabela de unidades
CREATE TABLE IF NOT EXISTS system_units (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    parent_id BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_parent_id (parent_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabela pivô: usuários-grupos
CREATE TABLE IF NOT EXISTS system_user_group (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    system_user_id BIGINT NOT NULL,
    system_group_id BIGINT NOT NULL,
    UNIQUE KEY uk_user_group (system_user_id, system_group_id),
    FOREIGN KEY (system_user_id) REFERENCES system_users(id) ON DELETE CASCADE,
    FOREIGN KEY (system_group_id) REFERENCES system_groups(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabela pivô: usuários-papéis
CREATE TABLE IF NOT EXISTS system_user_role (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    system_user_id BIGINT NOT NULL,
    system_role_id BIGINT NOT NULL,
    UNIQUE KEY uk_user_role (system_user_id, system_role_id),
    FOREIGN KEY (system_user_id) REFERENCES system_users(id) ON DELETE CASCADE,
    FOREIGN KEY (system_role_id) REFERENCES system_roles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabela pivô: usuários-unidades
CREATE TABLE IF NOT EXISTS system_user_unit (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    system_user_id BIGINT NOT NULL,
    system_unit_id BIGINT NOT NULL,
    UNIQUE KEY uk_user_unit (system_user_id, system_unit_id),
    FOREIGN KEY (system_user_id) REFERENCES system_users(id) ON DELETE CASCADE,
    FOREIGN KEY (system_unit_id) REFERENCES system_units(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabela pivô: grupos-programas
CREATE TABLE IF NOT EXISTS system_group_program (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    system_group_id BIGINT NOT NULL,
    system_program_id BIGINT NOT NULL,
    UNIQUE KEY uk_group_program (system_group_id, system_program_id),
    FOREIGN KEY (system_group_id) REFERENCES system_groups(id) ON DELETE CASCADE,
    FOREIGN KEY (system_program_id) REFERENCES system_programs(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabela pivô: grupos-papéis
CREATE TABLE IF NOT EXISTS system_group_role (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    system_group_id BIGINT NOT NULL,
    system_role_id BIGINT NOT NULL,
    UNIQUE KEY uk_group_role (system_group_id, system_role_id),
    FOREIGN KEY (system_group_id) REFERENCES system_groups(id) ON DELETE CASCADE,
    FOREIGN KEY (system_role_id) REFERENCES system_roles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabela pivô: programas-métodos-papéis (permissões granulares)
CREATE TABLE IF NOT EXISTS system_program_method_role (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    system_program_id BIGINT NOT NULL,
    method_name VARCHAR(100) NOT NULL,
    system_role_id BIGINT NOT NULL,
    UNIQUE KEY uk_program_method_role (system_program_id, method_name, system_role_id),
    FOREIGN KEY (system_program_id) REFERENCES system_programs(id) ON DELETE CASCADE,
    FOREIGN KEY (system_role_id) REFERENCES system_roles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabela de preferências globais
CREATE TABLE IF NOT EXISTS system_preferences (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    attribute VARCHAR(100) NOT NULL UNIQUE,
    value TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
