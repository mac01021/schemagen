package schemagen

import (
	"fmt"
	"strconv"
)

type Parser struct {
	Tables map[string]*TableSpec
}

func NewParser() *Parser {
	it := new(Parser)
	it.Tables = make(map[string]*TableSpec)
	return it
}

func filtered(tokens chan *Token) chan *Token {
	out := make(chan *Token)
	go func() {
		defer close(out)
		for tok := range tokens {
			if tok.Type != WHITE {
				out <- tok
			}
		}
	}()
	return out
}

func get(ttype TokType, tokens chan *Token, discard ...TokType) (*Token, error) {
	for tok := range tokens {
		if tok.Type == ttype {
			return tok, nil
		}
		if !tok.Type.in(discard) {
			return tok, fmt.Errorf("Unexpected token (%s) on line %d. Looking for %s",
				strconv.Quote(tok.Text), tok.Line, typenames[ttype])
		}
	}
	return nil, EOF
}

func collectUpTo(tokens chan *Token, ttypes ...TokType) ([]*Token, TokType, error) {
	var toks []*Token
	for tok := range tokens {
		if tok.Type.in(ttypes) {
			return toks, tok.Type, nil
		}
		toks = append(toks, tok)
	}
	return toks, 0, EOF
}

func newTableSpec(name string) *TableSpec {
	tab := new(TableSpec)
	tab.Name = name
	tab.Columns = make(map[string]*ColumnSpec)
	return tab
}

func (p *Parser) finishTable(tab *TableSpec, tokens chan *Token) error {
	var nbCols int = 0
	for {
		col, end, err := p.getColumn(tab, tokens)
		nbCols++
		if err != nil {
			return err
		}
		if col != nil {
			tab.Columns[col.Name] = col
			if col.IsPK {
				tab.Key = append(tab.Key, col)
			}
		}
		if end {
			break
		}
	}
	if nbCols == 0 {
		return fmt.Errorf("Table (%s) has no columns", tab.Name)
	}
	return nil
}

func (p *Parser) getColumn(tab *TableSpec, tokens chan *Token) (*ColumnSpec, bool, error) {
	buf, endl, err := collectUpTo(tokens, ENDL, RBRACE)
	if err != nil {
		return nil, false, err
	}
	done := endl == RBRACE
	col, err := p.makeColumn(buf)
	return col, done, err
}

func invalid(tok *Token) error {
	return fmt.Errorf("Invalid column description on line %d", tok.Line)
}

func (p *Parser) makeColumn(buf []*Token) (*ColumnSpec, error) {
	if len(buf) == 0 {
		return nil, nil
	}
	if len(buf) < 2 {
		return nil, invalid(buf[0])
	}
	if buf[0].Type != IDENT || buf[1].Type != IDENT {
		return nil, invalid(buf[0])
	}
	col := new(ColumnSpec)
	col.Name = buf[0].Text
	col.Type = coltype(buf[1])
	buf = buf[2:]
	var err error = nil
	if len(buf) > 0 {
		err = col.getSubtypeAndModifiers(buf)
	}
	return col, err
}

func coltype(tok *Token) ColType {
	if tok.Type != IDENT {
		return -1
	}
	if tok.Text == "fk" {
		return FK
	}
	if tok.Text == "string" {
		return STRING
	}
	if tok.Text == "int" {
		return INTEGER
	}
	if tok.Text == "timestamp" {
		return TIMESTAMP
	}
	if tok.Text == "date" {
		return DATE
	}
	if tok.Text == "binary" {
		return BINARY
	}
	if tok.Text == "bool" {
		return BOOLEAN
	}
	panic("unrecognized column type")
}

func (col *ColumnSpec) getSubtypeAndModifiers(buf []*Token) error {
	if len(buf) > 2 && buf[0].Type == LPAREN {
		if buf[2].Type != RPAREN {
			return invalid(buf[2])
		}
		typtok := buf[1]
		err := col.setSubtype(typtok)
		if err != nil {
			return err
		}
		buf = buf[3:]
	}
	return col.setModifiers(buf)
}

func (col *ColumnSpec) setSubtype(tok *Token) error {
	if col.Type == FK && tok.Type == IDENT {
		col.fkTargetName = tok.Text
		return nil
	} else if tok.Type == NUMBER {
		col.Size = tok.number()
		return nil
	} else {
		return invalid(tok)
	}
	panic("unreached")
}

func (col *ColumnSpec) setModifiers(buf []*Token) error {
	if len(buf) == 0 {
		return nil
	}
	if len(buf) > 1 || buf[0].Type != IDENT {
		return invalid(buf[0])
	}
	if buf[0].Text != "null" && buf[0].Text != "pk" {
		return invalid(buf[0])
	}
	if buf[0].Text == "null" {
		col.IsNullable = true
	}
	if buf[0].Text == "pk" {
		col.IsPK = true
	}
	return nil
}

func (p *Parser) Parse(tokens chan *Token) (map[string]*TableSpec, error) {
	tokens = filtered(tokens)
	for {
		id, err := get(IDENT, tokens, ENDL)
		if err == EOF {
			break
		} else if err != nil {
			return nil, err
		}
		tok := <-tokens
		if tok.Type != LBRACE {
			return nil, fmt.Errorf("Expected brace after table name on line %d", tok.Line)
		}
		tab := newTableSpec(id.Text)
		err = p.finishTable(tab, tokens)
		if err != nil {
			return nil, err
		}
		p.Tables[tab.Name] = tab
	}
	err := p.TypeCheck()
	return p.Tables, err
}

func (p *Parser) TypeCheck() error {
	for tabName, tab := range p.Tables {
		for colName, col := range tab.Columns {
			if col.Type < 0 {
				return fmt.Errorf("Column [%s] in table [%s] has invalid type", colName, tabName)
			}
			if col.Type == FK {
				target, ok := p.Tables[col.fkTargetName]
				if !ok {
					return fmt.Errorf("There is no table [%s] to be the target of FK [%s] in table [%s]", col.fkTargetName, colName, tabName)
				}
				col.FKTarget = target
			}
		}
	}
	return nil
}
