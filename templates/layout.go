package templates

import "../assets"

var Layout = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1, shrink-to-fit=no"
  />
  <script>
    window.onload = function(e){
      {{range .widgetFunctions}}
        {{.}}
      {{end}}
    }
  </script>
    <style>` + assets.BootstrapCss + `</style>
    <title>{{.title}}</title>
  </head>
  <body>
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark navbar-fixed-top">
      <a href="{{.homeUrl}}" class="navbar-brand">{{.title}}</a>
      <ul class="navbar-nav flex-row ml-md-auto d-none d-md-flex" style="color:white">
        {{range .headerWidgets}}
          <li class="nav-item">
            {{widget .}}
          </li>
        {{end}}
      </ul>
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
</html>`
