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
	amail       string
	isMail      bool
}

type acl struct {
	listname string
	ips      string
}

type option struct {
	recursion      bool
	quick_sync     bool
	isEdns         bool
	isMreduce      bool
	memory_size    int
	maintenance_ip string
	rndc_port      string
	isResolver     bool
	isSlave        bool
	forward_ip     string
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
	// contents-slave switch-toggle
	isResolver, _ := strconv.ParseBool(c.QueryParam("isResolver"))
	fmt.Println("isResolver", isResolver)
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
		z.amail = c.QueryParam("zone" + si + "[]amail")
		z.isMail, _ = strconv.ParseBool(c.QueryParam("zone" + si + "[isMailServer]"))
		fmt.Printf("zone", si, "domain:", z.domain_name)
		fmt.Printf(" ip:", z.ip)
		fmt.Printf(" cname:", z.cname)
		fmt.Printf(" rmail:", z.amail)
		fmt.Printf(" smail:", z.isMail)
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
	//
	// option
	//
	var o option
	o.isEdns, _ = strconv.ParseBool(c.QueryParam("option[isEdns]"))
	o.memory_size, _ = strconv.Atoi(c.QueryParam("option[memory]"))
	o.isMreduce, _ = strconv.ParseBool(c.QueryParam("option[isMreduce]"))
	o.rndc_port = c.QueryParam("option[port]")
	o.maintenance_ip = c.QueryParam("option[ip]")
	o.quick_sync, _ = strconv.ParseBool(c.QueryParam("option[isQsync]"))
	o.isResolver, _ = strconv.ParseBool(c.QueryParam("option[isResolver]"))
	o.isSlave, _ = strconv.ParseBool(c.QueryParam("option[isSlave]"))
	o.forward_ip = c.QueryParam("option[forwardIp]")
	fmt.Println("option", o)

	return c.String(http.StatusOK, buildConf(zones, acls, o))
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

func buildOptionContents(options option) string {
	s := "options {\n"
	s = s + "  recursion no"
	s = s + "  version \"\" \n"
	s = s + "  directory \"/var/named\" \n"
	s = s + "  edns no;\n"
	s = s + "  allow-new-zones no;\n"
	s = s + "  auth-nxdomain no;\n"
	s = s + "  memstatistics yes;\n"
	s = s + "  dialup no;\n"
	s = s + "  notify-to-soa no;\n"
	s = s + "  also-notyfy {slave;};\n"
	if options.quick_sync {
		s = s + "  transfers_out 20;\n"
		s = s + "  max-transfer-time-out 60;\n"
		s = s + "  tcp-clients 100;　\n"
		s = s + "  min-refresh-time 900;　\n"
		s = s + "  max-refresh-time 43200;　\n"
		s = s + "  min-retry-time 900;　\n"
		s = s + "  max-retry-time 1800;　\n"
	} else {
		s = s + "  transfers_out 10;\n"
		s = s + "  max-transfer-time-out 120;\n"
		s = s + "  tcp-clients 100;　\n"
		s = s + "  min-refresh-time 1800;　\n"
		s = s + "  max-refresh-time 86400;　\n"
		s = s + "  min-retry-time 900;　\n"
		s = s + "  max-retry-time 1800;　\n"
	}
	if options.isEdns {
		s = s + "  max-udp-size 4096;　\n"
	} else {
		s = s + "  max-udp-size 512;　\n"
	}
	s = s + "  request-nsid no;\n"
	s = s + "  multi-master no;\n"
	s = s + "  allow-query-cache {localhost;}\n" //default
	s = s + "  allow-update {none;}\n"
	s = s + "  query-source address * port *;\n" //default
	s = s + "  listen-on !" + options.maintenance_ip + ";\n"
	s = s + "}\n"
	s = s + `server 0.0.0.0/0{
		edns no;
	};
	`
	fmt.Println("build options:", s)
	return s
}

func buildOptionSlave(options option) string {
	s := "options {\n"
	s = s + "  directory \"/var/named\" \n"
	s = s + "  recursion no"
	s = s + "  version \"\" \n"
	s = s + "  allow-new-zones no;\n"
	s = s + "  auth-nxdomain no;\n"
	s = s + "  memstatistics yes;\n"
	s = s + "  dialup no;\n"
	s = s + "  minimal-responses yes;\n"
	s = s + "  notify-to-soa no;\n"
	if options.quick_sync {
		s = s + "  transfer-per-ns 5;\n"
		s = s + "  transfers_in 20;\n"
		s = s + "  max-transfer-time-in 60;\n"
		s = s + "  tcp-clients 200;　\n"
	} else {
		s = s + "  transfer-per-ns 2;\n"
		s = s + "  transfers_in 10;\n"
		s = s + "  max-transfer-time-in 120;\n"
		s = s + "  tcp-clients 100;　\n"
	}
	s = s + "  request-nsid no;\n"
	s = s + "  multi-master no;\n"
	s = s + "  allow-query-cache {localhost;}\n" //default
	s = s + "  allow-update {none;}\n"
	s = s + "  query-source address * port *;\n" //default
	s = s + "  listen-on !" + options.maintenance_ip + ";\n"
	s = s + "}\n"

	s = s + `server 0.0.0.0/0{
		edns no;
	};
	`
	fmt.Println("build options:", s)
	return s
}

func buildOptionResolver(options option) string {
	rclients := (options.memory_size * 1000 / 20 / 3)
	mcache := (options.memory_size / 3)
	s := "options {\n"
	s = s + "  directory \"/var/named\" \n"
	s = s + "  recursion yes\n"
	s = s + "  version \"\" \n"
	s = s + "  allow-recusion {internal;};\n"
	s = s + "  forward first;\n"
	s = s + "  forwarders {" + options.forward_ip + "}\n"
	if options.isEdns {
		s = s + "  max-udp-size 4096;\n"
	} else {
		s = s + "  max-udp-size 512;\n"
	}
	s = s + "  forwarders no;\n"
	s = s + "  auth-nxdomain no;\n"
	s = s + "  memstatistics yes;\n"
	s = s + "  dialup no;\n"
	s = s + "  lame-ttl 600;\n"
	s = s + "  max-cache-size " + strconv.Itoa(mcache) + ";\n"
	if options.quick_sync {
		s = s + "  max-ncache-ttl 300;\n"
		s = s + "  max-cache-ttl 3600;\n"
		s = s + "  cleaning-interval 60;\n"
	} else {
		s = s + "  max-ncache-ttl 300;\n"
		s = s + "  max-cache-ttl 2000;\n"
		s = s + "  cleaning-interval 60;\n"
	}
	s = s + "  minimal-responses yes;\n"
	s = s + "  notify-to-soa no;\n"
	s = s + "  request-nsid no;\n"
	s = s + "  recursive-clients " + strconv.Itoa(rclients) + ";\n"
	s = s + "  allow-query-cache {localhost;}\n" //default
	s = s + "  query-source address * port *;\n" //default
	s = s + "}\n"
	fmt.Println("build options:", s)

	if options.isEdns {
		s = s + `server 0.0.0.0/0{
			edns no;
		};
		`
	} else {
		s = s + `server 0.0.0.0/0{
			edns yes;
		};
		`
	}
	return s
}

func buildConf(zones []zone, acls []acl, options option) string {
	zonsString := "\n"
	for i := 0; i < len(zones); i++ {
		zonsString = zonsString + buildZone(zones[i].domain_name)
	}
	aclsString := "\n"
	for i := 0; i < len(acls); i++ {
		aclsString = aclsString + buildAcl(acls[i].listname, acls[i].ips)
	}
	optionsString := "\n"
	if options.isResolver {
		optionsString = buildOptionResolver(options)
	} else {
		if options.isSlave {
			optionsString = buildOptionSlave(options)
		} else {
			optionsString = buildOptionContents(options)
		}
	}

	configBase := `key "rndc-key" {
    algorithm hmac-md5;
    secret "qpBt4e8VSpFLWdVkD0AlRRq=";
};
controls {
    inet 127.0.0.1 port 953
    allow { 127.0.0.1; } keys { "rndc-key"; };
};
`
	configBase = aclsString + configBase
	configBase = configBase + optionsString
	configBase = configBase + zonsString
	return configBase
}
