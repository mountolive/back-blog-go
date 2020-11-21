// Holds basic lookup operations to be used by any regular *Store
package storefilter

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
	OrderBy  OrderBy
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

type Order int

const (
	DESC Order = iota
	ASC
)

type OrderBy struct {
	FieldName string
	Order     Order
}
