package database

import (
	"errors"
)

// Image contains metadata about an image
type Image struct {
	ID     string `json:"id"`
	Author string `json:"author"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

type Space struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	DefaultLocale string `json:"defaultLocale"`
}

// Provider is an interface for listing and retrieving images
type Provider interface {
	Get(id string) (*Image, error)
	GetRandom() (id string, err error)
	ListAll() ([]Image, error)
	List(offset, limit int) ([]Image, error)
	Shutdown()
}

type SpaceProvider interface {
	GetSpaceById(spaceId int64) (*Space, error)
	GetSpaceList() ([]Space, error)
}

// Errors
var (
	ErrNotFound = errors.New("Image does not exist")
)
