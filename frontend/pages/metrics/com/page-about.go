package com;
// pxnMetrics Frontend - about page

import(
	HTTP "net/http"
	TPL  "html/template"
	HTML "github.com/PoiXson/pxnGoCommon/utils/html"
);



func (pages *Pages) PageWeb_About(out HTTP.ResponseWriter, in *HTTP.Request) {
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
		Page:  "About",
		Title: "title",
	};
	out.Write([]byte(build.RenderTop()));
	tpl.ExecuteTemplate(out, "main.tpl", vars);
//	tpl.ExecuteTemplate(out, "about.tpl", vars);
	out.Write([]byte("About Page goes here"));
	out.Write([]byte(build.RenderBottom()));
}
