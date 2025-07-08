package com;
// pxnMetrics Frontend - status page

import(
	HTTP   "net/http"
	PxnWeb "github.com/PoiXson/pxnGoCommon/net/web"
);

import _ "embed";



//go:embed page-status.tpl
var TPL_Status []byte;



type PageStatus struct {
	Pages   *Pages
	Builder *PxnWeb.Builder
}



func (pages *Pages) NewPageStatus() *PageStatus {
	builder := pages.Builder.Clone().
		AddRawTPL(TPL_Status);
	return &PageStatus{
		Pages:   pages,
		Builder: builder,
	};
}



func (page *PageStatus) RenderWeb(out HTTP.ResponseWriter, in *HTTP.Request) {
	out.Header().Set("Content-Type", PxnWeb.Mime_HTML);
	tags := page.Builder.CloneTags();
	tags["Page" ] = "status";
	tags["Title"] = "pxnMetrics Status"
	page.Builder.TPL.ExecuteTemplate(out, "website", tags);
}



func (page *PageStatus) RenderAPI(out HTTP.ResponseWriter, in *HTTP.Request) {
	out.Header().Set("Content-Type", PxnWeb.Mime_JSON);
//	shards := make([]API.JSON_StatusShard, heart.NumShards);
//	for index:=uint8(0); index<heart.NumShards; index++ {




//	}


}
