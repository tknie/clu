# Clu REST API (_draft_)

## Introduction

The `github.com/tknie/clu` package contains an REST server to get and modify database data.
No complex infrastructure inbetween, a direct connection to the database.

There should be no difference if the database is a SQL or a NoSQL database.

A list of data recrods should be able to be inserted or updated in one call. May be as transaction or in an atomic matter.

In advance real main database functionality should be contained like:

* Create database tables
* Read, Search, Insert, Update and Delete of data records
* Search by passing database specific queries to be flexible for complex queries
* One-to-One mapping of `Golang` structures to database tables
* Large object support of big data
* Support creating batch jobs for database-specific tasks like SQL scripts
* Create index or other enhancements on database configuration

For details have a look at the API documentation. It can be referenced here: <https://godoc.org/github.com/tknie/flynn>

## Authentication

There are several possibilities to use authentication

* Create a realm file containing user and password
* Using LDAP or system password authentication
* Using a SQL query to authenticate to SQL database

All these configuration are located in <https://github.com/tknie/services>.

## Authorization

You can use SQL authorization limits if you use SQL users.
CLU provides a query entry role management to obmit user specific queries for read- or write-access.

## Example of Clu usage

### Query records in database

```http
Accept: application/json
Authorization: Base <base64>
GET http://localhost:8030/rest/view/Albums/ID,Title,published?limit=0&orderby=published:ASC
```

### Update records in database

```http
Accept: application/json
Authorization: Base <base64>
PUT http://localhost:8030/rest/view/Albums/ID,Title,published?limit=0&orderby=published:ASC
 {
  "Records": [
    {
      "id": "18",
      "title": "Der Ostergruss"
    },
 }
```

### Insert records in database

```http
Accept: application/json
Authorization: Base <base64>
PUSH http://localhost:8030/rest/view/Albums/ID,Title,published?limit=0&orderby=published:ASC
```

## Check List

Feature | Ready-State | Description
---------|----------|---------
 Login database | :heavy_check_mark: | Draft
 Query record | :heavy_check_mark: | Draft
 Search record | :heavy_check_mark: | Draft
 Insert record | :heavy_check_mark: | Draft
 Delete record | :heavy_check_mark: | Draft
 Load images out of database | :heavy_check_mark: | Draft
 Load videos out of database | :heavy_check_mark: | Draft
 Load binaries out of database |:heavy_check_mark: | Draft
 Insert Large Object (Image, binary or others) | :heavy_check_mark: | Draft(not recommended if images is to big)
 Create table |  | Draft
 Insert database |  | Draft
 Work with batch queries | :heavy_check_mark: | draft
 Complex search queries (common to SQL or NonSQL databases) |  | planned
