package widgets

import (
	"math/rand"
	"testing"
	"time"

	copan "github.com/rowdyroad/go-http-control-panel"
)

func TestDashboard(t *testing.T) {
	cp := copan.NewControlPanel(copan.Config{
		Listen: ":9999",
		Title:  "config",
	})

	db := NewDashboard(cp, time.Second)

	go func() {
		for {
			for i := 'Z'; i >= 'A'; i-- {
				db.Set(string(i), rand.Int63())
			}
			time.Sleep(time.Second)
		}
	}()

	cp.AddContentPage("/", "Test", "{{element .}}", db.ID())

	cp.Run()
}
