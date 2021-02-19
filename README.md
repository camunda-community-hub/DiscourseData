<img src="https://img.shields.io/badge/Camunda%20DevRel%20Project-Created%20by%20the%20Camunda%20Developer%20Relations%20team-0Ba7B9">

# DiscourseData

This is a project to read data from the Discourse API. The ultimate goal is to leverage that data into Orbit in order to surface community members that might need extra attention, recognition, etc.

## Reading from Discourse

This requires installing the [discourse-data-explorer](https://github.com/discourse/discourse-data-explorer) plugin. There are a lot of included, and many more open source, queries that can be run against your Discourse instance. I chose 4 that suited our needs. You will see them defined in the source:

```go
const lurkers = "5"
const top50 = "6"
const top50Posters = "4"
const mostSolutions = "7"

var queries = [...]string{lurkers, top50, top50Posters, mostSolutions}
```
Those numerical values are important, and the way to determine them is to load one of the queries you want to run in your browser, then look in the URL bar for the query number.

I gave them human readable names.

Next you will need a Discourse API Key. You will add that to the source file along with your API username.

These queries will run, sequentially, in an endless loop. In order to not overwhelm your discourse server there are delays built in. Each of these queries will get a certain about of data, but then it will run through the results of each one, pick out the User IDs, and do another query per User ID, to get detailed information about the users identified.

The queries that we run return the following JSON object of data:

```go
type QueryResult struct {
	Success     bool          `json:"success"`
	Errors      []interface{} `json:"errors"`
	Duration    float64       `json:"duration"`
	ResultCount int           `json:"result_count"`
	Params      struct {
                Page string `json:"page"`
              } `json:"params"`
	Columns      []string        `json:"columns"`
	DefaultLimit int             `json:"default_limit"`
	Rows         [][]interface{} `json:"rows"`
}
```
Ultimately all of the queried data is stored in a super-object:

```go
// Results the actual data in the table
type Results struct {
	ID           int       `json:"user_id"`
	Username     string    `json:"username"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	LastSeen     time.Time `json:"last_seen_at"`
	IPAddress    string    `json:"ip_address"`
	PostCount    int       `json:"post_count"`
	PostsCreated int       `json:"posts_created"`
	PostsRead    int       `json:"posts_read"`
	SolvedCount  int       `json:"solved_count"`
	AverageScore float64   `json:"average_score"`
}
```

All of this data is finally returned as an array of these objects.

## Putting the data somewhere

I chose to put this data into InfluxDB v2.0 at this point. Much of it is used as `fields` rather than `tags` because we want to be able to `group-by` some of this data in order to do anomaly detection.

In order to put this data into InfluxDB you will need:

1) And InfluxDB server accessible from wherever this script is run
2) An InfluxDB `write-token` in order to insert data
3) The `bucket` and `organization` names where you want to store the data


## Anomaly detection

Coming soon

## Anomaly notification

Coming soon

## Points assignments

Coming Soon

## Integration with Camunda Platform

Coming soon

## Integration with the Orbit Model

Coming soon
