package content

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pondok-tani-backend/config"
	contentmodels "pondok-tani-backend/models/content"

	"github.com/gin-gonic/gin"
)

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, " ", "-")
	value = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, value)
	value = strings.Trim(value, "-")
	if value == "" {
		return fmt.Sprintf("article-%d", time.Now().UnixNano())
	}
	return value
}

func saveArticleImage(c *gin.Context, fieldKey string) (string, error) {
	fileHeader, err := c.FormFile(fieldKey)
	if err != nil {
		return "", nil
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return "", fmt.Errorf("unsupported file extension")
	}

	if err := os.MkdirAll("uploads/articles", 0o755); err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("article_%d%s", time.Now().UnixNano(), ext)
	fullPath := filepath.Join("uploads", "articles", fileName)
	if err := c.SaveUploadedFile(fileHeader, fullPath); err != nil {
		return "", err
	}

	return "/uploads/articles/" + fileName, nil
}

func ListArticles(c *gin.Context) {
	category := strings.TrimSpace(c.DefaultQuery("category", ""))
	search := strings.TrimSpace(c.DefaultQuery("search", ""))
	status := strings.TrimSpace(c.DefaultQuery("status", "publish"))

	query := config.DB.Model(&contentmodels.Article{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if search != "" {
		query = query.Where("LOWER(title) LIKE ? OR LOWER(content) LIKE ? OR LOWER(excerpt) LIKE ?", "%"+strings.ToLower(search)+"%", "%"+strings.ToLower(search)+"%", "%"+strings.ToLower(search)+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var articles []contentmodels.Article
	if err := query.Order("published_at desc, created_at desc").Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles"})
		return
	}

	c.JSON(http.StatusOK, articles)
}

func GetArticle(c *gin.Context) {
	id := c.Param("id")
	var article contentmodels.Article
	if err := config.DB.First(&article, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}
	c.JSON(http.StatusOK, article)
}

func CreateArticle(c *gin.Context) {
	var req struct {
		Title       string `form:"title" binding:"required"`
		Category    string `form:"category"`
		Excerpt     string `form:"excerpt"`
		Content     string `form:"content" binding:"required"`
		Status      string `form:"status"`
		Author      string `form:"author"`
		PublishedAt string `form:"published_at"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article payload"})
		return
	}

	imageURL, err := saveArticleImage(c, "image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status == "" {
		status = "draft"
	}
	if status != "publish" && status != "draft" {
		status = "draft"
	}

	publishedAt := (*time.Time)(nil)
	if status == "publish" {
		if strings.TrimSpace(req.PublishedAt) != "" {
			parsed, err := time.Parse(time.RFC3339, req.PublishedAt)
			if err == nil {
				publishedAt = &parsed
			}
		} else {
			now := time.Now().UTC()
			publishedAt = &now
		}
	}

	article := contentmodels.Article{
		Title:       strings.TrimSpace(req.Title),
		Slug:        slugify(req.Title),
		Category:    strings.TrimSpace(req.Category),
		Excerpt:     strings.TrimSpace(req.Excerpt),
		Content:     strings.TrimSpace(req.Content),
		ImageURL:    imageURL,
		Status:      status,
		Author:      strings.TrimSpace(req.Author),
		PublishedAt: publishedAt,
	}

	if article.Category == "" {
		article.Category = "Umum"
	}
	if article.Author == "" {
		article.Author = "Admin"
	}

	if err := config.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
		return
	}

	c.JSON(http.StatusCreated, article)
}

func UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	var article contentmodels.Article
	if err := config.DB.First(&article, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	var req struct {
		Title       string `form:"title"`
		Category    string `form:"category"`
		Excerpt     string `form:"excerpt"`
		Content     string `form:"content"`
		Status      string `form:"status"`
		Author      string `form:"author"`
		PublishedAt string `form:"published_at"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article payload"})
		return
	}

	imageURL, err := saveArticleImage(c, "image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if imageURL != "" {
		article.ImageURL = imageURL
	}

	if strings.TrimSpace(req.Title) != "" {
		article.Title = strings.TrimSpace(req.Title)
		article.Slug = slugify(req.Title)
	}
	if strings.TrimSpace(req.Category) != "" {
		article.Category = strings.TrimSpace(req.Category)
	}
	article.Excerpt = strings.TrimSpace(req.Excerpt)
	article.Content = strings.TrimSpace(req.Content)
	if strings.TrimSpace(req.Status) != "" {
		status := strings.ToLower(strings.TrimSpace(req.Status))
		if status != "publish" && status != "draft" {
			status = "draft"
		}
		article.Status = status
	}
	if strings.TrimSpace(req.Author) != "" {
		article.Author = strings.TrimSpace(req.Author)
	}

	if article.Status == "publish" && article.PublishedAt == nil {
		now := time.Now().UTC()
		article.PublishedAt = &now
	}

	if err := config.DB.Save(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
		return
	}

	c.JSON(http.StatusOK, article)
}

func DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&contentmodels.Article{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Article deleted"})
}
