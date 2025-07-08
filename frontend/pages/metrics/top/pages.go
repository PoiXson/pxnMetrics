package top;
// minecraftmetrics.top

import(
	Gorilla  "github.com/gorilla/mux"
	PxnWeb   "github.com/PoiXson/pxnGoCommon/net/web"
	WebLink  "github.com/PoiXson/pxnMetrics/frontend/weblink"
	PagesCom "github.com/PoiXson/pxnMetrics/frontend/pages/metrics/com"
);



type Pages struct {
	Link    *WebLink.WebLink
	Builder *PxnWeb.Builder
	PageTop *PageTop
}



func New(weblink *WebLink.WebLink, router *Gorilla.Router) *Pages {
	pages := Pages{
		Link: weblink,
		Builder: PxnWeb.NewBuilder().
			AddRawTPL(PagesCom.TPL_MenuTop),
	};
	pages.PageTop = pages.NewPageTop();
	PxnWeb.AddStaticRoute(router);
	router.HandleFunc("/", pages.PageTop.RenderWeb);
	return &pages;
}
