package dkutils

type Comparison interface{}{
	Compare(a interface{}, b interface{}) (interface{}, error)
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
func DeepCompare(a interface{}, b interface{}, c Comparison)
