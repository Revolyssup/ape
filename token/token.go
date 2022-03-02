package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"for":    FOR,
	"return": RETURN,
}

const (
	//keywords
	LET      = "LET"
	FUNCTION = "FUNCTION"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	RETURN   = "RETURN"
	IF       = "IF"
	ELSE     = "ELSE"
	FOR      = "FOR"
	//Operators
	PLUS      = "+"
	MINUS     = "-"
	SLASH     = "/"
	ASTERIK   = "*"
	LESS_THAN = "<"
	GRTR_THAN = ">"
	BANG      = "!"
	ASSIGN    = "="
	EQUAL     = "=="
	NOT_EQUAL = "!="
	//delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LEFT_BRACKET        = "("
	RIGHT_BRACKET       = ")"
	LEFT_BRACE          = "{"
	RIGHT_BRACE         = "}"
	LEFT_OBJECT_BRACE   = "{{"
	RIGHT_OBJECT_BRACE  = "}}"
	LEFT_LARGE_BRACKET  = "["
	RIGHT_LARGE_BRACKET = "]"
	KEY_VAL_SEP         = ":"
	INCREMENT           = "++"
	DECREMENT           = "--"
	//identifier
	IDENTIFIER = "IDENT"

	//literal
	INTEGER = "INT"
	STRING  = "STRING"
	//special
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
)

//To check if given token is keyword or an identifier
func IdentOrKeyword(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
