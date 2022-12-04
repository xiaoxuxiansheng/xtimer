package mysql

import (
	"errors"
	"fmt"

	"github.com/xiaoxuxiansheng/xtimer/common/conf"

	mysql2 "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DuplicateEntryErrCode = 1062

// Client 客户端
type Client struct {
	*gorm.DB
}

// GetClient 获取一个数据库客户端
func GetClient(confProvider *conf.MysqlConfProvider) (*Client, error) {
	db, err := gorm.Open(mysql.Open(confProvider.Get().DSN), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to connect database, err: %w", err))
	}
	_db, err := db.DB()
	if err != nil {
		panic(err)
	}
	_db.SetMaxOpenConns(151) //设置数据库连接池最大连接数
	_db.SetMaxIdleConns(50)  //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭
	return &Client{DB: db}, nil
}

func NewClient(db *gorm.DB) *Client {
	return &Client{
		DB: db,
	}
}

func IsDuplicateEntryErr(err error) bool {
	var mysqlErr *mysql2.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == DuplicateEntryErrCode
}
