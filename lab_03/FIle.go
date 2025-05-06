package main

import "time"

// Struktura definiująca Plik
type File struct {
	name       string
	path       string
	data       []byte
	createdAt  time.Time
	modifiedAt time.Time
}

// Metoda tworząca nowy Plik
func createFile(name, path string) *File {
	return &File{
		name:       name,
		path:       path,
		data:       nil,
		createdAt:  time.Now(),
		modifiedAt: time.Now(),
	}
}

func (f *File) Name() string          { return f.name }
func (f *File) Path() string          { return f.path }
func (f *File) Size() int64           { return int64(len(f.data)) }
func (f *File) CreatedAt() time.Time  { return f.createdAt }
func (f *File) ModifiedAt() time.Time { return f.modifiedAt }

func (f *File) UpdateContent(newData []byte) {
	f.data = newData
	f.modifiedAt = time.Now()
}

func (f *File) Read(p []byte) (n int, err error) {
	copy(p, f.data)
	return len(f.data), nil
}

func (f *File) Write(p []byte) (n int, err error) {
	f.data = p
	f.modifiedAt = time.Now()
	return len(p), nil
}

func (f *File) Append(p []byte) (n int, err error) {
	f.data = append(f.data, p...)
	f.modifiedAt = time.Now()
	return len(p), nil
}
