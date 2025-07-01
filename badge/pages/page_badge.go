package pages;
// badge.minecraftmetrics.com

import(
	HTTP "net/http"
);



func (pages *Pages) Page_Badge(out HTTP.ResponseWriter, in *HTTP.Request) {
	out.Write([]byte("OK"));
}
