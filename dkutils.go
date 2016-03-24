package dkutils

import (
	"fmt"
	"github.com/dankozitza/go-convert"
	"reflect"
)

type ErrDkutilsGeneric string

func (e ErrDkutilsGeneric) Error() string {
	return "dkutils." + string(e)
}

type niltype struct{}
type invalidtype struct{}

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
			//panic(ErrDkutilsGeneric(msg))
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
	Check(expected interface{}, variable interface{}) (interface{}, error)
}

// DeepTypeCheck
//
// This function is used to crawl through an expected or 'template' data
// structure and compare it to an unknown data structure to check their contents
// against one another. When non-iterable data types are found they are passed
// to a Checker. If a Checker returns an error DeepTypeCheck stops crawling
// and returns the error to the caller.
//
func DeepTypeCheck(
	expected interface{},
	variable interface{},
	c Checker) (ret interface{}, err error) {

	if expected == nil && variable == nil {
		return nil, fmt.Errorf("expected and variable are nil")
	}
	if expected == nil && variable != nil {
		return nil, fmt.Errorf("expected is nil")
	}
	if expected != nil && variable == nil {
		return nil, fmt.Errorf("variable is nil")
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

	//fmt.Println("Before pointer removal: etype: " + etype.String() + ", vtype: " + vtype.String())

	// check if etype is a pointer and set evalue
	var evalue reflect.Value
	if etype.Kind() == reflect.Ptr {
		//fmt.Println("removed pointer to expected value")
		evalue = reflect.Indirect(reflect.ValueOf(expected))
		etype = reflect.TypeOf(evalue.Interface())

	} else {
		evalue = reflect.ValueOf(expected)
	}

	// same for vtype
	var vvalue reflect.Value
	if vtype.Kind() == reflect.Ptr {
		//fmt.Println("removed pointer to variable value")
		vvalue = reflect.Indirect(reflect.ValueOf(variable))
		vtype = reflect.TypeOf(vvalue.Interface())

	} else {
		vvalue = reflect.ValueOf(variable)
	}

	//fmt.Println("After pointer removal: etype: ", etype, ", vtype: ", vtype)

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
			return nil, ErrDkutilsGeneric(msg)

		default:
			ret, err = c.Check(evalue.Interface(), vvalue.Interface())
			if err != nil {
				return nil, err
			}
			//fmt.Println("DeepTypeCheck: setting vvalue")

			// this causes panic when vvalue is unaddressable
			//vvalue.Set(reflect.ValueOf(r))

			return ret, nil
		}
	}

	switch etype.Kind() {
	case reflect.Map:

		// make sure the return is the correct type
		ret = map[string]interface{}{}

		newexp, ok := evalue.Interface().(map[string]interface{})
		if !ok {
			msg = "DeepTypeCheck: maps are expected to be of type " +
				"map[string]interface{} not type " + etype.String()
			return nil, ErrDkutilsGeneric(msg)
		}
		newvar, ok := vvalue.Interface().(map[string]interface{})
		if !ok {
			msg = "DeepTypeCheck: maps are expected to be of type " +
				"map[string]interface{} not type " + vtype.String()
			return nil, ErrDkutilsGeneric(msg)
		}

		//fmt.Println("etype.Kind():", etype.Kind(), "vtype.Kind():", vtype.Kind())

		for k, _ := range newexp {

			// this is hardcoded for jobdist
			if k == "links" {
				ret.(map[string]interface{})[k] = newexp[k]
				continue
			}

			//fmt.Println("-\ngoing deeper into expected key:", k, ", val:", newexp[k])
			r, err := DeepTypeCheck(newexp[k], newvar[k], c)
			if err != nil {
				return nil, ErrDkutilsGeneric("DeepTypeCheck: While checking " +
					"expected key " + k + ": " + err.Error())
			}
			ret.(map[string]interface{})[k] = r
		}

	case reflect.Slice:

		// make sure the return is the correct type
		ret = []interface{}{}

		newexp, ok := evalue.Interface().([]interface{})
		if !ok {
			msg = "DeepTypeCheck: slices are expected to be of type " +
				"[]interface{} not type " + etype.String()
			return nil, ErrDkutilsGeneric(msg)
		}
		newvar, ok := vvalue.Interface().([]interface{})
		if !ok {
			msg = "DeepTypeCheck: slices are expected to be of type " +
				"[]interface{} not type " + vtype.String()
			return nil, ErrDkutilsGeneric(msg)
		}

		// check that lengths are the same
		elen := len(newexp)
		vlen := len(newvar)
		if elen != vlen {
			msg = fmt.Sprint("DeepTypeCheck: expected slice of length ", elen,
				" but found slice of length ", vlen)
			return nil, ErrDkutilsGeneric(msg)
		}

		fmt.Println("etype.Kind():", etype.Kind(), "vtype.Kind:", vtype.Kind())

		for i, _ := range newexp {
			//fmt.Println("-\ngoing deeper into expected index:", i, ", val:", newexp[i])
			r, err := DeepTypeCheck(newexp[i], newvar[i], c)
			if err != nil {
				return nil, err
			}
			ret.([]interface{})[i] = r
		}

	default:

		ret, err = c.Check(evalue.Interface(), vvalue.Interface())
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

// Persuader
//
// Used by DeepTypePeruade() to convert v to whatever type e is. Covers string
// and int types.
//
type Persuader struct{}

func (p Persuader) Check(e interface{}, v interface{}) (converted interface{}, err error) {

	//fmt.Println("converting ", reflect.TypeOf(v), "to ", reflect.TypeOf(e))

	switch e.(type) {
	case int:
		r, err := convert.Int(v)
		if err != nil {
			return nil, ErrDkutilsGeneric("Persuader.Check: could not convert: " +
				err.Error())
		}
		return r, nil

	case int32:
		r, err := convert.Int32(v)
		if err != nil {
			return nil, ErrDkutilsGeneric("Persuader.Check: could not convert: " +
				err.Error())
		}
		return r, nil

	case int64:
		r, err := convert.Int64(v)
		if err != nil {
			return nil, ErrDkutilsGeneric("Persuader.Check: could not convert: " +
				err.Error())
		}
		return r, nil

	case float32:
		r, err := convert.Float32(v)
		if err != nil {
			return nil, ErrDkutilsGeneric("Persuader.Check: could not convert: " +
				err.Error())
		}
		return r, nil

	case float64:
		r, err := convert.Float64(v)
		if err != nil {
			return nil, ErrDkutilsGeneric("Persuader.Check: could not convert: " +
				err.Error())
		}
		return r, nil

	case string:
		r, err := convert.String(v)
		if err != nil {
			return nil, ErrDkutilsGeneric("Persuader.Check: could not convert: " +
				err.Error())
		}
		return r, nil

	default:
		return nil, ErrDkutilsGeneric("Persuader.Check: unknown type")

	}
}

func DeepTypePersuade(expected interface{}, variable interface{}) (interface{}, error) {
	var p Persuader
	return DeepTypeCheck(expected, variable, p)
}

// DeepTypeSprint
//
// Calls DeepTypeSprintDepthIndent with default values for depth and indent.
//
func DeepTypeSprint(i interface{}) (string, error) {
	return DeepTypeSprintDepthIndent(i, 0, "   ")
}

// DeepTypeSprintDepthIndent
//
// DeepTypeSprintDepthIndent will crawl through a data structure building a
// string that shows the types of each element. It will not print memory
// locations.
//
func DeepTypeSprintDepthIndent(i interface{}, depth int64, ind string) (ret string, err error) {

	var indent string
	for i := int64(0); i < depth; i++ {
		indent += ind
	}

	var printedtype string

	itype := reflect.TypeOf(i)

	if i == nil {
		return indent + fmt.Sprint(itype) + "(\"<nil>\")\n", nil
	}

	printedtype = itype.String()

	// check if itype is a pointer and set ivalue
	var ivalue reflect.Value
	if itype.Kind() == reflect.Ptr {
		ivalue = reflect.Indirect(reflect.ValueOf(i))
		if !ivalue.IsValid() {
			ivalue = reflect.ValueOf(invalidtype{})
		}
		itype = reflect.TypeOf(ivalue.Interface())

	} else {
		ivalue = reflect.ValueOf(i)
	}

	switch itype.Kind() {
	case reflect.Map:

		newi, ok := ivalue.Interface().(map[string]interface{})
		if !ok {

			// make a new map that is map[string]interface{}
			newi = map[string]interface{}{}

			mkeys := ivalue.MapKeys()

			for _, k := range mkeys {
				newi[k.Interface().(string)] = ivalue.MapIndex(k).Interface()
			}
		}

		ret += indent + printedtype + "{\n"

		for k, _ := range newi {

			ret += indent + "-key[" + k + "]:\n"

			r, err := DeepTypeSprintDepthIndent(newi[k], depth+1, ind)
			if err != nil {
				return ret, ErrDkutilsGeneric("DeepTypeSprintDepthIndent: While checking " +
					"key " + k + ": " + err.Error())
			}
			ret += r
		}

		ret += indent + "}\n"

	case reflect.Slice:

		newi, ok := ivalue.Interface().([]interface{})
		if !ok {

			// make a new slice that is []interface{}
			newi = make([]interface{}, ivalue.Len())

			for index := 0; index < ivalue.Len(); index++ {
				newi[index] = ivalue.Index(index).Interface()
			}
		}

		ret += indent + printedtype + "{\n"

		for index, _ := range newi {

			ret += indent + "-index[" + fmt.Sprint(index) + "]:\n"

			r, err := DeepTypeSprintDepthIndent(newi[index], depth+1, ind)
			if err != nil {
				return ret, ErrDkutilsGeneric("DeepTypeSprintDepthIndent: While checking " +
					"index " + fmt.Sprint(index) + ": " + err.Error())
			}
			ret += r
		}

		ret += indent + "}\n"

	case reflect.Struct:

		ret += indent + printedtype + "{\n"

		for index := 0; index < ivalue.NumField(); index++ {
			iStructField := itype.Field(index)

			ret += indent + "-field[" + iStructField.Name + "]:\n"

			// if the field name does not start with an upper case letter skip it
			if iStructField.PkgPath != "" {
				ret += indent + ind + "<hidden field>\n"
				continue
			}

			r, err := DeepTypeSprintDepthIndent(
				ivalue.Field(index).Interface(),
				depth+1,
				ind)
			if err != nil {
				return ret, ErrDkutilsGeneric("DeepTypeSprintDepthIndent: While checking " +
					"field " + iStructField.Name + ": " + err.Error())
			}
			ret += r
		}

		ret += indent + "}\n"

	default:

		ret += indent + printedtype + "(\"" + fmt.Sprint(ivalue.Interface()) + "\")\n"
	}

	return ret, nil
}
