package models

import "time"

// Post 게시글
type Post struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Title     string    `json:"title" gorm:"size:255;not null"`
    Content   string    `json:"content" gorm:"type:text;not null"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

// SearchHistory 검색 기록
type SearchHistory struct {
    ID         uint      `json:"id" gorm:"primaryKey"`
    Keyword    string    `json:"keyword" gorm:"size:255;not null;unique"`
    SearchedAt time.Time `json:"searchedAt"`
}

// CreatePostRequest 게시글 작성 요청
type CreatePostRequest struct {
    Title   string `json:"title" binding:"required"`
    Content string `json:"content" binding:"required"`
}

// UpdatePostRequest 게시글 수정 요청
type UpdatePostRequest struct {
    Title   string `json:"title" binding:"required"`
    Content string `json:"content" binding:"required"`
}

// Cursor 커서 정보
type Cursor struct {
    CreatedAt time.Time `json:"createdAt"`
    ID        uint      `json:"id"`
}

// PageResponse 웹용 페이지 기반 응답
type PageResponse[T any] struct {
    Items     []T   `json:"items"`
    Page      int   `json:"page"`
    Limit     int   `json:"limit"`
    Total     int64 `json:"total"`
    TotalPages int  `json:"totalPages"`
}

// CursorResponse 앱용 커서 기반 응답
type CursorResponse[T any] struct {
    Items      []T     `json:"items"`
    NextCursor *Cursor  `json:"nextCursor"`
    HasMore    bool     `json:"hasMore"`
}