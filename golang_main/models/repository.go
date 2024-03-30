package models
// Repository represents the structure of a code repository.
type Repository struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Branch    string `json:"branch"`
	RepoPath  string `json:"repoPath"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}