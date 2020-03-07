# GORM

## 原始 SQL 操作 - sql.DB

连接数据库：

```go
import (
	"database/sql"

	"github.com/pingcap/log"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
)

// Connection info, should come from config
const (
	userName = "root"
	password = ""
	network  = "tcp"
	server   = "127.0.0.1"
	port     = 4000
	database = "performance_schema"
)

func OpenTiDB(config *config.Config) *sql.DB {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", userName, password, network, server, port, database)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Error("Connect to db failed: ", zap.Error(err))
		return nil
	}
	// defer db.Close()

	// setting
	db.SetConnMaxLifetime(100 * time.Second)
	db.SetMaxOpenConns(100)

	return db
}
```

查询单行：`row = db.QueryRow(...)`

```go
	sql = `select
	query_sample_text,last_seen
	from cluster_events_statements_summary_by_digest_history
	where schema_name=?
	and summary_begin_time=?
	and summary_end_time=?
	and digest=?
	order by last_seen desc
	limit 1`
	row = db.QueryRow(sql, schema, beginTime, endTime, digest)
	if err := row.Scan(&detail.QuerySampleText, &detail.LastSeen); err != nil {
		return detail, err
	}
  return detail, nil
```

查询多行：`rows, err = db.Query(...)`

```go
func QueryStatementNodes(db *sql.DB, schema, beginTime, endTime, digest string) ([]*StatementNode, error) {
	sql := `select
	address,sum_latency,exec_count,avg_latency,max_latency,avg_mem,sum_backoff_times
	from cluster_events_statements_summary_by_digest_history
	where schema_name=?
	and summary_begin_time=?
	and summary_end_time=?
	and digest=?`
	rows, err := db.Query(sql, schema, beginTime, endTime, digest)
	nodes := []*StatementNode{}

	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		return nodes, err
	}

	for rows.Next() {
		node := new(StatementNode)
		err = rows.Scan(
			&node.Address,
			&node.SumLatency,
			&node.ExecCount,
			&node.AvgLatency,
			&node.MaxLatency,
			&node.AvgMem,
			&node.SumBackoffTimes,
		)
		if err != nil {
			return nodes, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}
```

## gorm.DB.Raw()

查询单行：用 `row = db.Raw(...).Row()` 替代原始的 `db.QueryRow(...)`，其余一致。

查询多行：用 `rows, err := db.Raw(...).Rows()` 替代原始的 `db.Query(...)`，其余一致。

如果只选取一列：

```go
schemas := []string{}
err := db.Raw(sql).Pluck("Database", &schemas).Error
```

示例：

```go
func QueryStatementNodes(db *gorm.DB, schema, beginTime, endTime, digest string) ([]*Node, error) {
	sql := `select
	address,sum_latency,exec_count,avg_latency,max_latency,avg_mem,sum_backoff_times
	from cluster_events_statements_summary_by_digest_history
	where schema_name=?
	and summary_begin_time=?
	and summary_end_time=?
	and digest=?`
	nodes := []*Node{}

  db.Exec(selectPerformanceDB)
	rows, err := db.Raw(sql, schema, beginTime, endTime, digest).Rows()
	if err != nil {
		return nodes, err
  }
  defer rows.Close()

	for rows.Next() {
		node := new(StatementNode)
		err = rows.Scan(
			&node.Address,
			&node.SumLatency,
			&node.ExecCount,
			&node.AvgLatency,
			&node.MaxLatency,
			&node.AvgMem,
			&node.SumBackoffTimes,
		)
		if err != nil {
			return nodes, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}
```

## gorm.DB.Select(...).Where(...).Group(...).Order(...).Find(...)

以上面的 QueryStatementNodes() 方法为例，简化以后如下所示：

```go
func QueryStatementNodes(db *gorm.DB, schema, beginTime, endTime, digest string) (result []*Node, err error) {
	err = db.
		Select(`
			address,
			sum_latency,
			exec_count,
			avg_latency,
			max_latency,
			avg_mem,
			sum_backoff_times
		`).
		Table("PERFORMANCE_SCHEMA.cluster_events_statements_summary_by_digest_history").
		Where("schema_name = ?", schema).
		Where("summary_begin_time = ? AND summary_end_time = ?", beginTime, endTime).
		Where("digest = ?", digest).
		Order("sum_latency DESC").
		Find(&result).Error
	return result, err
}
```

如果只查询单行记录，则用 `.Scan(&rusult).Error` 或 `.First(&result).Error` 替代 `.Find(&result).Error`。
