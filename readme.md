# cq - cypher queries for database/sql
A database/sql implementation for Cypher. Still in development, but already surprisingly usable. 

[![Build Status](https://travis-ci.org/wfreeman/cq.png?branch=master)](https://travis-ci.org/wfreeman/cq)
[![Coverage Status](https://coveralls.io/repos/wfreeman/cq/badge.png?branch=HEAD)](https://coveralls.io/r/wfreeman/cq?branch=HEAD)
[![Waffle](https://badge.waffle.io/wfreeman/cq.png?label=ready)](https://waffle.io/wfreeman/cq)

Thanks to [Baron](http://twitter.com/xaprb), [Mike](http://twitter.com/mikearpaia), and [Jason](https://github.com/jmcvetta) for the ideas/motivation to start on this project. Cypher is close enough to SQL that it seems to fit pretty well in the idiomatic database/sql implementation.

#### Other Go drivers for Neo4j that support Cypher
* [Neoism](https://github.com/jmcvetta/neoism) (a careful/complete REST API implementation)
* [GonormCypher](https://github.com/marpaia/GonormCypher) (a port of AnormCypher, to get up and running quickly)

## usage
See the [excellent database/sql tutorial](http://go-database-sql.org/index.html) from [VividCortex](https://vividcortex.com/), as well as the [package documentation for database/sql](http://golang.org/pkg/database/sql/) for an introduction to the idiomatic go database access.

One key thing to mention that's slightly different, at least so far. For now, you can (and should) use parameters, but the placeholders must be numbers in sequence, e.g. `{0}`, `{1}`, `{2}`, and then you must put them in order in the calls to `Query`/`Exec`. I hope to overcome this at some point, but I haven't figured out a good way to represent Cypher's named parameters in the parameterized query feature of database/sql.

## [minimum viable snippet](http://blog.fogus.me/2012/08/23/minimum-viable-snippet/)

```go
package main

import (
	"database/sql"
	"log"
	
	_ "github.com/wfreeman/cq"
)

func main() {
	db, err := sql.Open("neo4j-cypher", "http://localhost:7474")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare(`
		match (n:User)-[:FOLLOWS]->(m:User) 
		where n.screenName = {0} 
		return m.screenName as friend
		limit 10
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query("wefreema")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var friend string
	for rows.Next() {
		err := rows.Scan(&friend)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(friend)
	}
}
```

## feature support
I've programmed the basic functionality of `Prepare()` and `Query()`, so you can already get things done. Here is a list of features planned (crossed out items are complete):

* ~~`sql.Open()`~~
* ~~`db.Prepare()`~~
* ~~`stmt.Exec()`~~
* ~~`stmt.Query()`~~
* ~~support for primitive parameters and results~~
* ~~`db.Begin()`~~
* ~~`tx.Prepare()`~~
* ~~tx: `stmt.Prepare()`~~
* ~~tx: `stmt.Exec()`~~
* ~~keepalive for transactions~~
* ~~support for array parameters and results via ValueConverter~~
* ~~support for map parameters and results via ValueConverter~~
* tx: `stmt.Query()`

## transactional API
The transactional API using `db.Begin()` is optimized for sending many queries to the [transactional Cypher endpoint](http://docs.neo4j.org/chunked/milestone/rest-api-transactional.html), in that it will batch them up and send them in chunks by default. Currently only supports `stmt.Exec()` within a transaction, will work on supporting `stmt.Query()` next and queueing up results.

#### transactional API example
```go
func main() {
	db, err := sql.Open("neo4j-cypher", "http://localhost:7474")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	
	stmt, err := tx.Prepare("create (:User {screenName:{0}})")	
	if err != nil {
		log.Fatal(err)
	}
	
	stmt.Exec("wefreema")
	stmt.Exec("JnBrymn")
	stmt.Exec("technige")
	
	err := tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
```

#### transactional API benchmarks
Able to get sustained times of 20k+ cypher statements per second, even with multiple nodes per create... on a 2011 vintage macbook.

```
(master ✓) wes-macbook:cq go test -bench=".*Transaction.*" -test.benchtime=10s
PASS
BenchmarkTransactional10SimpleCreate	  100000	    150630 ns/op
BenchmarkTransactional100SimpleCreate	  500000	     39202 ns/op
BenchmarkTransactional1000SimpleCreate	 1000000	     27320 ns/op
BenchmarkTransactional10000SimpleCreate	  500000	     28524 ns/op
ok  	github.com/wfreeman/cq	79.973s
```

## license

MIT license. See license file.

