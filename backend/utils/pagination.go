package utils

import (
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaginationData struct {
	Limit      int         `json:"limit,omitempty;query:limit"`
	Page       int         `json:"page,omitempty;query:page"`
	Sort       string      `json:"sort,omitempty;query:sort"`
	Search     string      `json:"search,omitempty;query:search"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Rows       interface{} `json:"rows"`
}

func (p *PaginationData) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *PaginationData) GetLimit() int {
	if p.Limit <= 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *PaginationData) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *PaginationData) GetSort() string {
	if p.Sort == "" {
		p.Sort = "id desc"
	}
	return p.Sort
}

// GeneratePaginationFromRequest creates a PaginationData struct from Gin context query parameters
func GeneratePaginationFromRequest(c *gin.Context) PaginationData {
	limit := 10
	page := 1
	sort := "id desc"
	search := ""

	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
	}

	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}

	if s := c.Query("sort"); s != "" {
		sort = s
	}

	if q := c.Query("search"); q != "" {
		search = q
	}

	return PaginationData{
		Limit:  limit,
		Page:   page,
		Sort:   sort,
		Search: search,
	}
}

// Search scope for GORM
func Search(columns []string, query string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if query == "" || len(columns) == 0 {
			return db
		}

		var queryConditions []string
		var queryValues []interface{}

		for _, col := range columns {
			queryConditions = append(queryConditions, col+" ILIKE ?")
			queryValues = append(queryValues, "%"+query+"%")
		}

		queryString := strings.Join(queryConditions, " OR ")
		return db.Where(queryString, queryValues...)
	}
}

// Paginate applies pagination and executes the query
func Paginate(db *gorm.DB, pagination *PaginationData, result interface{}) *gorm.DB {
	// Count total rows without limit and offset
	db.Count(&pagination.TotalRows)
	
	pagination.TotalPages = int(math.Ceil(float64(pagination.TotalRows) / float64(pagination.GetLimit())))

	// Apply pagination, sorting, and find rows
	err := db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort()).Find(result).Error
	
	pagination.Rows = result
	
	if err != nil {
		db.Error = err
	}
	
	return db
}

// FormatPaginationResponse standardizes the pagination response format
func FormatPaginationResponse(pagination *PaginationData) map[string]interface{} {
	return map[string]interface{}{
		"data": pagination.Rows,
		"meta": map[string]interface{}{
			"limit":       pagination.GetLimit(),
			"page":        pagination.GetPage(),
			"sort":        pagination.GetSort(),
			"total_rows":  pagination.TotalRows,
			"total_pages": pagination.TotalPages,
		},
	}
}
