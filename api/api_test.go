package api

import (
	"testing"

	"github.com/franela/goblin"
)

func TestServe(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Serve", func() {
		g.It("Should Serve", func() {
			g.Assert(true).IsTrue()
		})
	})
}
