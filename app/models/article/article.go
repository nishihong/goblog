package article

import (
	"ch35/goblog/app/models"
	//"ch35/goblog/pkg/logger"
	"ch35/goblog/pkg/route"
)

// Article 文章模型
type Article struct {
	models.BaseModel

	Title string
	Body  string
}

// Link 方法用来生成文章链接
func (a Article) Link() string {
	return route.Name2URL("articles.show", "id", a.GetStringID())
}