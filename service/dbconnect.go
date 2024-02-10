package service

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type DBConn struct {
	Conn *pgxpool.Pool
	// Conn *pgx.Conn
	Ctx context.Context
}

func (dbi *DBConfig) String() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", dbi.Host, dbi.Port, dbi.DBName, dbi.User, dbi.Password)
}

func CreateDbConnection() *DBConn {
	dbInfo := DBConfig{
		Host:     "localhost",
		Port:     54321,
		User:     "postgres",
		Password: "password",
		DBName:   "pgvectest",
	}

	ctx := context.Background()
	// conn, err := pgx.Connect(ctx, dbInfo.String())
	conn, err := pgxpool.New(ctx, dbInfo.String())
	if err != nil {
		log.Fatal(err)
	}

	return &DBConn{Conn: conn, Ctx: ctx}
}

func (db *DBConn) InitPgVectest() {
	db.Conn.Exec(db.Ctx, "CREATE EXTENSION IF NOT EXISTS vector")
}

func (db *DBConn) CreateAllTables() {
	db.InitPgVectest()
	db.CreateTextDataTable()
	db.CreateVectorDataTable()
}
