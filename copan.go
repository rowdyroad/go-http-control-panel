package copan

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/google/uuid"

	"./templates"

	"github.com/gin-gonic/gin"
	forms "github.com/rowdyroad/go-web-forms"
)

type User struct {
	Username string
	Password string
}

type Config struct {
	Title  string `yaml:"title" json:"title"`
	Listen string `yaml:"listen" json:"listen"`
	Users  []User `yaml:"users" json:"users"`
}

type menuItem struct {
	URL   string
	Title string
}

type ControlPanel struct {
	config          Config
	route           *gin.Engine
	menu            []menuItem
	widgetFunctions map[string]interface{}
	headerWidgets   []string
	templateFuncs   template.FuncMap
}

func NewControlPanel(config Config) *ControlPanel {
	route := gin.Default()

	cp := &ControlPanel{
		config,
		route,
		[]menuItem{},
		map[string]interface{}{},
		[]string{},
		template.FuncMap{},
	}

	cp.templateFuncs["widget"] = func(id string) interface{} {
		if _, ok := cp.widgetFunctions[id]; ok {
			return template.HTML(fmt.Sprintf(`<div class="widget-%s"></div>`, id))
		}
		return ""
	}

	if len(config.Users) > 0 {
		accounts := gin.Accounts{}
		for _, user := range config.Users {
			accounts[user.Username] = user.Password
		}
		cp.route.Use(gin.BasicAuth(accounts))
	}

	layout := template.Must(template.New("layout").Funcs(cp.templateFuncs).Parse(templates.Layout))
	cp.route.SetHTMLTemplate(layout)
	return cp
}

func (cc *ControlPanel) AddWidget(refresh time.Duration, cb func() (string, error)) string {
	wid := strings.Replace(uuid.New().String(), "-", "", -1)

	cc.widgetFunctions[wid] = template.JS(fmt.Sprintf(`function widget%s() {
					fetch('%s')
						.then(function(response) {
							return response.text()
						})
						.then(function(html) {
							var widgets = document.getElementsByClassName('widget-%s');
							for (var i = 0; i < widgets.length; i++) {
								widgets[i].innerHTML = html;
							}
						})
				}
				setInterval(widget%s, %d);
				widget%s();
	`, wid, wid, wid, wid, refresh.Nanoseconds()/1000000, wid))
	cc.route.GET(wid, func(c *gin.Context) {
		if content, err := cb(); err == nil {
			c.Data(200, "text/html; charset=utf-8", []byte(content))
		} else {
			c.Data(500, "text/html; charset=utf-8", []byte(fmt.Sprintf(`<div class="alert alert-danger" role="alert">Error: %s</div>`, err.Error())))
		}
	})
	return wid
}

func (cc *ControlPanel) AddWidgetToHeader(wid string) {
	cc.headerWidgets = append(cc.headerWidgets, wid)
}

func (cc *ControlPanel) AddContentPage(url string, menu string, tmpl string, content func() interface{}) {
	if menu != "" {
		cc.menu = append(cc.menu, menuItem{url, menu})
	}
	pageTemplate := template.Must(template.New(url).Funcs(cc.templateFuncs).Parse(tmpl))

	cc.route.GET(url, func(c *gin.Context) {
		out := bytes.Buffer{}
		pageTemplate.Execute(&out, content())
		cc.render(c, template.HTML(out.String()))
	})
}

func (cc *ControlPanel) render(c *gin.Context, content interface{}) {
	c.HTML(200, "layout", gin.H{
		"menu":            cc.menu,
		"title":           cc.config.Title,
		"content":         content,
		"headerWidgets":   cc.headerWidgets,
		"widgetFunctions": cc.widgetFunctions,
		"location":        c.Request.URL.Path,
	})
}

func (cc *ControlPanel) showForm(url string, menu string, data interface{}, c *gin.Context) {
	x := bytes.Buffer{}
	if forms.MakeHTML(data, &x, nil) {
		cc.render(c, template.HTML("<h1>Settings</h1>"+x.String()))
	}
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
		ret := forms.ParseForm(c.Request.Form, data)
		if cb(ret) {
			data = ret
		}
		cc.showForm(url, menu, data, c)
	})
}

func (cc *ControlPanel) Run() {
	cc.route.Run(cc.config.Listen)
}
