package dkutils

import (
	"fmt"
	"testing"
)

func TestForceType(t *testing.T) {

	fmt.Println("ForceType: testing testvar1")
	var testvar1 interface{} = "str"
	err := ForceType(&testvar1, 3)
	if err != nil {
		fmt.Println("ForceType returned error: " + err.Error())

	} else {
		fmt.Println("ForceType did not return an error!")
		t.Fail()
	}

	if fmt.Sprint(testvar1) != fmt.Sprint(3) {
		fmt.Println("ForceType did not set testvar1 to default")
		fmt.Println("testvar:", testvar1, "should be:", 3)
		t.Fail()
	}

	fmt.Println("ForceType: testing testvar2")
	var testvar2 interface{} = 0
	err = ForceType(&testvar2, 3)
	if err != nil {
		fmt.Println("ForceType returned error: " + err.Error())
		t.Fail()
	}

	if fmt.Sprint(testvar2) != fmt.Sprint(0) {
		fmt.Println("ForceType should not have modified testvar2")
		fmt.Println("testvar:", testvar2, "should be:", 0)
		t.Fail()
	}

	fmt.Println("ForceType: testing testvar3")
	var testvar3 interface{} = float64(7.8192)
	var testdefault interface{} = 2
	err = ForceType(&testvar3, &testdefault)
	if err != nil {
		fmt.Println("ForceType returned error: " + err.Error())

	} else {
		fmt.Println("ForceType did not return an error!")
		t.Fail()
	}

	if fmt.Sprint(testvar3) != fmt.Sprint(2) {
		fmt.Println("ForceType should not have modified testvar3")
		fmt.Println("testvar:", testvar3, "should be:", 2)
		t.Fail()
	}

	fmt.Println("ForceType: testing testvar4")
	var testvar4 interface{}
	err = ForceType(&testvar4, 2)
	if err != nil {
		fmt.Println("ForceType returned error: " + err.Error())

	} else {
		fmt.Println("ForceType did not return an error!")
		t.Fail()
	}

	if fmt.Sprint(testvar4) != fmt.Sprint(2) {
		fmt.Println("ForceType should not have modified testvar4")
		fmt.Println("testvar:", testvar4, "should be:", 2)
		t.Fail()
	}
}

//func TestDeepTypeCheck(t *testing.T) {
//
//	fmt.Println("\nDeepTypeCheck: testing testvar1")
//	var testvar1 interface{} = "str"
//	err := DeepTypeCheck(testvar1, 3)
//	if err != nil {
//		fmt.Println("DeepTypeCheck returned error: " + err.Error())
//
//	} else {
//		fmt.Println("DeepTypeCheck did not return an error!")
//		t.Fail()
//	}
//
//	if fmt.Sprint(testvar1) != fmt.Sprint(3) {
//		fmt.Println("DeepTypeCheck did not set testvar1 to default")
//		fmt.Println("testvar:", testvar1, "should be:", 3)
//		t.Fail()
//	}
//
//	fmt.Println("\nDeepTypeCheck: testing testvar2")
//	var testvar2 interface{} = map[string]interface{}{
//		"this": int64(12)}
//		err = DeepTypeCheck(map[string]interface{}{"this": int(12)}, testvar2)
//	if err != nil {
//		fmt.Println("DeepTypeCheck returned error: " + err.Error())
//		t.Fail()
//	}
//}

func TestKind(t *testing.T) {
	var tc TestChecker

	fmt.Println("\nDeepTypeCheck: testing kind")
	var expected4 interface{} = []interface{}{
		new(string),
		new(string)}
	var testvar4 interface{} = "horse"
	err := DeepTypeCheck(expected4, testvar4, tc)
	if err != nil {
		fmt.Println("DeepTypeCheck returned error: " + err.Error())
	} else {
		t.Fail()
	}
}

func TestMapTypes(t *testing.T) {
	var tc TestChecker

	fmt.Println("\nDeepTypeCheck: testing map types")
	// this setup shows that testvar3 contains int32 values!?!?
	//var expected3 interface{} = map[string]interface{}{
	//	"twelve": int32(12),
	//	"three": int32(3)}
	//var testvar3 interface{} = map[string]interface{}{
	//	"twelve": int64(12),
	//	"three": int64(3)}

	var expected3 interface{} = map[string]interface{}{
		"twelve": int32(12),
		"three":  int32(3)}
	var testvar3 interface{} = map[string]int32{
		"twelve": 12,
		"three":  3}

	err := DeepTypeCheck(expected3, testvar3, tc)
	if err != nil {
		fmt.Println("DeepTypeCheck returned error: " + err.Error())
	} else {
		t.Fail()
	}
}

func TestDeepMapTypes(t *testing.T) {
	var tc TestChecker

	fmt.Println("\nDeepTypeCheck: testing deep map types")

	// this setup shows that testvar3 contains int32 values!?!?
	var expected3 interface{} = map[string]interface{}{
		"twelve": int32(12),
		"three":  int32(3)}
	var testvar3 interface{} = map[string]interface{}{
		"twelve": int64(12),
		"three":  int64(3)}

	err := DeepTypeCheck(expected3, testvar3, tc)
	if err != nil {
		fmt.Println("DeepTypeCheck returned error: " + err.Error())
		t.Fail()
	}
}

func TestSliceType(t *testing.T) {
	var tc TestChecker
	fmt.Println("\nDeepTypeCheck: testing slice type")
	var expected4 interface{} = []interface{}{
		new(string),
		new(string)}
	var testvar4 interface{} = []*string{
		new(string),
		new(string)}
	err := DeepTypeCheck(expected4, testvar4, tc)
	if err != nil {
		fmt.Println("DeepTypeCheck returned error: " + err.Error())
	} else {
		t.Fail()
	}
}

func TestSliceLength(t *testing.T) {
	var tc TestChecker
	fmt.Println("\nDeepTypeCheck: testing slice length")
	var expected4 interface{} = []interface{}{
		new(string),
		new(int),
		new(int64)}
	var testvar4 interface{} = []interface{}{}
	//	"this": new(float64),
	//	"these": []int{0, 1, 2}}
	err := DeepTypeCheck(expected4, testvar4, tc)
	if err != nil {
		fmt.Println("DeepTypeCheck returned error: " + err.Error())
	} else {
		t.Fail()
	}
}

type TestChecker struct{}

func (tc TestChecker) Check(expected interface{}, variable interface{}) error {

	fmt.Println("TestChecker: got expected: ", expected, " and variable: ", variable)

	return nil
}
