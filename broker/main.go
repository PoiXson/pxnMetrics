package main;

import App "github.com/PoiXson/pxnMetrics/broker/app";

const Version = "{{{VERSION}}}";

func main() {
	app := App.New();
	app.Main();
}
