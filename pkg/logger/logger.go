/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger 初始化全局的 Zap Logger
func InitLogger() {
	// 1. 配置 lumberjack 实现日志按大小切割
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/app_1.log", // 日志存放路径
		MaxSize:    10,                // 单个文件最大尺寸 (MB)
		MaxBackups: 5,                 // 最多保留 5 个备份文件
		MaxAge:     30,                // 最多保留 30 天
		Compress:   true,              // 是否开启 gzip 压缩旧日志
	})

	// 2. 设置日志的输出格式
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 可读的时间格式 (例如 2026-03-01T16:34:16.000+0800)
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 大写的日志级别 (如 INFO, ERROR)

	// 使用 Console 格式（开发阶段更友好），如果纯生产环境可换成 zapcore.NewJSONEncoder(encoderConfig)
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 3. 配置多路输出：同时输出到文件和标准输出 (控制台)
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, writeSyncer, zap.InfoLevel),                 // 写入文件的级别 (INFO及以上)
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.DebugLevel), // 控制台打印的级别
	)

	// 4. 生成 Logger，开启 Caller 记录（打印出是在哪一行代码记录的日志）
	logger := zap.New(core, zap.AddCaller())

	// 5. 替换 zap 库的全局实例
	// 后续在项目的任何地方，只需要调用 zap.L() 或 zap.S() 就可以直接打印日志
	zap.ReplaceGlobals(logger)
}
