---
slug: '/wiki/uql'
title: 'UQL'
previous_page_title: 'Configuration'
previous_page_slug: '/wiki/configuration'
next_page_title: 'JSON'
next_page_slug: '/wiki/json'
---

## UQL

UQL (Unstructured query language) is advance query format in infinity datasource which will consolidate JSON, CSV, XML, GraphQL formats. UQL also provides ability to customize the results.

UQL is an opinionated query language designed for in-memory operations. UQL query can be formed with list of commands joined by `|`, in a line each.
Most of the times, fields are referred within double quotes and string values are referred with single quotes. UQL was inspired by kusto query language and follows similar syntax.

> UQL is still in BETA.

if your data looks like this,

```json
[
  { "id": 1, "name": { "firstName": "john", "lastName": "doe" }, "dob": "1985-01-01", "city": "chennai" },
  { "id": 2, "name": { "firstName": "alice", "lastName": "bob" }, "dob": "1990-12-31", "city": "london" }
]
```

then the following UQL query

```sql
parse-json
| extend "full name"=strcat("name.firstName",' ',"name.lastName"), "dob"=todatetime("dob")
| project "id", "full name", "dob", "date of birth"="dob"
| order by "full name" asc
```

will produce four column table (id,full name, dob, date of birth).

### Basic UQL commands

following are the basic UQL commands. All these commands are available in all the version unless specified.

### project

`project` command is used to select the columns to include in the results. If you want to select a property inside a nested object, you can use dot notation. Optionally, you can also alias the fields.

```sql
parse-json
| project "id", "name.firstName", "date of birth"="dob"
```

### project-away

`project-away` command is exactly opposite as `project`. It just drops specific columns from the data. It doesn't support alias or dot notation selector.

```sql
parse-json
| project-away "id", "city"
```

### order by

`order by` command sorts the input based on any column. sort direction should be either `asc` or `desc`

```sql
parse-json
| order by "full name" asc
```

### extend

`extend` command is similar to `project`. but instead of selecting the columns, it just adds/replace columns in existing data. `extends` expects an alias and a function.

```sql
parse-json
| extend "dob"=todatetime("dob"), "city"=toupper("city")
```

following are some of the available functions

| function keyword                 | syntax                                  | description                                                                                          | available from |
| -------------------------------- | --------------------------------------- | ---------------------------------------------------------------------------------------------------- | -------------- |
| trim                             | trim("name")                            | trims the string                                                                                     | 0.8.0          |
| trim_start                       | trim_start("name")                      | removes the space before                                                                             | 0.8.0          |
| trim_end                         | trim_end("name")                        | removes the space after                                                                              | 0.8.0          |
| tonumber                         | tonumber("age")                         | converts a string into number                                                                        | 0.8.0          |
| tostring                         | tostring("age")                         | converts a number into string                                                                        | 0.8.0          |
| todatetime                       | todatetime("age")                       | converts a datetime string into datetime                                                             | 0.8.0          |
| unixtime_seconds_todatetime      | unixtime_seconds_todatetime("dob")      | converts unix epoch s timestamp to datetime                                                          | 0.8.0          |
| unixtime_nanoseconds_todatetime  | unixtime_nanoseconds_todatetime("dob")  | converts unix epoch ns timestamp to datetime                                                         | 0.8.0          |
| unixtime_milliseconds_todatetime | unixtime_milliseconds_todatetime("dob") | converts unix epoch ms timestamp to datetime                                                         | 0.8.0          |
| unixtime_microseconds_todatetime | unixtime_microseconds_todatetime("dob") | converts unix epoch microsecond timestamp to datetime                                                | 0.8.0          |
| format_datetime                  | format_datetime("dob",'DD/MM/YYYY')     | converts datetime to a specific format                                                               | 0.8.0          |
| add_datetime                     | add_datetime("dob",'-1d')               | adds duration to a datetime field                                                                    | 0.8.0          |
| startofminute                    | startofminute("dob")                    | rounds the datetime field to the starting minute                                                     | 0.8.0          |
| startofhour                      | startofhour("dob")                      | rounds the datetime field to the starting hour                                                       | 0.8.0          |
| startofday                       | startofday("dob")                       | rounds the datetime field to the starting day                                                        | 0.8.0          |
| startofmonth                     | startofmonth("dob")                     | rounds the datetime field to the starting month                                                      | 0.8.0          |
| startofweek                      | startofweek("dob")                      | rounds the datetime field to the starting week                                                       | 0.8.0          |
| startofyear                      | startofyear("dob")                      | rounds the datetime field to the starting year                                                       | 0.8.0          |
| sum                              | sum("col1","col2")                      | sum of two or more columns                                                                           | 0.8.0          |
| diff                             | diff("col1","col2")                     | difference between two columns                                                                       | 0.8.0          |
| mul                              | mul("col1","col2")                      | multiplication of two columns                                                                        | 0.8.0          |
| strcat                           | strcat("col1","col2")                   | concatenates two or more columns                                                                     | 0.8.0          |
| parse_url                        | parse_url("col1")                       | parses the col1 as URL                                                                               | 0.8.6          |
|                                  | parse_url("col1",'pathname')            | returns the `pathname` of the URL. Options are `host`,`hash`,`origin`,`href`,`protocol` and `search` | 0.8.6          |
|                                  | parse_url("col1",'search','key1')       | returns the query string value for `key1`. 2nd arg is always `search`                                | 0.8.6          |

For example, the data `[ { "a": 12, "b" : 20 }, { "a" : 6, "b": 32} ]` and the following uql query

```sql
parse-json
| project "a", "triple"=sum("a","a","a"),"thrice"=mul("a",3), sum("a","b"),  diff("a","b"), mul("a","b")
```

wil produce the following output

```csv
a,triple,thrice,sum,diff,mul
12,36,36,32,-8,240
6,18,18,38,-26,192
```

To apply multiple transformations over a field, repeat them with the same field name. For example, the uql query `extend "name"=tolower("name"), "name"=trim("name")` will apply tolower function and then trim function over the name field.

### summarize

`summarize` command aggregates the data by a string column. summarize command expects alias, summarize by fields and summarize function. Following are the valid summarize functions.

| function keyword | syntax            | description       | available from |
| ---------------- | ----------------- | ----------------- | -------------- |
| count            | count()           | count of values   | 0.8.0          |
| sum              | sum("age")        | sum of age        | 0.8.0          |
| min              | min("population") | min of population | 0.8.0          |
| max              | max("foo")        | max of foo        | 0.8.0          |
| mean             | mean("foo")       | mean of foo       | 0.8.0          |

For example, the following data

```json
[
  { "city": "tokyo", "country": "japan", "population": 200 },
  { "city": "newyork", "country": "usa", "population": 60 },
  { "city": "oslo", "country": "usa", "population": 40 },
  { "city": "new delhi", "country": "india", "population": 180 },
  { "city": "mumbai", "country": "india", "population": 150 }
]
```

and the following uql query

```sql
parse-json
| summarize "number of cities"=count(), "total population"=sum("population") by "country"
| extend "country"=toupper("country")
| order by "total population" desc
```

will produce the output table like this

```csv
country,number of cities,total population
INDIA,2,330
JAPAN,1,200
USA,2,100
```

### parse-json

`parse-json` is the first command that parses the output as json. This is optional in case if you use JSON/GraphQL URL.

### parse-csv

`parse-csv` is used to specify that the results are in csv format

### parse-xml

`parse-xml` is used to specify that the results are in xml format

### parse-yaml

`parse-yaml` is used to specify that the results are in xml format

### count

count gives the number of results.

```sql
parse-json
| count
```

### limit

`limit` command restricts the number of results returned. For example, below query returns only 10 results

```sql
parse-json
| limit 10
```

### scope

`scope` commands sets the context of the output data. It is useful when the results are insides nested json object.

example

```json
{
  "meta": { "last-updated": "2021-08-09" },
  "count": 2,
  "users": [{ "name": "foo" }, { "name": "bar" }]
}
```

and the following uql query just results the "users" and ignores the other root level properties.

```sql
parse-json
| scope "users"
```

### mv-expand

`mv-expand` expands multi-value properties into their own records. For example, the command `mv-expand "user"="users"` over following data

```json
[
  { "group": "A", "users": ["user a1", "user a2"] },
  { "group": "B", "users": ["user b1"] }
]
```

will produce results like

```json
[
  { "group": "A", "user": "user a1" },
  { "group": "A", "user": "user a2" },
  { "group": "B", "user": "user b1" }
]
```

`mv-expand` should also work for non string arrays.

### comments

Any new line that starts with `#` will be treated as comment.
