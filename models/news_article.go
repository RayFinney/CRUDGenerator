package models

import (
	"errors"
)

type NewsArticle struct {
	UUID         string  `json:"uuid"`
	Title        string  `json:"title"`
	Article      string  `json:"article"`
	Views        int64   `json:"views"`
	Rating       float64 `json:"rating"`
	Public       bool    `json:"public"`
	CreatedBy    string  `json:"created_by"`
	ModifiedBy   string  `json:"modified_by"`
	DeletedBy    string  `json:"deleted_by"`
	CreatedDate  string  `json:"created_date"`
	LastModified string  `json:"last_modified"`
	DeletedDate  string  `json:"deleted_date"`
}

func (v *NewsArticle) Validate() error {
	if v.UUID == "" {
		return errors.New("UUID required")
	}
	if v.Title == "" {
		return errors.New("Title required")
	}
	if len(v.Title) > 50 {
		return errors.New("Title to large (max 50)")
	}
	if v.Article == "" {
		return errors.New("Article required")
	}
	if v.Rating == 0.0 {
		return errors.New("Rating required")
	}
	if v.CreatedBy == "" {
		return errors.New("CreatedBy required")
	}
	if v.CreatedDate == "" {
		return errors.New("CreatedDate required")
	}
	return nil
}
