package copan

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/rowdyroad/go-http-control-panel/templates"

	"github.com/gin-gonic/gin"
	forms "github.com/rowdyroad/go-web-forms"
)

type User struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" htmlForm:"type: password"`
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
	router          *gin.Engine
	menu            []menuItem
	widgetFunctions map[string]interface{}
	headerWidgets   []string
	templateFuncs   template.FuncMap
	server          *http.Server
}

func NewControlPanel(config Config) *ControlPanel {
	router := gin.Default()

	cp := &ControlPanel{
		config,
		router,
		[]menuItem{},
		map[string]interface{}{},
		[]string{},
		template.FuncMap{},
		&http.Server{
			Addr:    config.Listen,
			Handler: router,
		},
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
		cp.router.Use(gin.BasicAuth(accounts))
	}

	layout := template.Must(template.New("layout").Funcs(cp.templateFuncs).Parse(templates.Layout))
	cp.router.SetHTMLTemplate(layout)
	return cp
}

func (cc *ControlPanel) AddWidget(refresh time.Duration, tmpl string, content func() (interface{}, error)) string {
	wid := strings.Replace(uuid.New().String(), "-", "", -1)

	widgetTemplate := template.Must(template.New(wid).Funcs(cc.templateFuncs).Parse(tmpl))

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
	cc.router.GET(wid, func(c *gin.Context) {
		var data interface{}
		if content != nil {
			var err error
			data, err = content()
			if err != nil {
				c.Data(500, "text/html; charset=utf-8", []byte(fmt.Sprintf(`<div class="alert alert-danger" role="alert">Error: %s</div>`, err.Error())))
				return
			}
		}
		out := bytes.Buffer{}
		widgetTemplate.Execute(&out, data)
		c.Data(200, "text/html; charset=utf-8", []byte(out.String()))
	})
	return wid
}

func (cc *ControlPanel) AddWidgetToHeader(wid string) {
	cc.headerWidgets = append(cc.headerWidgets, wid)
}

func (cc *ControlPanel) AddContentPage(url string, menu string, tmpl string, content func() (interface{}, error)) {
	if menu != "" {
		cc.menu = append(cc.menu, menuItem{url, menu})
	}
	pageTemplate := template.Must(template.New(url).Funcs(cc.templateFuncs).Parse(tmpl))

	cc.router.GET(url, func(c *gin.Context) {
		var data interface{}
		if content != nil {
			var err error
			data, err = content()
			if err != nil {
				c.Data(500, "text/html; charset=utf-8", []byte(fmt.Sprintf(`<div class="alert alert-danger" role="alert">Error: %s</div>`, err.Error())))
				return
			}
		}
		out := bytes.Buffer{}
		pageTemplate.Execute(&out, data)
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
		cc.render(c, template.HTML(x.String()))
	}
}

func (cc *ControlPanel) AddFormPage(url string, menu string, data interface{}, cb func(data interface{}) bool) {
	if menu != "" {
		cc.menu = append(cc.menu, menuItem{url, menu})
	}

	cc.router.GET(url, func(c *gin.Context) {
		cc.showForm(url, menu, data, c)
	})

	cc.router.POST(url, func(c *gin.Context) {
		c.Request.ParseForm()
		ret := forms.ParseForm(c.Request.Form, data)
		if cb(ret) {
			data = ret
		}
		cc.showForm(url, menu, data, c)
	})
}

func (cc *ControlPanel) Run() {
	if err := cc.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func (cc *ControlPanel) Stop() {
	cc.server.Shutdown(context.Background())
}
