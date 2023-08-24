package core

import (
	"fmt"
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/db"
	"mt-hosting-manager/types"
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

			//TODO: add amount to user balance
		}
	}

	return tx, nil
}
