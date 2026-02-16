
## Setup MySQL
```sh
sudo apt install -U mysql-server

# connect to it as the root user
sudo mysql
# mysql -u root -p
##mysql>
```

Establish a database:
```sql
SHOW DATABASES;
DROP DATABASE IF EXISTS snippetbox;

CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- switch database
USE snippetbox;

-- create table
CREATE TABLE snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

-- add an index on the `created` column
CREATE INDEX idx_snippets_created ON snippets(created);
```

Add some dummy records:
```sql
INSERT INTO snippets (title, content, created, expires) VALUES
(
  'Idea: Micro-SaaS',
  'Build a tiny SaaS that converts long videos into short clips automatically.',
  UTC_TIMESTAMP(),
  DATE_ADD(UTC_TIMESTAMP(), INTERVAL 90 DAY)
),
(
  'SQL Reminder',
  'Remember to always use indexes on foreign keys and frequently queried columns.',
  UTC_TIMESTAMP(),
  DATE_ADD(UTC_TIMESTAMP(), INTERVAL 30 DAY)
),
(
  'Go Learning Note',
  'Practice building REST APIs with net/http before jumping into frameworks.',
  UTC_TIMESTAMP(),
  DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
),
(
  'Startup Thought',
  'Distribution matters more than product in early stages.',
  UTC_TIMESTAMP(),
  DATE_ADD(UTC_TIMESTAMP(), INTERVAL 90 DAY)
),
(
  'Debugging Rule',
  'If stuck >30 min, reduce problem to smallest reproducible case.',
  UTC_TIMESTAMP(),
  DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
);


-- verify inserts
SELECT id, title, created, expires
FROM snippets
ORDER BY id DESC;

```

Create a new (non-root) user:
```sql
-- while still connected to the MySQL prompt as root:
CREATE USER 'web'@'localhost';

GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';

ALTER USER 'web'@'localhost' IDENTIFIED BY 'ChangeM3';

exit

```

Test the new user:
```sh
mysql -D snippetbox -u web -p
```

Verify the permission:
```sql
DROP TABLE snippets;
-- ERROR 1142 (42000): DROP command denied to user 'web'@'localhost' for table 'snippets'
```

## Install a db driver
The **db driver** acts as a middleman; translating commands between Go & the MySQL db itself. 

```sh
go get github.com/go-sql-driver/mysql@v1

# you should see `go.sum` file added
# to download the exact version of the packages need for the project: go mod download
# to ensure that nothing in those downloaded packages has been changed unexpectedly: go mod verify
# go mod tidy: automatically removes any unused packages from the go.mod & go.sum files

# alternatively, you can remove a package like so:
go get github.com/foo/bar@none
```

N.B. using the `-u` flag increases the risk of breakages when upgrading packages (just don't do it).
```sh
# e.g., this will update the package & ALL ITS depentencies to their latest versions.
# but it may well be that in `go.mod` the dependencies versions are older.
go get -u github.com/go-sql-driver/mysql
```


```go
// params: driver name, data source name aka connection string aka DSN (format varies based on db & the driver)
// parseTime=true (instructs the driver to convert SQL TIME & DATE fields to Go time.Time values)
// returns: a sql.DB value (not just a db connection; it's a pool of many connections, safe for concurrent access)
// the connection pool is intented to be long-lived.
// Do NOT call `sql.Open()` in a short-lived HTTP handler (it's a waste of memory & network resources)
db, err := sql.Open()
```