package model

import (
	"encoding/json"
)

type ESBookMarc struct {
	MARCRecNo string `json:"marc_rec_no"` //书籍编号
	CallNo    string `json:"call_no"`     //中图法分类号
	Title     string `json:"title"`       //书名
	Author    string `json:"author"`      //作者
	Publisher string `json:"publisher"`   //出版社
	//PubYear   time.Time `json:"pub_year"`    //出版年月
	ISBN string `json:"isbn"` //ISBN编码
}

func (marc ESBookMarc) ToJSON() string {
	rawData, _ := json.Marshal(marc)
	return string(rawData)
}

func (marc ESBookMarc) String() string {
	return marc.ToJSON()
}
