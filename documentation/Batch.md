# Batch store

## Introduction

The batch store definitions are used to define predefined sequenze of database calls.

The parameter given in REST call can be passed to the batch sequence.

## Example of batch SQL

If you `https://<url>/rest/view/picview/*/tagname=holiday`

```SQL
SELECT *
FROM pictures p,
	(
	SELECT checksumpicture
	FROM picturetags
	WHERE tagname LIKE '<tagname>') c
WHERE
	p.checksumpicture = c.checksumpicture
	AND markdelete = false
```
To be continued ...
