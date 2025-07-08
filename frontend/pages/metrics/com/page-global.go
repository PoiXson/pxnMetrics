package com;
// pxnMetrics Frontend - global metrics page

import(
	HTTP   "net/http"
	PxnWeb "github.com/PoiXson/pxnGoCommon/net/web"
);

import _ "embed";



//go:embed page-global.tpl
var TPL_Global []byte;



type PageGlobal struct {
	Pages   *Pages
	Builder *PxnWeb.Builder
}



func (pages *Pages) NewPageGlobal() *PageGlobal {
	builder := pages.Builder.Clone().
		AddRawTPL(TPL_Global);
	return &PageGlobal{
		Pages:   pages,
		Builder: builder,
	};
}



func (page *PageGlobal) RenderWeb(out HTTP.ResponseWriter, in *HTTP.Request) {
	out.Header().Set("Content-Type", PxnWeb.Mime_HTML);
	tags := page.Builder.CloneTags();
	tags["Page" ] = "global";
	tags["Title"] = "Global Metrics"
	page.Builder.TPL.ExecuteTemplate(out, "website", tags);
}



func (page *PageGlobal) RenderAPI(out HTTP.ResponseWriter, in *HTTP.Request) {
	out.Header().Set("Content-Type", PxnWeb.Mime_JSON);
	out.Write([]byte("{}"));
}
