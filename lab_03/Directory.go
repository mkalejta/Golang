package main

import (
	"errors"
	"fmt"
	"time"
)

type Dir struct {
	name       string
	path       string
	size       int64
	createdAt  time.Time
	modifiedAt time.Time
	data       []FileSystemItem
}

func (p *Dir) Name() string          { return p.name }
func (p *Dir) Path() string          { return p.path }
func (p *Dir) Size() int64           { return p.size }
func (p *Dir) CreatedAt() time.Time  { return p.createdAt }
func (p *Dir) ModifiedAt() time.Time { return p.modifiedAt }

func (k *Dir) AddItem(item FileSystemItem) error {
	k.data = append(k.data, item)
	return nil
}

func (k *Dir) RemoveItem(item FileSystemItem) error {
	for i, elem := range k.data {
		if elem.Name() == item.Name() {
			k.data = append(k.data[:i], k.data[i+1:]...)
			return nil
		}
	}
	return errors.New("item not found")
}

func (k *Dir) OpenItem(item FileSystemItem) error {
	for _, elem := range k.data {
		if elem.Name() == item.Name() {
			fmt.Println(elem)
			return nil
		}
	}
	return errors.New("item not found")
}

func (k *Dir) Items() []FileSystemItem {
	return k.data
}
