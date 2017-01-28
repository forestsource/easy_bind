package main

import (
	"fmt"
	"os"
)

func zoneFile(z zone) {
	fmt.Println("make zone file", z)
	filename := z.domain_name
	adminMail := ""
	file, _ := os.Create(`/` + filename)
	defer file.Close()
	zoneBase = `
  $TTL 3600
  @      IN   SOA`
	zoneBase = zoneBase + "  " + z.domain_name + "  " + adminMail + "(\n"
	zoneBase = zoneBase + `
      1001251650 ; Serial
      10800      ; Refresh 3 hour
      3600       ; Retry 1 hour
      3600000    ; Expire 1000 hours
      3600       ; Minium 1 hours
  );
       IN    NS    NS1.DO-REG.JP.
       IN    A     203.183.220.196
;
;
WWW    IN    A     203.183.220.196
`
}
