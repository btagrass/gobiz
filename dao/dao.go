package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/btagrass/gobiz/app"
	"github.com/btagrass/gobiz/utl"
	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	db *gorm.DB
)

func init() {
	var dialector gorm.Dialector
	dsn := viper.GetString("data.dsn")
	if dsn == "" {
		return
	}
	uri, err := url.Parse(dsn)
	if err != nil {
		logrus.Fatal(err)
	}
	if uri.Scheme == "sqlite" {
		dsn = strings.TrimPrefix(uri.String(), fmt.Sprintf("%s://", uri.Scheme))
		err = utl.MakeDir(filepath.Dir(dsn))
		if err != nil {
			logrus.Fatal(err)
		}
		dialector = sqlite.Open(dsn)
	} else if uri.Scheme == "mysql" {
		uri.Host = fmt.Sprintf("(%s)", uri.Host)
		dsn = strings.TrimPrefix(uri.String(), fmt.Sprintf("%s://", uri.Scheme))
		name := strings.TrimPrefix(uri.Path, "/")
		info, err := sql.Open("mysql", utl.Replace(dsn, name, "information_schema"))
		if err != nil {
			logrus.Fatal(err)
		}
		defer info.Close()
		_, err = info.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", name))
		if err != nil {
			logrus.Fatal(err)
		}
		dialector = mysql.New(mysql.Config{
			DSN:               dsn,
			DefaultStringSize: 100,
		})
	}
	db, err = gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.New(
			log.New(io.MultiWriter(os.Stdout, app.LogFile), "", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				IgnoreRecordNotFoundError: true,
				LogLevel:                  logger.LogLevel(logrus.GetLevel() - 1),
			},
		),
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	if uri.Scheme == "sqlite" {
		err = db.Exec("PRAGMA journal_mode=WAL;").Error
		if err != nil {
			logrus.Fatal(err)
		}
	}
	err = db.Callback().Create().Before("gorm:create").Register("gorm:id", func(d *gorm.DB) {
		if d.Statement.Schema == nil {
			return
		}
		id := d.Statement.Schema.LookUpField("Id")
		if id == nil {
			return
		}
		kind := d.Statement.ReflectValue.Kind()
		if kind == reflect.Array || kind == reflect.Slice {
			for i := 0; i < d.Statement.ReflectValue.Len(); i++ {
				_, zero := id.ValueOf(d.Statement.Context, d.Statement.ReflectValue.Index(i))
				if zero {
					id.Set(d.Statement.Context, d.Statement.ReflectValue.Index(i), utl.IntId())
				}
			}
		} else if kind == reflect.Struct {
			_, zero := id.ValueOf(d.Statement.Context, d.Statement.ReflectValue)
			if zero {
				id.Set(d.Statement.Context, d.Statement.ReflectValue, utl.IntId())
			}
		}
	})
	if err != nil {
		logrus.Fatal(err)
	}
	sqlDb, err := db.DB()
	if err != nil {
		logrus.Fatal(err)
	}
	sqlDb.SetMaxIdleConns(viper.GetInt("dsn.maxIdleConns"))
	sqlDb.SetMaxOpenConns(viper.GetInt("dsn.maxOpenConns"))
	sqlDb.SetConnMaxLifetime(viper.GetDuration("dsn.maxLifetime"))
}

func Migrate(mdls []any, sqls ...string) error {
	if db == nil {
		return nil
	}
	err := db.AutoMigrate(mdls...)
	if err != nil {
		return err
	}
	for _, s := range sqls {
		sql := s
		if db.Dialector.Name() == "sqlite" {
			sql = utl.Replace(s, "INSERT IGNORE INTO", "INSERT OR IGNORE INTO")
		} else if db.Dialector.Name() == "mysql" {
			sql = utl.Replace(s, "INSERT OR IGNORE INTO", "INSERT IGNORE INTO")
		}
		err = db.Exec(sql).Error
		if err != nil {
			logrus.Error(err)
		}
	}
	return nil
}

type Dao[M any] struct {
	Db *gorm.DB
}

func NewDao[M any]() *Dao[M] {
	return &Dao[M]{
		Db: db,
	}
}

func (d *Dao[M]) Get(conds ...any) (*M, error) {
	var m M
	err := d.Db.First(&m, conds...).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &m, nil
}

func (d *Dao[M]) List(conds ...any) ([]M, int64, error) {
	var ms []M
	var count int64
	db := d.Make(conds...).Find(&ms)
	_, ok := db.Statement.Clauses["LIMIT"]
	if ok {
		db = db.Limit(-1).Offset(-1).Count(&count)
	}
	return ms, count, db.Error
}

func (d *Dao[M]) Make(conds ...any) *gorm.DB {
	db := d.Db
	if len(conds) > 0 {
		index := 0
		length := len(conds)
		cond, ok := conds[index].(map[string]any)
		if ok {
			size, ok := cond["size"]
			if ok {
				db = db.Limit(cast.ToInt(size))
				delete(cond, "size")
			}
			current, ok := cond["current"]
			if ok {
				db = db.Offset(cast.ToInt(size) * (cast.ToInt(current) - 1))
				delete(cond, "current")
			}
			var keys []string
			var values []any
			for k, v := range cond {
				value, ok := v.(string)
				if ok {
					if value != "" {
						keys = append(keys, fmt.Sprintf("%s like ?", k))
						values = append(values, fmt.Sprintf("%%%s%%", v))
					}
					delete(cond, k)
				}
			}
			if len(keys) > 0 {
				db = db.Where(strings.Join(keys, " and "), values...)
			}
			index++
		}
		order, ok := conds[length-1].(string)
		if ok && strings.Contains(order, "order by ") {
			db = db.Order(utl.Replace(order, "order by ", ""))
			length--
		}
		if index < length {
			db = db.Where(conds[index], conds[index+1:length]...)
		}
	}
	return db
}

func (d *Dao[M]) Purge(conds ...any) error {
	if len(conds) == 0 {
		conds = append(conds, "id > 0")
	}
	return d.Db.Unscoped().Delete(new(M), conds...).Error
}

func (d *Dao[M]) Remove(conds ...any) error {
	if len(conds) == 0 {
		conds = append(conds, "id > 0")
	}
	return d.Db.Delete(new(M), conds...).Error
}

func (d *Dao[M]) Save(m M, confs ...clause.Expression) error {
	if len(confs) == 0 {
		confs = []clause.Expression{
			clause.OnConflict{
				UpdateAll: true,
			},
		}
	}
	return d.Db.Clauses(confs...).Create(&m).Error
}

func (d *Dao[M]) Saves(ms []M, confs ...clause.Expression) error {
	if len(ms) == 0 {
		return nil
	}
	if len(confs) == 0 {
		confs = []clause.Expression{
			clause.OnConflict{
				UpdateAll: true,
			},
		}
	}
	return d.Db.Clauses(confs...).Create(&ms).Error
}

func (d *Dao[M]) Trans(funcs ...func(tx *gorm.DB) error) error {
	return d.Db.Transaction(func(tx *gorm.DB) error {
		for _, f := range funcs {
			err := f(tx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (d *Dao[M]) Update(values map[string]any, conds ...any) error {
	return d.Make(conds...).Model(new(M)).Updates(values).Error
}
