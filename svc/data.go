package svc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/btagrass/gobiz/app"
	"github.com/btagrass/gobiz/utl"
	"github.com/glebarez/sqlite"
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
		slog.Error(err.Error())
		os.Exit(1)
	}
	if uri.Scheme == "sqlite" {
		dsn = strings.TrimPrefix(uri.String(), fmt.Sprintf("%s://", uri.Scheme))
		err = utl.MakeDir(filepath.Dir(dsn))
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		dialector = sqlite.Open(dsn)
	} else if uri.Scheme == "mysql" {
		uri.Host = fmt.Sprintf("(%s)", uri.Host)
		dsn = strings.TrimPrefix(uri.String(), fmt.Sprintf("%s://", uri.Scheme))
		name := strings.TrimPrefix(uri.Path, "/")
		sdb, err := sql.Open("mysql", utl.Replace(dsn, name, "information_schema"))
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		defer sdb.Close()
		_, err = sdb.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", name))
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		dialector = mysql.New(mysql.Config{
			DSN:               dsn,
			DefaultStringSize: 100,
		})
	}
	var level logger.LogLevel
	switch app.LogLevel {
	case slog.LevelWarn:
		level = logger.Warn
	case slog.LevelError:
		level = logger.Error
	default:
		level = logger.Info
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
				LogLevel:                  level,
			},
		),
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	if uri.Scheme == "sqlite" {
		err = db.Exec("PRAGMA journal_mode=WAL;").Error
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}
	err = db.Callback().Create().Before("gorm:create").Register("gorm:id", func(db *gorm.DB) {
		if db.Statement.Schema == nil {
			return
		}
		id := db.Statement.Schema.LookUpField("Id")
		if id == nil {
			return
		}
		kind := db.Statement.ReflectValue.Kind()
		if kind == reflect.Array || kind == reflect.Slice {
			for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
				_, zero := id.ValueOf(db.Statement.Context, db.Statement.ReflectValue.Index(i))
				if zero {
					id.Set(db.Statement.Context, db.Statement.ReflectValue.Index(i), utl.IntId())
				}
			}
		} else if kind == reflect.Struct {
			_, zero := id.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
			if zero {
				id.Set(db.Statement.Context, db.Statement.ReflectValue, utl.IntId())
			}
		}
	})
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	sdb, err := db.DB()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	sdb.SetMaxIdleConns(viper.GetInt("dsn.maxIdleConns"))
	sdb.SetMaxOpenConns(viper.GetInt("dsn.maxOpenConns"))
	sdb.SetConnMaxLifetime(viper.GetDuration("dsn.maxLifetime"))
}

func Migrate(sqls ...string) error {
	for _, s := range sqls {
		sql := s
		if db.Dialector.Name() == "sqlite" {
			sql = utl.Replace(s, "INSERT IGNORE INTO", "INSERT OR IGNORE INTO")
		} else if db.Dialector.Name() == "mysql" {
			sql = utl.Replace(s, "INSERT OR IGNORE INTO", "INSERT IGNORE INTO")
		}
		err := db.Exec(sql).Error
		if err != nil {
			return err
		}
	}
	return nil
}

type DataSvc[M any] struct {
	*Svc[M]
}

func NewDataSvc[M any](prefix string, mdls ...any) *DataSvc[M] {
	s := &DataSvc[M]{
		Svc: NewSvc[M](prefix),
	}
	err := db.AutoMigrate(new(M))
	if err != nil {
		slog.Error(err.Error())
	}
	for _, m := range mdls {
		err := db.AutoMigrate(m)
		if err != nil {
			slog.Error(err.Error())
		}
	}
	return s
}

func (s *DataSvc[M]) Exist(conds ...any) bool {
	m, _ := s.Get(conds...)
	return m != nil
}

func (s *DataSvc[M]) Get(conds ...any) (*M, error) {
	var m M
	err := db.Preload(clause.Associations).First(&m, conds...).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &m, nil
}

func (s *DataSvc[M]) GetLocal(expiration time.Duration, conds ...any) (*M, error) {
	var m *M
	key := s.GetFullKey(conds)
	v, ok := s.Local.Get(key)
	if ok {
		m = v.(*M)
	} else {
		var err error
		m, err = s.Get(conds...)
		if err != nil {
			return m, err
		}
		s.Local.Set(key, m, expiration)
	}
	return m, nil
}

func (s *DataSvc[M]) GetRedis(expiration time.Duration, conds ...any) (*M, error) {
	var m *M
	key := s.GetFullKey(conds)
	err := s.Redis.Get(context.Background(), key).Scan(&m)
	if err != nil {
		m, err = s.Get(conds...)
		if err != nil {
			return m, err
		}
		err = s.Redis.Set(context.Background(), key, m, expiration).Err()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func (s *DataSvc[M]) List(conds ...any) ([]M, int64, error) {
	var ms []M
	var count int64
	db := s.Make(conds...).Preload(clause.Associations).Find(&ms)
	_, ok := db.Statement.Clauses["LIMIT"]
	if ok {
		db = db.Limit(-1).Offset(-1).Count(&count)
	}
	return ms, count, db.Error
}

func (s *DataSvc[M]) ListRecursion(conds ...any) ([]M, int64, error) {
	var ms []M
	var count int64
	db := s.Make(conds...).Preload(clause.Associations).Preload("Children", s.recursion).Find(&ms)
	_, ok := db.Statement.Clauses["LIMIT"]
	if ok {
		db = db.Limit(-1).Offset(-1).Count(&count)
	}
	return ms, count, db.Error
}

func (s *DataSvc[M]) ListLocal(expiration time.Duration, conds ...any) ([]M, error) {
	var ms []M
	key := s.GetFullKey(conds)
	v, ok := s.Local.Get(key)
	if ok {
		ms = v.([]M)
	} else {
		var err error
		ms, _, err = s.List(conds)
		if err != nil {
			return ms, err
		}
		s.Local.Set(key, ms, expiration)
	}
	return ms, nil
}

func (s *DataSvc[M]) ListRedis(expiration time.Duration, conds ...any) ([]M, error) {
	var ms []M
	key := s.GetFullKey(conds)
	err := s.Redis.Get(context.Background(), key).Scan(&ms)
	if err != nil {
		ms, _, err = s.List(conds)
		if err != nil {
			return ms, err
		}
		err = s.Redis.Set(context.Background(), key, ms, expiration).Err()
		if err != nil {
			return ms, err
		}
	}
	return ms, nil
}

func (s *DataSvc[M]) Make(conds ...any) *gorm.DB {
	db := db
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
			var vals []any
			for k, v := range cond {
				value, ok := v.(string)
				if ok {
					if value != "" {
						keys = append(keys, fmt.Sprintf("%s like ?", k))
						vals = append(vals, fmt.Sprintf("%%%s%%", v))
					}
					delete(cond, k)
				}
			}
			if len(keys) > 0 {
				db = db.Where(strings.Join(keys, " and "), vals...)
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

func (s *DataSvc[M]) Purge(conds ...any) error {
	if len(conds) == 0 {
		conds = append(conds, "id > 0")
	}
	return db.Unscoped().Delete(new(M), conds...).Error
}

func (s *DataSvc[M]) Remove(conds ...any) error {
	if len(conds) == 0 {
		conds = append(conds, "id > 0")
	}
	return db.Delete(new(M), conds...).Error
}

func (s *DataSvc[M]) Save(m M, confs ...clause.Expression) error {
	if len(confs) == 0 {
		confs = []clause.Expression{
			clause.OnConflict{
				UpdateAll: true,
			},
		}
	}
	return db.Clauses(confs...).Create(&m).Error
}

func (s *DataSvc[M]) Saves(ms []M, confs ...clause.Expression) error {
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
	return db.Clauses(confs...).Create(&ms).Error
}

func (s *DataSvc[M]) Trans(funcs ...func(tx *gorm.DB) error) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, f := range funcs {
			err := f(tx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *DataSvc[M]) Update(values map[string]any, conds ...any) error {
	return s.Make(conds...).Model(new(M)).Updates(values).Error
}

func (s *DataSvc[M]) recursion(db *gorm.DB) *gorm.DB {
	return db.Preload("Children", s.recursion)
}
