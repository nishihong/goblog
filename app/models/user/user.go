package user

import (
	"ch35/goblog/app/models"
	"ch35/goblog/pkg/model"
	"ch35/goblog/pkg/password"
	"ch35/goblog/pkg/types"
)

// User 用户模型
type User struct {
	models.BaseModel

	Name     string `gorm:"type:varchar(255);not null;unique" valid:"name"`
	Email    string `gorm:"type:varchar(255);unique;" valid:"email"`
	Password string `gorm:"type:varchar(255)" valid:"password"`

	// gorm:"-" —— 设置 GORM 在读写时略过此字段，仅用于表单验证
	PasswordConfirm string `gorm:"-" valid:"password_confirm"`
}

// ComparePassword 对比密码是否匹配
func (u User) ComparePassword(_password string) bool {
	return password.CheckHash(_password, u.Password)
}

// Get 通过 ID 获取用户
func Get(idstr string) (User, error) {
	var user User
	id := types.StringToInt(idstr)
	if err := model.DB.First(&user, id).Error; err != nil {
		return user, err
	}

	return user, nil
}

// Get 通过 email 获取用户
func GetByEmail(email string) (User, error) {
	var user User
	if err := model.DB.First(&user, "email = ?", email).Error; err != nil {
		return user, err
	}

	return user, nil
}

// Link 方法用来生成用户链接
func (u User) Link() string {
	return ""
}