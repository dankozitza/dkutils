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
      msg = "ForceType: v is type " + vtype.String() +
         ". It must be a pointer!"
      panic(ErrDkutilsGeneric(msg))

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
   vvaluetype := reflect.TypeOf(vvalue.Interface())

   if vvaluetype.String() != dtype.String() {
      msg = "ForceType: expected type " + dtype.String() +
         " but found type " + vvaluetype.String() + ". setting default value"

      if !vvalue.CanSet() {
         msg = "ForceType: vvalue cannot be Set!"
         panic(ErrDkutilsGeneric(msg))
      }
      vvalue.Set(dvalue)

      return ErrDkutilsGeneric(msg)
   }

   return nil
}
