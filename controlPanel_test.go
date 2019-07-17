package controlPanel

import (
	"log"
	"testing"
)

type File struct {
	Name string
	Path string
}
type Settings struct {
	Number []int  `htmlForm:"name: nnnumber\nreadonly: true\nlabel: Hello ho ho ho\ntype: number"`
	ZZZ    []int  `htmlForm:"name: zzzz"`
	String string `htmlForm:"label: string"`
	Text   string `htmlForm:"textarea"`
	Check  bool   `htmlForm:"label: checking go"`
	File   []File `htmlForm:"label: Files\nitemLabel: File"`
}

func TestMain(t *testing.T) {
	cp := NewControlPanel(Config{
		Listen: ":9999",
		Title:  "config",
	})

	s := Settings{
		[]int{100, 101, 102, 103},
		[]int{100, 101, 102, 103},
		"hello",
		"text",
		true,
		[]File{
			File{
				"test",
				"/etc/test",
			},
			File{
				"test1",
				"/etc/test1",
			},
		},
	}
	cp.AddFormPage("/settings", "Settings", s, func(data interface{}) bool {
		log.Println("DATA:", data)
		return true
	})

	cp.Run()

}
