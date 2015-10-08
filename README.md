# quicklog

Simple log aggregator, parser.

Example configurations are in examples/

## Installing everything:

go install github.com/sheenobu/quicklog/cmd/...

## Running

File config:

	$ quicklog -configFile quicklog.josn

Etcd config:

	$ quicklog -etcdEndpoints "http://127.0.0.1:4001" -instanceName "myQuicklog"

## Commands

### quicklog

runs the log input, filter, and output. Can be configured via configfile or etcd.

### ql2etcd

Loads a JSON quicklog config file into etcd under a specific instance name

### qlsearch

proof-of-concept search client for the bleve output

## components

### ql

Main engine

## inputs 

 * stdin - Read 'message' field from standard input
 * tcp - Read 'message' field from TCP input
 * nats - Read the entire log entry from a nats queue

### filters

 * uuid - Add a randomly generated uuid to the message data, as a custom field
 * uppercase - Uppercases the 'message' field

### outputs

 * stdout - Writes the 'message' field to standard out
 * debug - Writes the 'message' field to standard out, plus each additional field on the log entry
 * bleve - Writes all the field data to a bleve index. (experimental)
 * nats - Writes the entire log entry entry to a nats queue

### config 

JSON or etcd based configuration

