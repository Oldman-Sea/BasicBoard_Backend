package handlers

import (
	"board/database"
	"board/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPosts 게시글 목록 (커서/오프셋 병행)
func GetPosts(c *gin.Context) {
    limitStr := c.DefaultQuery("limit", "20")
    cursorStr := c.Query("cursor")   // 앱: ?cursor=42
    offsetStr := c.Query("offset")   // 웹: ?offset=20

    limit, _ := strconv.Atoi(limitStr)
    if limit > 100 {
        limit = 100
    }

    var posts []models.Post
    var totalCount int64

    database.DB.Model(&models.Post{}).Count(&totalCount)

    // === 커서 기반 (앱 무한스크롤) ===
    if cursorStr != "" {
        cursor, _ := strconv.Atoi(cursorStr)

        // cursor보다 작은 ID만 조회 (이전 데이터)
        database.DB.Where("id < ?", cursor).
            Order("created_at DESC, id DESC").
            Limit(limit + 1). // +1개 더 가져와서 hasMore 판단
            Find(&posts)

        hasMore := len(posts) > limit
        if hasMore {
            posts = posts[:limit] // 실제론 limit개만 반환
        }

        // 다음 커서 (마지막 항목의 ID)
        var nextCursor *uint
        if hasMore && len(posts) > 0 {
            lastID := posts[len(posts)-1].ID
            nextCursor = &lastID
        }

        c.JSON(200, models.PostListResponse{
            Posts:      posts,
            NextCursor: nextCursor,
            HasMore:    hasMore,
            TotalCount: totalCount,
        })
        return
    }

    // === 오프셋 기반 (웹 페이징) ===
    offset := 0
    if offsetStr != "" {
        offset, _ = strconv.Atoi(offsetStr)
    }

    database.DB.Order("created_at DESC, id DESC").
        Offset(offset).
        Limit(limit).
        Find(&posts)

    hasMore := offset+limit < int(totalCount)

    c.JSON(200, models.PostListResponse{
        Posts:      posts,
        HasMore:    hasMore,
        TotalCount: totalCount,
    })
}

// GetPost 게시글 상세 (조회수 증가)
func GetPost(c *gin.Context) {
    id := c.Param("id")

    var post models.Post
    if err := database.DB.First(&post, id).Error; err != nil {
        c.JSON(404, gin.H{"error": "Post not found"})
        return
    }

    // 조회수 증가 (updated_at 안 바뀌게)
    database.DB.Model(&post).UpdateColumn("view_count", post.ViewCount+1)

    c.JSON(200, post)
}

// CreatePost 게시글 작성
func CreatePost(c *gin.Context) {
    var req models.CreatePostRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 작성자 기본값
    if req.Author == "" {
        req.Author = "익명"
    }

    post := models.Post{
        Title:   req.Title,
        Content: req.Content,
        Author:  req.Author,
    }

    database.DB.Create(&post)
    c.JSON(201, post)
}

// UpdatePost 게시글 수정
func UpdatePost(c *gin.Context) {
    id := c.Param("id")

    var post models.Post
    if err := database.DB.First(&post, id).Error; err != nil {
        c.JSON(404, gin.H{"error": "Post not found"})
        return
    }

    var req models.UpdatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    post.Title = req.Title
    post.Content = req.Content

    database.DB.Save(&post)
    c.JSON(200, post)
}

// DeletePost 게시글 삭제
func DeletePost(c *gin.Context) {
    id := c.Param("id")

    result := database.DB.Delete(&models.Post{}, id)
    if result.RowsAffected == 0 {
        c.JSON(404, gin.H{"error": "Post not found"})
        return
    }

    c.JSON(200, gin.H{"message": "Deleted"})
}