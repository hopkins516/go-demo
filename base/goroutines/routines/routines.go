// Copyright © 2022 UCloud. All rights reserved.

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if len(r.FormValue("case-two")) > 0 {
			fmt.Println("case two")
		} else {
			fmt.Println("case one start")
			time.Sleep(time.Second * 5)
			fmt.Println("case one end")
		}
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
