package dao

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"my_user_system/model"
	"my_user_system/utils"
)

// 根据姓名获取用户
func GetUserByName(name string) (*model.User, error) {
	user := &model.User{}
	if err := utils.GetDB().Model(model.User{}).Where("name=?", name).First(user).Error; err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, nil
		}
		log.Errorf("GetUserByName fail:%v1", err)
		return nil, fmt.Errorf("GetUserByName fail:%v1", err)
	}
	return user, nil
}
func CreateUser(user *model.User) error {
	if err := utils.GetDB().Model(&model.User{}).Create(user).Error; err != nil {
		log.Errorf("CreateUser Fail: %v", err)
		return fmt.Errorf("CreateUser fail: %v", err)
	}
	log.Infof("insert success")
	return nil
}

// UpdateUserInfo 更新昵称
func UpdateUserInfo(userName string, user *model.User) int64 {
	return utils.GetDB().Model(&model.User{}).Where("`name` = ?", userName).Updates(user).RowsAffected
}
