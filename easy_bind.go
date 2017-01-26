package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type server_info struct {
	ip   string
	port int
}

func main() {
	si := server_info{ip: "127.0.0.1", port: 8080}
	fmt.Println("start server")
	http.HandleFunc("/", home)
	http.HandleFunc("/make", make)
	http.Handle("/statics/", staticFile())
	http.ListenAndServe(si.ip+":"+strconv.Itoa(si.port), nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("statics/home.html")
	err = tmpl.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

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

type context struct {
	zone []zone
	acl  []acl
}

func make(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("=====")
	fmt.Printf("%#v", r.Form)
	fmt.Println("=====")
	fmt.Println("=====")
	zoneNum, _ := strconv.Atoi(r.FormValue("zone_num"))
	aclNum, _ := strconv.Atoi(r.FormValue("acl_num"))
	fmt.Println("zonenum", zoneNum)
	fmt.Println("aclnum", aclNum)
	var zones [10]zone
	//zones[0] = zone{r.FormValue("zone0[domain]"), r.FormValue("zone0[ip]"), r.FormValue("zone0[cname]"), r.FormValue("zone0[rmail]"), r.FormValue("zone0[sname]")}
	zones[0] = zone{r.FormValue("zone0[domain]"), r.FormValue("zone0[ip]"), "", "", ""}
	fmt.Println("zone0 domain", r.FormValue("zone0[domain]"))
	fmt.Println("zone0 ip", r.FormValue("zone0[ip]"))
	fmt.Println("=====")
	//fmt.Println(r.FormValue("zone0[cname]"))
	//fmt.Println(r.FormValue("zone0[rmail]"))
	//fmt.Println(r.FormValue("zone0[smail]"))
	for i := 0; i > 1; i++ {
		fmt.Println("i:%d", i)
		num := strconv.Itoa(i)
		//stringDomain := []string{"zone", "aaaaaaa", "domain"}
		fmt.Println("num:", num)
		//fmt.Println(strings.Join(stringDomain, ","))
		//z := zone{r.FormValue(stringDomain), r.FormValue("zone" + num + "[ip]"), r.FormValue("zone" + num + "[cname]"), r.FormValue("zone" + num + "[rmail]"), r.FormValue("zone" + num + "[sname]")}
		//fmt.Printf("%#V", z)
		//zones = append(zones, z)
	}
	//fmt.Printf("%#v", r.FormValue("zones[zone1][domain]"))
	//jsonに個数を持たせて、それを元にformvalueでとればいい。
	tmpl, err := template.ParseFiles("statics/config.html")
	//tmpv := IPs{"127.0.0.1", "127.0.0.2"}
	//err = tmpl.Execute(w, tmpv)
	err = tmpl.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func staticFile() http.Handler {
	return http.StripPrefix("/statics/", http.FileServer(http.Dir("statics")))
}
