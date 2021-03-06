package pop_test

import (
	"os"
	"testing"

	"github.com/markbates/pop"
	"github.com/markbates/pop/slices"
	"github.com/stretchr/testify/require"
)

type Cake struct {
	Int    slices.Int    `db:"int_slice"`
	Float  slices.Float  `db:"float_slice"`
	String slices.String `db:"string_slice"`
}

func Test_String(t *testing.T) {
	if os.Getenv("SODA_DIALECT") == "postgres" {
		transaction(func(tx *pop.Connection) {
			r := require.New(t)

			c := &Cake{
				String: slices.String{"a", "b", "c"},
			}
			err := tx.Create(c)
			r.NoError(err)

			err = tx.Reload(c)
			r.Equal(slices.String{"a", "b", "c"}, c.String)
		})
	} else {
		t.SkipNow()
	}
}

func Test_Int(t *testing.T) {
	if os.Getenv("SODA_DIALECT") == "postgres" {
		transaction(func(tx *pop.Connection) {
			r := require.New(t)

			c := &Cake{
				Int: slices.Int{1, 2, 3},
				// Float: pop.Float{1.0, 2.1, 3.2},
			}
			err := tx.Create(c)
			r.NoError(err)

			err = tx.Reload(c)
			r.Equal(slices.Int{1, 2, 3}, c.Int)
		})
	} else {
		t.SkipNow()
	}
}

func Test_Float(t *testing.T) {
	if os.Getenv("SODA_DIALECT") == "postgres" {
		transaction(func(tx *pop.Connection) {
			r := require.New(t)

			c := &Cake{
				Float: slices.Float{1.0, 2.1, 3.2},
			}
			err := tx.Create(c)
			r.NoError(err)

			err = tx.Reload(c)
			r.Equal(slices.Float{1.0, 2.1, 3.2}, c.Float)
		})
	} else {
		t.SkipNow()
	}
}
