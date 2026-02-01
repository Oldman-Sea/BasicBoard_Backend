package main

import (
	"board/database"
	"board/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
    // DB 초기화
    database.InitDB()

    // Gin 라우터
    r := gin.Default()

    // CORS (프론트엔드 연동용)
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        // Json 요청 허락
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// 실제 요청 전에 OPTIONS 요청을 먼저 보낼 때 '이 API 써도 돼?' 묻는 용도
        if c.Request.Method == "OPTIONS" {
			// 그냥 204(No Content) 응답
            c.AbortWithStatus(204)
            return
        }

		// 해당 미들웨어 이후의 실제 API로직 실행
        c.Next()
    })

    // API 라우트
    api := r.Group("/api")
    {
        // 게시글
        api.GET("/posts", handlers.GetPosts)           // 목록 (웹: page, limit / 앱: limit, cursorCreatedAt, cursorId)
        api.GET("/posts/:id", handlers.GetPost)        // 상세
        api.POST("/posts", handlers.CreatePost)        // 작성
        api.PUT("/posts/:id", handlers.UpdatePost)     // 수정
        api.DELETE("/posts/:id", handlers.DeletePost)  // 삭제

        // 검색
        api.GET("/search", handlers.SearchPosts)                      // 검색 (웹: q, page, limit / 앱: q, limit, cursorCreatedAt, cursorId)
		// 해당 기능 있는 게 맞나 확인
        api.GET("/search/history", handlers.GetSearchHistory)         // 최근 검색어
        api.DELETE("/search/history/:id", handlers.DeleteSearchHistory) // 검색어 삭제
        api.DELETE("/search/history", handlers.ClearSearchHistory)      // 전체 삭제
    }

    r.Run(":8080")
}