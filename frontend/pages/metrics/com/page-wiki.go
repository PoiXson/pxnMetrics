package com;
// pxnMetrics Frontend - wiki page

import(
	HTTP   "net/http"
	PxnWeb "github.com/PoiXson/pxnGoCommon/net/web"
);

import _ "embed";



//go:embed page-wiki.tpl
var TPL_Wiki []byte;



type PageWiki struct {
	Pages   *Pages
	Builder *PxnWeb.Builder
}



func (pages *Pages) NewPageWiki() *PageWiki {
	builder := pages.Builder.Clone().
		AddRawTPL(TPL_Wiki);
	return &PageWiki{
		Pages:   pages,
		Builder: builder,
	};
}



func (page *PageWiki) RenderWeb(out HTTP.ResponseWriter, in *HTTP.Request) {
	out.Header().Set("Content-Type", PxnWeb.Mime_HTML);
	tags := page.Builder.CloneTags();
	tags["Page" ] = "wiki";
	tags["Title"] = "pxnMetrics Wiki"
	page.Builder.TPL.ExecuteTemplate(out, "website", tags);
}
