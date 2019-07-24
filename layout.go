package copan

var layout = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1, shrink-to-fit=no"
  />
  <script>
    function loadElement(id) {
      fetch(id)
        .then(function(response) {
          return response.text()
        })
        .then(function(html) {
          var widgets = document.getElementsByClassName('element-' + id);
          for (var i = 0; i < widgets.length; i++) {
            var widget = widgets[i];
            while (widget.firstChild) {
              widget.removeChild(widget.firstChild);
            }
            var el = document.createElement('div');
            el.innerHTML = html;
            widget.appendChild(el);
            var scripts = el.querySelectorAll('script');
            for (var j = 0; j < scripts.length; j++) {
              if (scripts[j].type === '' || scripts[j].type === 'text/javascript') {
                eval(scripts[j].text);
              }
            }
          }
        })
    }
  </script>
    <style>` + bootstrapCSS + `</style>
    <title>{{.title}}</title>
  </head>
  <body>
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark navbar-fixed-top">
      <a href="{{.homeUrl}}" class="navbar-brand">{{.title}}</a>
      <ul class="navbar-nav flex-row ml-md-auto d-none d-md-flex" style="color:white">
        {{range .headerElements}}
          <li class="nav-item">
            {{element .}}
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
