package pop

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/markbates/going/defaults"
	"github.com/markbates/going/randx"
)

// Connections contains all of the available connections
var Connections = map[string]*Connection{}

// Connection represents all of the necessary details for
// talking with a datastore
type Connection struct {
	ID      string
	Store   Store
	Dialect Dialect
	Timings []time.Duration
	TX      *tX
}

func (c *Connection) String() string {
	return c.Dialect.URL()
}

func (c *Connection) URL() string {
	return c.String()
}

// NewConnection creates a new connection, and sets it's `Dialect`
// appropriately based on the `ConnectionDetails` passed into it.
func NewConnection(deets *ConnectionDetails) *Connection {
	c := &Connection{
		ID:      randx.String(30),
		Timings: []time.Duration{},
	}
	switch deets.Dialect {
	case "postgres":
		c.Dialect = NewPostgreSQL(deets)
	case "mysql":
		c.Dialect = NewMySQL(deets)
		// case "sqlite3":
		// 	c.Dialect = NewSQLite(deets)
	}
	return c
}

// Connect takes the name of a connection, default is "development", and will
// return that connection from the available `Connections`. If a connection with
// that name can not be found an error will be returned. If a connection is
// found, and it has yet to open a connection with its underlying datastore,
// a connection to that store will be opened.
func Connect(e string) (*Connection, error) {
	e = defaults.String(e, "development")
	c := Connections[e]
	if c == nil {
		return c, fmt.Errorf("Could not find connection named %s!", e)
	}
	if c.Store != nil {
		return c, nil
	}
	db, err := sqlx.Open(c.Dialect.Details().Dialect, c.Dialect.URL())
	if err == nil {
		c.Store = &dB{db}
	}
	return c, nil
}

// Transaction will start a new transaction on the connection. If the inner function
// returns an error then the transaction will be rolled back, otherwise the transaction
// will automatically commit at the end.
func (c *Connection) Transaction(fn func(tx *Connection) error) error {
	var dberr error
	cn, err := c.NewTransaction()
	if err != nil {
		return err
	}
	err = fn(cn)
	if err != nil {
		dberr = cn.TX.Rollback()
	} else {
		dberr = cn.TX.Commit()
	}
	if err != nil {
		return err
	}
	return dberr
}

func (c *Connection) NewTransaction() (*Connection, error) {
	var cn *Connection
	if c.TX == nil {
		tx, err := c.Store.Transaction()
		if err != nil {
			return cn, err
		}
		cn = &Connection{
			ID:      randx.String(30),
			Store:   tx,
			Dialect: c.Dialect,
			Timings: []time.Duration{},
			TX:      tx,
		}
	} else {
		cn = c
	}
	return cn, nil
}

// Rollback will open a new transaction and automatically rollback that transaction
// when the inner function returns, regardless. This can be useful for tests, etc...
func (c *Connection) Rollback(fn func(tx *Connection)) error {
	var cn *Connection
	if c.TX == nil {
		tx, err := c.Store.Transaction()
		if err != nil {
			return err
		}
		cn = &Connection{
			ID:      randx.String(30),
			Store:   tx,
			Dialect: c.Dialect,
			Timings: []time.Duration{},
			TX:      tx,
		}
	} else {
		cn = c
	}
	fn(cn)
	return cn.TX.Rollback()
}

// Q creates a new "empty" query for the current connection.
func (c *Connection) Q() *Query {
	return Q(c)
}

func (c *Connection) timeFunc(name string, fn func() error) error {
	now := time.Now()
	err := fn()
	c.Timings = append(c.Timings, time.Now().Sub(now))
	return err
}
