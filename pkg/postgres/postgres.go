/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package postgres

import (
	global "CyberMatchmaker/pkg"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sync"
)

var pgOnce sync.Once

// InitDB 初始化 PostgreSQL 连接池
func InitDB(dsn string) {
	pgOnce.Do(func() {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("PostgreSQL 初始化失败: %v", err)
		}

		// 获取底层的 sql.DB 以配置连接池属性
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("获取底层 sql.DB 失败: %v", err)
		}

		// 设置连接池参数以应对高并发
		sqlDB.SetMaxIdleConns(10)  // 核心空闲连接数
		sqlDB.SetMaxOpenConns(100) // 最大打开连接数

		global.DB = db
		log.Println("PostgreSQL 初始化成功")
	})
}
