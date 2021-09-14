package cssvalue

// Term represents a term in a CSS value grammar.
type Term interface {
	isTerm()
}

// Component Value Types

// Keyword represents keyword values (such as auto, disc, etc.), which appear literally,
// without quotes (e.g. auto).
type Keyword string

func (Keyword) isTerm() {}

// Range represents a CSS bracketed range.
type Range struct {
	Min float64
	Max float64
}

// BasicType represents CSS basic data types, which appear between < and > (e.g.,
// <length>, <percentage>, etc.). For numeric data types, this type notation can annotate
// any range restrictions using the bracketed range notation described below.
type BasicType struct {
	Name  string
	Range *Range
}

func (*BasicType) isTerm() {}

// PropertyType represents types that have the same range of values as a property bearing
// the same name (e.g., <'border-width'>, <'background-attachment'>, etc.)
type PropertyType string

func (PropertyType) isTerm() {}

// NonTerminal represents non-terminals that do not share the same name as a property.
type NonTerminal string

func (NonTerminal) isTerm() {}

// Literal represents a literal character.
type Literal string

func (Literal) isTerm() {}

// Component Value Combinators

// Seq represents a juxtaposition of components that means that all of them must occur, in
// the given order.
type Seq struct {
	Name     string
	Operands []Term
}

func (*Seq) isTerm() {}

// AllOf represents two or more components, all of which must occur, in any order.
type AllOf struct {
	Operands []Term
}

func (*AllOf) isTerm() {}

// AnyOf represents two or more options: one or more of them must occur, in any order.
type AnyOf struct {
	Operands []Term
}

func (*AnyOf) isTerm() {}

// OneOf represents two or more alternatives: exactly one of them must occur.
type OneOf struct {
	Operands []Term
}

func (*OneOf) isTerm() {}

// Group represents a group of components. If Required is true, the group must produce at
// least one value.
type Group struct {
	Name     string
	Required bool

	Operand Term
}

func (*Group) isTerm() {}

// Component Value Multipliers

// ZeroOrMore indicates that the operand occurs zero or more times.
type ZeroOrMore struct {
	Operand Term
}

func (*ZeroOrMore) isTerm() {}

// OneOrMore indicates that the operand occurs one or more times.
type OneOrMore struct {
	Operand Term
}

func (*OneOrMore) isTerm() {}

// ZeroOrOne indicates that the operand is optional (occurs zero or one times).
type ZeroOrOne struct {
	Operand Term
}

func (*ZeroOrOne) isTerm() {}

// Range indicates that the operand occurs at least Min and at most Max times. If Commas
// is true, values must be separated by commas.
type Repeat struct {
	Min    int
	Max    int
	Commas bool

	Operand Term
}

func (*Repeat) isTerm() {}
