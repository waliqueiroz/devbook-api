package repository_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/waliqueiroz/devbook-api/model"
	"github.com/waliqueiroz/devbook-api/repository"
	"github.com/waliqueiroz/devbook-api/test/mock"
)

func TestCreatePost(t *testing.T) {
	postJson, _ := ioutil.ReadFile("../test/resource/json/created_post.json")

	var post model.Post
	json.Unmarshal(postJson, &post)

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
			name: "Create post",
		},
		{
			name:           "Create post - error in prepare",
			errorInPrepare: true,
			err:            errors.New("some error"),
		},
		{
			name:        "Create post - error in exec",
			errorInExec: true,
			err:         errors.New("some error"),
		},
		{
			name:          "Create post - error in result",
			errorInResult: true,
			err:           errors.New("some error"),
		},
		{
			name:           "Create post - error in scan row",
			errorInScanRow: true,
			err:            errors.New("some error"),
		},
	}

	repository := repository.NewPostRepository(db)

	insertQuery := "insert into posts \\(title, content, author_id\\) values \\(\\?, \\?, \\?\\)"
	selectQuery := "select p.*, u.nick from posts p join users u on p.author_id = u.id where p.id = \\?"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInPrepare {
				mock.ExpectPrepare(insertQuery).WillReturnError(subTest.err)

				_, err := repository.Create(post)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInExec {
				prep := mock.ExpectPrepare(insertQuery)
				prep.ExpectExec().WithArgs(post.Title, post.Content, post.AuthorID).WillReturnError(subTest.err)

				_, err := repository.Create(post)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInResult {
				prep := mock.ExpectPrepare(insertQuery)
				prep.ExpectExec().WithArgs(post.Title, post.Content, post.AuthorID).WillReturnResult(sqlmock.NewErrorResult(subTest.err))

				_, err := repository.Create(post)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInScanRow {
				prep := mock.ExpectPrepare(insertQuery)
				prep.ExpectExec().WithArgs(post.Title, post.Content, post.AuthorID).WillReturnResult(sqlmock.NewResult(1, 1))

				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(-1)

				mock.ExpectQuery(selectQuery).WithArgs(post.ID).WillReturnRows(rows)

				_, err := repository.Create(post)
				assert.Error(t, err)
			} else {
				prep := mock.ExpectPrepare(insertQuery)
				prep.ExpectExec().WithArgs(post.Title, post.Content, post.AuthorID).WillReturnResult(sqlmock.NewResult(1, 1))

				rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "likes", "created_at", "nick"}).
					AddRow(post.ID, post.Title, post.Content, post.AuthorID, post.Likes, post.CreatedAt, post.AuthorNick)

				mock.ExpectQuery(selectQuery).WithArgs(post.ID).WillReturnRows(rows)

				createdPost, _ := repository.Create(post)
				assert.Equal(t, post, createdPost)
			}
		})
	}
}

func TestFindPostByID(t *testing.T) {
	postJson, _ := ioutil.ReadFile("../test/resource/json/created_post.json")

	var post model.Post
	json.Unmarshal(postJson, &post)

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

	repository := repository.NewPostRepository(db)

	query := "select p.*, u.nick from posts p join users u on p.author_id = u.id where p.id = \\?"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInExec {
				mock.ExpectQuery(query).WillReturnError(subTest.err)

				_, err := repository.FindByID(post.ID)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInScanRow {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(-1)

				mock.ExpectQuery(query).WithArgs(post.ID).WillReturnRows(rows)

				_, err := repository.FindByID(post.ID)
				assert.Error(t, err)
			} else {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "likes", "created_at", "nick"}).
					AddRow(post.ID, post.Title, post.Content, post.AuthorID, post.Likes, post.CreatedAt, post.AuthorNick)

				mock.ExpectQuery(query).WithArgs(post.ID).WillReturnRows(rows)

				createdPost, _ := repository.FindByID(post.ID)
				assert.Equal(t, post, createdPost)
			}
		})
	}
}

func TestIndexPost(t *testing.T) {
	postJson, _ := ioutil.ReadFile("../test/resource/json/created_post.json")

	var post model.Post
	json.Unmarshal(postJson, &post)

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

	repository := repository.NewPostRepository(db)

	query := "select distinct p.*, u.nick from posts p join users u on p.author_id = u.id join followers f on p.author_id = f.user_id where u.id = \\? or f.follower_id = \\? order by p.id desc"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInExec {
				mock.ExpectQuery(query).WillReturnError(subTest.err)

				_, err := repository.Index(post.AuthorID)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInScanRow {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(-1)

				mock.ExpectQuery(query).WithArgs(post.AuthorID, post.AuthorID).WillReturnRows(rows)

				_, err := repository.Index(post.AuthorID)
				assert.Error(t, err)
			} else {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "likes", "created_at", "nick"}).
					AddRow(post.ID, post.Title, post.Content, post.AuthorID, post.Likes, post.CreatedAt, post.AuthorNick)

				mock.ExpectQuery(query).WithArgs(post.AuthorID, post.AuthorID).WillReturnRows(rows)

				createdPosts, _ := repository.Index(post.AuthorID)
				assert.Equal(t, []model.Post{post}, createdPosts)
			}
		})
	}
}

func TestUpdatePost(t *testing.T) {
	postJson, _ := ioutil.ReadFile("../test/resource/json/created_post.json")

	var post model.Post
	json.Unmarshal(postJson, &post)

	db, mock := mock.NewDatabaseConnection()
	defer db.Close()

	subTests := []struct {
		name           string
		errorInPrepare bool
		errorInExec    bool
		err            error
	}{
		{
			name: "Update post",
		},
		{
			name:           "Update post - error in prepare",
			errorInPrepare: true,
			err:            errors.New("some error"),
		},
		{
			name:        "Update post - error in exec",
			errorInExec: true,
			err:         errors.New("some error"),
		},
	}

	repository := repository.NewPostRepository(db)

	query := "update posts set title = \\?, content = \\? where id = \\?"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInPrepare {
				mock.ExpectPrepare(query).WillReturnError(subTest.err)

				err := repository.Update(post.ID, post)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInExec {
				prep := mock.ExpectPrepare(query)
				prep.ExpectExec().WithArgs(post.Title, post.Content, post.ID).WillReturnError(subTest.err)

				err := repository.Update(post.ID, post)
				assert.ErrorIs(t, err, subTest.err)
			} else {
				prep := mock.ExpectPrepare(query)
				prep.ExpectExec().WithArgs(post.Title, post.Content, post.ID).WillReturnResult(sqlmock.NewResult(1, 1))

				err := repository.Update(post.ID, post)
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeletePost(t *testing.T) {
	postJson, _ := ioutil.ReadFile("../test/resource/json/created_post.json")

	var post model.Post
	json.Unmarshal(postJson, &post)

	db, mock := mock.NewDatabaseConnection()
	defer db.Close()

	subTests := []struct {
		name           string
		errorInPrepare bool
		errorInExec    bool
		err            error
	}{
		{
			name: "Delete post",
		},
		{
			name:           "Delete post - error in prepare",
			errorInPrepare: true,
			err:            errors.New("some error"),
		},
		{
			name:        "Delete post - error in exec",
			errorInExec: true,
			err:         errors.New("some error"),
		},
	}

	repository := repository.NewPostRepository(db)

	query := "delete from posts where id = \\?"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInPrepare {
				mock.ExpectPrepare(query).WillReturnError(subTest.err)

				err := repository.Delete(post.ID)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInExec {
				prep := mock.ExpectPrepare(query)
				prep.ExpectExec().WithArgs(post.ID).WillReturnError(subTest.err)

				err := repository.Delete(post.ID)
				assert.ErrorIs(t, err, subTest.err)
			} else {
				prep := mock.ExpectPrepare(query)
				prep.ExpectExec().WithArgs(post.ID).WillReturnResult(sqlmock.NewResult(1, 1))

				err := repository.Delete(post.ID)
				assert.NoError(t, err)
			}
		})
	}
}

func TestFindPostByUser(t *testing.T) {
	postJson, _ := ioutil.ReadFile("../test/resource/json/created_post.json")

	var post model.Post
	json.Unmarshal(postJson, &post)

	db, mock := mock.NewDatabaseConnection()
	defer db.Close()

	subTests := []struct {
		name           string
		errorInExec    bool
		errorInScanRow bool
		err            error
	}{
		{
			name:        "Find by user",
			errorInExec: false,
		},
		{
			name:        "Find by user - error in exec query",
			errorInExec: true,
			err:         errors.New("some error"),
		},
		{
			name:           "Find by user - error in scan row",
			errorInScanRow: true,
			err:            errors.New("some error"),
		},
	}

	repository := repository.NewPostRepository(db)

	query := "select distinct p.*, u.nick from posts p join users u on p.author_id = u.id where	u.id = \\?"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInExec {
				mock.ExpectQuery(query).WillReturnError(subTest.err)

				_, err := repository.FindByUser(post.AuthorID)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInScanRow {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(-1)

				mock.ExpectQuery(query).WithArgs(post.AuthorID).WillReturnRows(rows)

				_, err := repository.FindByUser(post.AuthorID)
				assert.Error(t, err)
			} else {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "likes", "created_at", "nick"}).
					AddRow(post.ID, post.Title, post.Content, post.AuthorID, post.Likes, post.CreatedAt, post.AuthorNick)

				mock.ExpectQuery(query).WithArgs(post.AuthorID).WillReturnRows(rows)

				createdPosts, _ := repository.FindByUser(post.AuthorID)
				assert.Equal(t, []model.Post{post}, createdPosts)
			}
		})
	}

}

func TestLikePost(t *testing.T) {
	db, mock := mock.NewDatabaseConnection()
	defer db.Close()

	subTests := []struct {
		name           string
		errorInPrepare bool
		errorInExec    bool
		err            error
	}{
		{
			name: "Like post",
		},
		{
			name:           "Like post - error in prepare",
			errorInPrepare: true,
			err:            errors.New("some error"),
		},
		{
			name:        "Like post - error in exec",
			errorInExec: true,
			err:         errors.New("some error"),
		},
	}

	repository := repository.NewPostRepository(db)

	query := "update posts set likes = likes \\+ 1 where id = \\?"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInPrepare {
				mock.ExpectPrepare(query).WillReturnError(subTest.err)

				err := repository.LikePost(1)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInExec {
				prep := mock.ExpectPrepare(query)
				prep.ExpectExec().WithArgs(1).WillReturnError(subTest.err)

				err := repository.LikePost(1)
				assert.ErrorIs(t, err, subTest.err)
			} else {
				prep := mock.ExpectPrepare(query)
				prep.ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

				err := repository.LikePost(1)
				assert.NoError(t, err)
			}
		})
	}
}

func TestDesikePost(t *testing.T) {
	db, mock := mock.NewDatabaseConnection()
	defer db.Close()

	subTests := []struct {
		name           string
		errorInPrepare bool
		errorInExec    bool
		err            error
	}{
		{
			name:           "Deslike post",
			errorInPrepare: false,
			errorInExec:    false,
		},
		{
			name:           "Deslike post - error in prepare",
			errorInPrepare: true,
			errorInExec:    false,
			err:            errors.New("some error"),
		},
		{
			name:           "Deslike post - error in exec",
			errorInPrepare: false,
			errorInExec:    true,
			err:            errors.New("some error"),
		},
	}

	repository := repository.NewPostRepository(db)

	query := "update posts set likes = case when likes \\> 0 then likes \\- 1 else 0 end where id = \\?"

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {

			if subTest.errorInPrepare {
				mock.ExpectPrepare(query).WillReturnError(subTest.err)

				err := repository.DeslikePost(1)
				assert.ErrorIs(t, err, subTest.err)
			} else if subTest.errorInExec {
				prep := mock.ExpectPrepare(query)
				prep.ExpectExec().WithArgs(1).WillReturnError(subTest.err)

				err := repository.DeslikePost(1)
				assert.ErrorIs(t, err, subTest.err)
			} else {
				prep := mock.ExpectPrepare(query)
				prep.ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

				err := repository.DeslikePost(1)
				assert.NoError(t, err)
			}
		})
	}

}
