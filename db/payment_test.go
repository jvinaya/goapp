package db

import (
	"context"
	"strconv"
	"testing"

	"github.com/jvinaya/goapp/utils"
	"github.com/stretchr/testify/require"
)

func createRandomPayment(t *testing.T) Payment {
	user := createRandomUser(t)
	loan := createRandomLoan(t)
	loanAmount, _ := strconv.Atoi(loan.Amount)
	termAmount := strconv.Itoa(loanAmount / int(loan.Term))
	arg := CreatePaymentParams{
		UserID:        user.ID,
		LoanID:        loan.ID,
		Amount:        termAmount,
		CreatedBy:     utils.RandomString(6, false),
		LastUpdatedBy: utils.RandomString(6, false),
		IpFrom:        "localhost vscode " + utils.RandomString(6, false),
		UserAgent:     utils.RandomString(4, false),
	}
	res, err := testQueries.CreatePayment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.UserID, res.UserID)
	require.Equal(t, arg.LoanID, res.LoanID)
	require.Equal(t, arg.CreatedBy, res.CreatedBy)
	require.Equal(t, arg.UserID, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, res.ID)
	return res
}
func TestCreatePayement(t *testing.T) {
	createRandomPayment(t)

}
func TestGetPayment(t *testing.T) {
	res := createRandomPayment(t)
	res2, err := testQueries.GetPayment(context.Background(), res.ID)
	require.NoError(t, err)
	require.Equal(t, res, res2)
}
func TestUpdatePayment(t *testing.T) {

	req := createRandomPayment(t)
	updatedby := utils.RandomOwner()
	arg := UpdatePaymentParams{

		ID:            req.ID,
		LastUpdatedBy: updatedby,
		UserID:        req.UserID,
		LoanID:        req.LoanID,
		Amount:        req.Amount,
	}

	updatedRes, err := testQueries.UpdatePayment(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updatedRes)
	require.Equal(t, req.ID, updatedRes.ID)
	require.Equal(t, updatedby, updatedRes.LastUpdatedBy)

}

func TestListPayement(t *testing.T) {
	listPayment := [20]Payment{}
	for i := 0; i < 20; i++ {
		pay := createRandomPayment(t)
		listPayment[i] = pay
	}
	arg := ListPaymentParams{
		Limit:  int32(utils.RandomInt(1, 5)),
		Offset: 5,
	}
	res, err := testQueries.ListPayment(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, int32(len(res)), arg.Limit)

}
func TestListDescPayment(t *testing.T) {
	paymentList := [20]Payment{}
	for i := 0; i < 20; i++ {
		pay := createRandomPayment(t)
		paymentList[i] = pay
	}
	arg := ListDescPaymentParams{
		Limit:  int32(utils.RandomInt(1, 5)),
		Offset: int32(utils.RandomInt(1, 5)),
	}
	fetchedRes, err := testQueries.ListDescPayment(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, int32(len(fetchedRes)), arg.Limit)
	j := 20 - arg.Offset - arg.Limit
	for i := arg.Limit - 1; i > 0; i-- {

		require.Equal(t, paymentList[j], fetchedRes[i])
		j++
	}

}
