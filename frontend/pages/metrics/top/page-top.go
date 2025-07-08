package top;
// pxnMetrics Frontend - top page

import(
	HTTP "net/http"
);

import _ "embed";



//go:embed page-top.tpl
var TPL_Top []byte;



type PageTop struct {
	Pages *Pages
}



func (pages *Pages) NewPageTop() *PageTop {
	return &PageTop{};
}



func (page *PageTop) RenderWeb(out HTTP.ResponseWriter, in *HTTP.Request) {
	out.Write([]byte("OK"));
}
