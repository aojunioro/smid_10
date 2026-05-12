package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/aojunioro/smid_10/backend/internal/config"
)

// Alias identifica cada um dos quatro bancos legados reutilizados pelo SMID 10.
type Alias string

const (
	AliasSmid          Alias = "smid"
	AliasPermission    Alias = "permission"
	AliasLog           Alias = "log"
	AliasCommunication Alias = "communication"
)

// Aliases lista, em ordem canônica, todos os aliases gerenciados pelo Pools.
var Aliases = []Alias{AliasSmid, AliasPermission, AliasLog, AliasCommunication}

// Pools agrupa as conexões ativas para cada alias. Cada `*sql.DB` é seguro
// para uso concorrente e mantém seu próprio pool interno.
type Pools struct {
	conns map[Alias]*sql.DB
}

// Open abre os quatro pools com os parâmetros do `config.DBConfig`. Em caso de
// erro em qualquer abertura, todas as conexões parciais são fechadas antes do
// retorno para evitar vazamento.
func Open(cfg config.DBConfig) (*Pools, error) {
	dsns := map[Alias]string{
		AliasSmid:          cfg.SmidDSN,
		AliasPermission:    cfg.PermissionDSN,
		AliasLog:           cfg.LogDSN,
		AliasCommunication: cfg.CommunicationDSN,
	}

	p := &Pools{conns: make(map[Alias]*sql.DB, len(dsns))}

	for _, alias := range Aliases {
		dsn := dsns[alias]
		conn, err := sql.Open("mysql", dsn)
		if err != nil {
			p.Close()
			return nil, fmt.Errorf("db: open %s: %w", alias, err)
		}
		conn.SetMaxOpenConns(cfg.MaxOpenConns)
		conn.SetMaxIdleConns(cfg.MaxIdleConns)
		conn.SetConnMaxLifetime(cfg.ConnMaxLifetime)
		p.conns[alias] = conn
	}

	return p, nil
}

// Get retorna o handle de conexão associado ao alias. Retorna erro se o alias
// não foi registrado.
func (p *Pools) Get(alias Alias) (*sql.DB, error) {
	conn, ok := p.conns[alias]
	if !ok {
		return nil, fmt.Errorf("db: alias desconhecido %q", alias)
	}
	return conn, nil
}

// Ping executa um ping no alias informado respeitando o timeout do contexto.
func (p *Pools) Ping(ctx context.Context, alias Alias) error {
	conn, err := p.Get(alias)
	if err != nil {
		return err
	}
	return conn.PingContext(ctx)
}

// PingAll executa um ping em todos os pools e retorna um mapa alias→erro.
// O contexto é compartilhado, então um timeout global se aplica ao conjunto.
func (p *Pools) PingAll(ctx context.Context) map[Alias]error {
	result := make(map[Alias]error, len(p.conns))
	for _, alias := range Aliases {
		c, ok := p.conns[alias]
		if !ok {
			result[alias] = fmt.Errorf("db: alias não inicializado %q", alias)
			continue
		}
		pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		result[alias] = c.PingContext(pingCtx)
		cancel()
	}
	return result
}

// Close encerra todos os pools, ignorando erros individuais (eles são raros e
// não bloqueiam o shutdown).
func (p *Pools) Close() {
	for _, c := range p.conns {
		_ = c.Close()
	}
	p.conns = map[Alias]*sql.DB{}
}
