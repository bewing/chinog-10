package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"net/netip"
	"os"
	"text/template"

	"github.com/hairyhenderson/gomplate/v3"
	"github.com/hairyhenderson/gomplate/v3/data"
)

type Node struct {
	RouterId string
}

type NodeData struct {
	Region int
	Site   int
	Layer  int
	Type   string
	ASN    string
}

var tmpl *template.Template

func loadNodeData(routerId string) (NodeData, error) {
	nd := NodeData{}
	ip, err := netip.ParseAddr(routerId)
	if err != nil {
		return nd, err
	}
	data := ip.AsSlice()[2]
	typeByte := data & 12 >> 2
	if typeByte^1 == 0 {
		nd.Type = "leaf"
		nd.Layer = 1
	} else if typeByte^2 == 0 {
		nd.Type = "spine"
		nd.Layer = 2
	} else if typeByte^3 == 0 {
		nd.Type = "superspine"
		nd.Layer = 3
	} else {
		nd.Type = "server"
		nd.Layer = 0
	}
	nd.Region = int(data & 192 >> 6)
	nd.Site = int(data & 48 >> 4)

	nd.ASN = fmt.Sprintf("65%d%d%d.%d", nd.Region, nd.Site, nd.Layer, int(ip.AsSlice()[2])*256+int(ip.AsSlice()[3]))
	return nd, nil
}

func IPv4ToInt(routerId string) (uint32, error) {
	ip, err := netip.ParseAddr(routerId)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(ip.AsSlice()), nil
}

func init() {
	gfuncs := gomplate.CreateFuncs(context.Background(), new(data.Data))
	myFuncs := template.FuncMap{}
	myFuncs["loadNodeData"] = loadNodeData
	myFuncs["IPv4ToInt"] = IPv4ToInt
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
