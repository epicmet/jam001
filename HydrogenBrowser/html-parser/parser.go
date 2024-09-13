package htmlparser

import (
	"fmt"
	"io"
	"os"
)

type Node struct {
	Name       string
	Attributes map[string]string
	Content    string
	Children   []*Node
}

type Parser struct {
	l         *Lexer
	lastToken Token
	lastLit   string
	size      int // 0 or 1
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{l: NewLexer(reader)}
}

func (p *Parser) scan() (Token, string) {
	// Last token havent been consumed
	if p.size != 0 {
		p.size = 0
		return p.lastToken, p.lastLit
	}

	t, lit := p.l.Lex()
	p.lastToken, p.lastLit = t, lit

	return t, lit
}

func (p *Parser) scanIgnoreWhitespace() (Token, string) {
	t, lit := p.scan()
	if t == WS {
		t, lit = p.scan()
	}

	return t, lit
}

func (p *Parser) unscan() { p.size = 1 }

func (p *Parser) Parse() *Node {
	t, lit := p.scanIgnoreWhitespace()

	if t == OPENING_TAG {
		p.unscan()
		return p.parseOpeningTag()
	} else {
		// TODO: Extract & Centrelize name (From tokeinzer maybe?)
		return &Node{Name: "Text", Attributes: nil, Content: lit, Children: nil}
	}
}

func (p *Parser) parseOpeningTag() *Node {
	openTagToken, _ := p.scanIgnoreWhitespace()
	openNameToken, tagName := p.scanIgnoreWhitespace()

	for {
		t, lit := p.scanIgnoreWhitespace()
		// />

		// foo="dfd"
		// ></
		// >foo</
		// ><...></
		if t == SLASH {
			closeNameToken, _ := p.scanIgnoreWhitespace()
			closeTagToken, _ := p.scanIgnoreWhitespace()
		}
	}
}

// This parser does not support the whole HTML spec!
// just the provided HTML on the jam
func ParseHTML(filePath string) (*[]Node, error) {
	htmlFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer htmlFile.Close()

	// parser := NewParser(htmlFile)

	var ts []Token // token stack
	lexer := NewLexer(htmlFile)
	for {
		t, lit := lexer.Lex()

		if t == UNKNOWN {
			fmt.Printf("Unknown Token => Lit: %v\n", lit)
		}

		if t == EOF {
			break
		}

		if t == OPENING_TAG {
			ts = append(ts, t)
		}
	}

	return nil, nil
}
