package storage

import "errors"

var (
	ErrURKNotFound = errors.New("URL isn`t found")
	ErrUrlExists   = errors.New("URL already exists")
)
