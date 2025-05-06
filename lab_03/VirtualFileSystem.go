package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type VirtualFileSystem struct {
	root *Dir
}

func NewVirtualFileSystem() *VirtualFileSystem {
	return &VirtualFileSystem{
		root: &Dir{
			name:       "root",
			path:       "/",
			size:       0,
			createdAt:  time.Now(),
			modifiedAt: time.Now(),
			data:       []FileSystemItem{},
		},
	}
}

func (vfs *VirtualFileSystem) findDirectory(path string) (*Dir, error) {
	if path == "/" {
		return vfs.root, nil
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")
	current := vfs.root

	for _, part := range parts {
		found := false
		for _, item := range current.data {
			if Dir, ok := item.(*Dir); ok && Dir.name == part {
				current = Dir
				found = true
				break
			}
		}
		if !found {
			return nil, ErrItemNotFound
		}
	}

	return current, nil
}

func (vfs *VirtualFileSystem) CreateFile(path, name string, readonly bool, b []byte) error {
	dir, err := vfs.findDirectory(path)
	if err != nil {
		return err
	}

	for _, item := range dir.data {
		if item.Name() == name {
			return ErrItemExists
		}
	}

	var file FileSystemItem
	if readonly {
		file = &ReadableFile{
			name:       name,
			path:       path + "/" + name,
			data:       b,
			createdAt:  time.Now(),
			modifiedAt: time.Now(),
		}
	} else {
		file = &File{
			name:       name,
			path:       path + "/" + name,
			data:       b,
			createdAt:  time.Now(),
			modifiedAt: time.Now(),
		}
	}
	dir.AddItem(file)
	return nil
}

func (vfs *VirtualFileSystem) CreateDirectory(path, name string) error {
	dir, err := vfs.findDirectory(path)
	if err != nil {
		return err
	}

	for _, item := range dir.data {
		if item.Name() == name {
			return ErrItemExists
		}
	}

	subDir := &Dir{
		name:       name,
		path:       path + "/" + name,
		createdAt:  time.Now(),
		modifiedAt: time.Now(),
		data:       []FileSystemItem{},
	}
	dir.AddItem(subDir)
	return nil
}

func (vfs *VirtualFileSystem) FindItem(path string) (FileSystemItem, error) {
	if path == "/" {
		return vfs.root, nil
	}
	return vfs.findItem(vfs.root, path)
}

func (vfs *VirtualFileSystem) findItem(dir *Dir, path string) (FileSystemItem, error) {
	if dir.Path() == path {
		return dir, nil
	}

	for _, item := range dir.data {
		if item.Path() == path {
			return item, nil
		}
		if Dir, ok := item.(*Dir); ok {
			if found, err := vfs.findItem(Dir, path); err == nil {
				return found, nil
			}
		}
	}
	return nil, ErrItemNotFound
}

func (vfs *VirtualFileSystem) WriteToFile(path string, data []byte) error {
	item, err := vfs.FindItem(path)
	if err != nil {
		return err
	}

	if file, ok := item.(*File); ok {
		_, err := file.Write(data)
		return err
	}
	return errors.New("cannot write to a read-only file or directory")
}

func (vfs *VirtualFileSystem) ReadFromFile(path string) ([]byte, error) {
	item, err := vfs.FindItem(path)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 1024)
	if file, ok := item.(*File); ok {
		n, err := file.Read(buffer)
		return buffer[:n], err
	} else if file, ok := item.(*ReadableFile); ok {
		n, err := file.Read(buffer)
		return buffer[:n], err
	}
	return nil, errors.New("cannot read from a directory")
}

func (vfs *VirtualFileSystem) OpenFile(path string) error {
	data, err := vfs.ReadFromFile(path)
	if err != nil {
		return err
	}
	fmt.Println("Opening file:", path)
	fmt.Println("Data:", string(data))
	return nil
}

func (vfs *VirtualFileSystem) DeleteItem(path string) error {
	if path == "/" {
		return ErrPermissionDenied
	}
	parentPath := path[:strings.LastIndex(path, "/")]
	name := path[strings.LastIndex(path, "/")+1:]

	dir, err := vfs.findDirectory(parentPath)
	if err != nil {
		return err
	}

	return dir.RemoveItem(&File{name: name})
}
