package dkutils

type Comparison interface {
	// Compare
	//
	// Takes a and b and returns an interface. What is returned may depend on
	// a specific implementation of Compare. In most cases it will return some
	// combination of a and b.
	//
	Compare(a interface{}, b interface{}) (interface{}, error)
}

type Differ map[string]interface{}
func (d Differ) Compare(a interface{}, b interface{}) (interface{}, error) {
	if a != b {
		return false, ErrDkutilsGeneric("a does not equal b")
	}
	return true, nil
}

// DeepCompare
//
// Ranges over the data structures a and b calling c.Compare() on every pair of
// elements.
//
// var this map[string]interface{}{"cats": []interface{}{"one", "two"}}
// DeepCompare(this, nil, c)
//    // c.Compare() would be called:
//       // c.Compare(this, nil),
//       // c.Compare(this["cats"], nil),
//       // c.Compare(this["cats"][0], nil),
//       // c.Compare(this["cats"][1], nil),
//
// Breakdown for comparison of nested data structures in Perl:
//	http://www.perlmonks.org/?node_id=752326
// 
func DeepCompare(a interface{}, b interface{}, c Comparison) {}
