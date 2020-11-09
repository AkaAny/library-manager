package dbimport

type BookItem struct {
	PropNo    string `gorm:"column:PROP_NO"`     //书本编号
	MARCRecNo string `gorm:"column:MARC_REC_NO"` //书籍编号
	BarCode   string `gorm:"column:BAR_CODE"`    //条码
	Price     string `gorm:"column:PRICE"`       //价格
	InDate    string `gorm:"column:IN_DATE"`     //收录时间
	CallNo    string `gorm:"column:CALL_NO"`     //中图法分类号
}

func (BookItem) TableName() string {
	return "view_item"
}

type MARC struct {
	MARCRecNo string `gorm:"column:MARC_REC_NO" json:",omitempty"` //书籍编号
	CallNo    string `gorm:"column:M_CALL_NO" json:",omitempty"`   //中图法分类号
	Tittle    string `gorm:"column:M_TITLE" json:",omitempty"`     //书名
	Author    string `gorm:"column:M_AUTHOR" json:",omitempty"`    //作者
	Publisher string `gorm:"column:M_PUBLISHER" json:",omitempty"` //出版社
	PubYear   string `gorm:"column:M_PUB_YEAR" json:",omitempty"`  //出版年月
	ISBN      string `gorm:"column:M_ISBN" json:",omitempty"`      //ISBN编码
}

func (MARC) TableName() string {
	return "view_marc"
}
