package models

import (
	"math"
	"time"

	"gorm.io/gorm"
)

type Feed struct {
	gorm.Model  `json:"-"`
	ID          uint      `json:"id" gorm:"primaryKey; autoIncrement; not null"`
	Title       string    `json:"title" gorm:"type:varchar(200); min:10; max:200; not null; unique; index" validate:"required"`
	Description string    `json:"description" gorm:"type:varchar(250); not null; max:250" validate:"required"`
	Link        string    `json:"link" gorm:"not null" validate:"required"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// pagination
type Pagination struct {
	Page   int `json:"page"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
	Pages  int `json:"pages"`
	Offset int `json:"offset"`
}

// write a generic function to paginate the results of any model query
func Paginate(db *gorm.DB, feeds *[]Feed, page int, limit int, query string) Pagination {

	var pagination Pagination

	pagination.Page = page
	pagination.Limit = limit
	pagination.Offset = (page - 1) * limit

	if query != "" {
		db = db.Where("title ILIKE ?", "%"+query+"%").Limit(limit).Offset(pagination.Offset).Find(&feeds)
	} else {
		db.Limit(limit).Offset(pagination.Offset).Find(&feeds)
	}

	// count total number of records
	var count int64
	db.Model(&Feed{}).Count(&count)

	pagination.Total = int(count)
	pagination.Pages = int(math.Ceil(float64(pagination.Total) / float64(pagination.Limit)))

	return pagination
}

// gin post-save hook for gorm.Model
func (feed *Feed) AfterCreate(tx *gorm.DB) (err error) {
	return tx.Model(feed).UpdateColumn("created_at", feed.CreatedAt).Error
}

// gin post-update hook for gorm.Model
func (feed *Feed) AfterUpdate(tx *gorm.DB) (err error) {
	return tx.Model(feed).UpdateColumn("updated_at", feed.UpdatedAt).Error
}
