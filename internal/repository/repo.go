package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Vladroon22/TestTask-ITK-Academy/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(p *pgxpool.Pool) *Repo {
	return &Repo{
		pool: p,
	}
}

func (r *Repo) GetBalance(c context.Context, uuid string) (entity.WalletData, error) {
	ctx, cancel := context.WithTimeout(c, time.Second*3)
	defer cancel()

	tx, errTx := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if errTx != nil {
		log.Println("Beg Tx (GetBalance):", errTx)
		return entity.WalletData{}, errors.New("bad response from database")
	}

	defer func() {
		errRb := tx.Rollback(ctx)
		if errRb != nil && !errors.Is(errRb, pgx.ErrTxClosed) {
			log.Println("Rollback Tx (GetBalance): ", errRb)
		}
	}()

	wallet := entity.WalletData{}

	query1 := "SELECT wallet_id, balance, last_operation_type, created_at FROM wallets WHERE wallet_id = $1"

	if err := tx.QueryRow(ctx, query1, uuid).Scan(&wallet.Uuid, &wallet.Balance, &wallet.Operation_type, &wallet.Created_At); err != nil {
		log.Println("Database error (QR):", err)
		return entity.WalletData{}, errors.New("bad response from database")
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("failed to commit tx (GetBalance):", err)
		return entity.WalletData{}, errors.New("bad response from database")
	}

	log.Println("successfull operation")
	return wallet, nil
}

func (r *Repo) WalletOperation(c context.Context, wallet entity.WalletData) error {
	ctx, cancel := context.WithTimeout(c, time.Second*3)
	defer cancel()

	tx, errTx := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if errTx != nil {
		log.Println("Beg Tx (WalletOperation):", errTx)
		return errors.New("bad response from database")
	}

	defer func() {
		errRb := tx.Rollback(ctx)
		if errRb != nil && !errors.Is(errRb, pgx.ErrTxClosed) {
			log.Println("Rollback Tx (WalletOperation): ", errRb)
		}
	}()

	switch wallet.Operation_type {
	case "withdraw":
		balance := 0.00
		query1 := "SELECT balance FROM wallets WHERE wallet_id = $1"
		if err := tx.QueryRow(ctx, query1, wallet.Uuid).Scan(&balance); err != nil {
			log.Println("Database error:", err)
			return errors.New("bad response from database")
		}

		if balance == 0.00 {
			log.Println("nothing to withdraw")
			return errors.New("nothing to withdraw")
		}

		currTime := time.Now().UTC()
		query2 := " UPDATE wallets SET balance = balance - $1, last_operation_type = $2, created_at = $3 WHERE wallet_id = $4"
		if _, err := tx.Exec(ctx, query2, wallet.Balance, wallet.Operation_type, currTime, wallet.Uuid); err != nil {
			log.Println("Database error:", err)
			return errors.New("bad response from database")
		}
	case "deposit":
		currTime := time.Now().UTC()
		query1 := "UPDATE wallets SET balance = balance + $1, last_operation_type = $2, created_at = $3 WHERE wallet_id = $4"
		if _, err := tx.Exec(ctx, query1, wallet.Balance, wallet.Operation_type, currTime, wallet.Uuid); err != nil {
			log.Println("Database error:", err)
			return errors.New("bad response from database")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("failed to commit tx (WalletOperation):", err)
		return errors.New("bad response from database")
	}

	log.Println("successfull operation: " + wallet.Operation_type)
	return nil
}
