package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/jackc/pgx"
	//"github.com/jackc/pgx/pgtype"
	//"github.com/jackc/pgx/stdlib"
	"github.com/jivalabs/picsum-photos/database"
	"github.com/jmoiron/sqlx"
)

// Provider implements a postgresql based storage
type Provider struct {
	db *sqlx.DB
}

type SpaceProvider struct {
	db *sqlx.DB
}

func NewSpace(address string) (*SpaceProvider, error) {
	db, err := sqlx.Connect("mysql", address)
	if err != nil {
		return nil, err
	}

	// Use Unsafe so that the app doesn't fail if we add new columns to the database
	return &SpaceProvider{
		db: db.Unsafe(),
	}, nil
}

// New returns a new Provider instance
func New(address string) (*Provider, error) {
	db, err := sqlx.Connect("mysql", address)
	if err != nil {
		return nil, err
	}

	// Use Unsafe so that the app doesn't fail if we add new columns to the database
	return &Provider{
		db: db.Unsafe(),
	}, nil
}

// Get returns the image data for an image id
func (p *Provider) Get(id string) (*database.Image, error) {
	i := &database.Image{}
	//p.db.Preparex("select * from image where id=?")
	a := fmt.Sprintf("select * from image where id = '%s'", id)
	err := p.db.Get(i, a)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, database.ErrNotFound
		}

		return nil, err
	}

	return i, nil
}

// GetRandom returns a random image ID
func (p *Provider) GetRandom() (id string, err error) {
	// This will be slow on large tables
	err = p.db.Get(&id, "select id from image order by rand() limit 1")
	return
}

// ListAll returns a list of all the images
func (p *Provider) ListAll() ([]database.Image, error) {
	i := []database.Image{}
	err := p.db.Select(&i, "select * from image order by id")

	if err != nil {
		return nil, err
	}

	return i, nil
}

// List returns a list of all the images with an offset/limit
func (p *Provider) List(offset, limit int) ([]database.Image, error) {
	i := []database.Image{}
	err := p.db.Select(&i, "select * from image order by id LIMIT ?, ?", offset, limit)

	if err != nil {
		return nil, err
	}

	return i, nil
}

// Shutdown shuts down the database client
func (p *Provider) Shutdown() {
	p.db.Close()
}

func (p *SpaceProvider) CreateSpace(name string, location string) (*database.Space, error) {
	stmsIns, err := p.db.Prepare("INSERT INTO space VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}

	defer stmsIns.Close()

	result, err := stmsIns.Exec(name, location, "")
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	space, err := p.GetSpaceById(id)
	if err != nil {
		return nil, err
	}
	return space, nil
}

func (p *SpaceProvider) GetSpaceList() ([]database.Space, error) {
	s := []database.Space{}
	err := p.db.Select(&s, "SELECT id, name, defaultLocale FROM space")
	return s, err
}

func (p *SpaceProvider) GetSpaceById(spaceId int64) (*database.Space, error) {
	s := &database.Space{}
	err := p.db.Get(s, "select * from space where id = ?", spaceId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, database.ErrNotFound
		}

		return nil, err
	}

	return s, nil
}
