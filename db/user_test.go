package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jvinaya/goapp/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUserParams() CreateUserParams {
	return CreateUserParams{
		Mobile:        sql.NullString{String: utils.RandomMobile(), Valid: true},
		Email:         utils.RandomEmail(),
		CreatedBy:     utils.RandomString(6, false),
		LastUpdatedBy: utils.RandomString(6, false),
		IpFrom:        "localhost vscode " + utils.RandomString(6, false),
		Name:          utils.RandomOwner(),
		UserAgent:     utils.RandomString(4, false),
	}
}
func createRandomUser(t *testing.T) User {
	hashedPassword, err := utils.HashPassword(utils.RandomString(6, false))
	require.NoError(t, err)
	arg := createRandomUserParams()
	arg.HashedPassword = hashedPassword
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Mobile, user.Mobile)
	require.Equal(t, arg.CreatedBy, user.CreatedBy)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.ID)
	require.True(t, user.PasswordChangedAt.IsZero())
	return user
}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)

}
func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user.ID)
	require.NoError(t, err)
	require.Equal(t, user, user2)
}
func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)
	arg := UpdateUserParams{
		Name:   utils.RandomOwner(),
		Mobile: sql.NullString{String: utils.RandomMobile(), Valid: true},
		Email:  utils.RandomEmail(),
		ID:     user.ID,
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.Equal(t, user.ID, user2.ID)
	require.NotEmpty(t, user2)
	require.Equal(t, arg.Mobile, user2.Mobile)
	require.Equal(t, arg.Email, user2.Email)

}
func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), user.ID)

	require.NoError(t, err)
	user2, err := testQueries.GetUser(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user2.IsActive, false)

}

// func TestBulkInsert(t *testing.T)  {
// 	listUser:=[]User{};
// 	for i := 0; i < 5000; i++ {
// 		listUser = append(listUser, createRandomUser(t))
// 	}
// 	require.Equal(t,len(listUser),5000)
// }
func TestListUser(t *testing.T) {
	listUser := [20]User{}
	for i := 0; i < 20; i++ {
		user := createRandomUser(t)
		listUser[i] = user
	}
	arg := ListUserParams{
		Limit:  int32(utils.RandomInt(1, 5)),
		Offset: 5,
	}
	users, err := testQueries.ListUser(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, int32(len(users)), arg.Limit)

}
func TestDescUsers(t *testing.T) {
	listUser := [20]User{}
	for i := 0; i < 20; i++ {
		user := createRandomUser(t)
		listUser[i] = user
	}
	arg := ListDescUserParams{
		Limit:  int32(utils.RandomInt(1, 5)),
		Offset: int32(utils.RandomInt(1, 5)),
	}
	users, err := testQueries.ListDescUser(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, int32(len(users)), arg.Limit)
	j := 20 - arg.Offset - arg.Limit
	for i := arg.Limit - 1; i > 0; i-- {

		require.Equal(t, listUser[j], users[i])
		j++
	}

}
