package core

import (
	"fmt"
	"mt-hosting-manager/api/coinbase"
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/notify"
	"mt-hosting-manager/types"
	"strconv"
)

// refunds the given transaction by all or the available amount
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
	if tx.AmountRefunded > 0 {
		return nil, fmt.Errorf("already refunded: '%d'", tx.AmountRefunded)
	}

	user, err := c.repos.UserRepo.GetByID(tx.UserID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch user '%s': %v", tx.UserID, err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found: '%s'", tx.UserID)
	}

	tx_id, err := strconv.ParseInt(tx.TransactionID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse tx id: %v", err)
	}

	refund_amount := tx.Amount
	if refund_amount > user.Balance {
		// use remaining balance
		refund_amount = user.Balance
	}

	crs, err := c.wc.CreateRefund(&wallee.CreateRefundRequest{
		Amount:     fmt.Sprintf("%.2f", float64(refund_amount)/100),
		ExternalID: tx.ID,
		Transaction: &wallee.CreateRefundRequestTransaction{
			ID: tx_id,
		},
		Type: wallee.RefundCustomerInitiatedManual,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create refund: %v", err)
	}
	if crs.State != wallee.CreateRefundSuccessful {
		return nil, fmt.Errorf("refund not successful, state: %s", crs.State)
	}

	tx.AmountRefunded = refund_amount
	err = c.repos.PaymentTransactionRepo.Update(tx)
	if err != nil {
		return nil, fmt.Errorf("tx update error: %v", err)
	}

	err = c.SubtractBalance(tx.UserID, refund_amount)
	if err != nil {
		return nil, fmt.Errorf("user balance update error: %v", err)
	}

	c.AddAuditLog(&types.AuditLog{
		Type:                 types.AuditLogPaymentRefunded,
		UserID:               tx.UserID,
		PaymentTransactionID: &tx.ID,
		Amount:               &refund_amount,
	})

	notify.Send(&notify.NtfyNotification{
		Title:    fmt.Sprintf("Payment refunded by %s (%.2f)", user.Mail, float64(refund_amount)/100),
		Message:  fmt.Sprintf("User: %s, EUR %.2f", user.Mail, float64(refund_amount)/100),
		Priority: 3,
		Tags:     []string{"coin", "recycle"},
	}, true)

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
		switch tx.Type {
		case types.PaymentTypeWallee:
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

				err = c.repos.UserRepo.AddBalance(tx.UserID, tx.Amount)
				if err != nil {
					return nil, fmt.Errorf("could not add balance '%d' to user '%s': %v", tx.Amount, tx.UserID, err)
				}

				c.AddAuditLog(&types.AuditLog{
					Type:                 types.AuditLogPaymentReceived,
					UserID:               tx.UserID,
					PaymentTransactionID: &tx.ID,
					Amount:               &tx.Amount,
				})

				notify.Send(&notify.NtfyNotification{
					Title:    fmt.Sprintf("Payment received by %s (%.2f)", user.Mail, float64(tx.Amount)/100),
					Message:  fmt.Sprintf("User: %s, EUR %.2f", user.Mail, float64(tx.Amount)/100),
					Priority: 3,
					Tags:     []string{"coin"},
				}, true)
			}

		case types.PaymentTypeCoinbase:
			charge, err := c.cbc.GetCharge(tx.TransactionID)
			if err != nil {
				return nil, err
			}

			for _, payment := range charge.Data.Payments {
				if payment.Status == coinbase.PaymentStatusConfirmend {
					//TODO
					fmt.Println("confirmed")
				}
			}
		}
	}

	return tx, nil
}
