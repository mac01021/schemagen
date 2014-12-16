package schemagen

import (
	"unicode"
)

type StateFn func() StateFn

func (lexer *Lexer) startToken() StateFn {
	r, ok := lexer.next()
	if !ok {
		return nil
	}
	if r == '{' {
		lexer.emit(LBRACE)
		return lexer.startToken
	}
	if r == '}' {
		lexer.emit(RBRACE)
		return lexer.startToken
	}
	if r == '(' {
		lexer.emit(LPAREN)
		return lexer.startToken
	}
	if r == ')' {
		lexer.emit(RPAREN)
		return lexer.startToken
	}
	if r == ';' || r == '\n' {
		lexer.emit(ENDL)
		return lexer.startToken
	}
	if unicode.IsDigit(r) {
		return lexer.finishNumber
	}
	if unicode.IsLetter(r) {
		return lexer.finishIdent
	}
	if unicode.IsSpace(r) {
		return lexer.finishSpace
	}
	panic("rune has no class")
}

func (lexer *Lexer) finishNumber() StateFn {
	for {
		r, ok := lexer.next()
		if !ok {
			break
		}
		if !unicode.IsDigit(r) {
			lexer.backup()
			break
		}
	}
	lexer.emit(NUMBER)
	if lexer.isDone() {
		return nil
	}
	return lexer.startToken
}

type marker struct{}

var m marker
var idRunes = map[rune]marker{
	'a': m, 'b': m, 'c': m, 'd': m, 'e': m, 'f': m, 'g': m, 'h': m,
	'i': m, 'j': m, 'k': m, 'l': m, 'm': m, 'n': m, 'o': m, 'p': m,
	'q': m, 'r': m, 's': m, 't': m, 'u': m, 'v': m, 'x': m, 'y': m,
	'z': m, '_': m,
	'A': m, 'B': m, 'C': m, 'D': m, 'E': m, 'F': m, 'G': m, 'H': m,
	'I': m, 'J': m, 'K': m, 'L': m, 'M': m, 'N': m, 'O': m, 'P': m,
	'Q': m, 'R': m, 'S': m, 'T': m, 'U': m, 'V': m, 'X': m, 'Y': m,
	'Z': m,
}

func (lexer *Lexer) finishIdent() StateFn {
	for {
		r, ok := lexer.next()
		if !ok {
			break
		}

		if _, ok := idRunes[r]; !ok {
			lexer.backup()
			break
		}
	}
	lexer.emit(IDENT)
	if lexer.isDone() {
		return nil
	}
	return lexer.startToken
}

func (lexer *Lexer) finishSpace() StateFn {
	for {
		r, ok := lexer.next()
		if !ok {
			break
		}
		if r == '\n' || !unicode.IsSpace(r) {
			lexer.backup()
			break
		}
	}
	lexer.emit(WHITE)
	if lexer.isDone() {
		return nil
	}
	return lexer.startToken
}
