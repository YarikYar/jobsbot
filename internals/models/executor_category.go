package models

type ExecutorCategory struct {
	ID   string `json:"id,omitempty" db:"id"`
	Name string `json:"name,omitempty" db:"name"`
	Data string `json:"data,omitempty" db:"data"`
}
