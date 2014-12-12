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
