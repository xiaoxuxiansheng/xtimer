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
func GetClient(conf *conf.MySQLConfig) (*Client, error) {
	db, err := gorm.Open(mysql.Open(conf.DSN), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to connect database, err: %w", err))
	}
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
