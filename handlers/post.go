package handlers

import (
	"board/database"
	"board/models"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

// GetPosts 게시글 목록 조회
// 웹: GET /posts?page=1&limit=5 -> PageResponse<Post>
// 앱: GET /posts?limit=20&cursorCreatedAt=...&cursorId=... -> CursorResponse<Post>
func GetPosts(c *gin.Context) {
    

    // 웹: page 기반
    pageStr := c.Query("page")
    if pageStr != "" {

		limit := 5
	    if limitStr := c.Query("limit"); limitStr != "" {
	        limit, _ = strconv.Atoi(limitStr)
	    }

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

	limit := 20
    if limitStr := c.Query("limit"); limitStr != "" {
        limit, _ = strconv.Atoi(limitStr)
    }
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

	var prevId *uint
    var nextId *uint

    // 이전 글 id
    var prev uint
    if err := database.DB.
        Model(&models.Post{}).
        Select("id").
        Where("id < ?", post.ID).
        Order("id desc").
        Limit(1).
        Scan(&prev).Error; err == nil && prev != 0 {
        prevId = &prev
    }

    // 다음 글 id
    var next uint
    if err := database.DB.
        Model(&models.Post{}).
        Select("id").
        Where("id > ?", post.ID).
        Order("id asc").
        Limit(1).
        Scan(&next).Error; err == nil && next != 0 {
        nextId = &next
    }

    c.JSON(200, gin.H{
        "post":   post,
        "prevId": prevId,
        "nextId": nextId,
    })
}

// validateTitleLength 제목 길이 검증 (한글 45자/영어 72자)
// 바이트 길이로 근사치 계산: 한글 3바이트, 영어 1바이트
func validateTitleLength(title string) bool {
    title = strings.TrimSpace(title)
    if title == "" {
        return false
    }

    // 바이트 길이 계산
    byteLen := len([]byte(title))
    
    // 한글 45자 = 135바이트, 영어 72자 = 72바이트
    // 실제로는 한글과 영어가 섞여있을 수 있으므로, 최대 바이트 길이를 135로 제한
    // (한글만 있는 경우 45자, 영어만 있는 경우 72자보다 훨씬 많을 수 있음)
    // 더 정확하게는 유니코드 문자 수를 세어야 하지만, 바이트 길이로 근사치 계산
    if byteLen > 135 {
        return false
    }

    // 유니코드 문자 수로도 체크 (한글 45자 제한)
    runeCount := utf8.RuneCountInString(title)
    if runeCount > 45 {
        return false
    }

    return true
}

// CreatePost 게시글 작성
func CreatePost(c *gin.Context) {
    var req models.CreatePostRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "요청값 잘못됨(누락)"})
        return
    }

    // 제목/본문 공백 제거 후 검증
    req.Title = strings.TrimSpace(req.Title)
    req.Content = strings.TrimSpace(req.Content)

    if req.Title == "" || req.Content == "" {
        c.JSON(400, gin.H{"error": "요청값 잘못됨(누락)"})
        return
    }

    // 제목 길이 검증 (한글 45자/영어 72자)
    if !validateTitleLength(req.Title) {
        c.JSON(400, gin.H{"error": "제목 길이 제한을 초과했습니다 (한글 45자/영어 72자)"})
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

    // 제목/본문 공백 제거 후 검증
    req.Title = strings.TrimSpace(req.Title)
    req.Content = strings.TrimSpace(req.Content)

    if req.Title == "" || req.Content == "" {
        c.JSON(400, gin.H{"error": "요청값 잘못됨(누락)"})
        return
    }

    // 제목 길이 검증 (한글 45자/영어 72자)
    if !validateTitleLength(req.Title) {
        c.JSON(400, gin.H{"error": "제목 길이 제한을 초과했습니다 (한글 45자/영어 72자)"})
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