package db

import (
	"context"
	"testing"

	"github.com/jvinaya/goapp/utils"
	"github.com/stretchr/testify/require"
)

func createRandomBorrower(t *testing.T) Borrower {
	user := createRandomUser(t)
	loan := createRandomLoan(t)
	arg := CreateBorrowerParams{
		UserID:        user.ID,
		LoanID:        loan.ID,
		CreatedBy:     utils.RandomString(6, false),
		LastUpdatedBy: utils.RandomString(6, false),
		IpFrom:        "localhost vscode " + utils.RandomString(6, false),
		UserAgent:     utils.RandomString(4, false),
	}
	res, err := testQueries.CreateBorrower(context.Background(), arg)
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
func TestCreateBorrower(t *testing.T) {
	createRandomBorrower(t)

}
func TestGetBorrower(t *testing.T) {
	res := createRandomBorrower(t)
	res2, err := testQueries.GetBorrower(context.Background(), res.ID)
	require.NoError(t, err)
	require.Equal(t, res, res2)
}
func TestUpdateBorrower(t *testing.T) {

	req := createRandomBorrower(t)
	updatedby := utils.RandomOwner()
	arg := UpdateBorrowerParams{

		ID:            req.ID,
		LastUpdatedBy: updatedby,
		UserID:        req.UserID,
		LoanID:        req.LoanID,
	}

	updatedRes, err := testQueries.UpdateBorrower(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updatedRes)
	require.Equal(t, req.ID, updatedRes.ID)
	require.Equal(t, updatedby, updatedRes.LastUpdatedBy)

}
func TestDeleteBorrower(t *testing.T) {
	req := createRandomBorrower(t)

	err := testQueries.DeleteBorrower(context.Background(), req.ID)

	require.NoError(t, err)
	res2, err := testQueries.GetBorrower(context.Background(), req.ID)
	require.NoError(t, err)
	require.NotEmpty(t, res2)
	require.Equal(t, false, res2.IsActive)

}

func TestListBorrower(t *testing.T) {
	listUser := [20]Borrower{}
	for i := 0; i < 20; i++ {
		user := createRandomBorrower(t)
		listUser[i] = user
	}
	arg := ListBorrowerParams{
		Limit:  int32(utils.RandomInt(1, 5)),
		Offset: 5,
	}
	res, err := testQueries.ListBorrower(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, int32(len(res)), arg.Limit)

}
func TestListDescBorrower(t *testing.T) {
	listRequest := [20]Borrower{}
	for i := 0; i < 20; i++ {
		request := createRandomBorrower(t)
		listRequest[i] = request
	}
	arg := ListDescBorrowerParams{
		Limit:  int32(utils.RandomInt(1, 5)),
		Offset: int32(utils.RandomInt(1, 5)),
	}
	fetchedBorrowerRes, err := testQueries.ListDescBorrower(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, int32(len(fetchedBorrowerRes)), arg.Limit)
	j := 20 - arg.Offset - arg.Limit
	for i := arg.Limit - 1; i > 0; i-- {

		require.Equal(t, listRequest[j], fetchedBorrowerRes[i])
		j++
	}

}
