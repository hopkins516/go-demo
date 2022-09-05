package _template

import (
	sync2 "demo_go/base/sync"
	"fmt"
	"html/template"
	"os"
)

func TemplateDemo()  {
	p := sync2.Person { Name: "hopkins", Age: 34 }
	tmpl, err := template.New("test").Parse(
		"Name: {{.Name}}, Age: {{.Age}}")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, p)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println(tmpl.Name())
}