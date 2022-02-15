package parser

import (
	"fmt"
	"strings"
	"testing"
)

var testSrc = `
! ------------------------------------
! This is just a sample blocklist using different formats
! ------------------------------------
# AgGuard-formatted domains
||94.69.23.168^some crap
|1.212.242.67^$important
||azerbaijan-tourism.com^
@@||google.com^

# some plain domain names
example.com
r2---sn-n4v7knlz.googlevideo.com
google.com
blaahvlaah.com

# some host-formatted domains
0.0.0.0 100.1qingdao.com
0.0.0.0 7minuteworkout.com
0.0.0.0 	addshoppers.com

 ## some host file-formatted blocks
0.0.0.0 google1.com
127.0.0.1       blaahvlaah.com

# some native adguard filters
@@||ipinfo.io^! Issue 21
@@||yospace.com^! Issue 12
||1-cl0ud.com^
||0----q.tumblr.com^
||0--liamariejohnson--0.nedrobin.net^
||0-0-adult-superstore.com^
||0-0wearingglassesnakedmen.tumblr.com^
`

var testResult = `@@||google.com^$ctag=ctag1|tag2,$client=1.2.3.4|2.3.4.5,$denyallow=domain1|domain2,$dnstype=-A|AAAA,$dnsrewrite=domain3;domain4,$important,$badfilter
@@||ipinfo.io^$ctag=ctag1|tag2,$client=1.2.3.4|2.3.4.5,$denyallow=domain1|domain2,$dnstype=-A|AAAA,$dnsrewrite=domain3;domain4,$important,$badfilter
@@||yospace.com^$ctag=ctag1|tag2,$client=1.2.3.4|2.3.4.5,$denyallow=domain1|domain2,$dnstype=-A|AAAA,$dnsrewrite=domain3;domain4,$important,$badfilter
|100.1qingdao.com^$ctag=ctag1|tag2,$client=1.2.3.4|2.3.4.5,$denyallow=domain1|domain2,$dnstype=-A|AAAA,$dnsrewrite=domain3;domain4,$important,$badfilter
|7minuteworkout.com^$ctag=ctag1|tag2,$client=1.2.3.4|2.3.4.5,$denyallow=domain1|domain2,$dnstype=-A|AAAA,$dnsrewrite=domain3;domain4,$important,$badfilter
|blaahvlaah.com^$ctag=ctag1|tag2,$client=1.2.3.4|2.3.4.5,$denyallow=domain1|domain2,$dnstype=-A|AAAA,$dnsrewrite=domain3;domain4,$important,$badfilter
|example.com^$ctag=ctag1|tag2,$client=1.2.3.4|2.3.4.5,$denyallow=domain1|domain2,$dnstype=-A|AAAA,$dnsrewrite=domain3;domain4,$important,$badfilter
|google.com^$ctag=ctag1|tag2,$client=1.2.3.4|2.3.4.5,$denyallow=domain1|domain2,$dnstype=-A|AAAA,$dnsrewrite=domain3;domain4,$important,$badfilter
|google1.com^$ctag=ctag1|tag2,$client=1.2.3.4|2.3.4.5,$denyallow=domain1|domain2,$dnstype=-A|AAAA,$dnsrewrite=domain3;domain4,$important,$badfilter
|r2---sn-n4v7knlz.googlevideo.com^$ctag=ctag1|tag2,$client=1.2.3.4|2.3.4.5,$denyallow=domain1|domain2,$dnstype=-A|AAAA,$dnsrewrite=domain3;domain4,$important,$badfilter
||0----q.tumblr.com^$ctag=ctag1|tag2,$client=1.2.3.4|2.3.4.5,$denyallow=domain1|domain2,$dnstype=-A|AAAA,$dnsrewrite=domain3;domain4,$important,$badfilter
||0--liamariejohnson--0.nedrobin.net^$ctag=ctag1|tag2,$client=1.2.3.4|2.3.4.5,$denyallow=domain1|domain2,$dnstype=-A|AAAA,$dnsrewrite=domain3;domain4,$important,$badfilter
||0-0-adult-superstore.com^$ctag=ctag1|tag2,$client=1.2.3.4|2.3.4.5,$denyallow=domain1|domain2,$dnstype=-A|AAAA,$dnsrewrite=domain3;domain4,$important,$badfilter
||0-0wearingglassesnakedmen.tumblr.com^$ctag=ctag1|tag2,$client=1.2.3.4|2.3.4.5,$denyallow=domain1|domain2,$dnstype=-A|AAAA,$dnsrewrite=domain3;domain4,$important,$badfilter
||1-cl0ud.com^$ctag=ctag1|tag2,$client=1.2.3.4|2.3.4.5,$denyallow=domain1|domain2,$dnstype=-A|AAAA,$dnsrewrite=domain3;domain4,$important,$badfilter
||azerbaijan-tourism.com^$ctag=ctag1|tag2,$client=1.2.3.4|2.3.4.5,$denyallow=domain1|domain2,$dnstype=-A|AAAA,$dnsrewrite=domain3;domain4,$important,$badfilter`

func TestGeneration(t *testing.T) {
	config := conf{
		Ctag:       "ctag1,tag2",
		Client:     "1.2.3.4,2.3.4.5",
		DenyAllow:  "domain1,domain2",
		DNSType:    "-A,AAAA",
		DNSRewrite: "domain3,domain4",
		Important:  true,
		BadFilter:  true,
	}

	domains, err := config.parseRaw(testSrc)
	if err != nil {
		t.Error(err)
	}

	output := strings.Join(domains, "\n")

	if output != testResult {
		fmt.Println(output)
		t.Fail()
	}
}
