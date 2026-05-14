-- Migration script para criar tabelas de log no banco de dados log
-- Este script foi aplicado manualmente na VPS em 2026-05-14

CREATE DATABASE IF NOT EXISTS log;

USE log;

-- Tabelas criadas manualmente com AUTO_INCREMENT
CREATE TABLE IF NOT EXISTS system_access_log (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    sessionid VARCHAR(256),
    login VARCHAR(256),
    login_time VARCHAR(20),
    login_year VARCHAR(4),
    login_month VARCHAR(2),
    login_day VARCHAR(2),
    logout_time VARCHAR(20),
    impersonated CHAR(1),
    access_ip VARCHAR(45),
    impersonated_by VARCHAR(200)
);

CREATE TABLE IF NOT EXISTS system_change_log (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    logdate VARCHAR(20),
    login VARCHAR(256),
    tablename VARCHAR(256),
    primarykey VARCHAR(256),
    pkvalue VARCHAR(256),
    operation VARCHAR(256),
    columnname VARCHAR(256),
    oldvalue TEXT,
    newvalue TEXT,
    access_ip VARCHAR(256),
    transaction_id VARCHAR(256),
    log_trace TEXT,
    session_id VARCHAR(256),
    class_name VARCHAR(256),
    php_sapi VARCHAR(256),
    log_year VARCHAR(4),
    log_month VARCHAR(2),
    log_day VARCHAR(2)
);

CREATE TABLE IF NOT EXISTS system_sql_log (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    logdate VARCHAR(20),
    login VARCHAR(256),
    database_name VARCHAR(256),
    sql_command TEXT,
    statement_type VARCHAR(256),
    access_ip VARCHAR(45),
    transaction_id VARCHAR(256),
    log_trace TEXT,
    session_id VARCHAR(256),
    class_name VARCHAR(256),
    php_sapi VARCHAR(256),
    request_id VARCHAR(256),
    log_year VARCHAR(4),
    log_month VARCHAR(2),
    log_day VARCHAR(2)
);

CREATE TABLE IF NOT EXISTS system_request_log (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    endpoint TEXT,
    logdate VARCHAR(256),
    log_year VARCHAR(4),
    log_month VARCHAR(2),
    log_day VARCHAR(2),
    session_id VARCHAR(256),
    login VARCHAR(256),
    access_ip VARCHAR(256),
    class_name VARCHAR(256),
    class_method VARCHAR(256),
    http_host VARCHAR(256),
    server_port VARCHAR(256),
    request_uri TEXT,
    request_method VARCHAR(256),
    query_string TEXT,
    request_headers TEXT,
    request_body TEXT,
    request_duration INT
);

CREATE TABLE IF NOT EXISTS system_access_notification_log (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    login VARCHAR(256),
    email VARCHAR(256),
    ip_address VARCHAR(256),
    login_time VARCHAR(256)
);

CREATE TABLE IF NOT EXISTS system_schedule_log (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    logdate VARCHAR(19),
    title VARCHAR(256),
    class_name VARCHAR(256),
    method VARCHAR(256),
    status VARCHAR(1),
    message TEXT
);

CREATE TABLE IF NOT EXISTS system_sql_changes (
    id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    db_name VARCHAR(200),
    sql_date VARCHAR(20),
    sql_hash VARCHAR(32),
    sql_command TEXT
);
