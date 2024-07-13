package transaction

import (
	"database/sql"
	"fmt"

	"gorm.io/gorm"
)

// BeginTxWithIsolation begins a new transaction with the specified isolation level.
func BeginTxWithIsolation(db *gorm.DB, level sql.IsolationLevel) (*gorm.DB, error) {
	tx := db.Begin(&sql.TxOptions{Isolation: level})
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

// WithTransaction executes a function within a transaction if have (or open a new one with default isolation level if not), committing if successful and rolling back if not.
func WithTransaction(db *gorm.DB, tx *gorm.DB, fn func(tx *gorm.DB) error) error {
	var localTx *gorm.DB
	var err error

	if tx == nil {
		localTx = db.Begin()
		if localTx.Error != nil {
			return fmt.Errorf("failed to begin transaction: %v", localTx.Error)
		}
		defer func() {
			if r := recover(); r != nil {
				localTx.Rollback()
				panic(r) // re-throw after rollback
			} else if err != nil {
				localTx.Rollback()
			} else {
				err = localTx.Commit().Error
			}
		}()
	} else {
		localTx = tx
	}

	err = fn(localTx)
	return err
}

// WithNewTransaction always starts a new transaction, regardless of any existing transaction.
func WithNewTransaction(db *gorm.DB, level sql.IsolationLevel, fn func(tx *gorm.DB) error) error {
	tx, err := BeginTxWithIsolation(db, level)
	if err != nil {
		return fmt.Errorf("failed to begin new transaction: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // re-throw after rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	err = fn(tx)
	return err
}
