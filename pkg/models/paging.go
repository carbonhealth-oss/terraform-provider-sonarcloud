package models

// Paging used in /search API endpoints
type Paging struct {
PageIndex int64 `json:"pageIndex"`
PageSize  int64 `json:"pageSize"`
Total     int64 `json:"total"`
}