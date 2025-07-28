package repository

type CreateOptions struct {
	Name        string
	Description string
	Alias       string
}

type UpdateOptions struct {
	ID          string
	Name        string
	Description string
	Alias       string
}
