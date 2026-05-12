SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;


CREATE TABLE `acompanhar` (
  `id` int(11) NOT NULL,
  `lead_id` int(11) NOT NULL,
  `observacoes` varchar(1000) NOT NULL,
  `login` varchar(100) NOT NULL,
  `criado_em` datetime NOT NULL DEFAULT current_timestamp(),
  `alterado_em` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `coeficientes` (
  `id` int(11) NOT NULL,
  `login` varchar(255) NOT NULL,
  `n_parc` int(11) DEFAULT NULL,
  `coeficientes` varchar(20) NOT NULL,
  `criado_em` datetime NOT NULL DEFAULT current_timestamp(),
  `alterado_em` datetime DEFAULT NULL ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `comissoes` (
  `id` int(11) NOT NULL,
  `ped_id` int(11) NOT NULL,
  `vlr_comissao` decimal(10,0) NOT NULL,
  `total_pago` decimal(10,0) NOT NULL,
  `vlr_saldo` decimal(10,0) NOT NULL,
  `dt_prevista` datetime NOT NULL,
  `stt_comis` char(1) DEFAULT NULL,
  `criado_em` datetime NOT NULL DEFAULT current_timestamp(),
  `excluido_em` datetime DEFAULT NULL,
  `alterado_em` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `comis_ped_item` (
  `id` int(11) NOT NULL,
  `comis_id` int(11) NOT NULL,
  `vlr_pago` decimal(10,0) NOT NULL,
  `dt_pgto` datetime NOT NULL,
  `obs_pgto` varchar(800) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;
DELIMITER $$
CREATE TRIGGER `comis_saldo_AI` AFTER INSERT ON `comis_ped_item` FOR EACH ROW UPDATE `comissoes` SET `total_pago` = `total_pago` + NEW.vlr_pago, `vlr_saldo` = `vlr_comissao` - `total_pago` WHERE `id` = NEW.comis_id
$$
DELIMITER ;

CREATE TABLE `compras` (
  `id` int(11) NOT NULL,
  `ped_id` int(11) NOT NULL,
  `fornec_id` int(11) NOT NULL,
  `status_id` int(11) NOT NULL DEFAULT 1,
  `transp_id` int(11) DEFAULT NULL,
  `login` varchar(255) NOT NULL,
  `dt_compr` datetime NOT NULL,
  `dt_coleta` datetime DEFAULT NULL,
  `dt_chegada` datetime DEFAULT NULL,
  `frete` decimal(10,0) DEFAULT NULL,
  `vlr_compr` decimal(10,0) DEFAULT NULL,
  `n_nf` int(11) DEFAULT NULL,
  `n_parcelas` int(11) DEFAULT NULL,
  `dt_pgto` date DEFAULT NULL,
  `criado_em` datetime NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `compr_fornec` (
  `id` int(11) NOT NULL,
  `pessoa_id` int(11) NOT NULL,
  `fornecedor` varchar(100) NOT NULL,
  `contato` varchar(100) DEFAULT NULL,
  `fone` varchar(11) DEFAULT NULL,
  `email` varchar(200) DEFAULT NULL,
  `ativo` char(1) DEFAULT 'S'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `compr_status` (
  `id` int(11) NOT NULL,
  `stt_compr` varchar(100) NOT NULL,
  `ordem` tinyint(2) NOT NULL,
  `cor` char(7) NOT NULL,
  `ativo` char(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `compr_transport` (
  `id` int(11) NOT NULL,
  `pessoa_id` int(11) NOT NULL,
  `transportadora` varchar(100) NOT NULL,
  `contato` varchar(100) DEFAULT NULL,
  `fone` varchar(11) DEFAULT NULL,
  `email` varchar(200) DEFAULT NULL,
  `ativa` char(1) DEFAULT '1'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `contas` (
  `id` int(11) NOT NULL,
  `tipo_conta_id` int(11) NOT NULL,
  `categ_id` int(11) NOT NULL,
  `login` varchar(255) NOT NULL,
  `cpgto_id` int(11) DEFAULT NULL,
  `fornec_id` int(11) DEFAULT NULL,
  `ped_id` int(11) DEFAULT NULL,
  `comis_id` int(11) DEFAULT NULL,
  `anexo_id` int(11) NOT NULL,
  `dt_venc` date NOT NULL,
  `dt_emis` date DEFAULT NULL,
  `dt_pgto` date DEFAULT NULL,
  `valor` decimal(10,0) DEFAULT NULL,
  `parcela` int(11) DEFAULT NULL,
  `obs` varchar(500) DEFAULT NULL,
  `criado_em` datetime NOT NULL DEFAULT current_timestamp(),
  `alterado_em` datetime DEFAULT NULL,
  `excluido_em` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `conta_anexo` (
  `id` int(11) NOT NULL,
  `conta_id` int(11) NOT NULL,
  `tipo_anexo_id` int(11) NOT NULL,
  `descricao` varchar(255) DEFAULT NULL,
  `arquivo` mediumblob DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `conta_anexo_tipo` (
  `id` int(11) NOT NULL,
  `anexo` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `conta_categ` (
  `id` int(11) NOT NULL,
  `categoria` varchar(100) NOT NULL,
  `tipo_conta_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `conta_tipo` (
  `id` int(11) NOT NULL,
  `tipo_conta` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `ddds` (
  `id` int(11) NOT NULL,
  `ddd` char(2) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `endereco` (
  `id` int(11) NOT NULL,
  `lead_id` int(11) NOT NULL,
  `CEP` varchar(8) DEFAULT NULL,
  `rua` varchar(400) NOT NULL,
  `numero` char(128) NOT NULL,
  `complemento` varchar(255) DEFAULT NULL,
  `bairro` varchar(255) NOT NULL,
  `cidade` varchar(255) NOT NULL,
  `uf` char(2) NOT NULL,
  `referencias` varchar(500) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `equipes` (
  `id` int(11) NOT NULL,
  `equipe` varchar(255) NOT NULL,
  `supervisao` varchar(100) DEFAULT NULL,
  `cor` char(7) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `estoque_saldo` (
  `id` int(11) NOT NULL,
  `prod_id` int(11) NOT NULL,
  `local_estoq_id` int(11) NOT NULL,
  `qtde` int(11) NOT NULL,
  `vlr_unit` decimal(10,2) NOT NULL DEFAULT 0.00
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `estoq_entrada` (
  `id` int(11) NOT NULL,
  `prod_id` int(11) NOT NULL,
  `local_estoq_id` int(11) NOT NULL,
  `qtde` int(11) NOT NULL,
  `vlr_unit` decimal(10,2) NOT NULL DEFAULT 0.00,
  `criado_em` datetime NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;
DELIMITER $$
CREATE TRIGGER `prodEntrada_AD` AFTER DELETE ON `estoq_entrada` FOR EACH ROW CALL AtualizaEstoque (old.id, old.qtde * -1, old.local_estoq_id)
$$
DELIMITER ;
DELIMITER $$
CREATE TRIGGER `prodEntrada_AI` AFTER INSERT ON `estoq_entrada` FOR EACH ROW CALL AtualizaEstoque (new.prod_id, new.qtde, new.local_estoq_id)
$$
DELIMITER ;
DELIMITER $$
CREATE TRIGGER `prodEntrada_AU` AFTER UPDATE ON `estoq_entrada` FOR EACH ROW CALL AtualizaEstoque (new.prod_id, new.qtde - old.qtde, new.local_estoq_id)
$$
DELIMITER ;

CREATE TABLE `estoq_local` (
  `id` int(11) NOT NULL,
  `local` varchar(128) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `estoq_saida` (
  `id` int(11) NOT NULL,
  `prod_id` int(11) NOT NULL,
  `local_estoq_id` int(11) NOT NULL,
  `qtde` int(11) NOT NULL,
  `vlr_unit` decimal(10,2) NOT NULL DEFAULT 0.00,
  `criado_em` datetime NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;
DELIMITER $$
CREATE TRIGGER `prodSaida _AD` AFTER DELETE ON `estoq_saida` FOR EACH ROW CALL AtualizaEstoque (old.prod_id, old.qtde, old.local_estoq_id)
$$
DELIMITER ;
DELIMITER $$
CREATE TRIGGER `prodSaida_AI` AFTER INSERT ON `estoq_saida` FOR EACH ROW CALL AtualizaEstoque (new.prod_id, new.qtde * -1, new.local_estoq_id)
$$
DELIMITER ;
DELIMITER $$
CREATE TRIGGER `prodSaida_AU` AFTER UPDATE ON `estoq_saida` FOR EACH ROW CALL AtualizaEstoque (new.prod_id, old.qtde - new.qtde, new.local_estoq_id)
$$
DELIMITER ;

CREATE TABLE `generos` (
  `id` int(11) NOT NULL,
  `genero` char(128) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `historicos` (
  `id` int(11) NOT NULL,
  `lead_id` int(11) DEFAULT NULL,
  `vis_id` int(11) NOT NULL,
  `login` varchar(255) NOT NULL,
  `motivo_id` int(11) DEFAULT NULL,
  `ocorr_id` int(11) DEFAULT NULL,
  `hist` varchar(1000) NOT NULL,
  `criado_em` datetime NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `hist_motivo` (
  `id` int(11) NOT NULL,
  `motivo` varchar(100) NOT NULL,
  `cor` char(7) DEFAULT NULL,
  `ordem` tinyint(1) DEFAULT NULL,
  `positivo` char(1) DEFAULT 'N'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `hist_ocorrido` (
  `id` int(11) NOT NULL,
  `ocorrido` varchar(100) NOT NULL,
  `cor` char(7) DEFAULT NULL,
  `ordem` tinyint(1) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `img_compro` (
  `id` int(11) NOT NULL,
  `nome_img_compro` varchar(100) DEFAULT NULL,
  `img_compro` mediumblob NOT NULL,
  `ped_id` int(11) NOT NULL,
  `criado_em` datetime NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `img_ped` (
  `id` int(11) NOT NULL,
  `nome_img_ped` varchar(100) DEFAULT NULL,
  `img_ped` mediumblob NOT NULL,
  `ped_id` int(11) NOT NULL,
  `criado_em` datetime NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `img_protoc` (
  `id` int(11) NOT NULL,
  `nome_ img_protoc` varchar(100) DEFAULT NULL,
  `img_protoc` mediumblob NOT NULL,
  `ped_id` int(11) NOT NULL,
  `criado_em` datetime NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `img_suport` (
  `id` int(11) NOT NULL,
  `nome_img_suport` varchar(100) DEFAULT NULL,
  `img_suport` mediumblob NOT NULL,
  `ped_id` int(11) NOT NULL,
  `criado_em` datetime NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `leads` (
  `id` int(11) NOT NULL,
  `fone1` varchar(11) NOT NULL,
  `starttime` datetime DEFAULT NULL,
  `fone2` varchar(11) DEFAULT NULL,
  `nome` varchar(100) DEFAULT 'novo lead',
  `profissao` varchar(100) DEFAULT NULL,
  `idade` int(2) DEFAULT NULL,
  `patologia` varchar(255) DEFAULT NULL,
  `nome_acomp` varchar(100) DEFAULT NULL,
  `profis_acomp` varchar(100) DEFAULT NULL,
  `idd_acomp` int(2) DEFAULT NULL,
  `pato_acomp` varchar(255) DEFAULT NULL,
  `midia_id` int(11) DEFAULT NULL,
  `tent_id` int(11) DEFAULT 1,
  `contato_ok` varchar(1) NOT NULL DEFAULT 'N',
  `status_id` int(11) DEFAULT 1,
  `unidd_id` int(11) DEFAULT NULL,
  `meio_id` int(11) DEFAULT NULL,
  `mot_pend_id` int(11) DEFAULT NULL,
  `mot_perd_id` int(11) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `obs_curta_lead` varchar(300) DEFAULT NULL,
  `login` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `login_recep` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `login_super` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `criado_em` datetime NOT NULL DEFAULT current_timestamp(),
  `alterado_em` datetime DEFAULT NULL ON UPDATE current_timestamp(),
  `excluido_em` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `lead_duplicados` (
  `id` int(11) NOT NULL,
  `fone1` varchar(11) DEFAULT NULL,
  `fone2` varchar(11) DEFAULT NULL,
  `starttime` datetime DEFAULT NULL,
  `duplicado_em` datetime NOT NULL DEFAULT current_timestamp(),
  `leads_id` bigint(20) DEFAULT NULL,
  `login` varchar(255) DEFAULT NULL,
  `login_recep` varchar(255) DEFAULT NULL,
  `status_dupli` char(1) DEFAULT 'N',
  `meio_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `lead_event_tent` (
  `id` int(11) NOT NULL,
  `lead_id` int(11) NOT NULL,
  `login` varchar(255) NOT NULL,
  `obs_curta_lead` varchar(300) NOT NULL,
  `tent_id` int(11) DEFAULT NULL,
  `mot_pend_id` int(11) DEFAULT NULL,
  `mot_perd_id` int(11) DEFAULT NULL,
  `contato_ok` char(1) DEFAULT NULL,
  `criado_em` datetime DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `lead_meio` (
  `id` int(11) NOT NULL,
  `meio` varchar(100) NOT NULL,
  `ativo` char(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `lead_motivo_pend` (
  `id` int(11) NOT NULL,
  `motivo_pend` varchar(100) NOT NULL,
  `cor` char(7) DEFAULT NULL,
  `ativa` char(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `lead_motivo_perd` (
  `id` int(11) NOT NULL,
  `motivo_perd` varchar(100) NOT NULL,
  `cor` char(7) DEFAULT NULL,
  `ativa` char(1) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `lead_status` (
  `id` int(11) NOT NULL,
  `stt_lead` varchar(100) NOT NULL,
  `cor` char(7) DEFAULT NULL,
  `ordem` int(11) DEFAULT NULL,
  `kanban` char(1) DEFAULT NULL,
  `stt_inicial` char(1) DEFAULT NULL,
  `stt_final` char(1) DEFAULT NULL,
  `perm_edit` char(1) DEFAULT NULL,
  `perm_del` char(1) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `mensagens_hash` (
  `id` int(11) NOT NULL,
  `hash` varchar(255) NOT NULL,
  `criado_em` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `midias` (
  `id` int(11) NOT NULL,
  `unidd_id` varchar(255) DEFAULT NULL,
  `midia` varchar(100) NOT NULL,
  `tipo_id` int(11) DEFAULT NULL,
  `hora_ini` time DEFAULT NULL,
  `hora_fim` time DEFAULT NULL,
  `cor` varchar(7) DEFAULT NULL,
  `ativa` char(1) DEFAULT 'S',
  `criado_em` timestamp NOT NULL DEFAULT current_timestamp(),
  `alterado_em` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE current_timestamp(),
  `excluido_em` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `midia_tipo` (
  `id` int(11) NOT NULL,
  `tipo` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `notif_inbox` (
  `inbox_id` int(11) NOT NULL,
  `notif_id` int(11) NOT NULL,
  `login` varchar(255) NOT NULL,
  `notif_dtsent` datetime NOT NULL DEFAULT current_timestamp(),
  `notif_ontop` int(11) NOT NULL DEFAULT 0,
  `notif_isread` int(11) NOT NULL DEFAULT 0,
  `notif_dtread` datetime DEFAULT NULL,
  `notif_tags` varchar(255) DEFAULT NULL,
  `notif_important` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `notif_notifications` (
  `notif_id` int(11) NOT NULL,
  `notif_title` varchar(255) NOT NULL,
  `notif_message` text NOT NULL,
  `notif_dtcreated` datetime NOT NULL DEFAULT current_timestamp(),
  `notif_ontop` int(11) NOT NULL DEFAULT 0,
  `notif_dtexpire` datetime DEFAULT NULL,
  `notif_categ` varchar(60) DEFAULT NULL,
  `notif_login_sender` varchar(255) NOT NULL,
  `notif_type` varchar(60) DEFAULT NULL,
  `notif_link` varchar(255) DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `notif_pref` (
  `login` varchar(255) NOT NULL,
  `receive_email` int(11) NOT NULL DEFAULT 0,
  `receive_sms` int(11) NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `notif_profiles` (
  `profile_id` int(11) NOT NULL,
  `profile_name` varchar(255) DEFAULT NULL,
  `profile_users` text DEFAULT NULL,
  `profile_groups` text DEFAULT NULL,
  `profile_public` int(11) NOT NULL DEFAULT 0,
  `profile_owner` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `notif_tags` (
  `tag_id` int(11) NOT NULL,
  `tag_title` varchar(50) NOT NULL,
  `login` varchar(255) NOT NULL,
  `tag_color` varchar(100) DEFAULT NULL,
  `tag_active` int(11) NOT NULL DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `notif_user_tags` (
  `user_tags_id` int(11) NOT NULL,
  `login` varchar(255) NOT NULL,
  `login_sender` varchar(255) NOT NULL,
  `tags` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `orcamentos` (
  `id` int(11) NOT NULL,
  `status_id` int(11) NOT NULL,
  `login` varchar(255) NOT NULL,
  `leads_id` int(11) NOT NULL,
  `ped_id` int(11) DEFAULT NULL,
  `item_ped_id` int(11) DEFAULT NULL,
  `coefic_id` int(11) DEFAULT NULL,
  `total_orcam` double DEFAULT NULL,
  `obs_orcam` varchar(1000) DEFAULT NULL,
  `criado_em` datetime DEFAULT current_timestamp(),
  `alterado_em` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `orcam_status` (
  `id` int(11) NOT NULL,
  `stt_orcam` varchar(100) NOT NULL,
  `cor` char(7) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `pedidos` (
  `id` int(11) NOT NULL,
  `lead_id` int(11) DEFAULT NULL,
  `login_repre` varchar(100) NOT NULL,
  `status_id` int(11) NOT NULL,
  `fpgto_id` int(11) NOT NULL,
  `cpgto_id` int(11) DEFAULT NULL,
  `login` varchar(100) NOT NULL,
  `canal_id` int(11) NOT NULL,
  `dt_ped` date NOT NULL,
  `dt_prev` date DEFAULT NULL,
  `dt_quit` date DEFAULT NULL,
  `mes` char(2) DEFAULT NULL,
  `ano` char(4) DEFAULT NULL,
  `n_ped` varchar(11) NOT NULL,
  `total_ped` decimal(10,0) NOT NULL,
  `entrada_ped` decimal(10,0) NOT NULL,
  `obs_ped` varchar(1000) DEFAULT NULL,
  `obs_ped_ger` varchar(1000) DEFAULT NULL,
  `criado_em` datetime DEFAULT current_timestamp(),
  `excluido_em` datetime DEFAULT NULL,
  `alterado_em` datetime DEFAULT NULL ON UPDATE current_timestamp(),
  `nome` varchar(100) DEFAULT NULL,
  `cpf` varchar(11) DEFAULT NULL,
  `img_ped` tinyblob DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `ped_canal` (
  `id` int(11) NOT NULL,
  `canal_ped` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `ped_cpgto` (
  `id` int(11) NOT NULL,
  `c_pgto` varchar(100) NOT NULL,
  `n_parcelas` int(11) DEFAULT NULL,
  `inicio` int(11) DEFAULT NULL,
  `intervalo` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `ped_fpgto` (
  `id` int(11) NOT NULL,
  `f_pgto` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `ped_prod_item` (
  `id` int(11) NOT NULL,
  `ped_id` int(11) NOT NULL,
  `prod_id` int(11) NOT NULL,
  `med_id` int(11) NOT NULL,
  `especial` varchar(100) NOT NULL,
  `qtdd_item` int(11) DEFAULT NULL,
  `vlr_item` decimal(10,0) DEFAULT NULL,
  `vlr_total_item` decimal(10,0) DEFAULT NULL,
  `criado_em` datetime NOT NULL DEFAULT current_timestamp(),
  `excluido_em` datetime DEFAULT NULL,
  `alterado_em` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `ped_status` (
  `id` int(11) NOT NULL,
  `stt_ped` varchar(100) NOT NULL,
  `cor` char(7) DEFAULT NULL,
  `ordem` int(11) DEFAULT NULL,
  `kanban` char(1) DEFAULT NULL,
  `stt_inicial` char(1) DEFAULT NULL,
  `stt_final` char(1) DEFAULT NULL,
  `perm_edit` char(1) DEFAULT NULL,
  `perm_del` char(1) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `ped_status_negados` (
  `id` int(11) NOT NULL,
  `stts_negados` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `pos_perg` (
  `id` int(11) NOT NULL,
  `pergunta` varchar(500) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `pos_perg_item` (
  `id` int(11) NOT NULL,
  `pos_visita_id` int(11) NOT NULL,
  `perg_pos_id` int(11) NOT NULL,
  `resposta` varchar(300) NOT NULL,
  `criado_em` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `pos_visita` (
  `id` int(11) NOT NULL,
  `vis_id` int(11) NOT NULL,
  `login` varchar(100) NOT NULL,
  `nota_repre` tinyint(2) DEFAULT NULL,
  `nota_prod` tinyint(2) NOT NULL,
  `nota_empre` tinyint(2) NOT NULL,
  `visitado` char(1) NOT NULL,
  `pontual` char(1) DEFAULT NULL,
  `jaleco` char(1) NOT NULL,
  `adquiriu` char(1) NOT NULL,
  `obs` varchar(300) NOT NULL,
  `criado_em` datetime NOT NULL DEFAULT current_timestamp(),
  `lead_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `produtos` (
  `id` int(11) NOT NULL,
  `fornec_id` int(11) DEFAULT NULL,
  `nome_prod` varchar(100) NOT NULL,
  `vlr_prod_compra` decimal(10,0) NOT NULL,
  `vlr_prod_venda` decimal(10,0) NOT NULL,
  `estoq_min` int(11) NOT NULL,
  `estoq_max` int(11) NOT NULL,
  `ativo` char(1) DEFAULT '1',
  `criado_em` datetime NOT NULL DEFAULT current_timestamp(),
  `excluido_em` datetime DEFAULT NULL,
  `alterado_em` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `prod_categ` (
  `id` int(11) NOT NULL,
  `categoria` varchar(100) NOT NULL,
  `cor` char(7) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `prod_medidas` (
  `id` int(11) NOT NULL,
  `medida` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `prod_modelos` (
  `id` int(11) NOT NULL,
  `modelo` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `seg_apps` (
  `app_name` varchar(128) NOT NULL,
  `app_type` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `seg_groups` (
  `group_id` int(11) NOT NULL,
  `description` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `seg_groups_apps` (
  `group_id` int(11) NOT NULL,
  `app_name` varchar(128) NOT NULL,
  `priv_access` varchar(1) DEFAULT NULL,
  `priv_insert` varchar(1) DEFAULT NULL,
  `priv_delete` varchar(1) DEFAULT NULL,
  `priv_update` varchar(1) DEFAULT NULL,
  `priv_export` varchar(1) DEFAULT NULL,
  `priv_print` varchar(1) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `seg_log` (
  `id` int(8) NOT NULL,
  `inserted_date` datetime DEFAULT NULL,
  `username` varchar(90) NOT NULL,
  `application` varchar(255) NOT NULL,
  `creator` varchar(30) NOT NULL,
  `ip_user` varchar(255) NOT NULL,
  `action` varchar(30) NOT NULL,
  `description` text DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `seg_logged` (
  `id` int(11) NOT NULL,
  `login` varchar(255) NOT NULL,
  `date_login` varchar(128) DEFAULT NULL,
  `sc_session` varchar(32) DEFAULT NULL,
  `ip` varchar(32) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `seg_settings` (
  `set_name` varchar(255) NOT NULL,
  `set_value` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `seg_users` (
  `login` varchar(255) NOT NULL,
  `pswd` varchar(255) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `active` varchar(1) DEFAULT 'N',
  `activation_code` varchar(32) DEFAULT NULL,
  `priv_admin` varchar(1) DEFAULT 'N',
  `mfa` varchar(255) DEFAULT NULL,
  `picture` longblob DEFAULT NULL,
  `role` varchar(128) DEFAULT NULL,
  `phone` varchar(64) DEFAULT NULL,
  `pswd_last_updated` timestamp NULL DEFAULT NULL,
  `mfa_last_updated` timestamp NULL DEFAULT NULL,
  `unidd_id` varchar(255) NOT NULL DEFAULT '1',
  `equipe_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `seg_users_groups` (
  `login` varchar(255) NOT NULL,
  `group_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `suportes` (
  `id` int(11) NOT NULL,
  `ped_id` int(11) NOT NULL,
  `status_id` int(11) NOT NULL,
  `fone_sup` varchar(11) DEFAULT NULL,
  `solicit_id` int(11) NOT NULL,
  `depart_id` int(11) NOT NULL,
  `login` varchar(100) NOT NULL,
  `atrib_login` varchar(100) DEFAULT NULL,
  `prioridade` int(11) DEFAULT NULL,
  `dt_sup` date DEFAULT NULL,
  `dt_limit` date DEFAULT NULL,
  `relato_cli` varchar(500) DEFAULT NULL,
  `dt_resol` date DEFAULT NULL,
  `relato_tec` varchar(500) DEFAULT NULL,
  `img_ordem_id` int(11) DEFAULT NULL,
  `criado_em` datetime NOT NULL DEFAULT current_timestamp(),
  `excluido_em` datetime DEFAULT NULL,
  `alterado_em` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `suport_depart` (
  `id` int(11) NOT NULL,
  `departamento` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `suport_prioridd` (
  `id` int(11) NOT NULL,
  `prioridade` char(100) NOT NULL,
  `cor` varchar(7) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `suport_solicit` (
  `id` int(11) NOT NULL,
  `solicitacao` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `suport_status` (
  `id` int(11) NOT NULL,
  `stt_sup` varchar(100) NOT NULL,
  `cor` char(7) DEFAULT NULL,
  `ordem` int(11) DEFAULT NULL,
  `kanban` char(1) DEFAULT NULL,
  `stt_inicial` char(1) DEFAULT NULL,
  `stt_final` char(1) DEFAULT NULL,
  `perm_edit` char(1) DEFAULT NULL,
  `permi_del` char(1) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `system_access_log` (
  `id` int(11) NOT NULL,
  `sessionid` varchar(256) DEFAULT NULL,
  `login` varchar(256) DEFAULT NULL,
  `login_time` varchar(20) DEFAULT NULL,
  `login_year` varchar(4) DEFAULT NULL,
  `login_month` varchar(2) DEFAULT NULL,
  `login_day` varchar(2) DEFAULT NULL,
  `logout_time` varchar(20) DEFAULT NULL,
  `impersonated` char(1) DEFAULT NULL,
  `access_ip` varchar(45) DEFAULT NULL,
  `impersonated_by` varchar(200) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_change_log` (
  `id` int(11) NOT NULL,
  `logdate` varchar(256) DEFAULT NULL,
  `login` varchar(256) DEFAULT NULL,
  `tablename` varchar(256) DEFAULT NULL,
  `primarykey` varchar(256) DEFAULT NULL,
  `pkvalue` varchar(256) DEFAULT NULL,
  `operation` varchar(256) DEFAULT NULL,
  `columnname` varchar(256) DEFAULT NULL,
  `oldvalue` text DEFAULT NULL,
  `newvalue` text DEFAULT NULL,
  `access_ip` varchar(45) DEFAULT NULL,
  `transaction_id` varchar(256) DEFAULT NULL,
  `log_year` varchar(4) DEFAULT NULL,
  `log_month` varchar(2) DEFAULT NULL,
  `log_day` varchar(2) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_document` (
  `id` int(11) NOT NULL,
  `system_user_id` int(11) DEFAULT NULL,
  `title` varchar(256) DEFAULT NULL,
  `description` text DEFAULT NULL,
  `submission_date` date DEFAULT NULL,
  `archive_date` date DEFAULT NULL,
  `filename` varchar(512) DEFAULT NULL,
  `in_trash` char(1) DEFAULT NULL,
  `system_folder_id` int(11) DEFAULT NULL,
  `content` text DEFAULT NULL,
  `content_type` varchar(100) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_document_bookmark` (
  `id` int(11) NOT NULL,
  `system_user_id` int(11) DEFAULT NULL,
  `system_document_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_document_group` (
  `id` int(11) NOT NULL,
  `document_id` int(11) DEFAULT NULL,
  `system_group_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_document_user` (
  `id` int(11) NOT NULL,
  `document_id` int(11) DEFAULT NULL,
  `system_user_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_folder` (
  `id` int(11) NOT NULL,
  `system_user_id` int(11) DEFAULT NULL,
  `created_at` varchar(20) DEFAULT NULL,
  `name` varchar(256) NOT NULL,
  `in_trash` char(1) DEFAULT NULL,
  `system_folder_parent_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_folder_bookmark` (
  `id` int(11) NOT NULL,
  `system_user_id` int(11) DEFAULT NULL,
  `system_folder_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_folder_group` (
  `id` int(11) NOT NULL,
  `system_folder_id` int(11) DEFAULT NULL,
  `system_group_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_folder_user` (
  `id` int(11) NOT NULL,
  `system_folder_id` int(11) DEFAULT NULL,
  `system_user_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_group` (
  `id` int(11) NOT NULL,
  `name` varchar(256) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_group_program` (
  `id` int(11) NOT NULL,
  `system_group_id` int(11) DEFAULT NULL,
  `system_program_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_message` (
  `id` int(11) NOT NULL,
  `system_user_id` int(11) DEFAULT NULL,
  `system_user_to_id` int(11) DEFAULT NULL,
  `subject` varchar(256) DEFAULT NULL,
  `message` text DEFAULT NULL,
  `dt_message` varchar(20) DEFAULT NULL,
  `checked` char(1) DEFAULT NULL,
  `removed` char(1) DEFAULT NULL,
  `viewed` char(1) DEFAULT NULL,
  `attachments` text DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_message_tag` (
  `id` int(11) NOT NULL,
  `system_message_id` int(11) NOT NULL,
  `tag` varchar(256) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_notification` (
  `id` int(11) NOT NULL,
  `system_user_id` int(11) DEFAULT NULL,
  `system_user_to_id` int(11) DEFAULT NULL,
  `subject` varchar(256) DEFAULT NULL,
  `message` text DEFAULT NULL,
  `dt_message` varchar(20) DEFAULT NULL,
  `action_url` text DEFAULT NULL,
  `action_label` varchar(256) DEFAULT NULL,
  `icon` varchar(100) DEFAULT NULL,
  `checked` char(1) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_preference` (
  `id` varchar(256) NOT NULL,
  `preference` text DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_program` (
  `id` int(11) NOT NULL,
  `name` varchar(256) DEFAULT NULL,
  `controller` varchar(256) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_request_log` (
  `id` int(11) NOT NULL,
  `endpoint` text DEFAULT NULL,
  `logdate` varchar(256) DEFAULT NULL,
  `log_year` varchar(4) DEFAULT NULL,
  `log_month` varchar(2) DEFAULT NULL,
  `log_day` varchar(2) DEFAULT NULL,
  `session_id` varchar(256) DEFAULT NULL,
  `login` varchar(256) DEFAULT NULL,
  `access_ip` varchar(256) DEFAULT NULL,
  `class_name` varchar(256) DEFAULT NULL,
  `class_method` varchar(256) DEFAULT NULL,
  `http_host` varchar(256) DEFAULT NULL,
  `server_port` varchar(256) DEFAULT NULL,
  `request_uri` text DEFAULT NULL,
  `request_method` varchar(256) DEFAULT NULL,
  `query_string` text DEFAULT NULL,
  `request_headers` text DEFAULT NULL,
  `request_body` text DEFAULT NULL,
  `request_duration` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_role` (
  `id` int(11) NOT NULL,
  `name` varchar(256) DEFAULT NULL,
  `custom_code` varchar(256) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_sql_log` (
  `id` int(11) NOT NULL,
  `logdate` varchar(256) DEFAULT NULL,
  `login` varchar(256) DEFAULT NULL,
  `database_name` varchar(256) DEFAULT NULL,
  `sql_command` text DEFAULT NULL,
  `statement_type` varchar(256) DEFAULT NULL,
  `access_ip` varchar(45) DEFAULT NULL,
  `transaction_id` varchar(256) DEFAULT NULL,
  `log_year` varchar(4) DEFAULT NULL,
  `log_month` varchar(2) DEFAULT NULL,
  `log_day` varchar(2) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_unit` (
  `id` int(11) NOT NULL,
  `name` varchar(256) DEFAULT NULL,
  `connection_name` varchar(256) DEFAULT NULL,
  `custom_code` varchar(256) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_users` (
  `id` int(11) NOT NULL,
  `name` varchar(256) DEFAULT NULL,
  `login` varchar(256) DEFAULT NULL,
  `password` varchar(256) DEFAULT NULL,
  `email` varchar(256) DEFAULT NULL,
  `frontpage_id` int(11) DEFAULT NULL,
  `system_unit_id` int(11) DEFAULT NULL,
  `active` char(1) DEFAULT NULL,
  `accepted_term_policy` char(1) DEFAULT NULL,
  `accepted_term_policy_at` varchar(20) DEFAULT NULL,
  `accepted_term_policy_data` text DEFAULT NULL,
  `phone` varchar(256) DEFAULT NULL,
  `address` text DEFAULT NULL,
  `function_name` varchar(256) DEFAULT NULL,
  `about` text DEFAULT NULL,
  `two_factor_secret` varchar(256) DEFAULT NULL,
  `two_factor_active` char(1) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_user_group` (
  `id` int(11) NOT NULL,
  `system_user_id` int(11) DEFAULT NULL,
  `system_group_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_user_program` (
  `id` int(11) NOT NULL,
  `system_user_id` int(11) DEFAULT NULL,
  `system_program_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `system_user_unit` (
  `id` int(11) NOT NULL,
  `system_user_id` int(11) DEFAULT NULL,
  `system_unit_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `sys_setting` (
  `id` char(2) NOT NULL,
  `sys_data` longtext DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `tarefas` (
  `id` int(11) NOT NULL,
  `tarefa` varchar(300) NOT NULL,
  `dt_tarefa` date NOT NULL,
  `hr_tarefa` time NOT NULL,
  `status` char(1) NOT NULL DEFAULT 'N',
  `login` varchar(255) NOT NULL,
  `lead_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `unidades` (
  `id` int(11) NOT NULL,
  `unidade` varchar(255) NOT NULL,
  `ativa` char(1) NOT NULL,
  `cor` varchar(10) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `unidd_ddd` (
  `unidd_id` int(11) NOT NULL,
  `ddd_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `visitas` (
  `id` int(11) NOT NULL,
  `lead_id` int(11) DEFAULT NULL,
  `status_id` int(11) DEFAULT NULL,
  `login_recep` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `login_repre` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `dt_visita` date NOT NULL,
  `hr_visita` time NOT NULL,
  `confirm` char(1) DEFAULT NULL,
  `login_conf` varchar(100) DEFAULT NULL,
  `dt_confirm` datetime DEFAULT NULL,
  `interesse` tinyint(1) DEFAULT NULL,
  `hist_feito` char(1) NOT NULL DEFAULT 'N',
  `pos_feito` char(1) NOT NULL DEFAULT 'N',
  `criado_em` datetime NOT NULL DEFAULT current_timestamp(),
  `alterado_em` datetime DEFAULT NULL ON UPDATE current_timestamp(),
  `excluido_em` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

CREATE TABLE `vis_status` (
  `id` int(11) NOT NULL,
  `stt_visita` varchar(100) NOT NULL,
  `cor` char(7) NOT NULL,
  `ordem` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;
CREATE TABLE `v_admin_users` (
`login` varchar(255)
,`name` varchar(255)
,`active` varchar(1)
,`unidd_id` varchar(255)
,`priv_admin` varchar(1)
,`group_id` int(11)
);
CREATE TABLE `v_contatoOK` (
`lead_id` int(11)
,`login` varchar(255)
,`obs_curta_lead` varchar(300)
,`tent_id` int(11)
,`contato_ok` char(1)
,`criado_em` datetime
);
CREATE TABLE `v_pedidos_INvalidos` (
`lead_id` int(11)
,`nome` varchar(100)
,`cpf` varchar(11)
,`login_repre` varchar(100)
,`dt_ped` date
,`total_ped` decimal(10,0)
);
CREATE TABLE `v_pedidos_validos` (
`lead_id` int(11)
,`nome` varchar(100)
,`cpf` varchar(11)
,`login_repre` varchar(100)
,`dt_ped` date
,`total_ped` decimal(10,0)
);
CREATE TABLE `v_visitasCalendario` (
`lead` int(11)
,`status_id` int(11)
,`stt_lead` varchar(100)
,`bairro` varchar(255)
,`cidade` varchar(255)
,`dt_visita` date
,`hr_visita` time
,`cor` varchar(10)
,`id` int(11)
);
CREATE TABLE `v_visitasCalendario_nv` (
`LEAD` int(11)
,`status_id` int(11)
,`stt_lead` varchar(100)
,`bairro` varchar(255)
,`cidade` varchar(255)
,`dt_visita` date
,`hr_visita` time
,`cor` varchar(10)
,`id` int(11)
);
CREATE TABLE `v_visitas_INvalidas` (
`dt_visita` date
,`login_repre` varchar(255)
,`login_recep` varchar(100)
,`lead_id` int(11)
);
CREATE TABLE `v_visitas_VAlidas` (
`dt_visita` date
,`login_repre` varchar(255)
,`login_recep` varchar(100)
,`lead_id` int(11)
);


ALTER TABLE `acompanhar`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `coeficientes`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `comissoes`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ind_ped` (`ped_id`),
  ADD KEY `ind_dt_prev` (`dt_prevista`);

ALTER TABLE `comis_ped_item`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_comis` (`comis_id`);

ALTER TABLE `compras`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_compr_ped` (`ped_id`);

ALTER TABLE `compr_fornec`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `compr_status`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `compr_transport`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `contas`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `conta_anexo`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `conta_anexo_tipo`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `conta_categ`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `conta_tipo`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `ddds`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `endereco`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ind_lead` (`lead_id`),
  ADD KEY `ind_cidade` (`cidade`);

ALTER TABLE `equipes`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `estoque_saldo`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `estoq_entrada`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `estoq_local`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `estoq_saida`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `generos`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `historicos`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ind_vis` (`vis_id`),
  ADD KEY `ind_lead` (`lead_id`),
  ADD KEY `idx_historicos_vis_id` (`vis_id`,`id`);

ALTER TABLE `hist_motivo`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `hist_ocorrido`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `img_compro`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `img_ped`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ind_ped` (`ped_id`);

ALTER TABLE `img_protoc`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ind_ped` (`ped_id`);

ALTER TABLE `img_suport`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ind_ped` (`ped_id`);

ALTER TABLE `leads`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `unic_fone` (`fone2`,`fone1`),
  ADD KEY `ind_fone1` (`fone1`),
  ADD KEY `ind_fone2` (`fone2`),
  ADD KEY `ind_starttime` (`starttime`),
  ADD KEY `ind_midia` (`midia_id`),
  ADD KEY `ind_status` (`status_id`),
  ADD KEY `ind_unidd` (`unidd_id`),
  ADD KEY `ind_recep` (`login_recep`),
  ADD KEY `ind_criado` (`criado_em`),
  ADD KEY `ind_login_super` (`login_super`),
  ADD KEY `idx_leads_unidd_criado_em_desc` (`unidd_id`,`criado_em`),
  ADD KEY `idx_leads_midia_criado_id` (`midia_id`,`criado_em`,`id`);

ALTER TABLE `lead_duplicados`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ind_lead` (`leads_id`),
  ADD KEY `ind_recep` (`login_recep`) USING BTREE;

ALTER TABLE `lead_event_tent`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ind_lead` (`lead_id`),
  ADD KEY `ind_tent` (`tent_id`),
  ADD KEY `ind_login` (`login`),
  ADD KEY `idx_tent_lead_criado` (`lead_id`,`criado_em`,`tent_id`);

ALTER TABLE `lead_meio`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `lead_motivo_pend`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `lead_motivo_perd`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `lead_status`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `mensagens_hash`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `midias`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `midia_tipo`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `notif_inbox`
  ADD PRIMARY KEY (`inbox_id`);

ALTER TABLE `notif_notifications`
  ADD PRIMARY KEY (`notif_id`);

ALTER TABLE `notif_pref`
  ADD PRIMARY KEY (`login`);

ALTER TABLE `notif_profiles`
  ADD PRIMARY KEY (`profile_id`);

ALTER TABLE `notif_tags`
  ADD PRIMARY KEY (`tag_id`);

ALTER TABLE `notif_user_tags`
  ADD PRIMARY KEY (`user_tags_id`);

ALTER TABLE `orcamentos`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `orcam_status`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `pedidos`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ind_lead` (`lead_id`),
  ADD KEY `ind_dt_ped` (`dt_ped`),
  ADD KEY `ind_stts_ped` (`status_id`),
  ADD KEY `ind_n_ped` (`n_ped`),
  ADD KEY `ind_criado` (`criado_em`);

ALTER TABLE `ped_canal`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `ped_cpgto`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `ped_fpgto`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `ped_prod_item`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_ped_ped` (`ped_id`),
  ADD KEY `fk_ped_prod` (`prod_id`);

ALTER TABLE `ped_status`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `ped_status_negados`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `pos_perg`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `pos_perg_item`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `pos_visita`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ind_vis` (`vis_id`);

ALTER TABLE `produtos`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `prod_categ`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `prod_medidas`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `prod_modelos`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `seg_apps`
  ADD PRIMARY KEY (`app_name`);

ALTER TABLE `seg_groups`
  ADD PRIMARY KEY (`group_id`);

ALTER TABLE `seg_groups_apps`
  ADD PRIMARY KEY (`group_id`,`app_name`),
  ADD KEY `ind_appname` (`app_name`);

ALTER TABLE `seg_log`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `seg_logged`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ind_login` (`login`);

ALTER TABLE `seg_settings`
  ADD PRIMARY KEY (`set_name`);

ALTER TABLE `seg_users`
  ADD PRIMARY KEY (`login`);

ALTER TABLE `seg_users_groups`
  ADD PRIMARY KEY (`login`,`group_id`),
  ADD KEY `ind_group` (`group_id`);

ALTER TABLE `suportes`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_suport_ped` (`ped_id`);

ALTER TABLE `suport_depart`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `suport_prioridd`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `suport_solicit`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `suport_status`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `system_access_log`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `system_change_log`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `system_document`
  ADD PRIMARY KEY (`id`),
  ADD KEY `system_folder_id` (`system_folder_id`);

ALTER TABLE `system_document_bookmark`
  ADD PRIMARY KEY (`id`),
  ADD KEY `system_document_id` (`system_document_id`);

ALTER TABLE `system_document_group`
  ADD PRIMARY KEY (`id`),
  ADD KEY `document_id` (`document_id`);

ALTER TABLE `system_document_user`
  ADD PRIMARY KEY (`id`),
  ADD KEY `document_id` (`document_id`);

ALTER TABLE `system_folder`
  ADD PRIMARY KEY (`id`),
  ADD KEY `system_folder_parent_id` (`system_folder_parent_id`);

ALTER TABLE `system_folder_bookmark`
  ADD PRIMARY KEY (`id`),
  ADD KEY `system_folder_id` (`system_folder_id`);

ALTER TABLE `system_folder_group`
  ADD PRIMARY KEY (`id`),
  ADD KEY `system_folder_id` (`system_folder_id`);

ALTER TABLE `system_folder_user`
  ADD PRIMARY KEY (`id`),
  ADD KEY `system_folder_id` (`system_folder_id`);

ALTER TABLE `system_group`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `system_group_program`
  ADD PRIMARY KEY (`id`),
  ADD KEY `system_group_id` (`system_group_id`),
  ADD KEY `system_program_id` (`system_program_id`);

ALTER TABLE `system_message`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `system_message_tag`
  ADD PRIMARY KEY (`id`),
  ADD KEY `system_message_id` (`system_message_id`);

ALTER TABLE `system_notification`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `system_preference`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `system_program`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `controller` (`controller`);

ALTER TABLE `system_request_log`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `system_role`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `system_sql_log`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `system_unit`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `system_users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `login` (`login`),
  ADD KEY `frontpage_id` (`frontpage_id`),
  ADD KEY `system_unit_id` (`system_unit_id`);

ALTER TABLE `system_user_group`
  ADD PRIMARY KEY (`id`),
  ADD KEY `system_user_id` (`system_user_id`),
  ADD KEY `system_group_id` (`system_group_id`);

ALTER TABLE `system_user_program`
  ADD PRIMARY KEY (`id`),
  ADD KEY `system_user_id` (`system_user_id`),
  ADD KEY `system_program_id` (`system_program_id`);

ALTER TABLE `system_user_unit`
  ADD PRIMARY KEY (`id`),
  ADD KEY `system_user_id` (`system_user_id`),
  ADD KEY `system_unit_id` (`system_unit_id`);

ALTER TABLE `tarefas`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ind_lead` (`lead_id`),
  ADD KEY `ind_login` (`login`),
  ADD KEY `ind_stts` (`status`);

ALTER TABLE `unidades`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `unidd_ddd`
  ADD PRIMARY KEY (`unidd_id`);

ALTER TABLE `visitas`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ind_lead` (`lead_id`),
  ADD KEY `ind_repre` (`login_repre`),
  ADD KEY `ind_dt_vis` (`dt_visita`),
  ADD KEY `ind_recep` (`login_recep`),
  ADD KEY `ind_stts_vis` (`status_id`),
  ADD KEY `ind_hora` (`hr_visita`),
  ADD KEY `idx_visitas_lead_dt_hr_repre` (`lead_id`,`dt_visita`,`hr_visita`,`login_repre`),
  ADD KEY `idx_visita_ordenacao` (`dt_visita`,`hr_visita`,`login_repre`);

ALTER TABLE `vis_status`
  ADD PRIMARY KEY (`id`);


ALTER TABLE `acompanhar`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `coeficientes`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `comissoes`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `comis_ped_item`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `compras`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `compr_fornec`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `compr_status`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `compr_transport`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `contas`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `conta_anexo`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `conta_anexo_tipo`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `conta_categ`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `conta_tipo`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `ddds`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `endereco`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `equipes`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `estoque_saldo`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `estoq_entrada`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `estoq_local`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `estoq_saida`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `generos`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `historicos`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `hist_motivo`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `hist_ocorrido`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `img_compro`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `img_ped`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `img_protoc`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `img_suport`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `leads`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `lead_duplicados`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `lead_event_tent`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `lead_meio`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `lead_motivo_pend`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `lead_motivo_perd`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `lead_status`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `mensagens_hash`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `midias`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `midia_tipo`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `notif_inbox`
  MODIFY `inbox_id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `notif_notifications`
  MODIFY `notif_id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `notif_profiles`
  MODIFY `profile_id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `notif_tags`
  MODIFY `tag_id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `notif_user_tags`
  MODIFY `user_tags_id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `orcamentos`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `orcam_status`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `pedidos`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `ped_canal`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `ped_cpgto`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `ped_fpgto`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `ped_prod_item`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `ped_status`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `ped_status_negados`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `pos_perg`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `pos_perg_item`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `pos_visita`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `produtos`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `prod_categ`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `prod_medidas`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `prod_modelos`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `seg_groups`
  MODIFY `group_id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `seg_log`
  MODIFY `id` int(8) NOT NULL AUTO_INCREMENT;

ALTER TABLE `seg_logged`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `suportes`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `suport_depart`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `suport_prioridd`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `suport_solicit`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `suport_status`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_access_log`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_change_log`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_document`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_document_bookmark`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_document_group`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_document_user`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_folder`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_folder_bookmark`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_folder_group`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_folder_user`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_group`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_group_program`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_message`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_message_tag`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_notification`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_program`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_request_log`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_role`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_sql_log`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_unit`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_users`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_user_group`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_user_program`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `system_user_unit`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `tarefas`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `unidades`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `unidd_ddd`
  MODIFY `unidd_id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `visitas`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `vis_status`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
DROP TABLE IF EXISTS `v_admin_users`;

CREATE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `v_admin_users`  AS SELECT `u`.`login` AS `login`, `u`.`name` AS `name`, `u`.`active` AS `active`, `u`.`unidd_id` AS `unidd_id`, `u`.`priv_admin` AS `priv_admin`, `g`.`group_id` AS `group_id` FROM (`seg_users` `u` join `seg_users_groups` `g` on(`u`.`login` = `g`.`login`)) ;
DROP TABLE IF EXISTS `v_contatoOK`;

CREATE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `v_contatoOK`  AS SELECT `lead_event_tent`.`lead_id` AS `lead_id`, `lead_event_tent`.`login` AS `login`, `lead_event_tent`.`obs_curta_lead` AS `obs_curta_lead`, `lead_event_tent`.`tent_id` AS `tent_id`, `lead_event_tent`.`contato_ok` AS `contato_ok`, `lead_event_tent`.`criado_em` AS `criado_em` FROM `lead_event_tent` WHERE `lead_event_tent`.`contato_ok` = 'S' ;
DROP TABLE IF EXISTS `v_pedidos_INvalidos`;

CREATE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `v_pedidos_INvalidos`  AS SELECT `p`.`lead_id` AS `lead_id`, `p`.`nome` AS `nome`, `p`.`cpf` AS `cpf`, `p`.`login_repre` AS `login_repre`, `p`.`dt_ped` AS `dt_ped`, `p`.`total_ped` AS `total_ped` FROM `pedidos` AS `p` WHERE `p`.`status_id` is null OR `p`.`status_id` in (1,7,8) ;
DROP TABLE IF EXISTS `v_pedidos_validos`;

CREATE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `v_pedidos_validos`  AS SELECT `p`.`lead_id` AS `lead_id`, `p`.`nome` AS `nome`, `p`.`cpf` AS `cpf`, `p`.`login_repre` AS `login_repre`, `p`.`dt_ped` AS `dt_ped`, `p`.`total_ped` AS `total_ped` FROM `pedidos` AS `p` WHERE `p`.`status_id` is null OR `p`.`status_id` in (2,3,4,5,6,9) ;
DROP TABLE IF EXISTS `v_visitasCalendario`;

CREATE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `v_visitasCalendario`  AS SELECT `l`.`id` AS `lead`, `l`.`status_id` AS `status_id`, `s`.`stt_lead` AS `stt_lead`, `e`.`bairro` AS `bairro`, `e`.`cidade` AS `cidade`, `v`.`dt_visita` AS `dt_visita`, `v`.`hr_visita` AS `hr_visita`, `u`.`cor` AS `cor`, `u`.`id` AS `id` FROM (((`visitas` `v` left join (`leads` `l` join `endereco` `e` on(`l`.`id` = `e`.`lead_id`)) on(`e`.`lead_id` = `v`.`lead_id`)) join `lead_status` `s` on(`l`.`status_id` = `s`.`id`)) join `unidades` `u` on(`l`.`unidd_id` = `u`.`id`)) WHERE `l`.`status_id` = 2 OR `l`.`status_id` = 3 OR `l`.`status_id` = 4 ;
DROP TABLE IF EXISTS `v_visitasCalendario_nv`;

CREATE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `v_visitasCalendario_nv`  AS SELECT `l`.`id` AS `LEAD`, `l`.`status_id` AS `status_id`, `s`.`stt_lead` AS `stt_lead`, `e`.`bairro` AS `bairro`, `e`.`cidade` AS `cidade`, `v`.`dt_visita` AS `dt_visita`, `v`.`hr_visita` AS `hr_visita`, `u`.`cor` AS `cor`, `u`.`id` AS `id` FROM ((((`visitas` `v` left join `leads` `l` on(`l`.`id` = `v`.`lead_id`)) join `endereco` `e` on(`l`.`id` = `e`.`lead_id`)) join `lead_status` `s` on(`l`.`status_id` = `s`.`id`)) join `unidades` `u` on(`l`.`unidd_id` = `u`.`id`)) WHERE `l`.`status_id` in (2,3,4) AND `v`.`id` = (select `v2`.`id` from `visitas` `v2` where `v2`.`lead_id` = `v`.`lead_id` order by `v2`.`dt_visita` desc,`v2`.`hr_visita` desc limit 1) ;
DROP TABLE IF EXISTS `v_visitas_INvalidas`;

CREATE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `v_visitas_INvalidas`  AS SELECT `v`.`dt_visita` AS `dt_visita`, `v`.`login_repre` AS `login_repre`, `v`.`login_recep` AS `login_recep`, `v`.`lead_id` AS `lead_id` FROM ((`visitas` `v` left join `leads` `l` on(`v`.`lead_id` = `l`.`id`)) left join `historicos` `h` on(`v`.`id` = `h`.`vis_id`)) WHERE `h`.`motivo_id` = 2 ;
DROP TABLE IF EXISTS `v_visitas_VAlidas`;

CREATE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `v_visitas_VAlidas`  AS SELECT `v`.`dt_visita` AS `dt_visita`, `v`.`login_repre` AS `login_repre`, `v`.`login_recep` AS `login_recep`, `v`.`lead_id` AS `lead_id` FROM ((`visitas` `v` left join `leads` `l` on(`v`.`lead_id` = `l`.`id`)) left join `historicos` `h` on(`v`.`id` = `h`.`vis_id`)) WHERE `l`.`status_id` in (2,3,4,7,8) AND (`h`.`motivo_id` is null OR `h`.`motivo_id` <> 2) ORDER BY `v`.`lead_id` DESC ;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
