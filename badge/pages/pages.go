package pages;
// badge.minecraftmetrics.com

import(
	Gorilla "github.com/gorilla/mux"
	PxnWeb  "github.com/PoiXson/pxnGoCommon/net/web"
	WebLink "github.com/PoiXson/pxnMetrics/frontend/weblink"
);



type Pages struct {
	Link *WebLink.WebLink
}



func New(weblink *WebLink.WebLink, router *Gorilla.Router) *Pages {
	pages := Pages{
		Link: weblink,
	};
	PxnWeb.AddStaticRoute(router);
	router.HandleFunc("/", pages.Page_Badge);
	return &pages;
}
