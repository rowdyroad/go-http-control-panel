package widgets

import (
	"errors"
	"sync"
	"time"

	copan "github.com/rowdyroad/go-http-control-panel"
)

var dashboard = `<div class="row">
						{{range .}}
							<div class="col-md-4">
								<div class="card mb-4 text-center shadow-sm">
									<div class="card-body">
										<h5 class="card-title">{{.Name}}</h5>
										<p class="card-text" style="padding:2em">{{.Value}}</p>
									</div>
								</div>
							</div>
						{{end}}
					</div>`

type value struct {
	Name  string
	Value interface{}
}
type Dashboard struct {
	sync.Mutex
	copan  *copan.ControlPanel
	values []value
	keys   map[string]int
	wid    string
}

func NewDashboard(copan *copan.ControlPanel, refresh time.Duration, pairs ...interface{}) *Dashboard {
	db := &Dashboard{
		copan:  copan,
		values: []value{},
		keys:   map[string]int{},
	}
	db.wid = copan.AddWidget(refresh, dashboard, func() (interface{}, error) {
		db.Lock()
		defer db.Unlock()
		return db.values, nil
	})
	if len(pairs) > 0 {
		db.Set(pairs)
	}

	return db
}

func (cp *Dashboard) Set(pairs ...interface{}) *Dashboard {
	cp.Lock()
	defer cp.Unlock()
	if len(pairs)%2 != 0 {
		panic(errors.New("Not even pairs"))
	}

	for i := 0; i < len(pairs); i += 2 {
		key := pairs[i].(string)
		if idx, has := cp.keys[key]; has {
			cp.values[idx].Value = pairs[i+1]
		} else {
			cp.values = append(cp.values, value{key, pairs[i+1]})
			cp.keys[key] = len(cp.values) - 1
		}
	}
	return cp
}

func (cp *Dashboard) Unset(names ...string) *Dashboard {
	cp.Lock()
	defer cp.Unlock()
	for _, name := range names {
		if idx, has := cp.keys[name]; has {
			cp.values = append(cp.values[:idx], cp.values[idx+1:]...)
			for key, i := range cp.keys {
				if i > idx {
					cp.keys[key] = i - 1
				}
			}
		}
	}
	return cp
}

func (cp *Dashboard) Clear() *Dashboard {
	cp.Lock()
	defer cp.Unlock()
	cp.values = []value{}
	return cp
}

func (cp *Dashboard) ID() string {
	return cp.wid
}
