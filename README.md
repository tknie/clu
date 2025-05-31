# Clu REST API

## Introduction

This application named CluAPI is a generic database interface REST server. It can be used to query, search and modify database data.
There is no complex infrastructure inbetween, a direct connection to the database. A mechanism of authentication and authorization exists inbetween REST query and database access. Some predefined database queries can be defined in an batch database record.  The **batch** queries can be defined providing extra rest parameters. That's different to standard database views.

There should be no difference if the database is a SQL or a NoSQL database.

A list of data records should be able to be inserted or updated in one REST API call.
**May be as transaction or in an atomic matter**.

Based on database layer [Flynn at https://github.com/tknie/flynn](https://github.com/tknie/flynn) the database is queryed and updated.

* Read, Search, Insert, Update and Delete of data records
* Search by passing database specific queries to be flexible for complex queries
* Large object support of big data
  * binary data queries
  * image data queries
  * video data queries
* Support creating batch jobs for database-specific tasks like SQL scripts or a complex query

For possible detailed information about the Flynn layer, have a look at the API documentation. It can be referenced here: <https://godoc.org/github.com/tknie/flynn>

## Release information

The version 1.0.0 is released. It is a stable version used for various projects giving access to audit data stored in the database. It is stable and no problem known.
You can download the Docker image at <https://hub.docker.com/r/thknie/cluapi> or download it in Docker/Podman with

```sh
docker pull thknie/cluapi
```

## Authentication

There are several possibilities to use authentication

* Create a realm file containing user and password
* Use system user and password to authenticate
* Using LDAP password authentication
* Using a SQL query to authenticate using the SQL database

All these configuration are used from <https://github.com/tknie/services>.

## Authorization

You can use database roles and users authorization restriction if you use database users management.
CluAPI provides a query entry role management to obmit user specific queries for read- or write-access.
In addition you can define authorization access to some resources like tables or views.

For example in `users.yaml` you can define read or write restriction for table, views or batch scripts

* all tables and views are restricted using to prefix
* the prefix ^ restrict to batch processing tasks or complex queries
* the prefix < allows read/download file permissions
* the prefix > allows write/upload file permissions

## Batch store usage

You can define some predefined queries which are generated using the parameters given in an REST request.

Detail documentation about Batch store definition is found [here](documentation/Batch.md).

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

## REST API information

To evaluate the REST API look at [Swagger Editor with clu REST API informations.](https://editor.swagger.io/?url=https://raw.githubusercontent.com/tknie/clu/refs/heads/main/swagger/openapi-restserver.yaml)

## Download/Upload files

To exchange files into local file system like logs or extra information it is possible to upload or download files from a specific location.

Below a download location is defined which can be referenced using the download path.

```yaml
fileTransfer:
  Admin:
    role: xxx
  directories:
    directory:
      - name: download
        location: ${HOME}/Downloads
      - name: tmp
        location: ${CURDIR}/tmp
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
