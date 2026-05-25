package models

import (
	"fmt"
	"repair/global/config"
	"repair/global/log"
	"sync"

	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// 전역 DB 풀 (싱글톤)
// *sql.DB 자체가 연결 풀이므로 앱 전체에서 한 개만 생성하고 공유한다.
// 매 요청마다 sql.Open()을 호출하면 풀이 무한정 늘어나 MySQL max_connections를 초과해 다운된다.
var (
	globalDB *sql.DB
	dbOnce   sync.Once
)

func initGlobalDB() {
	db, err := sql.Open(config.Database.TypeString, config.Database.ConnectionString)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	// 연결 풀 설정
	db.SetMaxOpenConns(50)                  // 동시 최대 오픈 연결
	db.SetMaxIdleConns(25)                  // 유휴 연결 유지 수
	db.SetConnMaxLifetime(10 * time.Minute) // 연결 최대 수명
	db.SetConnMaxIdleTime(5 * time.Minute)  // 유휴 연결 타임아웃

	globalDB = db
}

type PagingType struct {
	Page     int
	Pagesize int
}

type OrderingType struct {
	Order string
}

type LimitType struct {
	Limit int
}

type OptionType struct {
	Page     int
	Pagesize int
	Order    string
	Limit    int
}

type Where struct {
	Column  string
	Value   interface{}
	Compare string
}

type Custom struct {
	Query string
}

type Base struct {
	Query string
}

type Groupby struct {
	Value int `json:"value"`
	Count int `json:"count"`
}

func Paging(page int, pagesize int) PagingType {
	return PagingType{Page: page, Pagesize: pagesize}
}

func Ordering(order string) OrderingType {
	return OrderingType{Order: order}
}

func Limit(limit int) LimitType {
	return LimitType{Limit: limit}
}

type Connection struct {
	Conn        *sql.DB
	Tx          *sql.Tx
	Transaction bool
	Isolation   bool
}

// Close는 전역 풀을 닫지 않는다.
// Connection은 풀에 대한 핸들 wrapper일 뿐이므로,
// 트랜잭션이 살아있다면 롤백만 수행하고 풀(c.Conn)은 그대로 둔다.
// (예전 코드는 매 요청마다 풀을 닫아버려서 풀링 효과가 사라지고
//
//	다량 요청 시 MySQL 연결이 폭증해 일시 다운되는 문제가 있었다.)
func (c *Connection) Close() {
	if c.Transaction && c.Tx != nil {
		_ = c.Tx.Rollback()
		c.Transaction = false
		c.Tx = nil
	}
}

func (c *Connection) IsConnect() bool {
	return c.Conn != nil
}

func (c *Connection) Exec(query string, params ...interface{}) (sql.Result, error) {
	if c.Transaction {
		return c.Tx.Exec(query, params...)
	} else {
		return c.Conn.Exec(query, params...)
	}
}

func (c *Connection) Query(query string, params ...interface{}) (*sql.Rows, error) {
	if c.Transaction {
		return c.Tx.Query(query, params...)
	} else {
		return c.Conn.Query(query, params...)
	}
}

func (c *Connection) Begin() {
	if c.Transaction {
		return
	}

	c.Tx, _ = c.Conn.Begin()
	c.Transaction = true
	c.Isolation = true
}

func (c *Connection) Commit() error {
	c.Transaction = false
	return c.Tx.Commit()
}

func (c *Connection) Rollback() {
	if !c.Transaction {
		return
	}

	err := c.Tx.Rollback()
	if err != nil {
		log.Error().Msg(err.Error())
	}
	c.Transaction = false
}

// GetConnection은 전역 풀에 대한 새 Connection 핸들을 반환한다.
// 풀 자체는 앱 시작 시 단 한 번만 초기화된다.
func GetConnection() *Connection {
	dbOnce.Do(initGlobalDB)
	if globalDB == nil {
		return nil
	}

	return &Connection{
		Conn:        globalDB,
		Tx:          nil,
		Transaction: false,
	}
}

func NewConnection() *Connection {
	db := GetConnection()

	if db != nil {
		return db
	}

	time.Sleep(100 * time.Millisecond)

	db = GetConnection()

	if db != nil {
		return db
	}

	time.Sleep(500 * time.Millisecond)

	db = GetConnection()

	if db != nil {
		return db
	}

	time.Sleep(1 * time.Second)

	db = GetConnection()

	if db != nil {
		return db
	}

	time.Sleep(2 * time.Second)

	db = GetConnection()

	return db
}

func QueryArray(db *Connection, query string, items []interface{}) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error

	rows, err = db.Conn.Query(query, items...)
	return rows, err
}

func ExecArray(db *Connection, query string, items []interface{}) error {
	var err error

	_, err = db.Conn.Exec(query, items...)
	return err
}

func InitDate() string {
	return "1000-01-01 00:00:00"
}

type Double float64

func (c Double) MarshalJSON() ([]byte, error) {
	if float64(c) == float64(int(c)) {
		return []byte(fmt.Sprintf("%v.0", int64(c))), nil
	}

	return []byte(fmt.Sprintf("%v", float64(c))), nil
}
