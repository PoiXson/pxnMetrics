package com;
// pxnMetrics Frontend - status page

import(
	Fmt     "fmt"
	HTTP    "net/http"
	Context "context"
	Runtime "runtime"
	TPL     "html/template"
	GEmpty  "google.golang.org/protobuf/types/known/emptypb"
	HTML    "github.com/PoiXson/pxnGoCommon/utils/html"
);
//TODO
//	GRPC     "google.golang.org/grpc"
//	GZIP     "google.golang.org/grpc/encoding/gzip"
//	FrontAPI "github.com/PoiXson/pxnMetrics/api/front"



func (pages *Pages) PageInit_Status() {
	pages.tpl_status = TPL.Must(TPL.ParseFiles(
		"html/main.tpl",
		"html/pages/status.tpl",
	));
}

func (pages *Pages) PageWeb_Status(out HTTP.ResponseWriter, in *HTTP.Request) {
	HTML.SetContentType(out, "html");
	build := pages.GetBuilder().
		AddBotJS("/static/status.js");
	vars := struct {
		Page  string
		Title string
	}{
		Page:  "status",
		Title: "title",
	};
	out.Write([]byte(build.RenderTop()));
	pages.tpl_status.ExecuteTemplate(out, "main.tpl",   vars);
	pages.tpl_status.ExecuteTemplate(out, "status.tpl", vars);
	out.Write([]byte(build.RenderBottom()));
}



func (pages *Pages) PageAPI_Status(out HTTP.ResponseWriter, in *HTTP.Request) {
	reply, err := pages.link.API.FetchStatusJSON(
		Context.Background(),
		&GEmpty.Empty{},
	);
//TODO: make this into a function?
	if err != nil {
		trace := make([]byte, 1024);
		n := Runtime.Stack(trace, true);
		HTTP.Error(out,
			Fmt.Sprintf("%s\n%s", err.Error(), trace[:n]),
			HTTP.StatusInternalServerError,
		);
		return;
	}
//TODO: optional? only when not unix socket?
//		GRPC.UseCompressor(GZIP.Name),
	HTML.SetContentType(out, "json");
	out.WriteHeader(HTTP.StatusOK);
	out.Write(reply.Data);
}
