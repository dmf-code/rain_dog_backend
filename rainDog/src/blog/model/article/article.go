package article

import "blog/utils/model"

type Article struct {
	model.BaseModel
	MdCode		string `json:"mdCode" gorm:"column:md_code;"`
	HtmlCode	string `json:"htmlCode" gorm:"html_code;"`
	Title		string `json:"title"`
	CategoryIds string `json:"categoryIds" gorm:"column:category_ids;"`
	TagIds		string `json:"tagIds" gorm:"column:tag_ids;"`
}

type PostField struct {
	Title string `json:"title"`
	CheckedCategorys string `json:"checkedCategory" gorm:"column:category_ids;"`
	CheckedTags string `json:"checkedTag" gorm:"column:tag_ids;"`
	MdCode string `json:"mdCode" gorm:"column:md_code;"`
	HtmlCode string `json:"htmlCode" gorm:"html_code;"`
}

type PutField struct {
	Id string `json:"id"`
	Title string `json:"title"`
	CheckedCategorys string `json:"checkedCategory" gorm:"column:category_ids;"`
	CheckedTags string `json:"checkedTag" gorm:"column:tag_ids;"`
	MdCode string `json:"mdCode"`
	HtmlCode string `json:"htmlCode"`
}

type GetField struct {
	model.BaseModel
	Title string `json:"title"`
	CheckedCategorys string `json:"checked_categorys" gorm:"column:category_ids;"`
	CheckedTags string `json:"checked_tags" gorm:"column:tag_ids;"`
	MdCode string `json:"mdCode"`
	HtmlCode string `json:"htmlCode"`
}

type DeleteField struct {
	Id string `json:"id"`
}
