package controlPanel

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Title  string `yaml:"title" json:"title"`
	Listen string `yaml:"listen" json:"listen"`
}

type menuItem struct {
	URL   string
	Title string
}

type formFieldArray struct {
	Name   string
	Index  int
	Length int
	Last   bool
}

type formField struct {
	ID               string      `yaml:"-"`
	Name             string      `yaml:"name"`
	Description      string      `yaml:"description"`
	Type             string      `yaml:"type"`
	Label            string      `yaml:"label"`
	ItemLabel        string      `yaml:"itemLabel"`
	DeleteBtnCaption string      `yaml:"deleteBtnCaption"`
	AddBtnCaption    string      `yaml:"addBtnCaption"`
	Value            interface{} `yaml:"-"`
	Disabled         bool        `yaml:"disabled"`
	Rows             int         `yaml:"rows"`
	Placeholder      string      `yaml:"placeholder"`
	Indent           int         `yaml:"-"`
	Readonly         bool        `yaml:"readonly"`
	Skip             bool        `yaml:"-"`
	IsArrayItem      bool        `yaml:"-"`
	Length           int         `yaml:"-"`
}

func (t formField) Copy() formField {
	return formField{
		Name:             t.Name,
		Description:      t.Description,
		Type:             t.Type,
		Label:            t.Label,
		ItemLabel:        t.ItemLabel,
		DeleteBtnCaption: t.DeleteBtnCaption,
		AddBtnCaption:    t.AddBtnCaption,
		Disabled:         t.Disabled,
		Rows:             t.Rows,
		Placeholder:      t.Placeholder,
		Readonly:         t.Readonly,
	}
}

type ControlPanel struct {
	config Config
	route  *gin.Engine
	menu   []menuItem
}

func NewControlPanel(config Config) *ControlPanel {
	route := gin.Default()
	route.SetHTMLTemplate(layout)
	return &ControlPanel{
		config,
		route,
		[]menuItem{},
	}
}

func (cc *ControlPanel) AddContentPage(url string, menu string, content func() string) {
	if menu != "" {
		cc.menu = append(cc.menu, menuItem{url, menu})
	}
	cc.route.GET(url, func(c *gin.Context) {
		c.HTML(200, "layout", gin.H{"menu": cc.menu, "title": cc.config.Title, "content": content(), "location": c.Request.URL.Path})
	})
}

func parseTags(sf reflect.StructField) formField {
	var ff formField
	tag := sf.Tag.Get("htmlForm")
	if tag == "-" {
		ff.Skip = true
		return ff
	}
	yaml.Unmarshal([]byte(tag), &ff)
	ff.ID = uuid.New().String()
	if ff.Name == "" {
		ff.Name = sf.Name
	}
	if ff.Label == "" {
		ff.Label = sf.Name
	}
	return ff
}

func processField(x io.Writer, f reflect.Value, ff *formField) {
	tmpl := formInput
	switch f.Type().Kind() {
	case reflect.Array, reflect.Slice:
		if ff != nil {
			formArrayHeader.Execute(x, ff)
			if !ff.Disabled && !ff.Readonly {
				x.Write([]byte(fmt.Sprintf(`<script type="x-template" id="new-%s">`, ff.Name)))
				fieldTemplate := ff.Copy()
				fieldTemplate.ID = fmt.Sprintf("id-%s-new", ff.Name)
				fieldTemplate.Name = fmt.Sprintf("name-%s-new", ff.Name)
				fieldTemplate.IsArrayItem = true
				formArrayItemWrapperHeader.Execute(x, fieldTemplate)
				processField(x, reflect.Zero(reflect.TypeOf(f.Interface()).Elem()), &fieldTemplate)
				formArrayItemWrapperFooter.Execute(x, fieldTemplate)
				x.Write([]byte(`</script>`))
			}
		}
		for i := 0; i < f.Len(); i++ {
			field := ff.Copy()
			field.ID = uuid.New().String()
			field.Name = fmt.Sprintf("%s[%d]", ff.Name, i)
			field.Value = f.Index(i).Interface()
			field.IsArrayItem = true

			formArrayItemWrapperHeader.Execute(x, field)
			processField(x, f.Index(i), &field)
			formArrayItemWrapperFooter.Execute(x, field)
		}
		if ff != nil {
			ff.Length = f.Len()
			formArrayFooter.Execute(x, ff)
		}
		return
	case reflect.Bool:
		if ff.Type == "" {
			ff.Type = "checkbox"
		}
		tmpl = formCheckbox
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if ff.Type == "" {
			ff.Type = "number"
		}
	case reflect.Struct:
		indent := 0

		if ff != nil {
			formStructHeader.Execute(x, ff)
			indent += ff.Indent + 1
		}

		for i := 0; i < f.NumField(); i++ {
			fd := parseTags(f.Type().Field(i))
			if fd.Skip {
				continue
			}
			if ff != nil {
				fd.Name = ff.Name + "[" + fd.Name + "]"
			}
			fd.Value = f.Field(i).Interface()
			fd.Indent = indent

			processField(x, f.Field(i), &fd)
		}
		return
	default:
		if ff.Type == "" {
			ff.Type = "text"
		}
	}
	tmpl.Execute(x, ff)
}

func parseForm(value reflect.Value, src reflect.Value, c *gin.Context, name string, form *url.Values) {
	switch value.Type().Kind() {
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			ff := parseTags(value.Type().Field(i))
			if src.IsValid() && (ff.Skip || ff.Disabled || ff.Readonly) {
				value.Field(i).Set(src.Field(i))
				continue
			}
			if name != "" {
				ff.Name = name + "[" + ff.Name + "]"
			}
			if src.IsValid() {
				parseForm(value.Field(i), src.Field(i), c, ff.Name, form)
			} else {
				parseForm(value.Field(i), reflect.Value{}, c, ff.Name, form)
			}
		}
	case reflect.Array, reflect.Slice:
		i := 0
		re := regexp.MustCompile(fmt.Sprintf("%s\\[(\\d+)\\]", name))
		for key := range *form {
			if regs := re.FindStringSubmatch(key); len(regs) == 2 {
				id, _ := strconv.ParseInt(regs[1], 10, 64)
				value.Set(reflect.Append(value, reflect.Zero(value.Type().Elem())))
				if !src.IsValid() || int(id) >= src.Len() {
					parseForm(reflect.Indirect(value.Index(i)), reflect.Value{}, c, regs[0], form)
				} else {
					parseForm(reflect.Indirect(value.Index(i)), src.Index(int(id)), c, regs[0], form)
				}
				i++
				for kk := range *form {
					if strings.Index(kk, fmt.Sprintf("%s[%d]", name, id)) != -1 {
						form.Del(kk)
					}
				}
			}
		}

	case reflect.Bool:
		v, _ := strconv.ParseBool(c.PostForm(name))
		value.SetBool(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, _ := strconv.ParseInt(c.PostForm(name), 10, 64)
		value.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, _ := strconv.ParseUint(c.PostForm(name), 10, 64)
		value.SetUint(v)
	case reflect.String:
		value.SetString(c.PostForm(name))
	}
}

func (cc *ControlPanel) showForm(url string, menu string, data interface{}, c *gin.Context) {
	x := bytes.Buffer{}
	x.WriteString(fmt.Sprintf(`<form method="POST"><h1>%s</h1>`, menu))
	processField(&x, reflect.ValueOf(data), nil)
	x.WriteString(`<button type="submit" class="btn btn-primary">Submit</button></form>`)
	c.HTML(200, "layout", gin.H{"menu": cc.menu, "title": cc.config.Title, "content": template.HTML(x.String()), "location": c.Request.URL.Path})
}

func (cc *ControlPanel) AddFormPage(url string, menu string, data interface{}, cb func(data interface{}) bool) {
	if menu != "" {
		cc.menu = append(cc.menu, menuItem{url, menu})
	}

	cc.route.GET(url, func(c *gin.Context) {
		cc.showForm(url, menu, data, c)
	})

	cc.route.POST(url, func(c *gin.Context) {
		c.Request.ParseForm()
		src := reflect.ValueOf(data)
		nv := reflect.New(reflect.Indirect(src).Type()).Elem()
		parseForm(nv, src, c, "", &c.Request.Form)
		if cb(nv.Interface()) {
			data = nv.Interface()
		}
		cc.showForm(url, menu, data, c)
	})
}

func (cc *ControlPanel) Run() {
	cc.route.Run(cc.config.Listen)
}
