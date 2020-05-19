package main

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/suite"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

type GormSuit struct {
	suite.Suite
	repo *Repository
	mock sqlmock.Sqlmock
}

func TestGormSuit(t *testing.T) {
	suite.Run(t, new(GormSuit))
}

func (s *GormSuit) SetupTest() {
	db, mock, err := sqlmock.New()
	gdb, err := gorm.Open("mysql", db)
	assert.NoError(s.T(), err)
	repo := &Repository{db: gdb}
	s.repo = repo
	s.mock = mock
}

func (s *GormSuit) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *GormSuit) TestCreate() {
	// blog to insert
	blog := &Blog{
		Title:     "title",
		Content:   "hell sqlmock",
		CreatedAt: time.Now(),
	}

	const (
		sqlExec       = "INSERT INTO `blogs` (`title`,`content`,`created_at`) VALUES (?,?,?)"
		newID   int64 = 1
	)
	s.mock.ExpectBegin()
	result := sqlmock.NewResult(1, 1)
	s.mock.ExpectExec(regexp.QuoteMeta(sqlExec)).
		WithArgs(blog.Title, blog.Content, blog.CreatedAt).
		WillReturnResult(result)
	s.mock.ExpectCommit()
	assert.Equal(s.T(), int64(0), blog.ID)

	err := s.repo.Create(blog)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), newID, blog.ID)

}

func (s *GormSuit) TestUpdate() {
	// blog to update
	blog := &Blog{
		Title:   "newPost",
		Content: "newContent",
	}

	const sqlExec = "UPDATE `blogs` SET `content` = ?, `title` = ? WHERE (id = ?)"
	const id int64 = 1
	s.mock.ExpectBegin()
	result := sqlmock.NewResult(1, 1)
	s.mock.ExpectExec(regexp.QuoteMeta(sqlExec)).
		WithArgs(blog.Content, blog.Title, id).
		WillReturnResult(result)
	s.mock.ExpectCommit()

	err := s.repo.Update(id, blog)
	assert.NoError(s.T(), err)

}

func (s *GormSuit) TestGetByID() {
	wantBlog := &Blog{
		Title:     "post",
		Content:   "content",
		CreatedAt: time.Now(),
	}

	const sqlQuery = "SELECT * FROM `blogs` WHERE (id = ?) ORDER BY `blogs`.`id` ASC LIMIT 1"
	const id int64 = 1

	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at"}).
		AddRow(id, wantBlog.Title, wantBlog.Content, wantBlog.CreatedAt)
	s.mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
		WithArgs(id).WillReturnRows(rows)

	gotBlog, err := s.repo.GetByID(id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), wantBlog.Title, gotBlog.Title)
	assert.Equal(s.T(), wantBlog.Content, gotBlog.Content)

}

func (s *GormSuit) TestGetByTitle() {
	wantBlogs := []*Blog{
		&Blog{
			ID:        1,
			Title:     "post",
			Content:   "content1",
			CreatedAt: time.Now(),
		},
		&Blog{
			ID:        2,
			Title:     "post",
			Content:   "content2",
			CreatedAt: time.Now(),
		},
	}

	const sqlQuery = "SELECT * FROM `blogs` WHERE (title = ?)"
	const title = "post"

	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at"}).
		AddRow(wantBlogs[0].ID, wantBlogs[0].Title, wantBlogs[0].Content, wantBlogs[0].CreatedAt).
		AddRow(wantBlogs[1].ID, wantBlogs[1].Title, wantBlogs[1].Content, wantBlogs[1].CreatedAt)
	s.mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
		WithArgs(title).WillReturnRows(rows)

	gotBlogs, err := s.repo.GetByTitle(title)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), wantBlogs[0], gotBlogs[0])
	assert.Equal(s.T(), wantBlogs[1], gotBlogs[1])

}

func (s *GormSuit) TestDelete() {
	query := "DELETE FROM `blogs` WHERE (id = ?)"
	var id int64 = 1
	result := sqlmock.NewResult(1, 1)
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnResult(result)
	s.mock.ExpectCommit()
	err := s.repo.Delete(id)
	assert.NoError(s.T(), err)
}
