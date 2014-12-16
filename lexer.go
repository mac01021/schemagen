package schemagen

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

type Lexer struct {
	in       *bufio.Reader
	eof      bool
	buf      *Deque
	consumed []rune

	//track the location in the file
	line int
	col  int

	//output channel
	Tokens chan *Token
}

func NewLexer(in io.Reader) *Lexer {
	it := new(Lexer)
	it.in = bufio.NewReader(in)
	it.buf = NewDeque()
	it.Tokens = make(chan *Token)
	it.line = 1
	it.col = -1 //TODO: track the column as well as the line number
	it.start()
	return it
}

var EOF = errors.New("EOF")

func (lexer *Lexer) fillBuf() {
	// Only to be called if the buffer is empty!
	var nb int
	for {
		r, _, err := lexer.in.ReadRune()
		nb++
		if err == nil {
			lexer.buf.pushR(r)
			if r == '\n' {
				break
			}
		} else if err == io.EOF {
			lexer.eof = true
			break
		} else {
			err = fmt.Errorf("invalid rune at line %d, col %d", lexer.line, lexer.col+nb)
			panic(err)
		}
	}
}

func (lexer *Lexer) isDone() bool {
	return lexer.buf.isEmpty() && lexer.eof
}

func (lexer *Lexer) forget() {
	lexer.consumed = nil
}

func (lexer *Lexer) text() string {
	return string(lexer.consumed)
}

func (lexer *Lexer) next() (rune, bool) {
	if lexer.buf.isEmpty() {
		lexer.fillBuf()
	}
	if lexer.isDone() {
		return 0, false
	}
	r, _ := lexer.buf.popL()
	lexer.consumed = append(lexer.consumed, r)
	lexer.updateCursorCoords(r)
	return r, true
}

func (lexer *Lexer) backup() error {
	l := len(lexer.consumed)
	if l < 1 {
		return errors.New("no operation to undo")
	}
	r := lexer.consumed[l-1]
	lexer.consumed = lexer.consumed[:l-1]
	lexer.buf.pushL(r)
	lexer.downdateCursorCoords(r)
	return nil
}

func (lexer *Lexer) emit(ttype TokType) {
	tok := new(Token)
	tok.Type = ttype
	tok.Line = lexer.line
	tok.Col = lexer.col
	tok.Text = lexer.text()
	lexer.Tokens <- tok
	lexer.forget()
}

func (lexer *Lexer) start() {
	go func() {
		state := lexer.startToken
		for {
			if state == nil {
				close(lexer.Tokens)
				break
			}
			state = state()
		}
	}()
}

func (lexer *Lexer) updateCursorCoords(r rune) {
	if r == '\n' {
		lexer.line += 1
	}
	lexer.col = -1 //TODO
}

func (lexer *Lexer) downdateCursorCoords(r rune) {
	if r == '\n' {
		lexer.line -= 1
	}
	lexer.col = -1 //TODO
}
