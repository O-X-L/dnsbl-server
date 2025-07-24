package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type requestLog struct {
	Client   string `json:"cli"`
	Type     string `json:"type"`
	Request  string `json:"req"`
	Status   int    `json:"status"`
	Response string `json:"res"`
}

type requestLogTime struct {
	Time     string `json:"time"`
	Client   string `json:"cli"`
	Type     string `json:"type"`
	Request  string `json:"req"`
	Status   int    `json:"status"`
	Response string `json:"res"`
}

func logRequest(q string, s int, cli string, c *DNSBLRunningConfig, t int, res string) {
	if !c.Log {
		return
	}
	ts := "IP"
	if t == LOOKUP_DOMAIN {
		ts = "Domain"
	}

	if c.LogJSON {
		var err error
		var l []byte
		if c.LogTime {
			l, err = json.Marshal(requestLogTime{
				Time:     time.Now().Format(time.RFC3339),
				Client:   cli,
				Type:     ts,
				Request:  q,
				Status:   s,
				Response: res,
			})
		} else {
			l, err = json.Marshal(requestLog{
				Client:   cli,
				Type:     ts,
				Request:  q,
				Status:   s,
				Response: res,
			})
		}

		if err != nil {
			fmt.Printf("ERROR: Unable to log in JSON-format - %v", err)
			return
		}
		fmt.Println(string(l))

	} else {
		l := fmt.Sprintf("[%s] => %s: %s <= %d %s", cli, ts, q, s, res)
		if c.LogTime {
			log.Println(l)
		} else {
			fmt.Println(l)
		}
	}
}
