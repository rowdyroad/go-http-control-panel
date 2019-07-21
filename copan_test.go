package copan

import (
	"fmt"
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

// type Settings struct {
// 	Struct

// 	StructOne Struct
// 	Structs   []Struct

// 	Textarea    string `htmlForm:"type: textarea; rows: 10"`
// 	CustomLabel string `htmlForm:"label: Custom; description: Custom Description"`

// 	ReadonlyString string  `htmlForm:"readonly: true"`
// 	ReadonlyInt    int     `htmlForm:"readonly: true"`
// 	ReadonlyFloat  float64 `htmlForm:"readonly: true"`
// 	ReadonlyBool   bool    `htmlForm:"readonly: true"`

// 	ReadonlyStrings []string  `htmlForm:"readonly: true"`
// 	ReadonlyInts    []int     `htmlForm:"readonly: true"`
// 	ReadonlyFloats  []float64 `htmlForm:"readonly: true"`
// 	ReadonlyBools   []bool    `htmlForm:"readonly: true"`

// 	ReadonlyStruct  Struct   `htmlForm:"readonly: true"`
// 	ReadonlyStructs []Struct `htmlForm:"itemLabel: Struct"`

// 	ReadonlyTextarea    string `htmlForm:"type: textarea; rows: 20; readonly: true"`
// 	ReadonlyCustomLabel string `htmlForm:"label: Custom; description: Custom Description; readonly: true"`
// }

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
	wid := cp.AddWidget(time.Second, func() (string, error) {
		return fmt.Sprintf("<div>%v</div>", time.Now().Format("15:04:05")), nil
	})
	cp.AddWidgetToHeader(wid)

	cp.AddContentPage("/", "Home", "<h1>Hello {{widget .wid}}</h1>", func() interface{} {
		return map[string]string{
			"wid": wid,
		}
	})
	cp.AddFormPage("/settings", "Settings", s, func(data interface{}) bool {
		log.Println("DATA:", data)
		return true
	})

	cp.Run()

}
