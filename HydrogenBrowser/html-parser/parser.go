package htmlparser

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Node struct {
	Name       string
	Tag        Tag
	Attributes map[string]string
	Content    string
	Children   []*Node
}

func (n *Node) String() string {
	return n.prettyPrint(0)
}

func (n *Node) prettyPrint(indent int) string {
	indentStr := strings.Repeat("  ", indent)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s<%s", indentStr, n.Name))

	if len(n.Attributes) > 0 {
		for key, value := range n.Attributes {
			sb.WriteString(fmt.Sprintf(" %s=\"%s\"", key, value))
		}
	}
	sb.WriteString(">\n")

	if n.Content != "" {
		sb.WriteString(fmt.Sprintf("%s  %s\n", indentStr, n.Content))
	}

	for _, child := range n.Children {
		sb.WriteString(child.prettyPrint(indent + 1)) // Increase indentation
	}

	sb.WriteString(fmt.Sprintf("%s</%s>\n", indentStr, n.Name))

	return sb.String()
}

const (
	CUSTOM_TAG Tag = iota
	TEXT
	HEADER
	BODY
	TITLE
	HEADING1
	NEXTID
	ANCHOR
	PARAGRAPH
	DL
	DT
	DD
)

type Tag int

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
	// Last token hasn't been consumed
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

	if t == EOF {
		return nil
	}

	if t == LT_SIGN {
		node := &Node{Attributes: make(map[string]string)}
		_, tagName := p.scanIgnoreWhitespace()
		node.Name = tagName
		node.Tag = p.parseKnownTag(tagName)

		for {
			t, _ := p.scanIgnoreWhitespace()
			if t == IDENT {
				p.unscan()
				key, value := p.parseAttr()
				node.Attributes[key] = value
			}

			if t == GT_SIGN {
				nt, _ := p.scanIgnoreWhitespace()
				if nt != CLOSING_TAG {
					if nt == LT_SIGN {
						p.unscan()
						child := p.Parse()
						node.Children = append(node.Children, child)
					} else {
						p.unscan()
						node.Content = p.parseContentWithSpace()
					}
				} else {
					p.scanIgnoreWhitespace() // Closing tag name
					p.scanIgnoreWhitespace() // Greater than sign
					return node
				}
			}

			if t == CLOSING_TAG {
				p.scanIgnoreWhitespace() // Closing tag name
				p.scanIgnoreWhitespace() // Greater than sign
				return node
			}

			if t == LT_SIGN {
				p.unscan()
				child := p.Parse()
				node.Children = append(node.Children, child)
				return node
			}

			if t == EOF {
				return node
			}
		}
	} else {
		return &Node{Name: "Text", Tag: TEXT, Attributes: nil, Content: lit, Children: nil}
	}
}

func (p *Parser) parseContentWithSpace() string {
	str := ""

	for {
		t, lit := p.scan()

		if t == CLOSING_TAG || t == LT_SIGN {
			if t == LT_SIGN {
				// TODO: How to handle not closed tag
			}

			p.unscan()
			break
		}

		str = str + lit
	}

	return str
}

func (p *Parser) parseKnownTag(ident string) Tag {
	t := CUSTOM_TAG

	switch strings.ToLower(ident) {
	case "header":
		t = HEADER
	case "body":
		t = BODY
	case "title":
		t = TITLE
	case "h1":
		t = HEADING1
	case "nextid":
		t = NEXTID
	case "a":
		t = ANCHOR
	case "p":
		t = PARAGRAPH
	case "dl":
		t = DL
	case "dt":
		t = DT
	case "dd":
		t = DD
	}

	return t
}

func (p *Parser) parseAttr() (string, string) {
	_, key := p.scanIgnoreWhitespace()
	_, _ = p.scanIgnoreWhitespace() // Equal sign
	_, _ = p.scanIgnoreWhitespace() // Quote
	_, value := p.scanIgnoreWhitespace()
	_, _ = p.scanIgnoreWhitespace() // Quote

	return key, value
}

// This parser does not support the whole HTML spec!
// just the provided HTML on the jam
func ParseHTML(filePath string) ([]*Node, error) {
	htmlFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer htmlFile.Close()

	nodes := []*Node{}
	parser := NewParser(htmlFile)

	for {
		node := parser.Parse()
		if node == nil {
			break
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}
