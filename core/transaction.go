package core

import (
	"fmt"
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/db"
	"mt-hosting-manager/types"

	"github.com/bojanz/currency"
)

func CheckTransaction(repos *db.Repositories, wc *wallee.WalleeClient, id string) (*types.PaymentTransaction, error) {
	tx, err := repos.PaymentTransactionRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("fetch payment tx failed: %v", err)
	}
	if tx == nil {
		return nil, fmt.Errorf("payment tx not found: %s", id)
	}

	if tx.State == types.PaymentStatePending {
		// verify tx success
		txr := &wallee.TransactionSearchRequest{
			Filter: &wallee.TransactionSearchFilter{
				FieldName: "id",
				Operator:  wallee.FilterOperatorEquals,
				Type:      wallee.FilterTypeLeaf,
				Value:     tx.TransactionID,
			},
		}
		tx_list, err := wc.SearchTransaction(txr)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch transaction %s: %v", tx.ID, err)
		}
		if tx_list == nil || len(tx_list) != 1 {
			return nil, fmt.Errorf("transaction not found %s", tx.ID)
		}
		verfifed_tx := tx_list[0]
		if verfifed_tx.State == wallee.TransactionStateFulfilled {
			tx.State = types.PaymentStateSuccess
			err = repos.PaymentTransactionRepo.Update(tx)
			if err != nil {
				return nil, fmt.Errorf("failed to save transaction: %v", err)
			}

			user, err := repos.UserRepo.GetByID(tx.UserID)
			if err != nil {
				return nil, fmt.Errorf("could not fetch user '%s': %v", tx.UserID, err)
			}
			if user == nil {
				return nil, fmt.Errorf("user not found: '%s'", tx.UserID)
			}

			a, err := currency.NewAmount(user.Balance, types.DEFAULT_CURRENCY)
			if err != nil {
				return nil, fmt.Errorf("could not parse balance '%s': %v", user.Balance, err)
			}

			tx_amount, err := currency.NewAmount(tx.Amount, types.DEFAULT_CURRENCY)
			if err != nil {
				return nil, fmt.Errorf("could not parse tx amount '%s': %v", tx.Amount, err)
			}

			new_amount, err := a.Add(tx_amount)
			if err != nil {
				return nil, fmt.Errorf("could not add amounts: %v", err)
			}

			user.Balance = new_amount.String()
			err = repos.UserRepo.Update(user)
			if err != nil {
				return nil, fmt.Errorf("could not update user: %v", err)
			}
		}
	}

	return tx, nil
}
