package db

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

const LaonIsPaid = "loan is Paid no need to pay again"

type TransactionDetail struct {
	LoanDetails    Loan
	CurrentPayment Payment
}

func (store *Store) CreateLoanWithBorrower(ctx context.Context, arg CreateLoanParams) (Loan, error) {
	var result Loan
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		user, err := q.GetUserByEmail(ctx, arg.CreatedBy)
		if err != nil {
			return err
		}
		result, err = q.CreateLoan(ctx, arg)
		if err != nil {
			return err
		}
		argBorrower := CreateBorrowerParams{
			UserID:        user.ID,
			LoanID:        result.ID,
			CreatedBy:     arg.CreatedBy,
			LastUpdatedBy: arg.LastUpdatedBy,
			IpFrom:        arg.IpFrom,
			UserAgent:     arg.UserAgent,
		}
		q.CreateBorrower(ctx, argBorrower)
		if err != nil {
			return err
		}

		return err
	})

	return result, err

}

func (store *Store) CreatePaymentTerms(ctx context.Context, arg CreatePaymentParams) (TransactionDetail, error) {
	var result TransactionDetail

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		loan, err := q.GetLoan(ctx, arg.LoanID)
		if err != nil {
			return err
		}

		if loan.RepaymentStatus == EnumPaymentStatusPaid {

			return fmt.Errorf(LaonIsPaid)

		} else {

			pendingAmount, _ := strconv.ParseFloat(loan.AmountNeedToPay, 64)
			user, err := q.GetUserByEmail(ctx, arg.CreatedBy)
			if err != nil {
				return err
			}
			paymentAmount, err := strconv.ParseFloat(arg.Amount, 64)
			if err != nil {
				return err
			}
			loanAmount, _ := strconv.ParseFloat(loan.Amount, 64)
			SingleLoanTermAmount := loanAmount / float64(loan.Term)
			arg.UserID = user.ID

			//if this the last full and final payment of loan
			if paymentAmount == pendingAmount {

				result.CurrentPayment, err = q.CreatePayment(ctx, arg)
				if err != nil {
					return err
				}
				loan.AmountNeedToPay = "0"
				loan.RepaymentStatus = EnumPaymentStatusPaid
				updateArg := UpdateLoanParams{
					ID:              loan.ID,
					Amount:          loan.Amount,
					AmountNeedToPay: loan.AmountNeedToPay,
					Term:            loan.Term,
					ApprovalStatus:  loan.ApprovalStatus,
					RepaymentStatus: loan.RepaymentStatus,
					LastUpdatedBy:   arg.LastUpdatedBy,
					UpdatedAt:       time.Now(),
				}
				result.LoanDetails, err = q.UpdateLoan(ctx, updateArg)

				if err != nil {
					return err
				}

			} else if paymentAmount >= SingleLoanTermAmount {

				result.CurrentPayment, err = q.CreatePayment(ctx, arg)
				if err != nil {
					return err
				}
				pend := pendingAmount - paymentAmount
				loan.AmountNeedToPay = strconv.FormatFloat(pend, 'E', -1, 64)
				updateArg := UpdateLoanParams{
					ID:              loan.ID,
					Amount:          loan.Amount,
					AmountNeedToPay: loan.AmountNeedToPay,
					Term:            loan.Term,
					ApprovalStatus:  loan.ApprovalStatus,
					RepaymentStatus: loan.RepaymentStatus,
					LastUpdatedBy:   arg.LastUpdatedBy,
					UpdatedAt:       time.Now(),
				}
				result.LoanDetails, err = q.UpdateLoan(ctx, updateArg)

				return err

			} else {
				return fmt.Errorf("plese enter amount greater than or equal to :%v ", SingleLoanTermAmount)
			}
		}

		return err
	})

	return result, err

}
