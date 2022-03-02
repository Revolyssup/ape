package lexer

import (
	token "github.com/Revolyssup/ape/token"
)

type Lexer struct {
	input    string
	lastRead int
	readPos  int
	ch       byte
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	if l.ch == '#' {
		l.skipComment()
		l.skipWhitespace()
	}
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			firstchar := l.ch
			l.read()
			tok = token.Token{Type: token.EQUAL, Literal: string(firstchar) + string(l.ch)}
			break
		}
		tok = newToken(token.ASSIGN, l.ch)
	case '+':
		if l.peekChar() == '+' {
			l.read()
			tok = newToken(token.INCREMENT, l.ch)
			break
		}
		tok = newToken(token.PLUS, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LEFT_BRACKET, l.ch)
	case ')':
		tok = newToken(token.RIGHT_BRACKET, l.ch)
	case '{':
		if l.peekChar() == '{' {
			l.read()
			tok = newToken(token.LEFT_OBJECT_BRACE, l.ch)
			break
		}
		tok = newToken(token.LEFT_BRACE, l.ch)
	case '}':
		if l.peekChar() == '}' {
			l.read()
			tok = newToken(token.RIGHT_OBJECT_BRACE, l.ch)
			break
		}
		tok = newToken(token.RIGHT_BRACE, l.ch)
	case '-':
		if l.peekChar() == '-' {
			l.read()
			tok = newToken(token.DECREMENT, l.ch)
			break
		}
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERIK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '!':
		if l.peekChar() == '=' {
			firstchar := l.ch
			l.read()
			tok = token.Token{Type: token.NOT_EQUAL, Literal: string(firstchar) + string(l.ch)}
			break
		}
		tok = newToken(token.BANG, l.ch)
	case '<':
		tok = newToken(token.LESS_THAN, l.ch)
	case '>':
		tok = newToken(token.GRTR_THAN, '>')
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '[':
		tok = newToken(token.LEFT_LARGE_BRACKET, '[')
	case ']':
		tok = newToken(token.RIGHT_LARGE_BRACKET, ']')
	case ':':
		tok = newToken(token.KEY_VAL_SEP, ':')
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default: //handling identifiers
		if l.isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.IdentOrKeyword(tok.Literal) //check if the given literal exists on keyword map
			return tok
		} else if l.isNumber(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INTEGER
			return tok
		} else {

			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.read()
	return tok
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.read()
	return l
}

//utilities

func (l *Lexer) read() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.lastRead = l.readPos
	l.readPos += 1
}

func (l *Lexer) readIdentifier() string {
	pos := l.lastRead
	for l.isLetter(l.ch) {
		l.read()
	}
	return l.input[pos:l.lastRead]
}

func (l *Lexer) readNumber() string {
	pos := l.lastRead
	for l.isNumber(l.ch) {
		l.read()
	}
	return l.input[pos:l.lastRead]
}
func (l *Lexer) readString() string {
	start := l.readPos
	for {
		l.read()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[start:l.lastRead]
}
func newToken(tt token.TokenType, ch byte) token.Token {
	return token.Token{Type: tt, Literal: string(ch)}
}

//currently only supporiting ASCII
func (l *Lexer) isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}
func (l *Lexer) isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.read()
	}
}

func (l *Lexer) skipComment() {
	l.read()
	for l.ch != '#' {
		if l.ch == 0 {
			return
		}
		l.read()
	}
	l.read()

}

//for two character token
func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]

}
