package env

import (
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

func LoadEnv() {
	LoadDbEnv()
}

func loadConfigNode() gdb.ConfigNode {
	return gdb.ConfigNode{
		Host:                 os.Getenv("DB_HOST"),
		Port:                 os.Getenv("DB_PORT"),
		User:                 os.Getenv("DB_USER"),
		Pass:                 os.Getenv("DB_PASS"),
		Name:                 os.Getenv("DB_NAME"),
		Type:                 os.Getenv("DB_TYPE"),
		Link:                 os.Getenv("DB_LINK"),
		Extra:                os.Getenv("DB_EXTRA"),
		Role:                 os.Getenv("DB_ROLE"),
		Debug:                gconv.Bool(os.Getenv("DB_DEBUG")),
		Prefix:               os.Getenv("DB_PREFIX"),
		DryRun:               gconv.Bool(os.Getenv("DB_DRYRUN")),
		Weight:               gconv.Int(os.Getenv("DB_WEIGHT")),
		Charset:              os.Getenv("DB_CHARSET"),
		Protocol:             os.Getenv("DB_PROTOCOL"),
		Timezone:             os.Getenv("DB_TIMEZONE"),
		MaxIdleConnCount:     gconv.Int(os.Getenv("DB_MAX_IDLE_CONN_COUNT")),
		MaxOpenConnCount:     gconv.Int(os.Getenv("DB_MAX_OPEN_CONN_COUNT")),
		MaxConnLifeTime:      gconv.Duration(os.Getenv("DB_MAX_CONN_LIFE_TIME")),
		QueryTimeout:         gconv.Duration(os.Getenv("DB_QUERY_TIMEOUT")),
		ExecTimeout:          gconv.Duration(os.Getenv("DB_EXEC_TIMEOUT")),
		TranTimeout:          gconv.Duration(os.Getenv("DB_TRAN_TIMEOUT")),
		PrepareTimeout:       gconv.Duration(os.Getenv("DB_PREPARE_TIMEOUT")),
		CreatedAt:            "",
		UpdatedAt:            "",
		DeletedAt:            "",
		TimeMaintainDisabled: gconv.Bool(os.Getenv("DB_TIME_MAINTAIN_DISABLED")),
	}
}

func mergeDbConfig(firstConfig, secondConfig gdb.ConfigNode) gdb.ConfigNode {
	newConfigMap := gmap.NewStrAnyMapFrom(gconv.Map(firstConfig))
	newConfigMap.Merge(gmap.NewStrAnyMapFrom(gconv.Map(secondConfig)))

	result := gdb.ConfigNode{}

	_ = gconv.Struct(newConfigMap, &result)

	return result
}

const development = ".env.development"
const test = ".env.test"
const demo = ".env.demo"
const production = ".env.production"

// LoadDbEnv 加载数据库配置
// 加载优先级为：(全局 DB_ENV ｜.env.development | .env.test | .env.demo | .env.production) > .env
// 例如：如果未设置 DB_ENV 环境变量来指定自定义配置文件，如果 (.env.development 或 .env.test 或 .env.demo 或 .env.production) 则加载其配置文件，并将其合并至 .env 配置中
func LoadDbEnv() {
	err := godotenv.Load()
	if err != nil {
		println("加载 .env 配置文件失败")
	}

	dbConfig := loadConfigNode()

	// 加载环境变量 DB_ENV 值指定的自定义配置，一般在开发模式中，配置多个运行环境，实现快速切换环境配置
	dbEnv := os.Getenv("DB_ENV")

	if dbEnv != "" {
		envFile := ".env." + strings.ToLower(dbEnv)
		err = godotenv.Load(".env." + envFile)
		if err != nil {
			println("加载 " + envFile + " 配置文件失败")
		}

		// 合并自定义配置
		dbConfig = mergeDbConfig(dbConfig, loadConfigNode())

		println(envFile + "配置文件已加载")
	} else {
		// 未设置 DB_ENV 环境变量
		if gfile.Exists(".env." + development) {
			// 加载开发环境配置
			err = godotenv.Load(".env." + development)
			if err != nil {
				println("加载 .env." + development + " 配置文件失败")
			} else {
				// 合并开发环境配置
				dbConfig = mergeDbConfig(dbConfig, loadConfigNode())
			}
			println(".env." + development + " 配置文件已加载")
		} else if gfile.Exists(".env." + test) {
			// 加载测试环境配置
			err = godotenv.Load(".env." + test)
			if err != nil {
				println("加载 .env." + test + " 配置文件失败")
			} else {
				// 合并测试环境配置
				dbConfig = mergeDbConfig(dbConfig, loadConfigNode())
			}
			println(".env." + test + " 配置文件已加载")
		} else if gfile.Exists(".env." + demo) {
			// 加载演示境配置
			err = godotenv.Load(".env." + demo)
			if err != nil {
				println("加载 .env." + demo + " 配置文件失败")
			} else {
				// 合并演示境配置
				dbConfig = mergeDbConfig(dbConfig, loadConfigNode())
			}
			println(".env." + demo + " 配置文件已加载")
		} else if gfile.Exists(".env." + production) {
			// 加载生产境配置
			err = godotenv.Load(".env." + production)
			if err != nil {
				println("加载 .env." + production + " 配置文件失败")
			} else {
				// 合并生产境配置
				dbConfig = mergeDbConfig(dbConfig, loadConfigNode())
			}
			println(".env." + production + " 配置文件已加载")
		} else {
			println("未设置 DB_ENV 环境变量，使用 .env 默认配置")
		}
	}

	if strings.Contains(dbConfig.Link, ":") || (dbConfig.Name != "" && dbConfig.User != "") {
		if strings.HasPrefix(dbConfig.Link, "postgres") {
			dbConfig.Type = "pgsql"
		} else {
			dbConfig.Type = strings.Split(dbConfig.Link, ":")[0]
		}

		gdb.SetConfig(gdb.Config{
			gdb.DefaultGroupName: gdb.ConfigGroup{
				dbConfig,
			},
		})
	}
}
