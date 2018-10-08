package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"github.com/google/martian"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
	"github.com/google/martian/log"
)

var (
	skipTLSVerify = flag.Bool("skip-tls-verify", false, "skip TLS server verification; insecure")
)

type Router struct {
	port int
	host string
	url  string
}

var (
	logLevel = flag.Int("l", 0, "log level")
)

func createProxy(s string) (p *Router, err error) {
	params := strings.SplitN(s, ":", 3)
	if len(params) < 3 {
		return nil, errors.New("invalid string")
	}
	port, err := strconv.Atoi(params[1])
	if err != nil {
		return nil, err
	}
	return &Router{port, params[0], params[2]}, nil
}

func serveProxy(r *Router) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", r.host, r.port))
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	m := martian.NewProxy()
	u, err := url.Parse(r.url)
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		Proxy:                 http.ProxyURL(u),
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: *skipTLSVerify,
		},
	}
	m.SetDownstreamProxy(u)
	m.SetRoundTripper(tr)
	m.Serve(l)
}
func main() {
	flag.Parse()
	log.SetLevel(*logLevel)
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	log.Debugf("fwproxy: log level: %d", log.GetLevel())
	routers := make([]Router, 0)
	for idx := range flag.Args() {
		if p, err := createProxy(flag.Arg(idx)); err == nil {
			routers = append(routers, *p)
		}
	}
	if len(routers) != 0 {
		fmt.Printf("fwproxy: %d proxy servers\n", len(routers))
		for idx := range routers {
			go serveProxy(&routers[idx])
		}
		//routers[0].Run()
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt, os.Kill)

		<-sigc
		log.Infof("fwproxy: shutting down")
	} else {
		log.Infof("fwproxy: no routers")
	}
}
