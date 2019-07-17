package controlPanel

import "html/template"

var formArrayHeader = template.Must(template.New("form/arrayHeader").Parse(`
	{{if .Label}}
		<h4>
			{{.Label}}
			{{if .Description}}<small>{{.Description}}</small>{{end}}
		</h4>
	{{end}}
	<div id="array-{{.Name}}">
`))

var formArrayFooter = template.Must(template.New("form/arrayFooter").Parse(`
	</div>
	{{if not .Readonly}}
		<div style="margin:0.4em;margin-bottom:1em">
			<input type="button" class="btn btn-secondary" value="{{if .AddBtnCaption}}{{.AddBtnCaption}}{{else}}Add{{end}}" onclick="addArrayItem('{{.Name}}', {{.Length}})"/>
		</div>
	{{end}}
`))
var formArrayItemWrapperHeader = template.Must(template.New("form/arrayItemWrapperHeader").Parse(`
	<div id="item-{{.Name}}" class="form-group" style="padding-left: {{.Indent}}em">
`))

var formStructHeader = template.Must(template.New("form/structHeader").Parse(`
	{{if .IsArrayItem}}
		{{if .ItemLabel}}
			<h4>
				{{.ItemLabel}}
				{{if .Description}}<small>{{.Description}}</small>{{end}}
			</h4>
		{{end}}
	{{else}}
		{{if .Label}}
			<h4>
				{{.Label}}
				{{if .Description}}<small>{{.Description}}</small>{{end}}
			</h4>
		{{end}}
	{{end}}
`))

var formArrayItemWrapperFooter = template.Must(template.New("form/arrayItemWrapperFooter").Parse(`
	{{if not .Readonly}}
		<div style="text-align:right; padding:0.4em 0">
			<input type="button" class="btn btn-danger" onclick="javascript:document.getElementById('item-{{.Name}}').remove()" value="{{if .DeleteBtnCaption}}{{.DeleteBtnCaption}}{{else}}Delete{{end}}"/>
		</div>
	{{end}}
	</div>
`))

var formInput = template.Must(template.New("form/input").Parse(`
	{{if not .IsArrayItem }}
		<div class="form-group" style="padding-left: {{.Indent}}em">
		{{if .Label}} <label class="form-check-label" for="{{.ID}}">{{.Label}}</label>{{end}}
	{{end}}
		<input type="{{.Type}}"  name="{{.Name}}" {{if .Disabled}}disabled{{end}} {{if .Readonly}}readonly{{end}} class="form-control" id="{{.ID}}" value="{{.Value}}" placeholder="{{.Placeholder}}"/>
	{{if not .IsArrayItem }}
		</div>
	{{end}}
`))

var formTextarea = template.Must(template.New("form/textarea").Parse(`
	{{if not .IsArrayItem }}
		<div class="form-group" style="padding-left: {{.Indent}}em">
		{{if .Label}} <label class="form-check-label" for="{{.ID}}">{{.Label}}</label>{{end}}
	{{end}}
		<textarea class="form-control" name="{{.Name}}" {{if .Disabled}}disabled{{end}} {{if .Readonly}}readonly{{end}} placeholder="{{.Placeholder}}" id="{{.ID}}" rows="{{.Rows}}">{{.Value}}</textarea>
	{{if not .IsArrayItem }}
		</div>
	{{end}}


	`))

var formCheckbox = template.Must(template.New("form/checkbox").Parse(`
	<div class="form-check">
		<input class="form-check-input" type="{{.Type}}" name="{{.Name}}" {{if .Disabled}}disabled{{end}} {{if .Readonly}}readonly{{end}} {{if .Value}}checked{{end}}  class="form-control" id="{{.ID}}" value="" placeholder="{{.Placeholder}}"/>
		{{if not .IsArrayItem }} {{if .Label}} <label class="form-check-label" for="{{.ID}}">{{.Label}}</label>{{end}} {{end}}
	</div>
`))

var formSelect = template.Must(template.New("form/select").Parse(`
	{{if not .IsArrayItem }}
		<div class="form-group" style="padding-left: {{.Indent}}em">
		{{if .Label}} <label class="form-check-label" for="{{.ID}}">{{.Label}}</label>{{end}}
	{{end}}
	<select class="form-control" name="{{.Name}}" id="{{.ID}}" {{if .Disabled}}disabled{{end}} {{if .Readonly}}readonly{{end}} >
		{{range $ox, $ob := .Options}}
			<option {{if eq .Value $ox}}selected{{end}} value={{$ox}}>{{$ob}}</option>
		{{end}}
	</select>
	{{if not .IsArrayItem }}
		</div>
	{{end}}
`))
