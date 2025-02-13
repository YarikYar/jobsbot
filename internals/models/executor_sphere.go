package models

type ExecutorSphere struct {
	ID   string `json:"id,omitempty" db:"id"`
	Name string `json:"name,omitempty" db:"name"`
	Uniq string `json:"uniq,omitempty" db:"uniq"`
	Hash string `json:"hash,omitempty" db:"hash"`
}
