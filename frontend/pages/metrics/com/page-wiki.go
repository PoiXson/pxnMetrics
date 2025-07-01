package com;
// pxnMetrics Frontend - wiki page

import(
	HTTP "net/http"
	TPL  "html/template"
	HTML "github.com/PoiXson/pxnGoCommon/utils/html"
);



func (pages *Pages) PageWeb_Wiki(out HTTP.ResponseWriter, in *HTTP.Request) {
	HTML.SetContentType(out, "html");
	build := pages.GetBuilder();
//TODO
build.IsDev = true;
	tpl, err := TPL.ParseFiles(
		"html/main.tpl",
	);
	if err != nil { panic(err); }
	vars := struct {
		Page  string
		Title string
	}{
		Page:  "wiki",
		Title: "title",
	};
	out.Write([]byte(build.RenderTop()));
	tpl.ExecuteTemplate(out, "main.tpl", vars);
//TODO
//	tpl.ExecuteTemplate(out, "wiki.tpl", vars);
	out.Write([]byte("Wiki goes here"));
	out.Write([]byte(build.RenderBottom()));
}
