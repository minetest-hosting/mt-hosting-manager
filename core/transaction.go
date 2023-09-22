package core

import (
	"fmt"
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/types"
)

func (c *Core) RefundTransaction(id string) (*types.PaymentTransaction, error) {
	tx, err := c.repos.PaymentTransactionRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("fetch payment tx failed: %v", err)
	}
	if tx == nil {
		return nil, fmt.Errorf("payment tx not found: %s", id)
	}
	if tx.State != types.PaymentStateSuccess {
		return nil, fmt.Errorf("payment state invalid: %s", tx.State)
	}
	if tx.AmountRefunded != "0" {
		return nil, fmt.Errorf("already refunded: '%s'", tx.AmountRefunded)
	}

	user, err := c.repos.UserRepo.GetByID(tx.UserID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch user '%s': %v", tx.UserID, err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found: '%s'", tx.UserID)
	}

	//TODO

	return tx, nil
}

func (c *Core) CheckTransaction(id string) (*types.PaymentTransaction, error) {
	tx, err := c.repos.PaymentTransactionRepo.GetByID(id)
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
		tx_list, err := c.wc.SearchTransaction(txr)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch transaction %s: %v", tx.ID, err)
		}
		if tx_list == nil || len(tx_list) != 1 {
			return nil, fmt.Errorf("transaction not found %s", tx.ID)
		}
		verfifed_tx := tx_list[0]
		if verfifed_tx.State == wallee.TransactionStateFulfilled {
			tx.State = types.PaymentStateSuccess
			err = c.repos.PaymentTransactionRepo.Update(tx)
			if err != nil {
				return nil, fmt.Errorf("failed to save transaction: %v", err)
			}

			user, err := c.repos.UserRepo.GetByID(tx.UserID)
			if err != nil {
				return nil, fmt.Errorf("could not fetch user '%s': %v", tx.UserID, err)
			}
			if user == nil {
				return nil, fmt.Errorf("user not found: '%s'", tx.UserID)
			}

			//TODO: add to balance
		}
	}

	return tx, nil
}
