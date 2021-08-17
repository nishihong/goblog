package category

import (
	"ch35/goblog/app/models"
	"ch35/goblog/pkg/model"
	"ch35/goblog/pkg/route"
	"ch35/goblog/pkg/types"
)

// Category 文章分类
type Category struct {
	models.BaseModel

	//ID uint64
	Name string `gorm:"type:varchar(255);not null;" valid:"name"`
}

// Link 方法用来生成文章链接
func (c Category) Link() string {
	return route.Name2URL("categories.show", "id", c.GetStringID())
}

// Get 通过 ID 获取分类
func Get(idstr string) (Category, error) {
	var category Category
	id := types.StringToInt(idstr)
	if err := model.DB.First(&category, id).Error; err != nil {
		return category, err
	}

	return category, nil
}
