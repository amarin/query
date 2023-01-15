# query

Query implements SQL query builder for Go.
The aim is to provide fast and simple interface to database functions without need to write SQL code at all.

Main features:

- supporting SELECT (and SELECT COUNT() as subset of SELECT), INSERT, UPDATE and DELETE queries;
- supporting fields conditions to use in SELECT/UPDATE/DELETE queries;
- supporting field and table names aliasing;
- supporting tables JOIN's keeping syntax as close to SQL as possible;
- provides string types to wrap table and field names constants allows to keep all definitions in single place and avoid mistypings;
- generated queries could use either '?' or '$N' placeholders depending on your needs;
- supporting conditionals building over single or several joined tables using complex conditions
- allows to extend standard conditions library with new condition implementations types when required
- all query builders are immutable which allows to keep original complex query definitions and easily derive new ones

## Alternatives and related projects

Here is open source alternatives on 09.09.2022:

- [sqlf](https://github.com/leporo/sqlf) - A fast SQL query builder for Go, last v1.3.0 fired 05.01.2022, switches  versions ~ once per 3 month
- [dbr](https://github.com/gocraft/dbr) -  provides additions to Go's database/sql for super fast performance and convenience. Last version v2.7.3 launched 24.12.2021, switches versions ~ once per 6 month. 
- [Squirrel](https://github.com/Masterminds/squirrel) - simple and fluent SQL queries builder from composable parts. Last version v1.5.3 launched 20.05.2022. Development complete, bug fixes will still be merged (slowly). Bug reports are welcome, but I will not necessarily respond to them. If another fork (or substantially similar project) actively improves on what Squirrel does, let me know and I may link to it here.
- [SQLR](https://github.com/elgris/sqrl) - fat-free version of squirrel - fluent SQL generator for Go. Actually a merge of squirrel and dbr ideas. No published versions at all, last commit on 28.07.2021. Looks stale.
- [Gocu](https://github.com/doug-martin/goqu) -  An expressive SQL builder and executor. Last version v9.18.0 launched 17.10.2021, before this v9.18.0 about a version per 2 months.

Not really a builders but are the SQL helpers and ORM packages:

- [GORP](https://github.com/go-gorp/gorp) - provides a simple way to marshal Go structs to and from SQL databases. It uses the database/sql package, and should work with any compliant database/sql driver
- [GORM](https://github.com/go-gorm/gorm) - fantastic ORM library for Golang, aims to be developer friendly.

# Editing this README

When you're ready to make this README your own, just edit this file and use the handy template below (or feel free to structure it however you want - this is just a starting point!).  
