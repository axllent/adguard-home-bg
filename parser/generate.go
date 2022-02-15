package parser

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/axllent/adguard-home-bg/app"
)

type conf struct {
	Ctag       string
	Client     string
	DenyAllow  string
	DNSType    string
	DNSRewrite string
	Important  bool
	BadFilter  bool
}

// Config returns the blocklist configuration
var Config = conf{}

// Build returns a report in JSON format
func Build(w http.ResponseWriter, r *http.Request) {
	urlParams := r.URL.Query()

	uri := urlParams.Get("url")

	u, err := url.ParseRequestURI(uri)
	if err != nil {
		app.HTTPError(w, err.Error())
		app.Log().Critical(err.Error())
		return
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		app.HTTPError(w, "Only http & https links supported")
		app.Log().Criticalf("Only http & https links supported: %s", uri)
		return
	}

	conf := modifiersFromArgs(r)

	domains, err := conf.URLToBlocklist(uri)
	if err != nil {
		app.HTTPError(w, err.Error())
		app.Log().Critical(err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "filename=\"blocklist.txt\"")

	if _, err := w.Write([]byte(domains)); err != nil {
		app.Log().Errorf(err.Error())
	}
}

func (c conf) parseRaw(raw string) ([]string, error) {
	modifiers := c.Modifiers()

	domains := []string{}

	domainRegex := regexp.MustCompile(`^([@|]+)?(([a-z0-9\-\.]+)\.([a-z]{2,4}))`)

	lines := strings.Split(raw, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(strings.ToLower(l))
		if l == "" || strings.HasPrefix(l, "#") || strings.HasPrefix(l, "!") {
			continue
		}

		words := strings.Split(l, " ")
		for _, w := range words {
			if w == "" {
				continue
			}

			matches := domainRegex.FindAllStringSubmatch(w, -1)

			cnt := len(matches)
			if cnt > 0 {
				if matches[0][1] == "" {
					domains = append(domains, "|"+matches[0][2]+modifiers)
				} else {
					domains = append(domains, matches[0][1]+matches[0][2]+modifiers)
				}
			}
		}

	}

	sort.Strings(domains)

	return unique(domains), nil

	// return domains, nil

	// currentTime := time.Now()
	// date := currentTime.Format(time.UnixDate)
	// sec := "! --------------------------------------------------\n"

	// return fmt.Sprintf("%s! Source:  %s\n! Domains: %d\n! Updated: %s\n%s", sec, uri, len(domains), date, sec) + strings.Join(domains, "\n"), nil
}

// URLToBlocklist will return a formatted blocklist from a URL
func (c conf) URLToBlocklist(uri string) (string, error) {
	raw, err := app.DownloadToString(uri)
	if err != nil {
		return "", err
	}

	domains, err := c.parseRaw(raw)
	if err != nil {
		return "", err
	}

	currentTime := time.Now()
	date := currentTime.Format(time.UnixDate)
	sec := "! --------------------------------------------------\n"

	return fmt.Sprintf("%s! Source:  %s\n! Domains: %d\n! Updated: %s\n%s", sec, uri, len(domains), date, sec) + strings.Join(domains, "\n"), nil

	// domains := []string{}

	// domainRegex := regexp.MustCompile(`^([@|]+)?(([a-z0-9\-\.]+)\.([a-z]{2,4}))`)

	// lines := strings.Split(raw, "\n")
	// for _, l := range lines {
	// 	l = strings.TrimSpace(strings.ToLower(l))
	// 	if l == "" || strings.HasPrefix(l, "#") || strings.HasPrefix(l, "!") {
	// 		continue
	// 	}

	// 	words := strings.Split(l, " ")
	// 	for _, w := range words {
	// 		if w == "" {
	// 			continue
	// 		}

	// 		matches := domainRegex.FindAllStringSubmatch(w, -1)

	// 		cnt := len(matches)
	// 		if cnt > 0 {
	// 			if matches[0][1] == "" {
	// 				domains = append(domains, "|"+matches[0][2]+modifiers)
	// 			} else {
	// 				domains = append(domains, matches[0][1]+matches[0][2]+modifiers)
	// 			}
	// 		}
	// 	}

	// }

	// sort.Strings(domains)

	// domains = unique(domains)

	// currentTime := time.Now()
	// date := currentTime.Format(time.UnixDate)
	// sec := "! --------------------------------------------------\n"

	// return fmt.Sprintf("%s! Source:  %s\n! Domains: %d\n! Updated: %s\n%s", sec, uri, len(domains), date, sec) + strings.Join(domains, "\n"), nil
}

func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (c conf) Modifiers() string {
	mods := []string{}

	if c.Ctag != "" {
		mods = append(mods, "$ctag="+strings.Replace(c.Ctag, ",", "|", -1))
	}

	if c.Client != "" {
		mods = append(mods, "$client="+strings.Replace(c.Client, ",", "|", -1))
	}

	if c.DenyAllow != "" {
		mods = append(mods, "$denyallow="+strings.Replace(c.DenyAllow, ",", "|", -1))
	}

	if c.DNSType != "" {
		mods = append(mods, "$dnstype="+strings.Replace(c.DNSType, ",", "|", -1))
	}

	if c.DNSRewrite != "" {
		mods = append(mods, "$dnsrewrite="+strings.Replace(c.DNSRewrite, ",", ";", -1))
	}

	if c.Important {
		mods = append(mods, "$important")
	}

	if c.BadFilter {
		mods = append(mods, "$badfilter")
	}

	return "^" + strings.Join(mods, ",")
}

func modifiersFromArgs(r *http.Request) conf {
	// mods := []string{}
	// opts := []string{"ctag", "client", "denyallow", "dnstype", "dnsrewrite"}
	urlParams := r.URL.Query()

	config := conf{}

	config.Ctag = strings.TrimSpace(urlParams.Get("ctag")) //, ",", "|", -1))
	config.Client = strings.TrimSpace(urlParams.Get("client"))
	config.DenyAllow = strings.TrimSpace(urlParams.Get("denyallow"))
	config.DNSType = strings.TrimSpace(urlParams.Get("dnstype"))
	config.DNSRewrite = strings.TrimSpace(urlParams.Get("dnsrewrite"))
	config.Important = urlParams.Has("important")
	config.BadFilter = urlParams.Has("badfilter")

	return config // .Modifiers()
}
