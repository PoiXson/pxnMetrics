package com;
// minecraftmetrics.com

import(
	Gorilla "github.com/gorilla/mux"
	PxnWeb  "github.com/PoiXson/pxnGoCommon/net/web"
	WebLink "github.com/PoiXson/pxnMetrics/frontend/weblink"
);

import _ "embed";



//go:embed menu-top.tpl
var TPL_MenuTop []byte;



type Pages struct {
	Link       *WebLink.WebLink
	Builder    *PxnWeb.Builder
	PageStatus *PageStatus
	PageGlobal *PageGlobal
	PageWiki   *PageWiki
	PageAbout  *PageAbout
}



func New(weblink *WebLink.WebLink, router *Gorilla.Router) *Pages {
	builder := PxnWeb.NewBuilder().
		WithIncludes().
		WithBootstrap().
		WithBootsIcons().
		WithTooltips().
		AddRawTPL(TPL_MenuTop).
		SetTitle("pxnMetrics").
		AddFileCSS("/static/metrics.css").
		SetFavIcon("/static/line-chart.ico");
	pages := Pages{
		Link:    weblink,
		Builder: builder,
	};
	pages.PageStatus = pages.NewPageStatus();
	pages.PageGlobal = pages.NewPageGlobal();
	pages.PageWiki   = pages.NewPageWiki();
	pages.PageAbout  = pages.NewPageAbout();
	PxnWeb.AddStaticRoute(router);
	router.HandleFunc("/",            pages.PageGlobal.RenderWeb);
	router.HandleFunc("/status/",     pages.PageStatus.RenderWeb);
	router.HandleFunc("/api/status/", pages.PageStatus.RenderAPI);
	router.HandleFunc("/wiki/",       pages.PageWiki  .RenderWeb);
	router.HandleFunc("/about/",      pages.PageAbout .RenderWeb);
	return &pages;
}
