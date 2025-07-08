package com;
// pxnMetrics Frontend - about page

import(
	HTTP   "net/http"
	PxnWeb "github.com/PoiXson/pxnGoCommon/net/web"
);

import _ "embed";



//go:embed page-about.tpl
var TPL_About []byte;



type PageAbout struct {
	Pages   *Pages
	Builder *PxnWeb.Builder
}



func (pages *Pages) NewPageAbout() *PageAbout {
	builder := pages.Builder.Clone().
		AddRawTPL(TPL_About);
	return &PageAbout{
		Pages:   pages,
		Builder: builder,
	};
}



func (page *PageAbout) RenderWeb(out HTTP.ResponseWriter, in *HTTP.Request) {
	out.Header().Set("Content-Type", PxnWeb.Mime_HTML);
	tags := page.Builder.CloneTags();
	tags["Page" ] = "about";
	tags["Title"] = "About pxnMetrics"
	page.Builder.TPL.ExecuteTemplate(out, "website", tags);
}
