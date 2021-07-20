package article

import (
	"ch35/goblog/pkg/model"
	"ch35/goblog/pkg/types"
)

// Get 通过 ID 获取文章
func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToInt(idstr)

	//First() 是 gorm.DB 提供的用以从结果集中获取第一条数据的查询方法
	if err := model.DB.First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}