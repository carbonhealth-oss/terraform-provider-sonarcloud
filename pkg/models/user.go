package models
// User struct
type User struct {
	Login       string   `json:"login,omitempty"`
	Name        string   `json:"name,omitempty"`
	Email       string   `json:"email,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	IsActive    bool     `json:"active,omitempty"`
	IsLocal     bool     `json:"local,omitempty"`
}

// GetUser for unmarshalling response body where users are returned
type GetUser struct {
	Paging Paging `json:"paging"`
	Users  []User `json:"users"`
}

