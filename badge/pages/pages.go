package pages;
// badge.minecraftmetrics.com

import(
	TPL     "html/template"
	Gorilla "github.com/gorilla/mux"
	HTML    "github.com/PoiXson/pxnGoCommon/utils/html"
	PxnNet  "github.com/PoiXson/pxnGoCommon/utils/net"
	WebLink "github.com/PoiXson/pxnMetrics/frontend/weblink"
);



type Pages struct {
	tpl_home *TPL.Template
	link     *WebLink.WebLink
}



func New(weblink *WebLink.WebLink, router *Gorilla.Router) *Pages {
	pages := Pages{
		link: weblink,
	};
	PxnNet.AddStaticRoute(router);
	router.HandleFunc("/", pages.Page_Badge);
	router.HandleFunc("/favicon.ico",
		PxnNet.NewRedirect("/static/line-chart.ico"));
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
