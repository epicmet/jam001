package htmlparser

import (
	"bufio"
	"bytes"
	"io"
)

const (
	// Special tokens
	UNKNOWN Token = iota
	EOF
	WS

	// Characters
	DOT           // .
	SLASH         // /
	COMMA         // ,
	DASH          // -
	QUESTION_MARK // ?
	HASH_SIGN     // #
	EQUAL_SIGN    // =
	DOUBLE_QUOTE  // "
	SINGLE_QUOTE  // '
	LT_SIGN       // <
	GT_SIGN       // >
	CLOSING_TAG   // </
	OPENING_PAR   // (
	CLOSING_PAR   // )

	// Literals
	IDENT
	DIGIT
)

type Token int

const eof = rune(0)

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

type Lexer struct {
	r *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{r: bufio.NewReader(reader)}
}

func (l *Lexer) read() rune {
	ch, _, err := l.r.ReadRune()
	if err != nil {
		return eof
	}

	return ch
}

func (l *Lexer) unread() {
	_ = l.r.UnreadRune()
}

func (l *Lexer) Lex() (Token, string) {
	ch := l.read()

	if isWhitespace(ch) {
		l.unread()
		return l.lexWhiteSpace()
	} else if isLetter(ch) {
		l.unread()
		return l.lexIdent()
	} else if isDigit(ch) {
		l.unread()
		return l.lexDigit()
	}

	sch := string(ch)

	switch ch {
	case eof:
		return EOF, ""
	case '.':
		return DOT, sch
	case ',':
		return COMMA, sch
	case '/':
		return SLASH, sch
	case '<':
		nch := l.read()
		if nch == '/' {
			return CLOSING_TAG, sch + string(nch)
		} else {
			l.unread()
			return LT_SIGN, sch
		}
	case '>':
		return GT_SIGN, sch
	case '=':
		return EQUAL_SIGN, sch
	case '"':
		return DOUBLE_QUOTE, sch
	case '\'':
		return SINGLE_QUOTE, sch
	case '(':
		return OPENING_PAR, sch
	case ')':
		return CLOSING_PAR, sch
	case '-':
		return DASH, sch
	case '?':
		return QUESTION_MARK, sch
	case '#':
		return HASH_SIGN, sch
	}

	return UNKNOWN, sch
}

func (s *Lexer) lexDigit() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isDigit(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return DIGIT, buf.String()
}

func (s *Lexer) lexIdent() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	ident := buf.String()
	t := IDENT

	return t, ident
}

func (s *Lexer) lexWhiteSpace() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}
