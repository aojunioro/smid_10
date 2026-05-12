-- VIEW para calendário de visitas
-- Database: indigita_adianti_teste
-- Usa system_unit do banco indigita_adianti_permission (cross-database)

USE indigita_adianti_teste;

DROP VIEW IF EXISTS v_calendario_visitas;

CREATE VIEW v_calendario_visitas AS
SELECT
    v.id AS visita_id,
    l.id AS lead_id,
    l.nome AS lead_nome,
    l.status_id,
    s.stt_lead,
    s.cor AS status_cor,
    e.bairro,
    e.cidade,
    v.dt_visita,
    v.hr_visita,
    CONCAT(v.dt_visita, ' ', v.hr_visita) AS start_time,
    CONCAT(v.dt_visita, ' ', DATE_ADD(v.hr_visita, INTERVAL 1 HOUR)) AS end_time,
    l.unidd_id,
    u.name AS unidade_nome,
    COALESCE(u.cor, '#007bff') AS unidade_cor,
    v.login_repre,
    v.status_id AS status_visita
FROM visitas v
LEFT JOIN leads l ON l.id = v.lead_id
LEFT JOIN endereco e ON l.id = e.lead_id
LEFT JOIN lead_status s ON l.status_id = s.id
LEFT JOIN indigita_adianti_permission.system_unit u ON l.unidd_id = u.id
WHERE l.status_id IN (2, 3, 4)
  AND v.stts_lead IN (2, 3, 4)
  AND v.dt_visita BETWEEN CURDATE() - INTERVAL 4 MONTH AND CURDATE() + INTERVAL 2 MONTH;
