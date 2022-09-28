package main

import (
	"context"
	"flag"
	"net/netip"
	"os"
	"text/template"

	"github.com/hairyhenderson/gomplate/v3"
	"github.com/hairyhenderson/gomplate/v3/data"
)

type Node struct {
	RouterId string
}

var tmpl *template.Template

func nodeType(routerId string) (string, error) {
	ip, err := netip.ParseAddr(routerId)
	if err != nil {
		return "", err
	}
	typeByte := ip.AsSlice()[2]
	typeByte = typeByte >> 2
	typeByte = typeByte & 3
	if typeByte^1 == 0 {
		return "leaf", nil
	} else if typeByte^2 == 0 {
		return "spine", nil
	} else if typeByte^3 == 0 {
		return "superspine", nil
	}
	return "server", nil
}

func asnFromRouterId(routerId string) (int, error) {
	ip, err := netip.ParseAddr(routerId)
	if err != nil {
		return -1, err
	}
	asnBytes := ip.AsSlice()[2:]
	asn := int(asnBytes[0])*256 + int(asnBytes[1])
	return asn, nil
}

func init() {
	gfuncs := gomplate.CreateFuncs(context.Background(), new(data.Data))
	myFuncs := template.FuncMap{}
	myFuncs["nodeType"] = nodeType
	myFuncs["asnFromRouterId"] = asnFromRouterId
	delete(gfuncs, "slice")
	tmpl = template.Must(template.New("").Funcs(gfuncs).Funcs(myFuncs).ParseGlob("templates/*.tpl"))
}

func main() {

	var routerId string
	flag.StringVar(&routerId, "router-id", "172.18.4.2", "")
	flag.Parse()
	td := Node{routerId}
	if err := tmpl.ExecuteTemplate(os.Stdout, "DCS-7280SR2A-48YC6.tpl", td); err != nil {
		panic(err)
	}

}
