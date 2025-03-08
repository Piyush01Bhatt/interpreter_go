package scanner

import (
	"log"
	"strconv"
	"unicode"

	u "github.com/Piyush01Bhatt/interpreter_go/internal/utils"
)

type TokenType int

const (
	// Single-character tokens.
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	// End of File
	EOF
)

// Token represents a scanned token.
type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

// LexScanner represents a scanner to scan tokens.
type LexScanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func NewLexScanner(input string) *LexScanner {
	return &LexScanner{
		source:  input,
		tokens:  make([]Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

// Token type names (for debugging/logging).
var tokenTypeNames = [...]string{
	"LEFT_PAREN", "RIGHT_PAREN", "LEFT_BRACE", "RIGHT_BRACE",
	"COMMA", "DOT", "MINUS", "PLUS", "SEMICOLON", "SLASH", "STAR",
	"BANG", "BANG_EQUAL", "EQUAL", "EQUAL_EQUAL",
	"GREATER", "GREATER_EQUAL", "LESS", "LESS_EQUAL",
	"IDENTIFIER", "STRING", "NUMBER",
	"AND", "CLASS", "ELSE", "FALSE", "FUN", "FOR", "IF", "NIL", "OR",
	"PRINT", "RETURN", "SUPER", "THIS", "TRUE", "VAR", "WHILE",
	"EOF",
}

// Keywords map for fast lookup.
var keywordsMap = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"fun":    FUN,
	"for":    FOR,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

// String method for debugging.
func (t TokenType) String() string {
	if int(t) < len(tokenTypeNames) {
		return tokenTypeNames[t]
	}
	return "UNKNOWN"
}

func (t *Token) String() string {
	literalStr := "<nil>"
	if str, ok := t.Literal.(string); ok {
		literalStr = str
	}
	return t.Type.String() + " " + t.Lexeme + " " + literalStr
}

func (ls *LexScanner) ScanTokens() []Token {
	for !ls.isAtEnd() {
		// We are at the beginning of the next lexeme.
		ls.start = ls.current
		ls.scan()
	}

	ls.addToken(EOF, nil)
	return ls.tokens
}

func (ls *LexScanner) isAtEnd() bool {
	return ls.current >= len(ls.source)
}

func (ls *LexScanner) scan() {
	ch := ls.advance()
	switch ch {
	case '(':
		ls.addToken(LEFT_PAREN, nil)
	case ')':
		ls.addToken(RIGHT_PAREN, nil)
	case '{':
		ls.addToken(LEFT_BRACE, nil)
	case '}':
		ls.addToken(RIGHT_BRACE, nil)
	case ',':
		ls.addToken(COMMA, nil)
	case '.':
		ls.addToken(DOT, nil)
	case '-':
		ls.addToken(MINUS, nil)
	case '+':
		ls.addToken(PLUS, nil)
	case ';':
		ls.addToken(SEMICOLON, nil)
	case '*':
		ls.addToken(STAR, nil)
	case '!':
		ls.addToken(u.Ternary(ls.match('='), BANG_EQUAL, BANG), nil)
	case '=':
		ls.addToken(u.Ternary(ls.match('='), EQUAL_EQUAL, EQUAL), nil)
	case '>':
		ls.addToken(u.Ternary(ls.match('='), GREATER_EQUAL, GREATER), nil)
	case '<':
		ls.addToken(u.Ternary(ls.match('='), LESS_EQUAL, LESS), nil)
	case '/':
		if ls.match('=') {
			for ls.peek() != '\n' && !ls.isAtEnd() {
				ls.advance()
			}
		} else {
			ls.addToken(SLASH, nil)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace
	case '\n':
		ls.line++
	case '"':
		ls.readString()
	default:
		if unicode.IsDigit(rune(ch)) {
			ls.readNumber()
			return
		}
		if unicode.IsLetter(rune(ch)) {
			ls.readIdentifier()
			return
		}
		log.Fatalf("unexpected character at line: %d", ls.line)
	}
}

func (ls *LexScanner) advance() byte {
	char := ls.source[ls.current]
	ls.current++
	return char
}

func (ls *LexScanner) match(expected byte) bool {
	if ls.isAtEnd() {
		return false
	}
	if ls.source[ls.current] != expected {
		return false
	}
	ls.current++
	return true
}

func (ls *LexScanner) addToken(tokenType TokenType, literal any) {
	lexeme := ls.source[ls.start:ls.current]
	token := Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    ls.line,
	}
	ls.tokens = append(ls.tokens, token)
}

func (ls *LexScanner) peek() byte {
	if ls.isAtEnd() {
		return '0'
	}
	return ls.source[ls.current]
}

func (ls *LexScanner) peekNext() byte {
	next := ls.current + 1
	if next >= len(ls.source) {
		return '0'
	}
	return ls.source[next]
}

func (ls *LexScanner) readString() {
	for ls.peek() != '"' && !ls.isAtEnd() {
		if ls.peek() == '\n' {
			ls.line++
		}
		ls.advance()
	}
	if ls.isAtEnd() {
		log.Fatalf("Unterminated string at line: %d", ls.line)
		return
	}
	ls.advance()
	value := ls.source[ls.start+1 : ls.current-1]
	ls.addToken(STRING, value)
}

func (ls *LexScanner) readNumber() {
	for unicode.IsDigit(rune(ls.peek())) {
		ls.advance()
	}
	if ls.peek() == '.' && unicode.IsDigit(rune(ls.peekNext())) {
		ls.advance()
		for unicode.IsDigit(rune(ls.peek())) {
			ls.advance()
		}
	}
	lexeme := ls.source[ls.start:ls.current]
	value, _ := strconv.ParseFloat(lexeme, 64)
	ls.addToken(NUMBER, value)
}

func (ls *LexScanner) readIdentifier() {
	for unicode.IsLetter(rune(ls.peek())) || unicode.IsDigit(rune(ls.peek())) || ls.peek() == '_' {
		ls.advance()
	}
	val := ls.source[ls.start:ls.current]
	keyword, exists := keywordsMap[val]
	if exists {
		ls.addToken(keyword, val)
		return
	}
	ls.addToken(IDENTIFIER, val)
}
