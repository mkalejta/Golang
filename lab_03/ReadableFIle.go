package main

import "time"

// Struktura definiująca PlikDoOdczytu
type ReadableFile struct {
	name       string
	path       string
	data       []byte
	createdAt  time.Time
	modifiedAt time.Time
}

// Metoda tworząca nowy PlikDoOdczytu
func createReadableFile(name, path string) *ReadableFile {
	return &ReadableFile{
		name:       name,
		path:       path,
		data:       nil,
		createdAt:  time.Now(),
		modifiedAt: time.Now(),
	}
}

func (f *ReadableFile) Name() string          { return f.name }
func (f *ReadableFile) Path() string          { return f.path }
func (f *ReadableFile) Size() int64           { return int64(len(f.data)) }
func (f *ReadableFile) CreatedAt() time.Time  { return f.createdAt }
func (f *ReadableFile) ModifiedAt() time.Time { return f.modifiedAt }

func (f *ReadableFile) UpdateContent(newContent []byte) {
	f.data = newContent
	f.modifiedAt = time.Now()
}

func (f *ReadableFile) Read(p []byte) (n int, err error) {
	copy(p, f.data)
	return len(f.data), nil
}
