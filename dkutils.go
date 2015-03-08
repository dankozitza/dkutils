package dkutils

import (
	"fmt"
	"reflect"
)

type ErrDkutilsGeneric string

func (e ErrDkutilsGeneric) Error() string {
	return "dkutils." + string(e)
}

type niltype struct{}

// ForceType
//
// Takes var *interface{} and a default value. Checks that the var's type
// matches the default's type. Sets the var to the default if the
// var is not the correct type.
//
// This function requires that v be type *interface{}.
//
// example:
//
//    var test interface{} = "test"
//    err := ForceType(&test, 3) // test now contains the int 3
//
//    var test interface{} = 0
//    err := ForceType(&test, 3) // same type, test still contains int 0
//
func ForceType(v interface{}, d interface{}) error {

	var msg string
	vtype := reflect.TypeOf(v)
	dtype := reflect.TypeOf(d)

	// check that v is a pointer
	if vtype.Kind() != reflect.Ptr {
		msg = "ForceType: v must be a pointer!"
		panic(ErrDkutilsGeneric(msg))
	}

	// make a Value object for use with refect. this one expects v is a pointer
	ptrvalue := reflect.ValueOf(v)
	// make a Value object for the value ptrvalue points to
	vvalue := reflect.Indirect(ptrvalue)

	if !vvalue.CanAddr() {
		if !vvalue.CanSet() {
			msg = "ForceType: vvalue cannot be Set!"
			panic(ErrDkutilsGeneric(msg))
		}
		vvalue.Set(reflect.ValueOf(niltype{}))

	} else {
		if fmt.Sprint(vvalue.Addr()) != "<*interface {} Value>" {
			msg = "ForceType: vvalue.Addr(): " + fmt.Sprint(vvalue.Addr()) +
				" must be <*interface {} Value>"
			panic(ErrDkutilsGeneric(msg))
		}
	}

	// check if d is a pointer and set dvalue
	var dvalue reflect.Value
	if dtype.Kind() == reflect.Ptr {
		dvalue = reflect.Indirect(reflect.ValueOf(d))
		dtype = reflect.TypeOf(dvalue.Interface())

	} else {
		dvalue = reflect.ValueOf(d)
	}

	// if vvalue is nil set it to niltype
	if vvalue.IsNil() {
		if !vvalue.CanSet() {
			msg = "ForceType: vvalue cannot be Set!"
			panic(ErrDkutilsGeneric(msg))
		}
		vvalue.Set(reflect.ValueOf(niltype{}))
	}

	// get the true type of the value v points to from vvalue.Interface()
	vtype = reflect.TypeOf(vvalue.Interface())

	if vtype.String() != dtype.String() {
		msg = "ForceType: expected type " + dtype.String() +
			" but found type " + vtype.String() + ". setting default value"

		if !vvalue.CanSet() {
			msg = "ForceType: vvalue cannot be Set!"
			panic(ErrDkutilsGeneric(msg))
		}
		vvalue.Set(dvalue)

		return ErrDkutilsGeneric(msg)
	}

	return nil
}

// Checker
//
// A checker is used by DeepTypeCheck when a pair of non-iterable variables are
// found.
//
type Checker interface {

	// Check
	//
	// the Check function will be called on any non-iterable expected, variable
	// pair found.
	//
	Check(expected interface{}, variable interface{}) error
}

// DeepTypeCheck
//
//	This function is used to crawl through an expected or 'template' data
// structure and compare it to an unknown data structure to check their contents
// against one another. When non-iterable data types are found they are passed
// to a Checker. If a Checker returns an error DeepTypeCheck stops crawling
// and returns the error to the caller.
//
// TODO: A checker can also return a new variable which will replace
//       whatever 'variable' parameter was passed in.
//
func DeepTypeCheck(expected interface{}, variable interface{}, c Checker) error {

	if expected == nil && variable == nil {
		fmt.Println("expected and variable are nil")
		return nil
	}
	if expected == nil && variable != nil {
		fmt.Println("only expected is nil")
		return nil
	}
	if expected != nil && variable == nil {
		fmt.Println("only variable is nil")
		return nil
	}

	var msg string
	etype := reflect.TypeOf(expected)
	vtype := reflect.TypeOf(variable)

	//	if conf["dkutils_DeepTypeCheck_ignore_pointers"].(int) == 0 {
	//		if etype.String() != vtype.String() {
	//			msg = "DeepTypeCheck: expected type " + etype.String() +
	//				" but found type " + vtype.String() + "."
	//
	//			return ErrDkutilsGeneric(msg)
	//
	//		}
	//	}

	fmt.Println("Before pointer removal: etype: " + etype.String() + ", vtype: " + vtype.String())

	// check if etype is a pointer and set evalue
	var evalue reflect.Value
	if etype.Kind() == reflect.Ptr {
		fmt.Println("removed pointer to expected value")
		evalue = reflect.Indirect(reflect.ValueOf(expected))
		etype = reflect.TypeOf(evalue.Interface())

	} else {
		evalue = reflect.ValueOf(expected)
	}

	// same for vtype
	var vvalue reflect.Value
	if vtype.Kind() == reflect.Ptr {
		fmt.Println("removed pointer to variable value")
		vvalue = reflect.Indirect(reflect.ValueOf(variable))
		vtype = reflect.TypeOf(vvalue.Interface())

	} else {
		vvalue = reflect.ValueOf(variable)
	}

	fmt.Println("After pointer removal: etype: ", etype, ", vtype: ", vtype)

	// maybe only check kind here? check type in the default case of the kind
	// switch.
	//if etype != vtype {
	//	msg = "DeepTypeCheck: expected type " + etype.String() +
	//		" but found type " + vtype.String() + "."

	//	return ErrDkutilsGeneric(msg)

	//}

	if etype.Kind() != vtype.Kind() {

		// if the expected is a map or slice then return an error, these cannot
		// be converted
		switch etype.Kind() {
		case reflect.Map, reflect.Slice:
			msg = "DeepTypeCheck: expected reflect.Kind " +
				fmt.Sprint(etype.Kind()) + " but found " +
				fmt.Sprint(vtype.Kind())
			return ErrDkutilsGeneric(msg)

		default:
			return c.Check(evalue.Interface(), vvalue.Interface())
		}
	}

	switch etype.Kind() {
	case reflect.Map:

		newexp, ok := evalue.Interface().(map[string]interface{})
		if !ok {
			msg = "DeepTypeCheck: maps are expected to be of type " +
				"map[string]interface{} not type " + etype.String()
			return ErrDkutilsGeneric(msg)
		}
		newvar, ok := vvalue.Interface().(map[string]interface{})
		if !ok {
			msg = "DeepTypeCheck: maps are expected to be of type " +
				"map[string]interface{} not type " + vtype.String()
			return ErrDkutilsGeneric(msg)
		}

		fmt.Println("etype.Kind():", etype.Kind(), "vtype.Kind():", vtype.Kind())

		for k, _ := range newexp {
			fmt.Println("-\ngoing deeper into expected key:", k, ", val:", newexp[k])
			err := DeepTypeCheck(newexp[k], newvar[k], c)
			if err != nil {
				return err
			}
		}

	case reflect.Slice:

		newexp, ok := evalue.Interface().([]interface{})
		if !ok {
			msg = "DeepTypeCheck: slices are expected to be of type " +
				"[]interface{} not type " + etype.String()
			return ErrDkutilsGeneric(msg)
		}
		newvar, ok := vvalue.Interface().([]interface{})
		if !ok {
			msg = "DeepTypeCheck: slices are expected to be of type " +
				"[]interface{} not type " + vtype.String()
			return ErrDkutilsGeneric(msg)
		}

		// check that lengths are the same
		elen := len(newexp)
		vlen := len(newvar)
		if elen != vlen {
			msg = fmt.Sprint("DeepTypeCheck: expected slice of length ", elen,
				" but found slice of length ", vlen)
			return ErrDkutilsGeneric(msg)
		}

		fmt.Println("etype.Kind():", etype.Kind(), "vtype.Kind:", vtype.Kind())

		for i, _ := range newexp {
			fmt.Println("-\ngoing deeper into expected index:", i, ", val:", newexp[i])
			err := DeepTypeCheck(newexp[i], newvar[i], c)
			if err != nil {
				return err
			}
		}

	default:

		// now call c.Check()
		fmt.Println("etype.Kind():", etype.Kind(), "vtype.Kind:", vtype.Kind())
		fmt.Println("did the default")

		return nil
	}

	return nil
}
