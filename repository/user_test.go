package repository_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/devbook-api/model"
	"github.com/waliqueiroz/devbook-api/repository"
	"github.com/waliqueiroz/devbook-api/test/mock"
)

func TestCreateUser(t *testing.T) {
	userJson, _ := ioutil.ReadFile("../test/resource/json/created_user.json")

	var user model.User
	json.Unmarshal(userJson, &user)

	db, mock := mock.NewDatabaseConnection()
	defer db.Close()

	subTests := []struct {
		name           string
		errorInPrepare bool
		errorInExec    bool
		errorInResult  bool
		errorInScanRow bool
		err            error
	}{
		{
			name: "Create user",
		},
		{
			name:           "Create user - error in prepare",
			errorInPrepare: true,
			err:            errors.New("some error"),
		},
		{
			name:        "Create user - error in exec",
			errorInExec: true,
			err:         errors.New("some error"),
		},
		{
			name:          "Create user - error in result",
			errorInResult: true,
			err:           errors.New("some error"),
		},
		{
			name:           "Create user - error in scan row",
			errorInScanRow: true,
			err:            errors.New("some error"),
		},
	}

	repository := repository.NewUserRepository(db)

	insertQuery := "insert into users \\(name, nick, email, password\\) values \\(\\?, \\?, \\?, \\?\\)"
	selectQuery := "select id, name, nick, email, created_at from users where id = \\?"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInPrepare {
				mock.ExpectPrepare(insertQuery).WillReturnError(subTest.err)

				_, err := repository.Create(user)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInExec {
				prep := mock.ExpectPrepare(insertQuery)
				prep.ExpectExec().WithArgs(user.Name, user.Nick, user.Email, user.Password).WillReturnError(subTest.err)

				_, err := repository.Create(user)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInResult {
				prep := mock.ExpectPrepare(insertQuery)
				prep.ExpectExec().WithArgs(user.Name, user.Nick, user.Email, user.Password).WillReturnResult(sqlmock.NewErrorResult(subTest.err))

				_, err := repository.Create(user)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInScanRow {
				prep := mock.ExpectPrepare(insertQuery)
				prep.ExpectExec().WithArgs(user.Name, user.Nick, user.Email, user.Password).WillReturnResult(sqlmock.NewResult(1, 1))

				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(-1)

				mock.ExpectQuery(selectQuery).WithArgs(user.ID).WillReturnRows(rows)

				_, err := repository.Create(user)
				assert.Error(t, err)
			} else {
				prep := mock.ExpectPrepare(insertQuery)
				prep.ExpectExec().WithArgs(user.Name, user.Nick, user.Email, user.Password).WillReturnResult(sqlmock.NewResult(1, 1))

				rows := sqlmock.NewRows([]string{"id", "name", "nick", "email", "created_at"}).
					AddRow(user.ID, user.Name, user.Nick, user.Email, user.CreatedAt)

				mock.ExpectQuery(selectQuery).WithArgs(user.ID).WillReturnRows(rows)

				createduser, err := repository.Create(user)
				fmt.Println(err)
				assert.Equal(t, user, createduser)
			}
		})
	}
}

func TestFindUserByID(t *testing.T) {
	userJson, _ := ioutil.ReadFile("../test/resource/json/created_user.json")

	var user model.User
	json.Unmarshal(userJson, &user)

	db, mock := mock.NewDatabaseConnection()
	defer db.Close()

	subTests := []struct {
		name           string
		errorInExec    bool
		errorInScanRow bool
		err            error
	}{
		{
			name:        "Find by ID",
			errorInExec: false,
		},
		{
			name:        "Find by ID - error in exec query",
			errorInExec: true,
			err:         errors.New("some error"),
		},
		{
			name:           "Find by ID - error in scan row",
			errorInScanRow: true,
			err:            errors.New("some error"),
		},
	}

	repository := repository.NewUserRepository(db)

	query := "select id, name, nick, email, created_at from users where id = \\?"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInExec {
				mock.ExpectQuery(query).WillReturnError(subTest.err)

				_, err := repository.FindByID(user.ID)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInScanRow {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(-1)

				mock.ExpectQuery(query).WithArgs(user.ID).WillReturnRows(rows)

				_, err := repository.FindByID(user.ID)
				assert.Error(t, err)
			} else {
				rows := sqlmock.NewRows([]string{"id", "name", "nick", "email", "created_at"}).
					AddRow(user.ID, user.Name, user.Nick, user.Email, user.CreatedAt)

				mock.ExpectQuery(query).WithArgs(user.ID).WillReturnRows(rows)

				createdUser, _ := repository.FindByID(user.ID)
				assert.Equal(t, user, createdUser)
			}
		})
	}
}

func TestFindUserByNameOrNick(t *testing.T) {
	userJson, _ := ioutil.ReadFile("../test/resource/json/created_user.json")

	var user model.User
	json.Unmarshal(userJson, &user)

	nameOrNick := fmt.Sprintf("%%%s%%", user.Name)

	db, mock := mock.NewDatabaseConnection()
	defer db.Close()

	subTests := []struct {
		name           string
		errorInExec    bool
		errorInScanRow bool
		err            error
	}{
		{
			name:        "Index",
			errorInExec: false,
		},
		{
			name:        "Index - error in exec query",
			errorInExec: true,
			err:         errors.New("some error"),
		},
		{
			name:           "Index - error in scan row",
			errorInScanRow: true,
			err:            errors.New("some error"),
		},
	}

	repository := repository.NewUserRepository(db)

	query := "select id, name, nick, email, created_at from users where name like \\? or nick like \\?"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInExec {
				mock.ExpectQuery(query).WillReturnError(subTest.err)

				_, err := repository.FindByNameOrNick(user.Name)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInScanRow {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(-1)

				mock.ExpectQuery(query).WithArgs(nameOrNick, nameOrNick).WillReturnRows(rows)

				_, err := repository.FindByNameOrNick(user.Name)
				assert.Error(t, err)
			} else {
				rows := sqlmock.NewRows([]string{"id", "name", "nick", "email", "created_at"}).
					AddRow(user.ID, user.Name, user.Nick, user.Email, user.CreatedAt)

				mock.ExpectQuery(query).WithArgs(nameOrNick, nameOrNick).WillReturnRows(rows)

				createdUsers, _ := repository.FindByNameOrNick(user.Name)
				assert.Equal(t, []model.User{user}, createdUsers)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	userJson, _ := ioutil.ReadFile("../test/resource/json/created_user.json")

	var user model.User
	json.Unmarshal(userJson, &user)

	db, mock := mock.NewDatabaseConnection()
	defer db.Close()

	subTests := []struct {
		name           string
		errorInPrepare bool
		errorInExec    bool
		err            error
	}{
		{
			name: "Update user",
		},
		{
			name:           "Update user - error in prepare",
			errorInPrepare: true,
			err:            errors.New("some error"),
		},
		{
			name:        "Update user - error in exec",
			errorInExec: true,
			err:         errors.New("some error"),
		},
	}

	repository := repository.NewUserRepository(db)

	query := "update users set name = \\?, nick = \\?, email = \\? where id = \\?"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInPrepare {
				mock.ExpectPrepare(query).WillReturnError(subTest.err)

				err := repository.Update(user.ID, user)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInExec {
				prep := mock.ExpectPrepare(query)
				prep.ExpectExec().WithArgs(user.Name, user.Nick, user.Email, user.ID).WillReturnError(subTest.err)

				err := repository.Update(user.ID, user)
				assert.ErrorIs(t, err, subTest.err)
			} else {
				prep := mock.ExpectPrepare(query)
				prep.ExpectExec().WithArgs(user.Name, user.Nick, user.Email, user.ID).WillReturnResult(sqlmock.NewResult(1, 1))

				err := repository.Update(user.ID, user)
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	userJson, _ := ioutil.ReadFile("../test/resource/json/created_user.json")

	var user model.User
	json.Unmarshal(userJson, &user)

	db, mock := mock.NewDatabaseConnection()
	defer db.Close()

	subTests := []struct {
		name           string
		errorInPrepare bool
		errorInExec    bool
		err            error
	}{
		{
			name: "Delete user",
		},
		{
			name:           "Delete user - error in prepare",
			errorInPrepare: true,
			err:            errors.New("some error"),
		},
		{
			name:        "Delete user - error in exec",
			errorInExec: true,
			err:         errors.New("some error"),
		},
	}

	repository := repository.NewUserRepository(db)

	query := "delete from users where id = \\?"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInPrepare {
				mock.ExpectPrepare(query).WillReturnError(subTest.err)

				err := repository.Delete(user.ID)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInExec {
				prep := mock.ExpectPrepare(query)
				prep.ExpectExec().WithArgs(user.ID).WillReturnError(subTest.err)

				err := repository.Delete(user.ID)
				assert.ErrorIs(t, err, subTest.err)
			} else {
				prep := mock.ExpectPrepare(query)
				prep.ExpectExec().WithArgs(user.ID).WillReturnResult(sqlmock.NewResult(1, 1))

				err := repository.Delete(user.ID)
				assert.NoError(t, err)
			}
		})
	}
}

func TestFindUserByEmail(t *testing.T) {
	userJson, _ := ioutil.ReadFile("../test/resource/json/created_user.json")

	var user model.User
	json.Unmarshal(userJson, &user)

	user.Password = "$2a$10$finFsyhIR/7UK/8nKmlUu.kdN.Vw3AaHBHBMZlp1HiP3J2JpMgkI6"

	db, mock := mock.NewDatabaseConnection()
	defer db.Close()

	subTests := []struct {
		name           string
		errorInExec    bool
		errorInScanRow bool
		err            error
	}{
		{
			name:        "Find by Email",
			errorInExec: false,
		},
		{
			name:        "Find by Email - error in exec query",
			errorInExec: true,
			err:         errors.New("some error"),
		},
		{
			name:           "Find by Email - error in scan row",
			errorInScanRow: true,
			err:            errors.New("some error"),
		},
	}

	repository := repository.NewUserRepository(db)

	query := "select id, password from users where email = \\?"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInExec {
				mock.ExpectQuery(query).WillReturnError(subTest.err)

				_, err := repository.FindByEmail(user.Email)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInScanRow {
				rows := sqlmock.NewRows([]string{"id", "password"}).
					AddRow(-1, user.Password)

				mock.ExpectQuery(query).WithArgs(user.Email).WillReturnRows(rows)

				_, err := repository.FindByEmail(user.Email)
				fmt.Println("teste", err)
				assert.Error(t, err)
			} else {
				rows := sqlmock.NewRows([]string{"id", "password"}).
					AddRow(user.ID, user.Password)

				mock.ExpectQuery(query).WithArgs(user.Email).WillReturnRows(rows)

				createdUser, _ := repository.FindByEmail(user.Email)
				assert.Equal(t, user.ID, createdUser.ID)
				assert.Equal(t, user.Password, createdUser.Password)
			}
		})
	}
}

func TestFollowUser(t *testing.T) {
	userJson, _ := ioutil.ReadFile("../test/resource/json/created_user.json")

	var user model.User
	json.Unmarshal(userJson, &user)

	db, mock := mock.NewDatabaseConnection()
	defer db.Close()

	subTests := []struct {
		name           string
		errorInPrepare bool
		errorInExec    bool
		err            error
	}{
		{
			name: "Follow user",
		},
		{
			name:           "Follow user - error in prepare",
			errorInPrepare: true,
			err:            errors.New("some error"),
		},
		{
			name:        "Follow user - error in exec",
			errorInExec: true,
			err:         errors.New("some error"),
		},
	}

	repository := repository.NewUserRepository(db)

	query := "insert ignore into followers \\(user_id, follower_id\\) values \\(\\?, \\?\\)"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInPrepare {
				mock.ExpectPrepare(query).WillReturnError(subTest.err)

				err := repository.Follow(2, user.ID)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInExec {
				prep := mock.ExpectPrepare(query)
				prep.ExpectExec().WithArgs(2, user.ID).WillReturnError(subTest.err)

				err := repository.Follow(2, user.ID)
				assert.ErrorIs(t, err, subTest.err)
			} else {
				prep := mock.ExpectPrepare(query)
				prep.ExpectExec().WithArgs(2, user.ID).WillReturnResult(sqlmock.NewResult(1, 1))

				err := repository.Follow(2, user.ID)
				assert.NoError(t, err)
			}
		})
	}
}

func TestUnfollowUser(t *testing.T) {
	userJson, _ := ioutil.ReadFile("../test/resource/json/created_user.json")

	var user model.User
	json.Unmarshal(userJson, &user)

	db, mock := mock.NewDatabaseConnection()
	defer db.Close()

	subTests := []struct {
		name           string
		errorInPrepare bool
		errorInExec    bool
		err            error
	}{
		{
			name: "Unfollow user",
		},
		{
			name:           "Unfollow user - error in prepare",
			errorInPrepare: true,
			err:            errors.New("some error"),
		},
		{
			name:        "Unfollow user - error in exec",
			errorInExec: true,
			err:         errors.New("some error"),
		},
	}

	repository := repository.NewUserRepository(db)

	query := "delete from followers where user_id = \\? and follower_id = \\?"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInPrepare {
				mock.ExpectPrepare(query).WillReturnError(subTest.err)

				err := repository.Unfollow(2, user.ID)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInExec {
				prep := mock.ExpectPrepare(query)
				prep.ExpectExec().WithArgs(2, user.ID).WillReturnError(subTest.err)

				err := repository.Unfollow(2, user.ID)
				assert.ErrorIs(t, err, subTest.err)
			} else {
				prep := mock.ExpectPrepare(query)
				prep.ExpectExec().WithArgs(2, user.ID).WillReturnResult(sqlmock.NewResult(1, 1))

				err := repository.Unfollow(2, user.ID)
				assert.NoError(t, err)
			}
		})
	}
}

func TestSearchFollowers(t *testing.T) {
	userJson, _ := ioutil.ReadFile("../test/resource/json/created_user.json")

	var user model.User
	json.Unmarshal(userJson, &user)

	db, mock := mock.NewDatabaseConnection()
	defer db.Close()

	subTests := []struct {
		name           string
		errorInExec    bool
		errorInScanRow bool
		err            error
	}{
		{
			name:        "Search followers",
			errorInExec: false,
		},
		{
			name:        "Search followers - error in exec query",
			errorInExec: true,
			err:         errors.New("some error"),
		},
		{
			name:           "Search followers - error in scan row",
			errorInScanRow: true,
			err:            errors.New("some error"),
		},
	}

	repository := repository.NewUserRepository(db)

	query := "select u.id, u.name, u.nick, u.email, u.created_at from users u join followers f on u.id = f.follower_id where f.user_id = \\?"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInExec {
				mock.ExpectQuery(query).WillReturnError(subTest.err)

				_, err := repository.SearchFollowers(user.ID)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInScanRow {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(-1)

				mock.ExpectQuery(query).WithArgs(user.ID).WillReturnRows(rows)

				_, err := repository.SearchFollowers(user.ID)
				assert.Error(t, err)
			} else {
				rows := sqlmock.NewRows([]string{"id", "name", "nick", "email", "created_at"}).
					AddRow(user.ID, user.Name, user.Nick, user.Email, user.CreatedAt)

				mock.ExpectQuery(query).WithArgs(user.ID).WillReturnRows(rows)

				createdUsers, _ := repository.SearchFollowers(user.ID)
				assert.Equal(t, []model.User{user}, createdUsers)
			}
		})
	}
}

func TestSearchFollowing(t *testing.T) {
	userJson, _ := ioutil.ReadFile("../test/resource/json/created_user.json")

	var user model.User
	json.Unmarshal(userJson, &user)

	db, mock := mock.NewDatabaseConnection()
	defer db.Close()

	subTests := []struct {
		name           string
		errorInExec    bool
		errorInScanRow bool
		err            error
	}{
		{
			name:        "Search following",
			errorInExec: false,
		},
		{
			name:        "Search following - error in exec query",
			errorInExec: true,
			err:         errors.New("some error"),
		},
		{
			name:           "Search following - error in scan row",
			errorInScanRow: true,
			err:            errors.New("some error"),
		},
	}

	repository := repository.NewUserRepository(db)

	query := "select u.id, u.name, u.nick, u.email, u.created_at from users u join followers f on u.id = f.user_id where f.follower_id = \\?"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInExec {
				mock.ExpectQuery(query).WillReturnError(subTest.err)

				_, err := repository.SearchFollowing(user.ID)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInScanRow {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(-1)

				mock.ExpectQuery(query).WithArgs(user.ID).WillReturnRows(rows)

				_, err := repository.SearchFollowing(user.ID)
				assert.Error(t, err)
			} else {
				rows := sqlmock.NewRows([]string{"id", "name", "nick", "email", "created_at"}).
					AddRow(user.ID, user.Name, user.Nick, user.Email, user.CreatedAt)

				mock.ExpectQuery(query).WithArgs(user.ID).WillReturnRows(rows)

				createdUsers, _ := repository.SearchFollowing(user.ID)
				assert.Equal(t, []model.User{user}, createdUsers)
			}
		})
	}
}

func TestFindPassword(t *testing.T) {
	userJson, _ := ioutil.ReadFile("../test/resource/json/created_user.json")

	var user model.User
	json.Unmarshal(userJson, &user)

	user.Password = "$2a$10$finFsyhIR/7UK/8nKmlUu.kdN.Vw3AaHBHBMZlp1HiP3J2JpMgkI6"

	db, mock := mock.NewDatabaseConnection()
	defer db.Close()

	subTests := []struct {
		name           string
		errorInExec    bool
		errorInScanRow bool
		err            error
	}{
		{
			name:        "Find password",
			errorInExec: false,
		},
		{
			name:        "Find password - error in exec query",
			errorInExec: true,
			err:         errors.New("some error"),
		},
		{
			name:           "Find password - error in scan row",
			errorInScanRow: true,
			err:            errors.New("some error"),
		},
	}

	repository := repository.NewUserRepository(db)

	query := "select password from users where id = \\?"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInExec {
				mock.ExpectQuery(query).WillReturnError(subTest.err)

				_, err := repository.FindPassword(user.ID)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInScanRow {
				rows := sqlmock.NewRows([]string{"password", "name"}).
					AddRow(-1, -1)

				mock.ExpectQuery(query).WithArgs(user.ID).WillReturnRows(rows)

				_, err := repository.FindPassword(user.ID)
				assert.Error(t, err)
			} else {
				rows := sqlmock.NewRows([]string{"password"}).
					AddRow(user.Password)

				mock.ExpectQuery(query).WithArgs(user.ID).WillReturnRows(rows)

				password, _ := repository.FindPassword(user.ID)
				assert.Equal(t, user.Password, password)
			}
		})
	}
}

func TestUpdatePassword(t *testing.T) {
	userJson, _ := ioutil.ReadFile("../test/resource/json/created_user.json")

	var user model.User
	json.Unmarshal(userJson, &user)

	password := "$2a$10$finFsyhIR/7UK/8nKmlUu.kdN.Vw3AaHBHBMZlp1HiP3J2JpMgkI6"

	db, mock := mock.NewDatabaseConnection()
	defer db.Close()

	subTests := []struct {
		name           string
		errorInPrepare bool
		errorInExec    bool
		err            error
	}{
		{
			name: "Update password",
		},
		{
			name:           "Update password - error in prepare",
			errorInPrepare: true,
			err:            errors.New("some error"),
		},
		{
			name:        "Update password - error in exec",
			errorInExec: true,
			err:         errors.New("some error"),
		},
	}

	repository := repository.NewUserRepository(db)

	query := "update users set password = \\? where id = \\?"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInPrepare {
				mock.ExpectPrepare(query).WillReturnError(subTest.err)

				err := repository.UpdatePassword(user.ID, password)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInExec {
				prep := mock.ExpectPrepare(query)
				prep.ExpectExec().WithArgs(password, user.ID).WillReturnError(subTest.err)

				err := repository.UpdatePassword(user.ID, password)
				assert.ErrorIs(t, err, subTest.err)
			} else {
				prep := mock.ExpectPrepare(query)
				prep.ExpectExec().WithArgs(password, user.ID).WillReturnResult(sqlmock.NewResult(1, 1))

				err := repository.UpdatePassword(user.ID, password)
				assert.NoError(t, err)
			}
		})
	}
}
