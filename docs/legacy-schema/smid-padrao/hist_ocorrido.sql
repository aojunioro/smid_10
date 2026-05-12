--
-- Estrutura para tabela `hist_ocorrido`
--

CREATE TABLE `hist_ocorrido` (
  `id` int(11) NOT NULL,
  `ocorrido` varchar(100) NOT NULL,
  `cor` char(7) DEFAULT NULL,
  `ordem` tinyint(1) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;

--
-- Despejando dados para a tabela `hist_ocorrido`
--

INSERT INTO `hist_ocorrido` (`id`, `ocorrido`, `cor`, `ordem`) VALUES
(1, 'Endereço não encontrado', '#ffb878', 1),
(2, 'Cliente esqueceu e saiu', '#ff887c', 2),
(3, 'Rejeitou a visita', '#dc2127', 3),
(4, 'Conge não esta em casa ', '', 0);

--
-- Índices para tabelas despejadas
--

--
-- Índices de tabela `hist_ocorrido`
--
ALTER TABLE `hist_ocorrido`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT para tabelas despejadas
--

--
-- AUTO_INCREMENT de tabela `hist_ocorrido`
--
ALTER TABLE `hist_ocorrido`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;