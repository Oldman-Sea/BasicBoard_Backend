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
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    })

    // API 라우트
    api := r.Group("/api")
    {
        // 게시글
        api.GET("/posts", handlers.GetPosts)           // 목록
        api.GET("/posts/:id", handlers.GetPost)        // 상세
        api.POST("/posts", handlers.CreatePost)        // 작성
        api.PUT("/posts/:id", handlers.UpdatePost)     // 수정
        api.DELETE("/posts/:id", handlers.DeletePost)  // 삭제

        // 검색
        api.GET("/search", handlers.SearchPosts)                      // 검색
        api.GET("/search/history", handlers.GetSearchHistory)         // 최근 검색어
        api.DELETE("/search/history/:id", handlers.DeleteSearchHistory) // 검색어 삭제
        api.DELETE("/search/history", handlers.ClearSearchHistory)      // 전체 삭제
    }

    r.Run(":8080")
}