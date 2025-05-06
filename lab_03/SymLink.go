package main

import "time"

// Struktura definiująca Dowiązanie Symboliczne
type SymLink struct {
	name       string
	path       string
	createdAt  time.Time
	modifiedAt time.Time
	referencja FileSystemItem
}

func (p *SymLink) Name() string          { return p.name }
func (p *SymLink) Path() string          { return p.path }
func (p *SymLink) CreatedAt() time.Time  { return p.createdAt }
func (p *SymLink) ModifiedAt() time.Time { return p.modifiedAt }

func (s *SymLink) Size() int64 {
	if s.referencja != nil {
		return s.referencja.Size()
	}
	return 0
}
