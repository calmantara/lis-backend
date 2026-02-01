package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/Calmantara/lis-backend/internal/core/configurations"
	"github.com/Calmantara/lis-backend/internal/core/models"
	"github.com/Calmantara/lis-backend/internal/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	Client struct {
		config configurations.Database
		db     *gorm.DB
	}

	connection struct {
		master *gorm.DB
		slave  *gorm.DB
	}

	MySql interface {
		Master() *gorm.DB
		Slave() *gorm.DB
	}
)

func NewConnection(master *Client, slave *Client) MySql {
	return &connection{master: master.db, slave: slave.db}
}

func (con *connection) Master() *gorm.DB {
	return con.master
}

func (con *connection) Slave() *gorm.DB {
	return con.slave
}

func NewClient(config configurations.Database) *Client {
	cln := &Client{config: config}
	// connect to database
	_, err := cln.connect()
	if err != nil {
		panic(err)
	}

	return cln
}

func (g *Client) connect() (*sql.DB, error) {
	sqlDB, err := sql.Open(g.config.Driver, g.config.ConnectionString())
	if err != nil {
		return nil, err
	}
	// set connection setting
	sqlDB.SetMaxOpenConns(g.config.MaxConnection)
	sqlDB.SetMaxIdleConns(g.config.MaxIDleConnection)
	sqlDB.SetConnMaxLifetime(g.config.MaxLifeTime)
	// load gorm connection
	g.loadGorm(sqlDB)

	return sqlDB, nil
}

var (
	mapConditionSymbol = map[string]string{
		"lte":   "<=",
		"lt":    "<",
		"gte":   ">=",
		"ilike": "ilike",
		"wild":  "~*",
	}

	mapConditionOr = map[string]string{
		"or": "OR",
	}
)

func (g *Client) loadGorm(sqlDB *sql.DB) {
	// checking if connection to db has been established
	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		panic(err)
	}
	g.db, _ = gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

func pagination(tx *gorm.DB, page models.Pagination) *gorm.DB {
	if page.Limit == 0 || page.Page == 0 {
		return tx
	}
	offset := (page.Page - 1) * page.Limit

	return tx.Offset(offset).Limit(page.Limit)
}

func beforeFind(tx *gorm.DB, t any) *gorm.DB {
	valueOf := reflect.ValueOf(t)
	typeOf := reflect.TypeOf(t)

	for i := range make([]struct{}, typeOf.NumField()) {
		fieldType := valueOf.FieldByName(typeOf.Field(i).Name)
		// check param
		tag := typeOf.Field(i).Tag.Get("filter")

		// check or condition
		tags := strings.Split(tag, ";")
		if len(tags) > 2 && mapConditionOr[strings.ToLower(tags[1])] == "OR" {
			tn := *tx
			for i, t := range tags {
				if i == 0 {
					if tr := filterTag(t, fieldType, &tn, ""); tr != nil {
						tn = *tr
					}
				}
				if i > 0 && t != "or" {
					if tr := filterTag(t, fieldType, &tn, "OR"); tr != nil {
						tn = *tr
					}
				}
			}
			tx = tx.Where(&tn)
		} else {
			if tr := filterTag(tag, fieldType, tx, ""); tr != nil {
				tx = tr
			}
		}

		tag = typeOf.Field(i).Tag.Get("sort")
		if tr := sortTag(tag, fieldType, tx); tr != nil {
			tx = tr
		}
	}

	return tx
}

func filterTag(tag string, fieldType reflect.Value, tx *gorm.DB, or string) (res *gorm.DB) {
	if tag == "" {
		return
	}
	// separate and init variables
	tags := strings.Split(tag, ";")
	query := tag + " = ?"
	kindSlice := fieldType.Type().Kind() == reflect.Slice

	var condition string
	if len(tags) > 1 && mapConditionSymbol[strings.ToLower(tags[1])] != "" {
		var ok bool
		condition, ok = mapConditionSymbol[strings.ToLower(tags[1])]
		if !ok {
			return
		}
		query = tags[0] + " " + condition + " ?"
	} else {
		// check slice
		if kindSlice {
			if fieldType.IsNil() {
				return
			}
			query = tag + " IN ?"
		}
	}

	val := fieldType.Interface()
	if fieldType.IsZero() {
		return
	}
	// check for bidx fields
	if strings.Contains(tags[0], "_bidx") {
		val, _ = utils.EncryptAES(val.(string))
	}

	// build query
	tn := tx.Where
	if or == "OR" {
		tn = tx.Or
	}

	switch condition {
	case "ilike":
		res = tn(query, fmt.Sprintf("%%%v%%", val))
	case "~*":
		filter := val
		if kindSlice {
			// only work in string / uuid
			var vals []string
			utils.ObjectMapper(val, &vals)
			filter = strings.Join(vals, "|")
		}
		res = tn(query, filter)
	default:
		res = tn(query, val)
	}

	return
}

func sortTag(tag string, fieldType reflect.Value, tx *gorm.DB) (res *gorm.DB) {
	if tag == "" {
		return
	}

	// check string
	kindString := fieldType.Type().Kind() == reflect.String
	if !kindString {
		return
	}
	val := fieldType.Interface()
	if fieldType.IsZero() {
		return
	}
	sorts := strings.Split(val.(string), ",")
	for _, sort := range sorts {
		if sort == "" {
			continue
		}
		sorter := strings.Split(sort, ".")

		// build query
		query := sorter[0]
		if len(sorter) > 1 {
			query += " " + sorter[1]
		}
		res = tx.Order(query)
	}

	return
}

type GormTransactionKey struct{}

func begin(ctx context.Context, tx *gorm.DB) (context.Context, *gorm.DB) {
	txn := getTransaction(ctx)
	if txn != nil {
		// already has transaction
		return ctx, txn
	}
	txn = tx
	// store transaction to context
	ctx = context.WithValue(ctx, GormTransactionKey{}, txn)
	// transaction
	return ctx, txn
}

func getTransaction(ctx context.Context) *gorm.DB {
	txn, ok := ctx.Value(GormTransactionKey{}).(*gorm.DB)
	if ok {
		return txn
	}

	return nil
}

func commit(ctx context.Context) error {
	txn := getTransaction(ctx)
	if txn == nil {
		return nil
	}

	return txn.Commit().Error
}

func rollback(ctx context.Context) error {
	txn := getTransaction(ctx)
	if txn == nil {
		return nil
	}

	return txn.Rollback().Error
}
