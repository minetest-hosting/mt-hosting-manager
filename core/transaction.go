package core

import (
	"fmt"
	"mt-hosting-manager/api/coinbase"
	"mt-hosting-manager/api/wallee"
	"mt-hosting-manager/notify"
	"mt-hosting-manager/types"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (c *Core) CreateTransaction(userid string, create_tx_req *types.CreateTransactionRequest) (*types.PaymentTransaction, error) {

	user, err := c.repos.UserRepo.GetByID(userid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found: '%s'", userid)
	}

	if user.Balance+create_tx_req.Amount > int64(c.cfg.MaxBalance) {
		return nil, fmt.Errorf("max balance of %d exceeded", c.cfg.MaxBalance)
	}
	if create_tx_req.Amount < 500 {
		return nil, fmt.Errorf("min payment: EUR 5")
	}

	payment_tx_id := uuid.NewString()
	back_url := fmt.Sprintf("%s/#/finance/detail/%s", c.cfg.BaseURL, payment_tx_id)

	payment_tx := &types.PaymentTransaction{
		ID:             payment_tx_id,
		Type:           create_tx_req.Type,
		Created:        time.Now().Unix(),
		Expires:        time.Now().Add(time.Hour).Unix(),
		UserID:         userid,
		Amount:         create_tx_req.Amount,
		AmountRefunded: 0,
		State:          types.PaymentStatePending,
	}

	switch create_tx_req.Type {
	case types.PaymentTypeZahlsch:
		if !c.cfg.ZahlschEnabled {
			return nil, fmt.Errorf("zahlsch provider not enabled")
		}

		// use our own id
		payment_tx.TransactionID = payment_tx.ID
		payment_tx.PaymentURL = fmt.Sprintf("https://%s.zahls.ch/en/pay?invoice_amount=%.2f&custom_user_id=%s&custom_transaction_id=%s&tid=%s",
			c.cfg.ZahlschUser, float64(create_tx_req.Amount)/100, user.ID, payment_tx.ID, c.cfg.ZahlschPageID)

	case types.PaymentTypeWallee:
		if !c.cfg.WalleeEnabled {
			return nil, fmt.Errorf("wallee provider not enabled")
		}

		item := &wallee.LineItem{
			Name:               "Minetest hosting credits",
			Quantity:           1,
			AmountIncludingTax: float64(create_tx_req.Amount) / 100,
			Type:               wallee.LineItemTypeProduct,
			UniqueID:           payment_tx_id,
		}

		tx, err := c.wc.CreateTransaction(&wallee.TransactionRequest{
			Currency:   "EUR",
			LineItems:  []*wallee.LineItem{item},
			SuccessURL: back_url,
			FailedURL:  back_url,
		})
		if err != nil {
			return nil, fmt.Errorf("create transaction failed: %v", err)
		}

		url, err := c.wc.CreatePaymentPageURL(tx.ID)
		if err != nil {
			return nil, fmt.Errorf("create payment url failed: %v", err)
		}

		payment_tx.TransactionID = fmt.Sprintf("%d", tx.ID)
		payment_tx.PaymentURL = url

	case types.PaymentTypeCoinbase:
		if !c.cfg.CoinbaseEnabled {
			return nil, fmt.Errorf("coinbase provider not enabled")
		}

		charge, err := c.cbc.CreateCharge(&coinbase.CreateChargeRequest{
			Name:        "Minetest hosting",
			Description: "Minetest hosting payment",
			PricingType: coinbase.PricingTypeFixed,
			LocalPrice: &coinbase.LocalPrice{
				Amount:   fmt.Sprintf("%.2f", float64(create_tx_req.Amount)/100),
				Currency: coinbase.CURRENCY_EUR,
			},
			RedirectURL: back_url,
			CancelURL:   back_url,
		})
		if err != nil {
			return nil, err
		}

		payment_tx.TransactionID = charge.Data.Code
		payment_tx.PaymentURL = charge.Data.HostedURL

		return payment_tx, nil
	default:
		return nil, fmt.Errorf("payment type not implemented: %s", create_tx_req.Type)
	}

	err = c.repos.PaymentTransactionRepo.Insert(payment_tx)
	if err != nil {
		return nil, fmt.Errorf("payment tx insert failed: %v", err)
	}

	c.AddAuditLog(&types.AuditLog{
		Type:                 types.AuditLogPaymentCreated,
		UserID:               userid,
		PaymentTransactionID: &payment_tx_id,
		Amount:               &create_tx_req.Amount,
	})

	notify.Send(&notify.NtfyNotification{
		Title:    fmt.Sprintf("Transaction created by %s (%.2f)", user.Name, float64(create_tx_req.Amount)/100),
		Message:  fmt.Sprintf("User: %s\nEUR %.2f\nType: %s", user.Name, float64(create_tx_req.Amount)/100, payment_tx.Type),
		Priority: 3,
		Click:    &back_url,
		Tags:     []string{"credit_card", "new"},
	}, true)

	return payment_tx, nil

}

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

	switch tx.Type {
	case types.PaymentTypeWallee:
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
			Title:    fmt.Sprintf("Payment refunded by %s (%.2f)", user.Name, float64(refund_amount)/100),
			Message:  fmt.Sprintf("User: %s, EUR %.2f", user.Name, float64(refund_amount)/100),
			Priority: 3,
			Tags:     []string{"coin", "recycle"},
		}, true)

		return tx, nil
	default:
		return nil, fmt.Errorf("refund for payment type %s not implemented", tx.Type)
	}
}

func (c *Core) CheckTransaction(id string) (*types.PaymentTransaction, error) {
	tx, err := c.repos.PaymentTransactionRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("fetch payment tx failed: %v", err)
	}
	if tx == nil {
		return nil, fmt.Errorf("payment tx not found: %s", id)
	}

	user, err := c.repos.UserRepo.GetByID(tx.UserID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch user '%s': %v", tx.UserID, err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found: '%s'", tx.UserID)
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
					Title:    fmt.Sprintf("Payment received by %s (%.2f)", user.Name, float64(tx.Amount)/100),
					Message:  fmt.Sprintf("User: %s, EUR %.2f", user.Name, float64(tx.Amount)/100),
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
				if payment.Status == coinbase.PaymentStatusConfirmed {

					crypto_value := payment.Value[coinbase.PaymentValueCrypto]
					if crypto_value == nil {
						return nil, fmt.Errorf("crypto value not found")
					}

					local_value := payment.Value[coinbase.PaymentValueLocal]
					if local_value == nil {
						return nil, fmt.Errorf("local value not found")
					}
					if local_value.Currency != "EUR" {
						return nil, fmt.Errorf("invalid local currency: %s", local_value.Currency)
					}
					eur_value, err := strconv.ParseFloat(local_value.Amount, 64)
					if err != nil {
						return nil, fmt.Errorf("could not convert local amount '%s': %v", local_value.Amount, err)
					}
					if eur_value < 0 {
						return nil, fmt.Errorf("negative amount: %.2f", eur_value)
					}
					eurocent_value := int64(eur_value * 100) // TODO: proper currency calc or possible rounding error fix of 1 eurocent

					tx.Amount = eurocent_value
					tx.State = types.PaymentStateSuccess
					err = c.repos.PaymentTransactionRepo.Update(tx)
					if err != nil {
						return nil, fmt.Errorf("could not update tx: %v", err)
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
						Title: fmt.Sprintf("Crypto payment (%s %s, local: %d eurocent) received from %s",
							crypto_value.Amount, crypto_value.Currency, eurocent_value, user.Name),
						Message:  fmt.Sprintf("User: %s, EUR %.2f", user.Name, float64(tx.Amount)/100),
						Priority: 3,
						Tags:     []string{"coin"},
					}, true)

					break
				}
			}
		default:
			return nil, fmt.Errorf("check for payment type %s not implemented", tx.Type)
		}
	}

	// check expiration _after_ confirmation check
	if time.Now().After(time.Unix(tx.Expires, 0)) && tx.State != types.PaymentStateSuccess {
		// payment expired
		tx.State = types.PaymentStateExpired
		err = c.repos.PaymentTransactionRepo.Update(tx)
		return tx, err
	}

	return tx, nil
}
