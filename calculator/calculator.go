package calculator

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func Calc(r io.Reader) float64 {
	s := newScanner(r)
	p := newParser(s)

	n := p.parse()
	return n.Result()
}

func CalcString(s string) float64 {
	return Calc(strings.NewReader(s))
}

// input
// 0-9, +, -, *, /, (, )

type token interface{}

type (
	tkNumber  float64
	tkOp1     string
	tkOp2     string
	tkParenOp struct{}
	tkParenEn struct{}
)

// text => token
type scanner struct {
	r *bufio.Scanner
}

func newScanner(r io.Reader) *scanner {
	s := &scanner{}

	scanner := bufio.NewScanner(r)
	scanner.Split(s.split)
	s.r = scanner

	return s
}

func (s *scanner) split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// "     1   +   2"
	{
		var (
			i int
			b byte
		)
		for i, b = range data {
			if b != ' ' {
				break
			}
		}

		advance = i
		data = data[i:]
	}

	if atEOF {
		if len(data) == 0 {
			return advance, nil, io.EOF
		}
		return advance + len(data), data, nil
	}

	for i, b := range data {
		if b == ' ' {
			return advance + i, data[:i], nil
		}
	}

	return advance, nil, nil
}

func (s *scanner) Scan() bool {
	return s.r.Scan()
}

func (s *scanner) Token() token {
	tk := s.r.Text()

	switch tk {
	case "*", "/":
		return tkOp1(tk)
	case "+", "-":
		return tkOp2(tk)
	case "(":
		return tkParenOp{}
	case ")":
		return tkParenEn{}
	}

	if f, err := strconv.ParseFloat(tk, 64); err == nil {
		return tkNumber(f)
	}

	panic("invalid")
}

type node interface {
	Result() float64
}

type nodeNumber float64

func (n nodeNumber) Result() float64 {
	return float64(n)
}

func (n nodeNumber) String() string {
	return fmt.Sprintf("%f", n)
}

type nodeOp struct {
	op string
	l  node
	r  node
}

func (n nodeOp) Result() float64 {
	switch n.op {
	case "+":
		return n.l.Result() + n.r.Result()
	case "-":
		return n.l.Result() - n.r.Result()
	case "*":
		return n.l.Result() * n.r.Result()
	case "/":
		return n.l.Result() / n.r.Result()
	}

	panic("unreachable")
}

func (n nodeOp) String() string {
	return fmt.Sprintf("(%s %s %s)",
		n.op,
		n.l,
		n.r,
	)
}

type parser struct {
	s *scanner
}

func newParser(s *scanner) *parser {
	return &parser{
		s: s,
	}
}

func (p *parser) parse() node {
	var prevNode node

	for p.s.Scan() {
		tk := p.s.Token()

		switch tk := tk.(type) {
		case tkNumber:
			prevNode = nodeNumber(tk)
		case tkOp1:
			prevNode = nodeOp{
				op: string(tk),
				l:  prevNode,
				r:  p.parseNumber(),
			}
		case tkOp2:
			return nodeOp{
				op: string(tk),
				l:  prevNode,
				r:  p.parse(),
			}
		case tkParenOp:
			prevNode = p.parse()
		case tkParenEn:
			return prevNode
		default:
			panic("unreachable")
		}
	}

	return prevNode
}

func (p *parser) parseNumber() node {
	if !p.s.Scan() {
		panic("invalid")
	}

	tk := p.s.Token()
	switch tk := tk.(type) {
	case tkNumber:
		return nodeNumber(tk)
	case tkParenOp:
		return p.parse()
	}

	panic("invalid")
}
