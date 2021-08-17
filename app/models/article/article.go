package article

import (
	"ch35/goblog/app/models"
	"ch35/goblog/app/models/category"
	"ch35/goblog/app/models/user"
	"ch35/goblog/pkg/model"
	"ch35/goblog/pkg/pagination"
	"net/http"

	//"ch35/goblog/pkg/logger"
	"ch35/goblog/pkg/route"
)

// Article 文章模型
type Article struct {
	models.BaseModel

	Title string
	Body  string

	UserID uint64 `gorm:"not null;index"`
	User   user.User

	CategoryID uint64 `gorm:"not null;default:4;index"`
	Category   category.Category
}

// Link 方法用来生成文章链接
func (a Article) Link() string {
	return route.Name2URL("articles.show", "id", a.GetStringID())
}

// CreatedAtDate 创建日期
func (a Article) CreatedAtDate() string {
	return a.CreatedAt.Format("2006-01-02")
}

// GetByCategoryID 获取分类相关的文章
func GetByCategoryID(cid string, r *http.Request, perPage int) ([]Article, pagination.ViewData, error) {

	// 1. 初始化分页实例
	db := model.DB.Model(Article{}).Where("category_id = ?", cid).Order("created_at desc")
	_pager := pagination.New(r, db, route.Name2URL("categories.show", "id", cid), perPage)

	// 2. 获取视图数据
	viewData := _pager.Paging()

	// 3. 获取数据
	var articles []Article
	_pager.Results(&articles)

	return articles, viewData, nil
}