# quicklog

example configurations are in examples/

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

JSON or (soon) etcd based configuration

