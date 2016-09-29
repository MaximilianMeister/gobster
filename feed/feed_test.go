package feed

import (
	"reflect"
	"testing"

	"github.com/franela/goblin"
)

func TestOpen(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Open Database", func() {
		g.It("Should Return A Database", func() {
			db, err := Open()
			defer db.Close()
			g.Assert(err == nil).IsTrue()
			g.Assert(reflect.TypeOf(db).String()).Equal("*simplebolt.Database")
		})
	})
}

func TestSet(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Add Entry To List", func() {
		g.It("Should Return Successfully", func() {
			err := Set("somename", "This is a quote")
			g.Assert(err == nil).IsTrue()
		})
	})
}

func TestGet(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Get Random List Entry", func() {
		g.It("Should Return A Quote", func() {
			quote, err := Get("somename")
			g.Assert(err == nil).IsTrue()
			g.Assert(reflect.TypeOf(quote).String()).Equal("string")
		})
	})
}

func TestGetAll(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Get All List Entries", func() {
		g.It("Should Return A Numbered List of Quotes", func() {
			quote, err := GetAll("somename")
			g.Assert(err == nil).IsTrue()
			g.Assert(reflect.TypeOf(quote).String()).Equal("[]string")
		})
	})
}
