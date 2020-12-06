package dbimport

import (
	"regexp"
)

var sDateExpr = regexp.MustCompile("^\\d{4}\\.\\d{2}$")

func CheckPubYear(rawPubYear string) bool {
	return sDateExpr.MatchString(rawPubYear)
}

var sNewISBNExpr = regexp.MustCompile("^\\d{3}-\\d{1,5}-\\d{1,7}-\\d{1,6}-\\d{1}$")
var sOldISBNExpr = regexp.MustCompile("^\\d{1,5}-\\d{1,7}-\\d{1,6}-\\d{1}$")

func CheckISBN(isbn string) bool {
	//2007年1月1日之前，ISBN是10位的，之后是13位
	if sNewISBNExpr.MatchString(isbn) {
		return true
	}
	if sOldISBNExpr.MatchString(isbn) {
		return true
	}
	return false
}
