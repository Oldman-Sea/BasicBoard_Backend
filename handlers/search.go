package handlers

import (
	"board/database"
	"board/models"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SearchPosts 검색
// 웹: GET /search?q=검색어&page=1&limit=5 -> PageResponse<Post>
// 앱: GET /search?q=검색어&limit=5&cursorCreatedAt=...&cursorId=... -> CursorResponse<Post>
func SearchPosts(c *gin.Context) {
    keyword := c.Query("q")
    if keyword == "" {
        c.JSON(400, gin.H{"error": "요청값 잘못됨(누락)"})
        return
    }

    // 검색어에서 모든 공백 제거
    processedKeyword := strings.ReplaceAll(keyword, " ", "")
    processedKeyword = strings.TrimSpace(processedKeyword)

    // 공백만 입력한 경우 빈 결과 반환
    if processedKeyword == "" {
        // 웹: page 기반
        pageStr := c.Query("page")
        if pageStr != "" {
            page, _ := strconv.Atoi(pageStr)
            if page <= 0 {
                page = 1
            }
            limitStr := c.DefaultQuery("limit", "5")
            limit, _ := strconv.Atoi(limitStr)
            if limit <= 0 {
                limit = 5
            }

            c.JSON(200, models.PageResponse[models.Post]{
                Items:     []models.Post{},
                Page:      page,
                Limit:     limit,
                Total:     0,
                TotalPages: 0,
            })
            return
        }

        // 앱: cursor 기반
        limitStr := c.DefaultQuery("limit", "20")
        limit, _ := strconv.Atoi(limitStr)
        if limit <= 0 {
            limit = 20
        }

        c.JSON(200, models.CursorResponse[models.Post]{
            Items:      []models.Post{},
            NextCursor: nil,
            HasMore:    false,
        })
        return
    }

    // 검색 기록 저장 (공백 제거된 검색어로 저장)
    saveSearchHistory(processedKeyword)

    // 공백 제거된 검색어로 검색 패턴 생성
    searchPattern := "%" + processedKeyword + "%"

    // 웹: page 기반
    pageStr := c.Query("page")
    if pageStr != "" {
        limitStr := c.DefaultQuery("limit", "5")
        limit, _ := strconv.Atoi(limitStr)
        if limit > 100 {
            limit = 100
        }
        if limit <= 0 {
            limit = 5
        }

        page, _ := strconv.Atoi(pageStr)
        if page <= 0 {
            page = 1
        }

        var posts []models.Post
        var totalCount int64

        // 검색 조건으로 카운트
        database.DB.Model(&models.Post{}).
            Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern).
            Count(&totalCount)

        offset := (page - 1) * limit
        database.DB.Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern).
            Order("created_at DESC, id DESC").
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
    limitStr := c.DefaultQuery("limit", "5")
    limit, _ := strconv.Atoi(limitStr)
    if limit > 100 {
        limit = 100
    }
    if limit <= 0 {
        limit = 5
    }

    cursorCreatedAtStr := c.Query("cursorCreatedAt")
    cursorIdStr := c.Query("cursorId")

    var posts []models.Post

    query := database.DB.Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern).
        Order("created_at DESC, id DESC")

    // 커서가 있으면 그 이후 데이터 조회
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

// GetSearchHistory 최근 검색어 (최대 10개)
func GetSearchHistory(c *gin.Context) {
    var history []models.SearchHistory

    database.DB.Order("searched_at DESC").
        Limit(10).
        Find(&history)

    c.JSON(200, gin.H{"history": history})
}

// DeleteSearchHistory 검색어 삭제
func DeleteSearchHistory(c *gin.Context) {
    id := c.Param("id")

    database.DB.Delete(&models.SearchHistory{}, id)
    c.JSON(200, gin.H{"message": "Deleted"})
}

// ClearSearchHistory 전체 삭제
func ClearSearchHistory(c *gin.Context) {
    database.DB.Exec("TRUNCATE TABLE search_history")
    c.JSON(200, gin.H{"message": "All cleared"})
}

// 검색어 저장 헬퍼
func saveSearchHistory(keyword string) {
    keyword = strings.TrimSpace(keyword)
    if keyword == "" {
        return
    }

    var history models.SearchHistory

    // 기존 검색어가 있으면 시간만 업데이트 (UNIQUE 제약)
    result := database.DB.Where("keyword = ?", keyword).First(&history)

    if result.Error != nil {
        // 없으면 새로 추가
        history = models.SearchHistory{
            Keyword: keyword,
        }
        database.DB.Create(&history)
    } else {
        // 있으면 searched_at 업데이트
        database.DB.Model(&history).Update("searched_at", time.Now())
    }
}