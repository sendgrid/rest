package rest

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"
	"time"
)

func TestLicense(t *testing.T) {
	re, _ := regexp.Compile("[0-9]{4}")

	f, err := ioutil.ReadFile("LICENSE.txt")

	if err != nil {
		t.Errorf("could not read file: %v", err)
	}

	currentYear := fmt.Sprint(time.Now().Year())
	specifiedYear := string(re.Find(f))

	if specifiedYear != currentYear {
		t.Errorf("Specified year (%s) is different of the current year (%s)", specifiedYear, currentYear)
	}

}
