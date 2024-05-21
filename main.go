package main

import (
	"log"
)

func serve(srv *ProxySrv) {
	log.Fatal(srv.ListenAndServe())
}

func main() {
	canSrv := NewProxyServer("cfgs/can.yaml", can)
	gbrSrv := NewProxyServer("cfgs/gbr.yaml", gbr)
	indSrv := NewProxyServer("cfgs/ind.yaml", ind)
	jpSrv := NewProxyServer("cfgs/jp.yaml", jp)
	nldSrv := NewProxyServer("cfgs/nld.yaml", nld)
	rusSrv := NewProxyServer("cfgs/rus.yaml", rus)
	sgpSrv := NewProxyServer("cfgs/sgp.yaml", sgp)
	taiwanSrv := NewProxyServer("cfgs/taiwan.yaml", taiwan)
	usaSrv := NewProxyServer("cfgs/usa.yaml", usa)
	go serve(canSrv)
	go serve(gbrSrv)
	go serve(indSrv)
	go serve(jpSrv)
	go serve(nldSrv)
	go serve(rusSrv)
	go serve(sgpSrv)
	go serve(taiwanSrv)
	go serve(usaSrv)
	select {}
}
