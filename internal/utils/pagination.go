package utils

import (
	"net/http"
	"strconv"
)

// Page represents the pagination data
// @Description Pagination details
type Page struct {
	Limit   int `json:"per_page" example:"10"`
	Current int `json:"page_no" example:"1"`
	Total   int `json:"total" example:"486"`
}

// Offset returns the offset of the  page
func (p *Page) Offset() int {
	return (p.Current * p.Limit) - p.Limit
}

// NewPage is the factory function  a new page
func NewPage(r *http.Request) *Page {
	limit, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if limit < 1 {
		limit = 10
	}
	currentPageP := r.URL.Query().Get("page_no")
	currentPage, _ := strconv.Atoi(currentPageP)
	if currentPage <= 0 {
		currentPage = 1
	}

	return &Page{
		Limit:   limit,
		Current: currentPage,
	}
}
