package lexer

import (
	"testing"

	"github.com/Revolyssup/ape/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	let ten = 10;
	let add = fn(x, y) {
	x + y;
	};
	let result = add(five, ten);

	!-/*5;
	5 < 10 > 5==5!=67;
	"foobar"
	"foo bar"
	[1,]
	 `
	tests := []struct {
		Type    token.TokenType
		Literal string
	}{
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INTEGER, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INTEGER, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LEFT_BRACKET, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RIGHT_BRACKET, ")"},
		{token.LEFT_BRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RIGHT_BRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LEFT_BRACKET, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RIGHT_BRACKET, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERIK, "*"},
		{token.INTEGER, "5"},
		{token.SEMICOLON, ";"},
		{token.INTEGER, "5"},
		{token.LESS_THAN, "<"},
		{token.INTEGER, "10"},
		{token.GRTR_THAN, ">"},
		{token.INTEGER, "5"},
		{token.EQUAL, "=="},
		{token.INTEGER, "5"},
		{token.NOT_EQUAL, "!="},
		{token.INTEGER, "67"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.LEFT_LARGE_BRACKET, "["},
		{token.INTEGER, "1"},
		{token.COMMA, ","},
		{token.RIGHT_LARGE_BRACKET, "]"},

		{token.EOF, ""},
	}

	lex := New(input)

	for i, tt := range tests {
		tok := lex.NextToken()

		if tok.Type != tt.Type {
			t.Fatalf("test[%d]: Wrong Type Token. Expected %q--Got %q", i, tt.Type, tok.Type)

		}
		if tok.Literal != tt.Literal {
			t.Fatalf("test[%d]: Wrong Literal. Expected %q--Got %q", i, tt.Literal, tok.Literal)

		}

	}
}
