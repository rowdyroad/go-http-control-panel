package controlPanel

import "html/template"
import "./assets"

var layout = template.Must(template.New("layout").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1, shrink-to-fit=no"
	/>
    <style>` + assets.BootstrapCss + `</style>
    <script>
      function uuidv4() {
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
          var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
          return v.toString(16);
        });
      }

      var indexes = {};


      function addArrayItem(id, indexMax) {
        if (!indexes[id]) {
          indexes[id] = indexMax
        }
        var index = indexes[id];



        var cnt = document.getElementById('new-' + id)
                            .innerHTML
                            .replace(new RegExp('id-'+id+'-new', 'g'), uuidv4())
                            .replace(new RegExp('name-'+id+'-new', 'g'), id+'['+indexes[id]+']')
                            .replace(new RegExp('display:none;', 'g'), '');

        var e = document.createElement('div');
        e.innerHTML = cnt;

        document.getElementById('array-'+id).appendChild(e);
        indexes[id] = index + 1;
      }
    </script>
    <title>{{.title}}</title>
  </head>
  <body>
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark navbar-fixed-top">
      <a href="{{.homeUrl}}" class="navbar-brand">{{.title}}</a>
    </nav>
    <main class="container-fluid" style="margin-top:1em">
    <div class="row">
      <div class="col-sm-3 col-md-2">
          <div class="nav flex-column nav-pills" role="tablist" aria-orientation="vertical">
            {{$x := .}}
            {{range .menu}}
              <a class="nav-item nav-link {{if eq .URL $x.location}}active{{end}}" href="{{.URL}}">{{.Title}}</a>
            {{end}}
          </div>
      </div>
      <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2">
      {{.content}}
      </div>
    </div>
    </main>
 </body>
</html>`))
