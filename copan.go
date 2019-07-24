package copan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"

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
	config         Config
	router         *gin.Engine
	menu           []menuItem
	elements       map[string]string
	headerElements []string
	templateFuncs  template.FuncMap
	server         *http.Server
}

func NewControlPanel(config Config) *ControlPanel {
	router := gin.Default()

	cp := &ControlPanel{
		config,
		router,
		[]menuItem{},
		map[string]string{},
		[]string{},
		template.FuncMap{},
		&http.Server{
			Addr:    config.Listen,
			Handler: router,
		},
	}

	cp.templateFuncs["element"] = func(id string) interface{} {
		if data, ok := cp.elements[id]; ok {
			return template.HTML(data)
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

	layout := template.Must(template.New("layout").Funcs(cp.templateFuncs).Parse(layout))
	cp.router.SetHTMLTemplate(layout)
	return cp
}

func (cc *ControlPanel) AddWidget(refresh time.Duration, tmpl string, content func() (interface{}, error)) string {
	wid := strings.Replace(uuid.New().String(), "-", "", -1)

	widgetTemplate := template.Must(template.New(wid).Funcs(cc.templateFuncs).Parse(tmpl))

	cc.elements[wid] = fmt.Sprintf(`<div class="element-%s"></div>
		<script>
			if (!window.elements) {
				window.elements = {};
			}
			if (!window.elements['%s']) {
				window.elements['%s'] = true;
				document.addEventListener("DOMContentLoaded", function() {
					setInterval(function() {
						loadElement('%s');
					}, %d);
					loadElement('%s');
				});
			}
		</script>`, wid, wid, wid, wid, refresh.Nanoseconds()/1000000, wid)

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

func (cc *ControlPanel) AddForm(data interface{}, callback func(data interface{}) bool) string {
	wid := strings.Replace(uuid.New().String(), "-", "", -1)
	cc.elements[wid] = fmt.Sprintf(`<div class="element-%s"></div>
		<script>
			document.addEventListener("DOMContentLoaded", function() {
				loadElement('%s');
			});
		</script>`, wid, wid)

	cc.router.GET(wid, func(c *gin.Context) {
		log.Println("ret", data)
		forms.MakeHTML(wid, data, c.Writer)
	})

	cc.router.POST(wid, func(c *gin.Context) {
		srcBts, _ := json.Marshal(data)
		dst := reflect.New(reflect.TypeOf(data)).Interface()
		json.Unmarshal(srcBts, dst)
		if err := c.Bind(dst); err != nil {
			log.Println("Error:", err)
			return
		}
		dst = reflect.ValueOf(dst).Elem().Interface()
		if callback(dst) {
			data = dst
		}
	})

	return wid
}

func (cc *ControlPanel) AddElementToHeader(wid string) {
	cc.headerElements = append(cc.headerElements, wid)
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
		"menu":           cc.menu,
		"title":          cc.config.Title,
		"content":        content,
		"headerElements": cc.headerElements,
		"elements":       cc.elements,
		"location":       c.Request.URL.Path,
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
