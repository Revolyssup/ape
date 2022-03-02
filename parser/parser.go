package parser

import (
	"fmt"
	"strconv"

	"github.com/Revolyssup/ape/ast"
	"github.com/Revolyssup/ape/lexer"
	"github.com/Revolyssup/ape/token"
)

//To parse expressions using pratt parsing. Based on what type of expression that is, we will define different functins. Broadly they will be in two categories:
type (
	infixParsefunc  func(ast.Expression) ast.Expression //It takes in the expression before the operator/token
	prefixParsefunc func() ast.Expression
)

type Parser struct {
	l         *lexer.Lexer
	currToken token.Token
	peekToken token.Token
	errors    []string
	//Each token type will have some parse function associated with it.
	infixParsefuncns  map[token.TokenType]infixParsefunc
	prefixParsefuncns map[token.TokenType]prefixParsefunc
}

// These are the precedence of operators which would be passed in function call to specific parseExpression functions.
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // ><
	SUMSUB      // +
	PRODUCT     // *
	PREFIX      // -X and !X
	CALL        // func(x)
	INDEX
)

//mapping each token to its appropriate precedence
var precedence = map[token.TokenType]int{
	token.EQUAL:              EQUALS,
	token.NOT_EQUAL:          EQUALS,
	token.LESS_THAN:          LESSGREATER,
	token.GRTR_THAN:          LESSGREATER,
	token.MINUS:              SUMSUB,
	token.PLUS:               SUMSUB,
	token.SLASH:              PRODUCT,
	token.ASTERIK:            PRODUCT,
	token.BANG:               PREFIX,
	token.RIGHT_BRACKET:      LOWEST,
	token.LEFT_BRACKET:       CALL,
	token.LEFT_LARGE_BRACKET: INDEX,
	token.LEFT_OBJECT_BRACE:  INDEX,
}

//functins to compare precedences of tokens
func (p *Parser) currPrecedence() int {
	if p, ok := precedence[p.currToken.Type]; ok {
		return p
	}
	return LOWEST
}
func (p *Parser) peekPrecedence() int {
	if p, ok := precedence[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

//Parsing expressions
func (p *Parser) parseExpression(precedence int) ast.Expression {

	prefix := p.prefixParsefuncns[p.currToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}

	leftExp := prefix()
	// fmt.Println("LEFT EXP IS: " + leftExp.String())
	for p.peekToken.Type != token.SEMICOLON && p.peekPrecedence() > precedence {
		// fmt.Println("Coming in loop cuz peektoken is : " + p.peekToken.Literal)
		infix := p.infixParsefuncns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.NextToken()
		leftExp = infix(leftExp)
		// fmt.Println("LEFT EXP after innfix IS: " + leftExp.String())
	}
	return leftExp
}
func (p *Parser) parsePrefixExpression() ast.Expression {
	pexp := &ast.PrefixExpression{Token: p.currToken, Operator: p.currToken.Literal}

	p.NextToken()
	pexp.RightExpression = p.parseExpression(PREFIX)
	return pexp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	iexp := &ast.InfixExpression{Token: p.currToken, LeftExpression: left, Operator: p.currToken.Literal}
	precedence := p.currPrecedence()
	p.NextToken()
	iexp.RightExpression = p.parseExpression(precedence)
	return iexp
}

//Creating instance of the parser.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.NextToken()
	p.NextToken()
	p.prefixParsefuncns = make(map[token.TokenType]prefixParsefunc)
	p.infixParsefuncns = make(map[token.TokenType]infixParsefunc)
	//registering parseExpressinoFunctions

	p.registerPrefixParse(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefixParse(token.INTEGER, p.parseIntegerLiteral)
	p.registerPrefixParse(token.TRUE, p.parseBoolean)
	p.registerPrefixParse(token.FALSE, p.parseBoolean)
	p.registerPrefixParse(token.BANG, p.parsePrefixExpression)
	p.registerPrefixParse(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixParse(token.LEFT_BRACKET, p.parseGroupedExpression)
	p.registerPrefixParse(token.IF, p.parseIfExpression)
	p.registerPrefixParse(token.FOR, p.parseForExpression)
	p.registerPrefixParse(token.FUNCTION, p.parseFunctionLiterals)
	p.registerPrefixParse(token.STRING, p.parseStringLiteral)
	p.registerPrefixParse(token.LEFT_LARGE_BRACKET, p.parseArray)
	p.registerPrefixParse(token.LEFT_OBJECT_BRACE, p.parseObject)

	p.registerInfixParse(token.PLUS, p.parseInfixExpression)
	p.registerInfixParse(token.MINUS, p.parseInfixExpression)
	p.registerInfixParse(token.SLASH, p.parseInfixExpression)
	p.registerInfixParse(token.ASTERIK, p.parseInfixExpression)
	p.registerInfixParse(token.EQUAL, p.parseInfixExpression)
	p.registerInfixParse(token.NOT_EQUAL, p.parseInfixExpression)
	p.registerInfixParse(token.LESS_THAN, p.parseInfixExpression)
	p.registerInfixParse(token.GRTR_THAN, p.parseInfixExpression)
	p.registerInfixParse(token.LEFT_BRACKET, p.parseFunctionCall)
	p.registerInfixParse(token.LEFT_LARGE_BRACKET, p.parseArrObjElement)
	p.registerInfixParse(token.LEFT_OBJECT_BRACE, p.parseArrObjElement)
	return p
}

func (p *Parser) NextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil { //we will get some sort of parsed statement
			program.Statements = append(program.Statements, stmt)
		}
		p.NextToken()
	}
	return program
}
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekErrors(t token.TokenType) {
	msg := fmt.Sprintf("Expected token type %s. Got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		{
			return p.parseLetStatement()
		}
	case token.RETURN:
		{
			return p.parseReturnStatement()
		}

	default:
		{
			return p.parseExpressionStatement()
		}
	}
}

//parsing different types of statements.

func (p *Parser) parseLetStatement() *ast.LetStatement {
	letstmt := &ast.LetStatement{Token: p.currToken}
	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	letstmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.NextToken()
	letstmt.Value = p.parseExpression(LOWEST)
	for p.peekToken.Type == token.SEMICOLON {
		p.NextToken()
	}

	return letstmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	retstmt := &ast.ReturnStatement{Token: p.currToken}
	p.NextToken()
	retstmt.ReturnValue = p.parseExpression(LOWEST)
	for p.peekToken.Type == token.SEMICOLON {
		p.NextToken()
	}
	return retstmt
}

//Parsing expressionns using pratt parser technique.
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON && p.currToken.Type != token.EOF { //Semicolon is not mandatory
		p.NextToken()
	}
	return stmt
}

//utilities
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t { //good to go
		p.NextToken()
		return true
	}
	p.peekErrors(t)
	return false
}

func (p *Parser) expectCurr(t token.TokenType) bool {
	if p.currToken.Type == t { //good to go
		p.NextToken()
		return true
	}
	return false
}

func (p *Parser) registerPrefixParse(t token.TokenType, f prefixParsefunc) {
	p.prefixParsefuncns[t] = f
}

func (p *Parser) registerInfixParse(t token.TokenType, f infixParsefunc) {
	p.infixParsefuncns[t] = f
}

// different types of parseExpressionfunc based on token type

func (p *Parser) parseIdentifier() ast.Expression { //For token.IDENT
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	intexp := &ast.IntegerLiteral{Token: p.currToken}
	val, err := strconv.ParseInt(p.currToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as int64", intexp)
		p.errors = append(p.errors, msg)
		return nil
	}
	intexp.Value = val
	return intexp
}

func (p *Parser) parseStringLiteral() ast.Expression {
	stringexp := &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Literal}
	return stringexp
}
func (p *Parser) parseBoolean() ast.Expression {
	boolexp := &ast.Boolean{Token: p.currToken}

	val, err := strconv.ParseBool(p.currToken.Literal)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as nool", boolexp)
		p.errors = append(p.errors, msg)
		return nil
	}
	boolexp.Value = val
	return boolexp
}

func (p *Parser) parseArray() ast.Expression { //Enter with currtoken set as '['
	arr := &ast.ArrayLiteral{Token: p.currToken}
	exp := []ast.Expression{}
	p.NextToken()
	for p.currToken.Type != token.RIGHT_LARGE_BRACKET && p.currToken.Type != token.EOF {
		tempExp := p.parseExpression(LOWEST)

		exp = append(exp, tempExp)
		if p.peekToken.Type != token.COMMA {
			if p.peekToken.Type == token.RIGHT_LARGE_BRACKET {
				p.NextToken()
				arr.Value = exp
				return arr
			}
			p.errors = append(p.errors, "No comma after element in array.")
			return arr
		}
		p.NextToken()
		p.NextToken()
	}

	p.NextToken()
	arr.Value = exp
	return arr
}

func (p *Parser) parseArrObjElement(id ast.Expression) ast.Expression {
	arrele := &ast.ArrObjElement{Token: p.currToken, Name: id}
	p.NextToken()
	arrele.Index = p.parseExpression(LOWEST)
	if p.peekToken.Type != token.RIGHT_LARGE_BRACKET {
		return nil
	}
	p.NextToken()
	return arrele
}
func (p *Parser) parseObject() ast.Expression { //Enter with currtoken set as '{'
	obj := &ast.ObjectLiteral{Token: p.currToken}
	exp := map[ast.Expression]ast.Expression{}
	p.NextToken()
	for p.currToken.Type != token.RIGHT_OBJECT_BRACE && p.currToken.Type != token.EOF {
		keyExp := p.parseExpression(LOWEST)
		p.NextToken()
		if p.currToken.Type != token.KEY_VAL_SEP {
			p.errors = append(p.errors, "No seperator found between key-values")
			return obj
		}
		p.NextToken()
		valueExp := p.parseExpression(LOWEST)
		exp[keyExp] = valueExp
		if p.peekToken.Type != token.COMMA {
			if p.peekToken.Type == token.RIGHT_OBJECT_BRACE {
				p.NextToken()
				obj.Value = exp
				return obj
			}
			p.errors = append(p.errors, "No comma after element in object, found "+p.peekToken.Literal)
			return obj
		}
		p.NextToken()
		p.NextToken()
	}

	p.NextToken()
	obj.Value = exp
	return obj
}

//For parenthesis(grouped expressions)

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.NextToken()

	//This will go on recursively parsing the expression untill just before the right parenthesis for this parent expression is encountered.
	exp := p.parseExpression(LOWEST)
	if p.peekToken.Type != token.RIGHT_BRACKET {
		return nil
	}
	p.NextToken()
	return exp

}

//Parsing Blocks
func (p *Parser) parseBlockStatements() *ast.BlockStatement { //Will enter with currToken at `{`
	bs := &ast.BlockStatement{Token: p.currToken}
	bs.Stmts = []ast.Statement{}

	p.NextToken()

	for p.currToken.Type != token.RIGHT_BRACE && p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			bs.Stmts = append(bs.Stmts, stmt)
		}
		p.NextToken()
	}
	return bs //Exit with currToken either `}` or file ends
}

//IF_ELSE are expressions in monkey as they produce a value. Hence, if (x>3) 2; is equivalent to if (x>3) return 2;
func (p *Parser) parseIfExpression() ast.Expression {
	ife := &ast.IfExpression{Token: p.currToken}
	if p.peekToken.Type != token.LEFT_BRACKET {
		return nil
	}
	p.NextToken()

	ife.Condition = p.parseExpression(LOWEST)
	if p.peekToken.Type != token.LEFT_BRACE {
		return nil
	}

	p.NextToken()
	ife.MainStmt = p.parseBlockStatements()
	//In case there is an else statement

	if p.peekToken.Type == token.ELSE {
		p.NextToken()
		if p.peekToken.Type != token.LEFT_BRACE {
			return nil
		}
		p.NextToken()
		ife.AltStmt = p.parseBlockStatements()
	}
	return ife
}

//Parsing For expressions-Looks exactly like If expressions
func (p *Parser) parseForExpression() ast.Expression {
	fore := &ast.ForExpression{Token: p.currToken}
	if p.peekToken.Type != token.LEFT_BRACKET {
		return nil
	}
	p.NextToken()

	fore.Condition = p.parseExpression(LOWEST)
	if p.peekToken.Type != token.LEFT_BRACE {
		return nil
	}
	p.NextToken()
	fore.Stmt = p.parseBlockStatements()
	return fore
}

//Parsing functino literals.Function declarations in go are just like expressions. fn(..params){body}
func (p *Parser) parseFunctionLiterals() ast.Expression {
	fl := &ast.FunctionLiteral{Token: p.currToken}
	if p.peekToken.Type != token.LEFT_BRACKET {
		return nil
	}
	p.NextToken()
	fl.Params = p.parseParameters()
	p.NextToken()
	fl.Body = p.parseBlockStatements()
	return fl
}

//Parsing all the parameters inside of function declaration.
func (p *Parser) parseParameters() []*ast.Identifier { //Current token will be  `(` when we enter this function
	params := []*ast.Identifier{}

	if p.peekToken.Type == token.RIGHT_BRACKET {
		p.NextToken()
		return params
	}
	p.NextToken()
	param := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	params = append(params, param)
	for p.peekToken.Type == token.COMMA {
		p.NextToken() //Will go to next comma
		p.NextToken() //Will reach to next param
		param := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		params = append(params, param)
	}
	if p.peekToken.Type != token.RIGHT_BRACKET {
		return nil
	}
	p.NextToken()
	return params //Leave the function with currToken `)`
}

//Parsing function calls.
//Function calls- <expression>(args). expression can be either an identifier pointing to a function literal or a function literal itself.And args is also expression.
//We can have nested funciton literals inside of out function call as arguments.

func (p *Parser) parseFunctionCall(function ast.Expression) ast.Expression { //While entering: currtoken would be `(` before the args
	fc := &ast.FunctionCall{Token: p.currToken, Function: function}
	fc.Arguments = p.parseArgs()
	return fc
}

//Returning all expressions inside function call
func (p *Parser) parseArgs() []ast.Expression {
	args := []ast.Expression{}
	if p.peekToken.Type == token.RIGHT_BRACKET {
		p.NextToken()
		return args
	}
	p.NextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekToken.Type == token.COMMA {
		p.NextToken()
		p.NextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if p.peekToken.Type != token.RIGHT_BRACKET {
		return nil
	}
	p.NextToken() //Leaves at RIGHT BRACKER
	return args
}
