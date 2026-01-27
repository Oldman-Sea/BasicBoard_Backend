package models

import "time"

// Post 게시글
type Post struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Title     string    `json:"title" gorm:"size:255;not null"`
    Content   string    `json:"content" gorm:"type:text;not null"`
    Author    string    `json:"author" gorm:"size:50;default:익명"`
    ViewCount uint      `json:"viewCount" gorm:"default:0"`
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
    Author  string `json:"author"`
}

// UpdatePostRequest 게시글 수정 요청
type UpdatePostRequest struct {
    Title   string `json:"title" binding:"required"`
    Content string `json:"content" binding:"required"`
}

// PostListResponse 목록 응답
type PostListResponse struct {
    Posts      []Post `json:"posts"`
    NextCursor *uint  `json:"nextCursor,omitempty"` // 앱: 무한스크롤용
    HasMore    bool   `json:"hasMore"`
    TotalCount int64  `json:"totalCount"`           // 웹: 페이지 계산용
}