package core

import (
	"fmt"
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/types"
	"strconv"

	"github.com/bojanz/currency"
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

	balance, err := currency.NewAmount(user.Balance, types.DEFAULT_CURRENCY)
	if err != nil {
		return nil, fmt.Errorf("could not parse balance '%s': %v", user.Balance, err)
	}

	tx_amount, err := currency.NewAmount(tx.Amount, types.DEFAULT_CURRENCY)
	if err != nil {
		return nil, fmt.Errorf("could not parse tx.Amount '%s': %v", tx.Amount, err)
	}

	tx_refunded_amount, err := currency.NewAmount(tx.AmountRefunded, types.DEFAULT_CURRENCY)
	if err != nil {
		return nil, fmt.Errorf("could not parse tx.AmountRefunded '%s': %v", tx.Amount, err)
	}

	// remaining amount refundable on this transaction
	tx_amount_remaining, err := tx_amount.Sub(tx_refunded_amount)
	if err != nil {
		return nil, fmt.Errorf("could not subtract amount - refund: %v", err)
	}

	refund_amount := tx_amount_remaining
	cmp, err := balance.Cmp(tx_amount_remaining)
	if err != nil {
		return nil, fmt.Errorf("could not compare balance->tx_amount_remaining: %v", err)
	}
	if cmp != 1 {
		// balance less or equal than tx-refund rest
		// use entire remaining balance as refund
		refund_amount = balance
	}

	tx_id, err := strconv.ParseInt(tx.TransactionID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse tx id: %v", err)
	}

	refund_resp, err := c.wc.CreateRefund(&wallee.CreateRefundRequest{
		Amount:     refund_amount.Number(),
		ExternalID: tx.ID,
		Transaction: &wallee.CreateRefundRequestTransaction{
			ID: tx_id,
		},
		Type: wallee.RefundCustomerInitiatedManual,
	})
	if err != nil {
		return nil, fmt.Errorf("refund api error: %v", err)
	}
	if refund_resp.State != wallee.CreateRefundSuccessful {
		return nil, fmt.Errorf("refund not successful, state: %s", refund_resp.State)
	}

	tx.AmountRefunded = refund_amount.Number()
	err = c.repos.PaymentTransactionRepo.Update(tx)
	if err != nil {
		return nil, fmt.Errorf("could not update tx: %v", err)
	}

	remaining_amount, err := balance.Sub(refund_amount)
	if err != nil {
		return nil, fmt.Errorf("could not subtract balance-refund_amount amount: %v", err)
	}

	user.Balance = remaining_amount.Number()
	err = c.repos.UserRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("could not update user balance: %v", err)
	}

	currency := types.DEFAULT_CURRENCY
	c.AddAuditLog(&types.AuditLog{
		Type:                 types.AuditLogPaymentRefunded,
		UserID:               tx.UserID,
		PaymentTransactionID: &tx.ID,
		Amount:               &tx.AmountRefunded,
		Currency:             &currency,
	})

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

			user.Balance = new_amount.Number()
			err = c.repos.UserRepo.Update(user)
			if err != nil {
				return nil, fmt.Errorf("could not update user: %v", err)
			}

			currency := types.DEFAULT_CURRENCY
			c.AddAuditLog(&types.AuditLog{
				Type:                 types.AuditLogPaymentReceived,
				UserID:               tx.UserID,
				PaymentTransactionID: &tx.ID,
				Amount:               &tx.Amount,
				Currency:             &currency,
			})
		}
	}

	return tx, nil
}
