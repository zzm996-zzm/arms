package demo

import (
	"database/sql"
	"github.com/zzm996-zzm/arms/framework/contract"
	"time"

	"github.com/zzm996-zzm/arms/framework/gin"
	"github.com/zzm996-zzm/arms/framework/provider/orm"
)

func (api *DemoApi) DemoOrm(c *gin.Context) {
	logger := c.MustMakeLog()
	logger.Info(c, "request start", nil)

	//初始化一个 orm.DB

	gormService := c.MustMakeOrm()
	db, err := gormService.GetDB(orm.WithConfigPath("database.default"))
	if err != nil {
		logger.Error(c, err.Error(), nil)
		// c.AbortWithError(50001, err)
		return
	}

	db.WithContext(c)
	db.AutoMigrate(&User{})
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	logger.Info(c, "migrate ok", nil)

	//插入一条数据
	email := "zzmhaoshuai@icloud.com"
	name := "zzm"
	age := uint8(18)
	birthday := time.Date(1999, 8, 6, 0, 0, 0, 0, time.Local)
	user := &User{
		Name:         name,
		Email:        &email,
		Age:          age,
		Birthday:     &birthday,
		MemberNumber: sql.NullString{},
		ActivatedAt:  sql.NullTime{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = db.Create(user).Error
	logger.Info(c, "insert user", map[string]interface{}{"id": user.ID, "err": err})

	//更新一条数据
	user.Name = "bb"
	err = db.Save(user).Error
	logger.Info(c, "update user", map[string]interface{}{"err": err, "id": user.ID})

	//查询一条数据
	queryUser := &User{ID: user.ID}
	err = db.First(queryUser).Error
	logger.Info(c, "query user", map[string]interface{}{"err": err, "name": queryUser.Name})

	// 删除一条数据
	err = db.Delete(queryUser).Error
	logger.Info(c, "delete user", map[string]interface{}{"err": err, "id": user.ID})

	c.JSON(200, "ok")

}




// DemoCache cache的简单例子
func (api *DemoApi) DemoCache(c *gin.Context) {
	logger := c.MustMakeLog()
	logger.Info(c, "request start", nil)
	// 初始化cache服务
	cacheService := c.MustMake(contract.CacheKey).(contract.Cache)
	// 设置key为foo
	err := cacheService.Set(c, "foo", "bar", 1*time.Hour)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	// 获取key为foo
	val, err := cacheService.Get(c, "foo")
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	logger.Info(c, "cache get", map[string]interface{}{
		"val": val,
	})
	// 删除key为foo
	if err := cacheService.Del(c, "foo"); err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, "ok")
}