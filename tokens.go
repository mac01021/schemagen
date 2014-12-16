package schemagen

import (
	"strconv"
)

type Token struct {
	Line int
	Col  int
	Type TokType
	Text string
}

func (tok *Token) number() int {
	i, _ := strconv.Atoi(tok.Text)
	return i
}

type TokType int

const (
	LPAREN TokType = iota
	RPAREN
	LBRACE
	RBRACE
	IDENT
	WHITE
	ENDL
	NUMBER
)

var typenames = map[TokType]string{
	LPAREN: "LPAREN",
	RPAREN: "RPAREN",
	LBRACE: "LBRACE",
	RBRACE: "RBRACE",
	IDENT:  "IDENTIFIER",
	WHITE:  "WHITEPACE",
	ENDL:   "ENDL",
	NUMBER: "NUMBER",
}

func TypeName(tok *Token) string {
	name, _ := typenames[tok.Type]
	return name
}

func (ttype TokType) in(ttypes []TokType) bool {
	for _, t := range ttypes {
		if ttype == t {
			return true
		}
	}
	return false
}
