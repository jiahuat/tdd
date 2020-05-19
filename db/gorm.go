package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Repository struct {
	db *gorm.DB
}

type Blog struct {
	ID        int64
	Title     string
	Content   string
	CreatedAt time.Time
}

func (p *Repository) Create(b *Blog) error {
	return p.db.Create(b).Error
}

func (p *Repository) Update(id int64, b *Blog) error {
	return p.db.Table("blogs").Where("id = ?", id).Update(b).Error
}

func (p *Repository) Delete(id int64) error {
	b := &Blog{}
	return p.db.Where("id = ?", id).Delete(b).Error
}

func (p *Repository) GetByID(id int64) (*Blog, error) {
	blog := &Blog{}
	err := p.db.Model(blog).Where("id = ?", id).First(blog).Error
	return blog, err
}

func (p *Repository) GetByTitle(title string) ([]*Blog, error) {
	var blogs []*Blog
	err := p.db.Table("blogs").Where("title = ?", title).Find(&blogs).Error
	return blogs, err
}

func main() {

}
