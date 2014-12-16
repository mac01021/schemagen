package main

import "fmt"
import "strings"
import "schemagen"
import "strconv"

func main() {
	var text = `

customers {
	id int pk
	name string(42)
	phone string(512) null
}

invoices {
	id int pk; sent_on date; filled_on date  null
}

`

	in := strings.NewReader(text)
	scanner := schemagen.NewLexer(in)
	for tok := range scanner.Tokens {
		printToken(tok)
	}
}

func printToken(tok *schemagen.Token) {
	fmt.Printf("%s[%s] on line %d\n", schemagen.TypeName(tok),
		strconv.Quote(tok.Text), tok.Line)
}
