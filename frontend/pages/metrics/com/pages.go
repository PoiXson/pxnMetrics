package com;
// minecraftmetrics.com

import(
	TPL     "html/template"
	Gorilla "github.com/gorilla/mux"
	HTML    "github.com/PoiXson/pxnGoCommon/utils/html"
	PxnNet  "github.com/PoiXson/pxnGoCommon/utils/net"
	WebLink "github.com/PoiXson/pxnMetrics/frontend/weblink"
);



type Pages struct {
	link *WebLink.WebLink
	tpl_status *TPL.Template
}



func New(weblink *WebLink.WebLink, router *Gorilla.Router) *Pages {
	pages := Pages{
		link: weblink,
	};
	PxnNet.AddStaticRoute(router);
	router.HandleFunc("/",            pages.PageWeb_Global);
	router.HandleFunc("/wiki/",       pages.PageWeb_Wiki  );
	router.HandleFunc("/status/",     pages.PageWeb_Status);
	router.HandleFunc("/api/status/", pages.PageAPI_Status);
	router.HandleFunc("/about/",      pages.PageWeb_About );
	router.HandleFunc("/favicon.ico", PxnNet.NewRedirect("/static/line-chart.ico"));
	pages.PageInit_Status();
	return &pages;
}



func (pages *Pages) GetBuilder() *HTML.Builder {
	return HTML.NewBuilder().
		WithBootstrap().
		WithBootstrapIcons().
		WithBootstrapTooltips().
		SetFavIcon("/static/line-chart.ico").
		AddCSS("/static/metrics.css").
		SetTitle("pxnMetrics");
}
