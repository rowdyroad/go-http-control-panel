package copan

import (
	"log"
	"testing"
	"time"
)

type Struct struct {
	String        string
	Int           int
	Float         float64
	Bool          bool
	Duration      time.Duration
	Strings       []string
	Bools         []bool
	Ints          []int
	Floats        []float64
	Configuration Config
}

func TestMain(t *testing.T) {
	cp := NewControlPanel(Config{
		Listen: ":8011",
		Title:  "config",
		Users: []User{
			User{
				"admin",
				"1q2w3e4r",
			},
		},
	})

	s := Struct{
		String: "string",
		Int:    1,
		Float:  0.5,
		Bool:   true,

		Strings: []string{"a", "b", "c"},
		Ints:    []int{10, 20, 30},
		Floats:  []float64{1.2, 1.5, 2.5},
		Bools:   []bool{true, true, false},
		Configuration: Config{
			Users: []User{
				User{
					"admin",
					"1q2w3e4r",
				},
			},
		},
	}
	wid := cp.AddWidget(time.Second, `<div>{{.}}</div>`, func() (interface{}, error) {
		return time.Now().Format("15:04:05"), nil
	})

	fid := cp.AddForm(s, func(data interface{}) bool {
		log.Println("HHH", data)
		return true
	})

	cp.AddElementToHeader(wid)

	cp.AddContentPage("/", "Home",
		`<h1>Hello</h1>
		<div>
			{{element .widget}}
			{{element .form}}
		</div>`, func() (interface{}, error) {
			return map[string]string{"widget": wid, "form": fid}, nil
		})

	cp.AddContentPage("/home", "Home", "<h1>Hello {{.}}</h1>", "World")

	cp.Run()

}
