// Holds basic lookup operations to be used by any regular *Store
package storehelper

// Single lookup condition
type Lookup struct {
	FieldName  string
	Comparator Comparator
}

// Composed search criteria
// Unifies Lookups either by AND or by ORs
type Criteria struct {
	Lookups  []Lookup
	Operator LogicalOperator
}

type Comparator int

// EQ   ->  =
// NEQ  ->  !=
// LET  ->  <=
// LT   ->  <
// GET  ->  >=
// GT   ->  >
const (
	EQ Comparator = iota
	NEQ
	LET
	LT
	GET
	GT
)

type LogicalOperator int

const (
	AND LogicalOperator = iota
	OR
)
