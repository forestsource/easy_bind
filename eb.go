package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

type zone struct {
	domain_name string
	ip          string
	cname       string
	rmail       string
	smail       string
}

type acl struct {
	listname string
	ips      string
}

type option struct {
	recursion bool
	allow-transfer string
	allow-query string
	max-cache-size string
	maintenance-ip string
	transfers-out string
	transfers-in string

}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("statics/*.html")),
	}
	e.Renderer = t
	e.File("/", "statics/home.html")
	e.GET("/make", make)
	e.Static("/statics", "statics")
	e.Logger.Fatal(e.Start("localhost:8080"))
}

func home(c echo.Context) error {
	return c.Render(http.StatusOK, "home", "")
}

func make(c echo.Context) error {
	fmt.Println("=====")
	fmt.Printf("%#v", c.QueryParams())
	fmt.Println("=====")
	fmt.Println("=====")
	zoneNum, _ := strconv.Atoi(c.QueryParam("zone_num"))
	aclNum, _ := strconv.Atoi(c.QueryParam("acl_num"))
	fmt.Println("zonenum", zoneNum)
	fmt.Println("aclnum", aclNum)
	fmt.Println("=====")
	//
	// zone loop
	//
	var zones []zone
	for i := 0; i < zoneNum; i++ {
		si := strconv.Itoa(i)
		fmt.Println("zone num:", si)
		var z zone
		z.domain_name = c.QueryParam("zone" + si + "[domain]")
		z.ip = c.QueryParam("zone" + si + "[ip]")
		z.cname = c.QueryParam("zone" + si + "[cname]")
		z.rmail = c.QueryParam("zone" + si + "[rmail]")
		z.smail = c.QueryParam("zone" + si + "[smail]")
		fmt.Printf("zone", si, "domain:", z.domain_name)
		fmt.Printf(" ip:", z.ip)
		fmt.Printf(" cname:", z.cname)
		fmt.Printf(" rmail:", z.rmail)
		fmt.Printf(" smail:", z.smail)
		fmt.Println("")
		zones = append(zones, z)
	}

	//
	//acl loop
	//
	var acls []acl
	for i := 0; i < aclNum; i++ {
		si := strconv.Itoa(i)
		fmt.Println("acl num:", si)
		var a acl
		a.listname = c.QueryParam("acl" + si + "[listname]")
		a.ips = c.QueryParam("acl" + si + "[ips]")
		fmt.Printf(" listname:", a.listname)
		fmt.Printf(" ips:", a.ips)
		fmt.Println("")
		acls = append(acls, a)
	}
	return c.String(http.StatusOK, buildConf(zones, acls))
}

func staticFile() http.Handler {
	return http.StripPrefix("/statics/", http.FileServer(http.Dir("statics")))
}

func buildZone(zoneName string) string {
	s := "zone \"" + zoneName + "\" {\n"
	s = s + "  type hint;\n"
	s = s + "  file '" + zoneName + "';\n"
	s = s + "};\n"
	fmt.Println("build zone:", s)
	return s
}

func buildAcl(aclName string, aclIp string) string {
	s := "acl \"" + aclName + "\" {\n"
	for {
		i := strings.Index(aclIp, ",")
		fmt.Println("index ", i)
		if i == -1 {
			s = s + "  " + aclIp + ";\n"
			break
		}
		s = s + "  " + aclIp[:i] + ";\n"
		aclIp = aclIp[i+1:]
	}
	s = s + "}\n"
	fmt.Println("build acl:", s)
	return s
}
func buildOptionContents(){
	s = "options {"
	s = s + "  edns no;\n"
	s = s + "  allow-new-zones no;\n"
	s = s + "  auth-nxdomain no;\n"
	s = s + "  memstatistics yes;\n"
	s = s + "  dialup no;\n"
	s = s + "  minimal-responses yes;\n"
	s = s + "  notify-to-soa no;\n"
	s = s + "  request-nsid no;\n"
	s = s + "  request-sit no;\n"//default
	s = s + "  nosit-udp-size 512;\n"//default
	s = s + "  multi-master no;\n"
	s = s + "  allow-query-cache {localhost;};\n"//default
	s = s + "  allow-update none;\n"
	s = s + "  query-source address * port *;\n"//default
	s = s + "}\n"
	fmt.Println("build options:", s)
	return s
}

func buildOptionResolver(){
	s = "options {"
	s = s + "  edns no;\n"
	s = s + "  allow-new-zones yes;\n"
	s = s + "  auth-nxdomain no;\n"
	s = s + "  memstatistics yes;\n"
	s = s + "  dialup no;\n"
	s = s + "  cleaning-interval 60;\n"
	s = s + "  lame-ttl 600;\n"
	s = s + "  max-cache-size unlimited;\n"
	s = s + "  max-ncache-ttl 300;\n"
	s = s + "  minimal-responses yes;\n"
	s = s + "  notify-to-soa no;\n"
	s = s + "  request-nsid no;\n"
	s = s + "  request-sit no;\n"//default
	s = s + "  nosit-udp-size 512;\n"//default
	s = s + "  multi-master no;\n"
	s = s + "  allow-query-cache {localhost;};\n"//default
	s = s + "  allow-update none;\n"
	s = s + "  query-source address * port *;\n"//default
	s = s + "}\n"
	fmt.Println("build options:", s)
	return s
}

func buildConf(zones []zone, acls []acl) string {
	zonsString := "\n"
	for i := 0; i < len(zones); i++ {
		zonsString = zonsString + buildZone(zones[i].domain_name)
	}
	aclsString := "\n"
	for i := 0; i < len(acls); i++ {
		aclsString = aclsString + buildAcl(acls[i].listname, acls[i].ips)
	}
	configBase := `key "rndc-key" {
    algorithm hmac-md5;
    secret "qpBt4e8VSpFLWdVkD0AlRRq=";
};

controls {
    inet 127.0.0.1 port 953
    allow { 127.0.0.1; } keys { "rndc-key"; };
};

options {
    directory "/var/named";
    allow-transfer { none}
        {{.Ip1}};
    allow-query {
      any;
    };
    version "none";
    recursion{no};
    max-cache-size {{.chache-size}};
};`
	configBase = aclsString + configBase
	configBase = configBase + zonsString
	return configBase
}
