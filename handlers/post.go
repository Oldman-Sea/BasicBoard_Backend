package handlers

import (
	"board/database"
	"board/models"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetPosts 게시글 목록 조회
// 웹: GET /posts?page=1&limit=20 -> PageResponse<Post>
// 앱: GET /posts?limit=20&cursorCreatedAt=...&cursorId=... -> CursorResponse<Post>
func GetPosts(c *gin.Context) {
    limitStr := c.DefaultQuery("limit", "5")
    limit, _ := strconv.Atoi(limitStr)
    if limit > 100 {
        limit = 100
    }
    if limit <= 0 {
        limit = 5
    }

    // 웹: page 기반
    pageStr := c.Query("page")
    if pageStr != "" {
        page, _ := strconv.Atoi(pageStr)
        if page <= 0 {
            page = 1
        }

        var posts []models.Post
        var totalCount int64

        database.DB.Model(&models.Post{}).Count(&totalCount)

        offset := (page - 1) * limit
        database.DB.Order("created_at DESC, id DESC").
            Offset(offset).
            Limit(limit).
            Find(&posts)

        totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

        c.JSON(200, models.PageResponse[models.Post]{
            Items:     posts,
            Page:      page,
            Limit:     limit,
            Total:     totalCount,
            TotalPages: totalPages,
        })
        return
    }

    // 앱: cursor 기반
    cursorCreatedAtStr := c.Query("cursorCreatedAt")
    cursorIdStr := c.Query("cursorId")

    var posts []models.Post

    query := database.DB.Order("created_at DESC, id DESC")

    // 커서가 있으면 그 이후 데이터 조회
	// 이전에 본 게시글보다 더 오래된 것만 조회
    if cursorCreatedAtStr != "" && cursorIdStr != "" {
        cursorCreatedAt, err := time.Parse(time.RFC3339, cursorCreatedAtStr)
        if err == nil {
            cursorId, _ := strconv.Atoi(cursorIdStr)
            query = query.Where("(created_at < ?) OR (created_at = ? AND id < ?)",
                cursorCreatedAt, cursorCreatedAt, cursorId)
        }
    }

    // limit+1개 가져와서 hasMore 판단
    query.Limit(limit + 1).Find(&posts)

    hasMore := len(posts) > limit
    if hasMore {
        posts = posts[:limit]
    }

    var nextCursor *models.Cursor
    if hasMore && len(posts) > 0 {
        lastPost := posts[len(posts)-1]
        nextCursor = &models.Cursor{
            CreatedAt: lastPost.CreatedAt,
            ID:        lastPost.ID,
        }
    }

    c.JSON(200, models.CursorResponse[models.Post]{
        Items:      posts,
        NextCursor: nextCursor,
        HasMore:    hasMore,
    })
}

// GetPost 게시글 상세 조회
func GetPost(c *gin.Context) {
    id := c.Param("id")

    var post models.Post
    if err := database.DB.First(&post, id).Error; err != nil {
        c.JSON(404, gin.H{"error": "대상없음(없는id)"})
        return
    }

    c.JSON(200, post)
}

// CreatePost 게시글 작성
func CreatePost(c *gin.Context) {
    var req models.CreatePostRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "요청값 잘못됨(누락)"})
        return
    }

    post := models.Post{
        Title:   req.Title,
        Content: req.Content,
    }

    if err := database.DB.Create(&post).Error; err != nil {
        c.JSON(500, gin.H{"error": "서버 내부 오류"})
        return
    }

    c.JSON(201, post)
}

// UpdatePost 게시글 수정
func UpdatePost(c *gin.Context) {
    id := c.Param("id")

    var post models.Post
    if err := database.DB.First(&post, id).Error; err != nil {
        c.JSON(404, gin.H{"error": "대상없음(없는id)"})
        return
    }

    var req models.UpdatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "요청값 잘못됨(누락)"})
        return
    }

    post.Title = req.Title
    post.Content = req.Content

    if err := database.DB.Save(&post).Error; err != nil {
        c.JSON(500, gin.H{"error": "서버 내부 오류"})
        return
    }

    c.JSON(201, post)
}

// DeletePost 게시글 삭제
func DeletePost(c *gin.Context) {
    id := c.Param("id")

    result := database.DB.Delete(&models.Post{}, id)
    if result.RowsAffected == 0 {
        c.JSON(404, gin.H{"error": "대상없음(없는id)"})
        return
    }

    if result.Error != nil {
        c.JSON(500, gin.H{"error": "서버 내부 오류"})
        return
    }

    c.JSON(204, nil)
}