package main

import "fmt"
import "strings"
import "schemagen"

func main() {
	var text = `

customers {
	id int pk
	name string(42)
	phone string(512) null
}

invoices {
	id int pk; sent_on date; filled_on date  null; customer fk(customers)
}


blocks { id int(64) pk; content binary(512) }

modifications {
	at timestamp pk
	target fk(blocks) pk
}


editors {
	id int(16) pk
	is_admin  bool
	most_recent_edit fk(modifications) null
}

`

	in := strings.NewReader(text)
	scanner := schemagen.NewLexer(in)
	parser := schemagen.NewParser()
	schema, err := parser.Parse(scanner.Tokens)
	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Success")
	for name, tab := range schema {
		fmt.Println(name + ":")
		fmt.Println(schemagen.CreateStatement(tab))
		fmt.Println()
	}
}
