package dbimport

import (
	"library-manager/logger"
	"library-manager/matcher"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var dateMatcher matcher.RegexMatcher

func initDateMatcher() {
	//除了这条规则，其它的规则匹配到的数据多少有点问题，需要重新核实修改，优先豆瓣，豆瓣没有且比较老的书可以去孔夫子旧书网核查信息
	dateMatcher.MustAddRule("year.month",
		"^\\d{4}\\.\\d{1,2}$",
		func(expr *regexp.Regexp, raw string) (interface{}, error) {
			var dateStrs = strings.Split(raw, ".")
			year, err := parseInt64(dateStrs[0])
			if err != nil {
				return nil, err
			}
			month, err := parseInt64(dateStrs[1])
			var yearAndMonth = time.Date(int(year), time.Month(month),
				1, 0, 0, 0, 0, time.Local)
			return yearAndMonth, nil
		})
	//0000013676,J292.12/050.2,巴尔扎克名言硬笔书法字帖,高杨选编,北京广播学院出版社,1992.4.,7-81004-338-2
	//u1s1,学校当年买这个本明显水的不行而且借了没法用的字帖书感觉有点和书商py交易恰回扣的味道
	dateMatcher.MustAddRule("year.month.", "^\\d{2,4}\\.\\d{1,2}\\.$", func(expr *regexp.Regexp, raw string) (interface{}, error) {
		raw = strings.Trim(raw, ".")
		return dateMatcher.TryParse(raw) //交给"year.month"
	})
	//0000000357,F120/020,阶段 特征 战略 社会主义初级阶段经济问题,方生主编,经济日报出版社,1989 10,7-80036
	dateMatcher.MustAddRule("year month", "^\\d{4} \\d{1,2}$", func(expr *regexp.Regexp, raw string) (interface{}, error) {
		raw = strings.Replace(raw, " ", ".", 1)
		return dateMatcher.TryParse(raw)
	})
	//把这些来自"未来"的书集合起来，做一个大赏
	//0000002529,D922.2/421,会计法规专论,杨纪琬主编,东北财经大学出版社,18989.6,7-81005
	//0000011513,F740.2/402,国际期货市场,李文贤，张春荣编著,中国对外经济贸易出版社,19990.8,7-80004
	//0000008624,C93/222,多目标决策分析及其在工程和经济中的应用,（美）乔伊科奇等著；王寅初译,航空工业出版社,19878.5,7-80046
	//0000014512,,English（英语） v1,北京外国语学院英语系,商务印书馆,19789.6,7-100
	//0000015838,F713.5/240,市场营销调研学,（美）鲁克等著；《市场营销调研学》翻译小组译,福建人民出版社,19889.3,7-211
	dateMatcher.MustAddRule("wrongFiveCharYear.Month", "^\\d{5}\\.\\d{1,2}$", func(expr *regexp.Regexp, raw string) (interface{}, error) {
		var dateStrs = strings.Split(raw, ".")
		var yearStr = dateStrs[0]
		yearStr = strings.Replace(yearStr, "8", "", 1) //18989，其实这个规则有问题，如果是19889的话可能会得到错误的出版日期
		if strings.Contains(yearStr, "99") {           //19990
			yearStr = strings.Replace(yearStr, "9", "", 1)
		}
		dateStrs[0] = yearStr
		raw = strings.Join(dateStrs, ".")
		return dateMatcher.TryParse(raw)
	})
	//0000016967,F714.1/243,价格管理学,伍世安，顾晓燕主编,北京经济学院出版社,19925,7-5638-0299-1
	//这么久远的未来还有这么多财会管理的书，看来communism还没实现
	dateMatcher.MustAddRule("wrongFiveCharYearOnly", "^\\d{5}$", func(expr *regexp.Regexp, raw string) (interface{}, error) {
		raw = raw + ".1" //这样就能匹配"wrongFiveCharYear.Month"的规则了（其实用一个原则，更改->适应而不是创建更多细致但有重复逻辑的规则）
		return dateMatcher.TryParse(raw)
	})

	//0000004119,F271/040,我国各地试办农工商联合企业情况汇编,辛农编,新华出版社,1981.,7-5011
	//0000016192,F270/554,透视未来企业∶革命性的未来企业新观念,托夫勒著；潘祖铭译,志文出版社,1985..,W-1015
	dateMatcher.MustAddRule("wrongYear.+", "^\\d{4}\\.+$", func(expr *regexp.Regexp, raw string) (interface{}, error) {
		raw = strings.TrimRight(raw, ".")
		return dateMatcher.TryParse(raw)
	})
	//0000003635,F091.9/262,不平等交换:对帝国主义贸易的研究,（希腊）伊曼纽尔著；文贯中科译,中国对外经济贸易出版社,988.5,7-80004
	//本来和"wrongTwoIntYear.Month"规则可以合并，但考虑到学校可能会有18xx年的书（虽然我觉得没有）
	dateMatcher.MustAddRule("wrongThreeIntYear.Month", "^\\d{3}\\.\\d{1,2}$", func(expr *regexp.Regexp, raw string) (interface{}, error) {
		raw = "1" + raw
		return dateMatcher.TryParse(raw)
	})
	//0000002785,D922.29/030.2,一九八五年商业政策法规汇编,商业部办公厅编,中国商业出版社,86.3,7-5044
	dateMatcher.MustAddRule("wrongTwoIntYear.Month", "^\\d{2}\\.\\d{1,2}$", func(expr *regexp.Regexp, raw string) (interface{}, error) {
		raw = "19" + raw
		return dateMatcher.TryParse(raw)
	})
	//0000003751,A41/234,毛泽东选集: v1,毛泽东著,人民出版社,52—77,7-01
	dateMatcher.MustAddRule("wrongTwoCharEstimatedYear", "^\\d{2}—\\d{2}$", func(expr *regexp.Regexp, raw string) (interface{}, error) {
		var years = strings.Split(raw, "—")
		return dateMatcher.TryParse(years[0])
	})
	//[0000004029 G250/462 外国图书馆学名著选读 袁咏秋，李主编 北京大学出版社 19 7-301]

	dateMatcher.MustAddRule("wrongTwoIntEstimatedYear", "^\\d{2,3}$", func(expr *regexp.Regexp, raw string) (interface{}, error) {
		if len(raw) == 2 {
			raw = "19" + raw
		}
		if len(raw) == 3 {
			raw = raw + "0"
		}
		return dateMatcher.TryParse(raw)
	})
	//0000005344,D67/232,香港政务官阶层的构成,穆迈伦著 ；杨立信，罗绍熙译,上海翻译出版公司,198 .11,7-80514
	//0000009301,C93/450,管理的数量概念,（美）埃本，高尔德著；赵国士译,机械工业出版社,19  .8,7-111
	//我理解这可能是后边的年份不详，但是这一手明显就是踢皮球给后面的人的操作符合一个图书馆工作人员的负责态度吗
	dateMatcher.MustAddRule("wrongEstimatedYearMissingLast.Month", "^\\d{2,4} +\\.\\d{1,2}$", func(expr *regexp.Regexp, raw string) (interface{}, error) {
		var dateStrs = strings.Split(raw, ".")
		var yearStr = dateStrs[0]
		yearStr = strings.TrimSpace(yearStr)
		if len(yearStr) == 3 {
			yearStr = strings.TrimSpace(yearStr) + "0"
		}
		if len(yearStr) == 2 {
			yearStr = strings.TrimSpace(yearStr) + "00"
		}
		dateStrs[0] = yearStr
		raw = strings.Join(dateStrs, ".")
		return dateMatcher.TryParse(raw) //交给year.month处理
	})
	dateMatcher.MustAddRule("yearOnly",
		"^\\d{4}$",
		func(expr *regexp.Regexp, raw string) (interface{}, error) {
			year, err := parseInt64(raw)
			if err != nil {
				return nil, err
			}
			var yearOnly = time.Date(int(year),
				time.January, 1, 0, 0, 0, 0, time.Local)
			return yearOnly, nil
		})
	//0000000775,D912.29/960,经济法概要,肖明编,航空工业出版社,.11,7-80046
	dateMatcher.MustAddRule(".monthOnly", "^\\.\\d{1,2}$", func(expr *regexp.Regexp, raw string) (interface{}, error) {
		raw = "0000" + raw
		return dateMatcher.TryParse(raw) //交给"year.month"处理
	})

	dateMatcher.MustAddRule("nullOrWhiteSpace", "^$", func(expr *regexp.Regexp, raw string) (interface{}, error) {
		return time.Time{}, nil //空字符串输入，直接返回最早的日期，确保能够用time.After()查找的到
	})
}

var authorMatcher matcher.RegexMatcher

func initAuthorMatcher() {

}

func init() {
	initDateMatcher()
}

func parseInt64(str string) (int64, error) {
	val, err := strconv.ParseInt(str, 10, 64)
	return val, err
}

func removeUnMatched(expr *regexp.Regexp, str string) string {
	var result string
	str = expr.ReplaceAllStringFunc(str, func(s string) string {
		result += s
		return s
	})
	return result
}

// 这里更好的实现是正则匹配器
func ParsePubYear(rawPubYear string) (time.Time, error) {
	//预处理
	rawPubYear = strings.TrimSpace(rawPubYear) //已经trim过了，后续无需再TrimSpace
	//0000005039,A56/31,马克思 恩格斯 列宁 斯大林论自然辩证法,河南师范大学马列主义教研室编辑,,1983.1m,
	//0000015573,H31/558.3,星期日广播英语选 编v1,申葆青编,,"1984,.7",
	//0000007331,F112/633,世界经济参考资料:v1-v3,国家计委经济研究所世界经济组,,FIHI,
	//替换所有非空格+非.+非数字字符（包括纯字母序列）
	rawPubYear = removeUnMatched(regexp.MustCompile("( +)|(\\.)|(\\d*)"), rawPubYear)

	//匹配处理
	obj, err := dateMatcher.TryParse(rawPubYear)
	if err != nil {
		return time.Time{}, err
	}
	if obj == nil {
		logger.Error.Printf("no rules matched for raw pub year:\"%s\"", rawPubYear)
		return time.Time{}, nil
	}
	var timeResult = obj.(time.Time)
	return timeResult, nil
}

func ParseAuthor(rawAuthorStr string) {

}
