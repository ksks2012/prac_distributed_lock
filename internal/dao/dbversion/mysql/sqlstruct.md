# Heading Code

* `keep-empty-line`

```go
package mysql

import (
  "database/sql"
  "fmt"
  "context"

  metastore "github.com/semeqetjsakatayza/go-metastore-mysql"
)

```

# DistributedLockMeta (distributed-lock-meta) r.1

* `strip-spaces`

```sql
```

## Routines

### prepare fetch revision

```go
metaStoreInst := metastore.MetaStore{
  TableName: metaStoreTableName,
  Ctx:       m.ctx,
  Conn:      m.conn,
}
```

### fetch revision

```go
if ${SCHEMA_REV_VAR}, _, err = metaStoreInst.FetchRevision(${SCHEMA_REV_KEY}); nil != err {
  return nil, err
}
```

### update revision

```go
metaStoreInst := metastore.MetaStore{
  TableName: metaStoreTableName,
  Ctx:       m.ctx,
  Conn:      m.conn,
}
err = metaStoreInst.StoreRevision(key, rev)
```

# ExclusiveLocks (exclusive-locks) r.1

* `const`: `sqlCreate`
* `strip-spaces`

```sql
CREATE TABLE `lock_exclusive_lock` (
   `id` int(4) NOT NULL AUTO_INCREMENT COMMENT 'Primary key',
   `resource_name` varchar(64) NOT NULL DEFAULT '' COMMENT 'Locked resource name',
   `owner` varchar(64) NOT NULL DEFAULT '' COMMENT 'lock owner',
   `desc` varchar(1024) NOT NULL DEFAULT 'Remarks',
   `modified_on` int(10) unsigned DEFAULT '0' COMMENT 'Save data time, automatically generated',
   PRIMARY KEY (`id`),
   UNIQUE KEY `uidx_resource_name` (`resource_name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Resource in exclusive lock';
```
