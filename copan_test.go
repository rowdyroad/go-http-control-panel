package copan

import (
	"log"
	"testing"
	"time"
)

type Struct struct {
	String  string
	Int     int
	Float   float64
	Bool    bool
	Strings []string
	Bools   []bool
	Ints    []int
	Floats  []float64
}

func TestMain(t *testing.T) {
	cp := NewControlPanel(Config{
		Listen: ":9999",
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
	}
	wid := cp.AddWidget(time.Second, `<div>{{.}}</div>`, func() (interface{}, error) {
		return time.Now().Format("15:04:05"), nil
	})

	fid := cp.AddForm(s, func(data interface{}) bool {
		log.Println("HHH", data)
		return true
	})

	cp.AddElementToHeader(wid)

	cp.AddContentPage("/", "Home", "<h1>Hello {{element .widget}} {{element .form}}</h1>", func() (interface{}, error) {
		return map[string]string{"widget": wid, "form": fid}, nil
	})

	cp.Run()

}
