//line parser.go.y:2

// Parser is a type to analyse scheme source's syntax.
// It embeds Lexer to generate tokens from a source code.
// Parser.Parse() does syntactic analysis and returns scheme object pointer.

package scheme

import __yyfmt__ "fmt"

//line parser.go.y:6
//line parser.go.y:9
type yySymType struct {
	yys     int
	objects []Object
	object  Object
	token   string
}

const IDENTIFIER = 57346
const NUMBER = 57347
const BOOLEAN = 57348
const STRING = 57349

var yyToknames = []string{
	"IDENTIFIER",
	"NUMBER",
	"BOOLEAN",
	"STRING",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line parser.go.y:78

type Parser struct {
	*Lexer
}

func NewParser(source string) *Parser {
	return &Parser{NewLexer(source)}
}

func (p *Parser) Parse(parent Object) []Object {
	p.ensureAvailability()
	if yyParse(p.Lexer) != 0 {
		panic("parse error")
	}

	for _, r := range p.results {
		r.setParent(parent)
	}
	return p.results
}

func (p *Parser) parseObject(parent Object) Object {
	tokenType := p.TokenType()
	token := p.NextToken()

	switch tokenType {
	case '(':
		return p.parseApplication(parent)
	case '\'':
		return p.parseSingleQuote(parent)
	case IntToken:
		return NewNumber(token, parent)
	case IdentifierToken:
		return NewVariable(token, parent)
	case BooleanToken:
		return NewBoolean(token, parent)
	case StringToken:
		return NewString(token[1:len(token)-1], parent)
	default:
		return nil
	}
	return nil
}

// This is for parsing syntax sugar '*** => (quote ***)
func (p *Parser) parseSingleQuote(parent Object) Object {
	if len(p.PeekToken()) == 0 {
		runtimeError("unterminated quote")
	}
	application := NewApplication(parent)
	application.procedure = NewVariable("quote", application)
	application.arguments = NewList(application, p.parseObject(application))
	return application
}

// This function returns *Pair of first object and list from second.
// Scanner position ends with the next of close parentheses.
func (p *Parser) parseList(parent Object) Object {
	pair := NewPair(parent)
	pair.Car = p.parseObject(pair)
	if pair.Car == nil {
		return pair
	}
	pair.Cdr = p.parseList(pair).(*Pair)
	return pair
}

func (p *Parser) parseApplication(parent Object) Object {
	if p.PeekToken() == ")" {
		p.NextToken()
		return Null
	}
	application := NewApplication(parent)
	application.procedure = p.parseObject(application)
	application.arguments = p.parseList(application)

	return application
}

func (p *Parser) parseQuotedObject(parent Object) Object {
	tokenType := p.TokenType()
	token := p.NextToken()

	switch tokenType {
	case '(':
		return p.parseQuotedList(parent)
	case '\'':
		return p.parseSingleQuote(parent)
	case IntToken:
		return NewNumber(token, parent)
	case IdentifierToken:
		return NewSymbol(token)
	case BooleanToken:
		return NewBoolean(token, parent)
	case ')':
		return nil
	default:
		runtimeError("unterminated quote")
		return nil
	}
}

func (p *Parser) parseQuotedList(parent Object) Object {
	pair := NewPair(parent)
	pair.Car = p.parseQuotedObject(pair)
	if pair.Car == nil {
		return pair
	}
	pair.Cdr = p.parseQuotedList(pair).(*Pair)
	return pair
}

func (p *Parser) ensureAvailability() {
	// Error message will be printed by interpreter
	recover()
}

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 13
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 23

var yyAct = []int{

	4, 7, 8, 9, 5, 6, 12, 13, 4, 7,
	8, 9, 5, 6, 14, 15, 2, 3, 1, 0,
	10, 11, 16,
}
var yyPact = []int{

	-1000, 4, -1000, -1000, -1000, 4, -4, -1000, -1000, -1000,
	-1000, 4, -1000, 5, 4, -1000, -1000,
}
var yyPgo = []int{

	0, 18, 7, 14, 17,
}
var yyR1 = []int{

	0, 1, 1, 2, 2, 3, 3, 3, 3, 4,
	4, 4, 4,
}
var yyR2 = []int{

	0, 0, 2, 0, 2, 1, 1, 2, 4, 1,
	1, 1, 2,
}
var yyChk = []int{

	-1000, -1, -3, -4, 4, 8, 9, 5, 6, 7,
	-3, -3, 10, -2, -3, 10, -2,
}
var yyDef = []int{

	1, -2, 2, 5, 6, 0, 0, 9, 10, 11,
	7, 3, 12, 0, 3, 8, 4,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 8,
	9, 10,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(yyToknames) {
		if yyToknames[c-4] != "" {
			return yyToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(c), uint(char))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		//line parser.go.y:28
		{
			yyVAL.objects = []Object{}
		}
	case 2:
		//line parser.go.y:32
		{
			yyVAL.objects = append(yyS[yypt-1].objects, yyS[yypt-0].object)
			if l, ok := yylex.(*Lexer); ok {
				l.results = yyVAL.objects
			}
		}
	case 3:
		//line parser.go.y:40
		{
			yyVAL.object = Null
		}
	case 4:
		//line parser.go.y:42
		{
			pair := NewPair(nil)
			pair.Car = yyS[yypt-1].object
			pair.Car.setParent(pair)
			pair.Cdr = yyS[yypt-0].object
			pair.Cdr.setParent(pair)
			yyVAL.object = pair
		}
	case 5:
		//line parser.go.y:53
		{
			yyVAL.object = yyS[yypt-0].object
		}
	case 6:
		//line parser.go.y:55
		{
			yyVAL.object = NewVariable(yyS[yypt-0].token, nil)
		}
	case 7:
		//line parser.go.y:57
		{
			yyVAL.object = yyS[yypt-0].object
		}
	case 8:
		//line parser.go.y:59
		{
			app := NewApplication(nil)
			app.procedure = yyS[yypt-2].object
			app.procedure.setParent(app)
			app.arguments = yyS[yypt-1].object
			app.arguments.setParent(app)
			yyVAL.object = app
		}
	case 9:
		//line parser.go.y:70
		{
			yyVAL.object = NewNumber(yyS[yypt-0].token)
		}
	case 10:
		//line parser.go.y:72
		{
			yyVAL.object = NewBoolean(yyS[yypt-0].token)
		}
	case 11:
		//line parser.go.y:74
		{
			yyVAL.object = NewString(yyS[yypt-0].token[1 : len(yyS[yypt-0].token)-1])
		}
	case 12:
		//line parser.go.y:76
		{
			yyVAL.object = Null
		}
	}
	goto yystack /* stack new state and value */
}
