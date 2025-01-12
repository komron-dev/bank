package db

import (
	"context"
	"database/sql"
	"github.com/komron-dev/bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)

	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomOwner()
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: oldUser.HashedPassword,
			Valid:  true,
		},
		Email: sql.NullString{
			String: oldUser.Email,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, updatedUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, updatedUser.Email, oldUser.Email)
	require.Equal(t, newFullName, updatedUser.FullName)

	require.NotEqual(t, updatedUser.FullName, oldUser.FullName)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := util.RandomEmail()
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: oldUser.HashedPassword,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: oldUser.FullName,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, updatedUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, updatedUser.FullName, oldUser.FullName)
	require.Equal(t, newEmail, updatedUser.Email)

	require.NotEqual(t, updatedUser.Email, oldUser.Email)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword := util.RandomString(7)
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newPassword,
			Valid:  true,
		},
		Email: sql.NullString{
			String: oldUser.Email,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: oldUser.FullName,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, updatedUser.Email, oldUser.Email)
	require.Equal(t, updatedUser.FullName, oldUser.FullName)
	require.Equal(t, newPassword, updatedUser.HashedPassword)

	require.NotEqual(t, updatedUser.HashedPassword, oldUser.HashedPassword)
}

func TestUpdateUserFull(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword := util.RandomString(7)
	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newPassword,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, newPassword, updatedUser.HashedPassword)

	require.NotEqual(t, updatedUser.FullName, oldUser.FullName)
	require.NotEqual(t, updatedUser.Email, oldUser.Email)
	require.NotEqual(t, updatedUser.HashedPassword, oldUser.HashedPassword)
}
