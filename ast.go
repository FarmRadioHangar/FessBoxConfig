package config

import (
	"errors"
	"fmt"
	"io"
)

type ast struct {
	*node
	sections []*nodeSection
}

type node struct {
	values []*nodeIdent
}

type nodeSection struct {
	name   string
	line   int
	values []*nodeIdent
}

type nodeIdent struct {
	key   string
	value string
	line  int
}

type parser struct {
	tokens  []*Token
	ast     *ast
	currPos int
}

func newParser(src io.Reader) (*parser, error) {
	s := NewScanner(src)
	var toks []*Token
	var err error
	var tok *Token
	for err == nil {
		tok, err = s.Scan()
		if err != nil {
			if err.Error() != io.EOF.Error() {
				return nil, err
			}
			break
		}
		if tok != nil {
			if tok.Type != Comment || tok.Type != WHiteSpace {
				toks = append(toks, tok)
			}
		}
	}
	return &parser{tokens: toks, ast: &ast{}}, nil
}

func (p *parser) parse() (*ast, error) {
	var err error
	if err != nil {
		return nil, err
	}
END:
	for {
		tok := p.next()
		if tok.Type == EOF {
			break END
		}
		switch tok.Type {
		case OpenBrace:
			p.rewind()
			err = p.parseSection()
			if err != nil {
				break END
			}
		}
	}
	return p.ast, err
}

func (p *parser) next() *Token {
	if p.currPos >= len(p.tokens)-1 {
		return &Token{Type: EOF}
	}
	t := p.tokens[p.currPos]
	p.currPos++
	return t
}

func (p *parser) seek(at int) {
	p.currPos = at
}

func (p *parser) parseSection() (err error) {
	left := p.next()
	if left.Type != OpenBrace {
		return errors.New("bad token")
	}
	ns := &nodeSection{}
	completeName := false
END:
	for {
	BEGIN:
		tok := p.next()
		if tok.Type == EOF {
			p.rewind()
			break END
		}
		if !completeName {
			switch tok.Type {
			case Ident:
				ns.name = ns.name + tok.Text
				goto BEGIN
			case ClosingBrace:
				completeName = true
				goto BEGIN
			}
		}
		switch tok.Type {
		case NewLine:
			n1 := p.next()
			if n1.Type == NewLine {
				n2 := p.next()
				if n2.Type == NewLine {
					fmt.Println("HERE")
					break END
				}
				p.rewind()
				goto BEGIN
			}
			p.rewind()
			goto BEGIN
		case Ident:
			p.rewind()
			err = p.parseIdent(ns)
			if err != nil {
				break END
			}
		default:
			break END
		}
	}
	if err == nil {
		p.ast.sections = append(p.ast.sections, ns)
	}
	return
}

func (p *parser) rewind() {
	p.currPos--
}

func (p *parser) parseIdent(sec *nodeSection) (err error) {
	fmt.Printf("parsing ident for %s -- ", sec.name)
	n := &nodeIdent{}
	doneKey := false
END:
	for {
	BEGIN:
		tok := p.next()
		if tok.Type == EOF {
			p.rewind()
			break END
		}

		if !doneKey {
			switch tok.Type {
			case Ident:
				n.key = n.key + tok.Text
				goto BEGIN
			case Operand:
				doneKey = true
				goto BEGIN
			default:
				err = errors.New("some fish")
				break END
			}

		}
		switch tok.Type {
		case Ident:
			n.value = n.value + tok.Text
			goto BEGIN
		case NewLine:
			break END
		default:
			err = errors.New("some fish")
			break END
		}
	}
	if err == nil {
		sec.values = append(sec.values, n)
	}
	fmt.Println("done")

	return
}
