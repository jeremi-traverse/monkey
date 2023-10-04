package lexer

import (
	"github.com/jeremi-traverse/monkey/token"
)

// Making the current char a byte makes our lexer supports only ASCII characters
type Lexer struct {
	input           string // source code
	currentPosition int    // current position in input (points to current char)
	nextPosition    int    // current reading position in input (position + 1, next char)
	currentChar     byte   // current char being examined
}

func New(input string) *Lexer {
	l := &Lexer{input: input}

	// Read the first character of the lexer's input
	l.readChar()
	return l
}

// consume the current char
func (l *Lexer) readChar() {
	if l.nextPosition >= len(l.input) {
		l.currentChar = 0
	} else {
		l.currentChar = l.input[l.nextPosition]
	}

	// Next char to examine becomes the current char
	l.currentPosition = l.nextPosition
	// Move the reading pointer by 1
	l.nextPosition += 1
}

// Returns a token containing information about the current char
// and advance the lexer to the next char
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	switch l.currentChar {
	case '=':
		if l.peekChar() == '=' {
			currentCharSnapshot := l.currentChar
			// Consume next character
			l.readChar()
			litteral := string(currentCharSnapshot) + string(l.currentChar)
			tok = token.Token{Type: token.EQ, Literal: litteral}
		} else {
			tok = newToken(token.ASSIGN, l.currentChar)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.currentChar)
	case '(':
		tok = newToken(token.LPAREN, l.currentChar)
	case ')':
		tok = newToken(token.RPAREN, l.currentChar)
	case '{':
		tok = newToken(token.LBRACE, l.currentChar)
	case '}':
		tok = newToken(token.RBRACE, l.currentChar)
	case ',':
		tok = newToken(token.COMMA, l.currentChar)
	case '+':
		tok = newToken(token.PLUS, l.currentChar)
	case '-':
		tok = newToken(token.MINUS, l.currentChar)
	case '!':
		if l.peekChar() == '=' {
			currentCharSnapshot := l.currentChar
			// Consume next character
			l.readChar()
			litteral := string(currentCharSnapshot) + string(l.currentChar)
			tok = token.Token{Type: token.NOT_EQ, Literal: litteral}
		} else {
			tok = newToken(token.BANG, l.currentChar)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.currentChar)
	case '/':
		tok = newToken(token.SLASH, l.currentChar)
	case '<':
		tok = newToken(token.LT, l.currentChar)
	case '>':
		tok = newToken(token.GT, l.currentChar)
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(l.currentChar) {
			tok.Literal = l.readIndentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			// early returns to skip the readChar after the switch statement
			return tok
		} else if isDigit(l.currentChar) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			// early returns to skip the readChar after the switch statement
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.currentChar)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, literal byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(literal)}
}

func (l *Lexer) readIndentifier() string {
	initialPosition := l.currentPosition

	for isLetter(l.currentChar) {
		l.readChar()
	}

	return l.input[initialPosition:l.currentPosition]
}

func (l *Lexer) readNumber() string {
	initialPosition := l.currentPosition

	for isDigit(l.currentChar) {
		l.readChar()
	}

	return l.input[initialPosition:l.currentPosition]
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextPosition]
	}
}

func isLetter(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func (l *Lexer) skipWhiteSpace() {
	for l.currentChar == ' ' || l.currentChar == '\t' || l.currentChar == '\n' || l.currentChar == '\r' {
		l.readChar()
	}
}
