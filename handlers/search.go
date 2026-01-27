package handlers

import (
	"board/database"
	"board/models"

	"github.com/gin-gonic/gin"
)

// SearchPosts 통합 검색 (제목 + 본문)
func SearchPosts(c *gin.Context) {
    keyword := c.Query("q")
    if keyword == "" {
        c.JSON(400, gin.H{"error": "검색어를 입력하세요"})
        return
    }

    // 검색 기록 저장
    saveSearchHistory(keyword)

    var posts []models.Post

    // FULLTEXT 검색 (MySQL)
    database.DB.Where("MATCH(title, content) AGAINST(? IN NATURAL LANGUAGE MODE)", keyword).
        Order("created_at DESC").
        Limit(50).
        Find(&posts)

    c.JSON(200, gin.H{
        "keyword": keyword,
        "posts":   posts,
        "count":   len(posts),
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
        database.DB.Model(&history).Update("searched_at", "NOW()")
    }
}