package tutorial

import "app/model"

type ContentTutorial struct {
	model.BaseModel
	MenuTutorialId int `json:"menu_tutorial_id" gorm:"type:int;column:menu_tutorial_id;comment:'菜单id'"`
	MdCode		string `json:"mdCode" gorm:"type:longtext;column:md_code;comment:'markdown代码'"`
	HtmlCode	string `json:"htmlCode" gorm:"type:longtext;column:html_code;comment:'html代码'"`
}