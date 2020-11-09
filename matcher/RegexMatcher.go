package matcher

import (
	"library-manager/logger"
	"regexp"
)

type RuleCallback func(expr *regexp.Regexp, raw string) (interface{}, error)
type RegexMatcher struct {
	mQueue PriorityList //优先匹配更复杂的规则
}

type Rule struct {
	Alias     string
	mExpr     *regexp.Regexp
	Callback  RuleCallback
	mPriority int //长度越大，优先级越高
}

func (r Rule) Less(other ISortable) bool {
	var otherRule = other.(Rule)
	return r.mPriority < otherRule.mPriority
}

func (r Rule) Equal(other ISortable) bool {
	var otherRule = other.(Rule)
	var lengthEqual = r.mPriority == otherRule.mPriority
	var exprEqual = r.mExpr.String() == otherRule.mExpr.String()
	return lengthEqual && exprEqual
}

func (matcher *RegexMatcher) MustAddRule(alias string, exprStr string, callback RuleCallback) {
	var rule = Rule{Alias: alias, Callback: callback}
	var expr = regexp.MustCompile(exprStr)
	rule.mExpr = expr
	rule.mPriority = len(expr.String())
	err := matcher.mQueue.AddZ(rule)
	if err != nil {
		logger.Error.Fatalln(err)
		return
	}
}

func (matcher RegexMatcher) TryParse(rawString string) (interface{}, error) {
	var result interface{} = nil
	var err error
	matcher.mQueue.EnumerateZ(func(index int, item ISortable) bool {
		var rule = item.(Rule)
		if !rule.mExpr.MatchString(rawString) {
			return false
		}
		//logger.Info.Printf("hit rule:%s",rule.Alias)
		result, err = rule.Callback(rule.mExpr, rawString)
		return true
	})
	return result, err
}
