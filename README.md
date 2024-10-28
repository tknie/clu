# Clu REST API

## Introduction

The `github.com/tknie/clu` package contains an REST server to get and modify database data.
No complex infrastructure inbetween, a direct connection to the database.

There should be no difference if the database is a SQL or a NoSQL database.

A list of data records should be able to be inserted or updated in one REST API call.
**May be as transaction or in an atomic matter**.

Based on database layer at `https://github.com/tknie/flynn` the database should be accessed

* Read, Search, Insert, Update and Delete of data records
* Search by passing database specific queries to be flexible for complex queries
* Large object support of big data
  * binary data queries
  * image data queries
  * video data queries
* Support creating batch jobs for database-specific tasks like SQL scripts

For possible future details have a look at the API documentation. It can be referenced here: <https://godoc.org/github.com/tknie/flynn>

## Authentication

There are several possibilities to use authentication

* Create a realm file containing user and password
* Using LDAP or system password authentication
* Using a SQL query to authenticate to SQL database

All these configuration are used from <https://github.com/tknie/services>.

## Authorization

You can use SQL database authorization restriction if you use SQL users.
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
 Insert Large Object (Image, binary or others) | :heavy_check_mark: | Draft (not recommended if images is to big)
 Create table |  | Draft
 Insert database |  | Draft
 Work with predefined batch queries | :heavy_check_mark: | Draft
 Complex search queries (common to SQL or NonSQL databases) |  | planned
