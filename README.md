# quicklog

[![Join the chat at https://gitter.im/sheenobu/quicklog](https://badges.gitter.im/sheenobu/quicklog.svg)](https://gitter.im/sheenobu/quicklog?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Simple log aggregator, parser.

Example configurations are in examples/

## Installing everything:

	$ go install github.com/sheenobu/quicklog/cmd/...

## Embedding

You can embed quicklog, and only include the parts you need. See /examples/embedded/main.go

## Running

File config:

	$ quicklog -configFile quicklog.json

Etcd config:

	$ quicklog -etcdEndpoints "http://127.0.0.1:4001" -instanceName "myQuicklog"

## Commands

### quicklog

runs the log input, filter, and output. Can be configured via configfile or etcd.

### ql2etcd

Loads a JSON quicklog config file into etcd under a specific instance name

## components

### ql

Main engine

## inputs

 * stdin - Read from standard input
 * tcp - Read from TCP input
 * udp - Read from UDP input
 * nats - Read the entire log entry from a nats queue

## parsers

 * plain - parse input as a plain message
 * json - parse input as json
 * otto - parse input via javascript function
 * csv - parse input as a CSV entry (TODO)

### filters

 * uuid - Add a randomly generated uuid to the message data, as a custom field
 * uppercase - Uppercases the 'message' field
 * rename\_field - Renames a field
 * hostname - Add the current hostname as a field

### outputs

 * stdout - Writes the 'message' field to standard out
 * debug - Writes the 'message' field to standard out, plus each additional field on the log entry
 * nats - Writes the entire log entry entry to a nats queue
 * elasticsearch-http - Writes the entire JSON to elasticsearch, under the given index/type

### config

JSON or etcd based configuration, see /examples

