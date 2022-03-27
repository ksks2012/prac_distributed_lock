package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	metastore "github.com/semeqetjsakatayza/go-metastore-mysql"
	mysqlroundrobinconnector "github.com/yinyin/go-mysql-round-robin-connector"

	interfaces "github.com/distributed_lock/interfaces"
)

//go:generate go-literal-code-gen -do-not-edit -in sqlstruct.md -out sqlstruct.go -sqlschema

// ErrServerAddressRequired indicates server address is not given
var ErrServerAddressRequired = errors.New("database server address is required")

const dialTimeout = time.Second * 10
const readTimeout = time.Second * 10
const writeTimeout = time.Second * 10

const defaultRoundRobinLocationName = "dlock-tab-1"
const mysqlDialerContextName = "dlock-conn"

const metaStoreTableName = "distributed_lock.meta"

type mySQLStorageEngine struct {
	DSN  string
	Conn *sql.DB
}

// NewMySQLStorageEngine creates an instance of MySQL-based storage engine.
func NewMySQLStorageEngine(userName, password, networkAddress, socketPath, databaseName string) (storage interfaces.StorageEngine, err error) {
	cfg := mysqldriver.NewConfig()
	cfg.User = userName
	cfg.Passwd = password
	if networkAddress != "" {
		cfg.Net = "tcp"
		cfg.Addr = networkAddress
	} else if socketPath != "" {
		cfg.Net = "unix"
		cfg.Addr = socketPath
	} else {
		return nil, ErrServerAddressRequired
	}
	cfg.DBName = databaseName
	cfg.Timeout = dialTimeout
	cfg.ReadTimeout = readTimeout
	cfg.WriteTimeout = writeTimeout
	// cfg.ParseTime = true
	return &mySQLStorageEngine{
		DSN: cfg.FormatDSN(),
	}, nil
}

func makeMySQLRoundRobinLocations(networkAddresses []string, socketPath string) (locations []mysqlroundrobinconnector.Location) {
	if socketPath != "" {
		loc := mysqlroundrobinconnector.Location{
			Network: "unix",
			Address: socketPath,
		}
		locations = append(locations, loc)
	}
	for _, networkAddress := range networkAddresses {
		if "" == networkAddress {
			continue
		}
		loc := mysqlroundrobinconnector.Location{
			Network: "tcp",
			Address: networkAddress,
		}
		locations = append(locations, loc)
	}
	return
}

// NewMySQLRoundRobinStorageEngine creates an instance of MySQL-based storage engine with round-robin connector.
func NewMySQLRoundRobinStorageEngine(userName, password string, networkAddresses []string, socketPath, databaseName string) (storage interfaces.StorageEngine, err error) {
	locations := makeMySQLRoundRobinLocations(networkAddresses, socketPath)
	if err = mysqlroundrobinconnector.RegisterLocations(defaultRoundRobinLocationName, locations); nil != err {
		return
	}
	mysqldriver.RegisterDialContext(mysqlDialerContextName, mysqlroundrobinconnector.RoundRobinDialContext) // TODO: go-sql-mysql 1.4.2+
	// mysqldriver.RegisterDial(mysqlDialerContextName, mysqlroundrobinconnector.RoundRobinDial)
	cfg := mysqldriver.NewConfig()
	cfg.User = userName
	cfg.Passwd = password
	cfg.Net = mysqlDialerContextName
	cfg.Addr = defaultRoundRobinLocationName
	cfg.DBName = databaseName
	cfg.Timeout = dialTimeout
	cfg.ReadTimeout = readTimeout
	cfg.WriteTimeout = writeTimeout
	// cfg.ParseTime = true
	return &mySQLStorageEngine{
		DSN: cfg.FormatDSN(),
	}, nil
}

func (eng *mySQLStorageEngine) prepareSchema() (err error) {
	ctx := context.Background()
	metaStoreInst := metastore.MetaStore{
		TableName: metaStoreTableName,
		Ctx:       ctx,
		Conn:      eng.Conn,
	}
	if _, err = metaStoreInst.PrepareSchema(); nil != err {
		return
	}
	mgmt := schemaManager{
		ctx:  ctx,
		conn: eng.Conn,
	}
	schemaRev, err := mgmt.FetchSchemaRevision()
	if nil != err {
		return err
	}
	if schemaChanged, err := mgmt.UpgradeSchema(schemaRev); nil != err {
		return err
	} else if schemaChanged {
		if schemaRev, err = mgmt.FetchSchemaRevision(); nil != err {
			return err
		}
	}
	if !schemaRev.IsUpToDate() {
		return fmt.Errorf("DistributedLocks schema not up-to-date: %#v", schemaRev)
	}
	return
}

func (eng *mySQLStorageEngine) Open() (err error) {
	if eng.Conn != nil {
		return nil
	}
	if eng.Conn, err = sql.Open("mysql", eng.DSN); nil != err {
		eng.Conn = nil
	} else if err = eng.prepareSchema(); nil != err {
		eng.Close()
	}
	return
}

func (eng *mySQLStorageEngine) Close() (err error) {
	if eng.Conn == nil {
		return nil
	}
	conn := eng.Conn
	eng.Conn = nil
	return conn.Close()
}

func (eng *mySQLStorageEngine) FetchMetaInt64(ctx context.Context, metaKey string, defaultValue int64) (value int64, err error) {
	metaStoreInst := metastore.MetaStore{
		TableName: metaStoreTableName,
		Ctx:       ctx,
		Conn:      eng.Conn,
	}
	value, _, err = metaStoreInst.FetchInt64(metaKey, defaultValue)
	return
}
