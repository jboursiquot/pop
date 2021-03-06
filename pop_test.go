package pop_test

import (
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/markbates/going/nulls"
	"github.com/markbates/pop"
	_ "github.com/mattes/migrate/migrate"
	// _ "github.com/mattn/go-sqlite3"
)

var PDB *pop.Connection

func init() {
	pop.Debug = false
	pop.AddLookupPaths("./")

	dialect := os.Getenv("SODA_DIALECT")

	var err error
	PDB, err = pop.Connect(dialect)
	if err != nil {
		log.Panic(err)
	}

	pop.MapTableName("Friend", "good_friends")
	pop.MapTableName("Friends", "good_friends")
}

func transaction(fn func(tx *pop.Connection)) {
	PDB.Rollback(func(tx *pop.Connection) {
		fn(tx)
	})
}

func ts(s string) string {
	return PDB.Dialect.TranslateSQL(s)
}

type User struct {
	ID        int           `db:"id"`
	Name      nulls.String  `db:"name"`
	Alive     nulls.Bool    `db:"alive"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
	BirthDate nulls.Time    `db:"birth_date"`
	Bio       nulls.String  `db:"bio"`
	Price     nulls.Float64 `db:"price"`
	FullName  nulls.String  `db:"full_name" select:"name as full_name"`
}

type Users []User

type Friend struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

type Friends []Friend

type Enemy struct {
	A string
}
