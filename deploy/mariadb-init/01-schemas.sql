-- Cria os 4 bancos canônicos exigidos pelo SMID 10
-- (espelha os bancos legados do SMID 8.x para testes funcionais).
-- Executado pelo entrypoint do MariaDB apenas na inicialização do volume vazio.

CREATE DATABASE IF NOT EXISTS `smid`
  CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE DATABASE IF NOT EXISTS `permission`
  CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE DATABASE IF NOT EXISTS `log`
  CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE DATABASE IF NOT EXISTS `communication`
  CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Usuário de aplicação (a senha vem da variável MYSQL_USER_PASSWORD no entrypoint).
-- Concede acesso só aos 4 bancos.
GRANT ALL PRIVILEGES ON `smid`.*          TO 'smid10'@'%';
GRANT ALL PRIVILEGES ON `permission`.*    TO 'smid10'@'%';
GRANT ALL PRIVILEGES ON `log`.*           TO 'smid10'@'%';
GRANT ALL PRIVILEGES ON `communication`.* TO 'smid10'@'%';
FLUSH PRIVILEGES;
