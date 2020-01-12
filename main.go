package main

import (
	"fmt"
	"os"

	"Email/pop3"
	"Email/shared/tr"
)

func main() {
	tr.Init()

	addr := os.Getenv("BSMAIL_SERVER")
	user := os.Getenv("BSMAIL_USER")
	pass := os.Getenv("BSMAIL_PASS")
	if addr == "" || user == "" || pass == "" {
		return
	}

	if p3 := pop3.New(addr, -1, pop3.TLS); p3 != nil {
		defer func() {
			if q := p3.Quit(); q != nil {
				fmt.Println(q.Data)
			}
			p3.Close()
		}()

		if rsp := p3.Read(); rsp != nil && rsp.IsOK() {
			fmt.Println(rsp.Data)

			if _, ok := p3.Auth(user, pass); ok {
				if r, ok := p3.Stat(); ok {
					fmt.Println("STAT:", r.Data)

					if items := p3.List(); len(items) > 0 {
						for _, item := range items {
							p3.ReadEmail(item)
							fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
						}
					}
				}
			}
		}
	}
}
