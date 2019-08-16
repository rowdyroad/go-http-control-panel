package widgets

import (
	"log"
	"math/rand"
	"testing"
	"time"

	copan "github.com/rowdyroad/go-http-control-panel"
)

func TestSetUnset(t *testing.T) {
	cp := copan.NewControlPanel(copan.Config{
		Listen: ":9999",
		Title:  "config",
	})
	db := NewDashboard(cp, time.Second)
	db.Set("A", 1, "B", 2, "C", 3)
	db.Unset("A")
	db.Set("A", 4, "B", 5, "C", 6)
	db.Unset("B")
	db.Set("A", 7, "B", 8, "C", 9)
	db.Unset("C")
	db.Set("A", 10, "B", 11, "C", 12)

	log.Println(db.keys, db.values)
}
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
