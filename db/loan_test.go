package db

import (
	"context"
	"testing"

	"github.com/jvinaya/goapp/utils"
	"github.com/stretchr/testify/require"
)

func createLoanParams() CreateLoanParams {
	amount := utils.RandomString(5, true)
	return CreateLoanParams{
		Amount:          amount,
		AmountNeedToPay: amount,
		Term:            int32(utils.RandomInt(1, 5)),
		ApprovalStatus:  EnumApprovalStatusPending,
		RepaymentStatus: EnumPaymentStatusUnpaid,
		CreatedBy:       utils.RandomString(6, false),
		LastUpdatedBy:   utils.RandomString(6, false),
		IpFrom:          "localhost vscode " + utils.RandomString(6, false),
		UserAgent:       utils.RandomString(4, false),
	}
}
func createRandomLoan(t *testing.T) Loan {

	arg := createLoanParams()
	loan, err := testQueries.CreateLoan(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, loan)
	require.Equal(t, arg.CreatedBy, loan.CreatedBy)
	require.Equal(t, arg.RepaymentStatus, loan.RepaymentStatus)
	require.NotZero(t, loan.CreatedAt)
	require.NotZero(t, loan.ID)

	return loan
}
func TestCreateLoan(t *testing.T) {
	createRandomLoan(t)

}
func TestGetLoan(t *testing.T) {
	res := createRandomLoan(t)
	res2, err := testQueries.GetLoan(context.Background(), res.ID)
	require.NoError(t, err)
	require.Equal(t, res, res2)
}
func TestUpdateLoan(t *testing.T) {

	req := createRandomLoan(t)
	updatedby := utils.RandomOwner()
	updateAmount := utils.RandomString(8, true)
	arg := UpdateLoanParams{

		ID:              req.ID,
		Amount:          updateAmount,
		AmountNeedToPay: updateAmount,
		LastUpdatedBy:   updatedby,
		Term:            req.Term,
		ApprovalStatus:  EnumApprovalStatusApproved,
		RepaymentStatus: EnumPaymentStatusPaid,
	}

	updatedRes, err := testQueries.UpdateLoan(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updatedRes)
	require.Equal(t, req.ID, updatedRes.ID)
	require.Equal(t, updatedby, updatedRes.LastUpdatedBy)
	require.Equal(t, updateAmount, updatedRes.Amount)
	require.Equal(t, EnumApprovalStatusApproved, updatedRes.ApprovalStatus)
	require.Equal(t, EnumPaymentStatusPaid, updatedRes.RepaymentStatus)

}
func TestDeleteLoan(t *testing.T) {
	req := createRandomLoan(t)

	err := testQueries.DeleteLoan(context.Background(), req.ID)

	require.NoError(t, err)
	res2, err := testQueries.GetLoan(context.Background(), req.ID)
	require.NoError(t, err)
	require.NotEmpty(t, res2)
	require.Equal(t, false, res2.IsActive)

}
func TestListLoan(t *testing.T) {
	listLoan := [20]Loan{}
	for i := 0; i < 20; i++ {
		loan := createRandomLoan(t)
		listLoan[i] = loan
	}
	arg := ListLoanParams{
		Limit:  int32(utils.RandomInt(1, 5)),
		Offset: int32(utils.RandomInt(1, 5)),
	}
	fetchedListRes, err := testQueries.ListLoan(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedListRes)
	require.Equal(t, arg.Limit, int32(len(fetchedListRes)))

}
func TestListDescLoan(t *testing.T) {
	listLoan := [20]Loan{}
	for i := 0; i < 20; i++ {
		loan := createRandomLoan(t)
		listLoan[i] = loan
	}
	arg := ListDescLoanParams{
		Limit:  int32(utils.RandomInt(1, 5)),
		Offset: int32(utils.RandomInt(1, 5)),
	}
	fetchedLoanResponse, err := testQueries.ListDescLoan(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.Limit, int32(len(fetchedLoanResponse)))
	j := 20 - arg.Offset - arg.Limit
	for i := arg.Limit - 1; i > 0; i-- {

		require.Equal(t, listLoan[j], fetchedLoanResponse[i])
		j++
	}

}
