# Azure EventHub Tools
This repo has a couple of tools to interact with Azure EventHub and was created because I
needed to stress test a couple of systems in the project I'm currently working on.

My initial goal was to read at least 1 million messages/day and write at least 100 messages/second,
so I can do more meaningful tests.

An important note: Although BadgerDb can handle several terabytes of data and these tools have a decent performance,
I don't believe it's ready to be used as some production-ready tool. Great to do stress tests, if your application uses
Azure EventHub.

Another thing: When starting up, this tools will try to create the required folders (all defined in the 
configuration). If it does not have enough permissions, please create the folders first and run the app later.

# Tools
## Hub Send
This tool will read a directory and send a message for every file there.
Currently, there's two send modes:
1. Buffered: Will read and save every message on database and only after processing all of them, will start sending. 
2. Unbuffered (default): As soon as the tool reads a file, will try to send a message.

The main advantage of buffering messages is that you can re-send them how many times you need.


## Hub Read
This tool reads messages from an EventHub. For this to work, you also need to inform a consumer 
group - event if it's ```$Default```.

Every message read will be saved logged in the database. You can also dump the messages to disk as 
soon as they are read (but doing this will greatly reduce the read speed).

A third option is to read the messages and dump to file just the ones that contains any of the filters you pass. 
(more about it on the configuration section)


## Hub Export
This last tool is more of a companion to **Hub Send**. It will read the database and save to disk the messages 
that were logged. Everytime you run it, will resume exporting the messages and/or export again any files that were 
deleted.

By default, every exported message is saved under a sub-folder with the current date.


# Benchmark
All benchmarks were made using about 500k messages of 7kb each.

## Hub Send
1. Buffered: 765 messages/second (about 66 million messages/day)
2. Unbuffered: 445 messages/second (about 38.5 million messages/day)

## Hub Read
1. Read:330 messages/second (about 28.5 million messages/day)
2. Read + Export: 25 messages/second (about 2.1 million messages/day)

## Hub Export
1. Export: 25 messages/second (about 2.1 million messages/day)


# How to use
To use any of these tools, you must have a configuration file. It can be the same configuration file for all of them.
If the configuration file is named ```default.config.json```, you don't need to specify the name. The tool will try 
to use it.

## Command line syntax
```shell
hubsend.exe [-config=<configuration file>]
```
```shell
hubread.exe [-config=<configuration file>]
```
```shell
hubexport.exe [-config=<configuration file>]
```

## To (implicitly) use the default config file
```shell
hubsend.exe
```
```shell
hubread.exe
```
```shell
hubexport.exe
```


## To (explicitly) use the default config file:
```shell
hubsend.exe -config=./default.config.json
```
```shell
hubread.exe -config=./default.config.json
```
```shell
hubexport.exe -config=./default.config.json
```

## To specify a configuration file:
```shell
hubsend.exe -config=./hubsend-specific.config.json
```
```shell
hubread.exe -config=./hubread-specific.config.json
```
```shell
hubexport.exe -config=./hubexport-specific.config.json
```

# Configuration
Every optional parameter can be omitted in the config file, and you can create a single file ready to be used by 
all 3 tools.

Note that all paths that have default values are relative to the application location and all optional parameters can be omitted.


## Full configuration file with default value:
```json
{
  "EventHubConnString": "<REQUIRED! NO DEFAULT!>",
  "entityPath": "<REQUIRED! NO DEFAULT!>",
  "skipGetRuntimeInfo": false,
  "badgerConfig": {
    "verboseMode": false,
    "badgerSkipCompactL0OnClose": false,
    "badgerValueLogFileSize": 10485760,
    "outboundBaseDir": ".app-data/badger-db/.outbound",
    "outboundDir": ".app-data/badger-db/.outbound/dir",
    "outboundValueDir": ".app-data/badger-db/.outbound/value-dir",
    "inboundBaseDir": ".app-data/badger-db/.inbound",
    "inboundDir": ".app-data/badger-db/.inbound/dir",
    "inboundValueDir": ".app-data/badger-db/.inbound/value-dir"
  },
  "inboundConfig": {
    "consumerGroup": "<REQUIRED! NO DEFAULT!>",
    "partitionId": 0,
    "inboundFolder": ".inbound",    
    "readToFile": false,
    "ignoreCheckpoint": false,
    "dumpContentOnly": false,
    "dumpFilter": []
  },
  "outboundConfig": {
    "outboundFolder": ".outbound",
    "buffered": false,
    "justSendBuffered": false,
    "ignoreStatus": false
  }
}
```

## Minimum required configuration file
### For reading from EventHub
```json
{
  "EventHubConnString": "<REQUIRED! NO DEFAULT!>",
  "entityPath": "<REQUIRED! NO DEFAULT!>",
  "inboundConfig": {
    "consumerGroup": "<REQUIRED! NO DEFAULT!>"
  }
}  
```

### For sending messages to EventHub
```json
{
  "EventHubConnString": "<REQUIRED! NO DEFAULT!>",
  "entityPath": "<REQUIRED! NO DEFAULT!>"
} 
```

## Other configuration file examples
### Sending messages to EventHub using a custom outbound folder.
#### Sending each message only 1x.
```json
{
  "EventHubConnString": "<REQUIRED! NO DEFAULT!>",
  "entityPath": "<REQUIRED! NO DEFAULT!>",
  "outboundConfig": {
    "outboundFolder": "c:\\messages_to_send"
  }  
} 
```

#### Sending every message everytime hubsend is executed.
```json
{
  "EventHubConnString": "<REQUIRED! NO DEFAULT!>",
  "entityPath": "<REQUIRED! NO DEFAULT!>",
  "outboundConfig": {
    "outboundFolder": "c:\\messages_to_send",
    "ignoreStatus": true
  }  
} 
```


## Understanding the config file.
### Root section
1. **EventHubConnString**: This is required for **Hub Read** and **Hub Send**. It's the connection string that will be used to connect to EventHub
2. **entityPath**: This is required for **Hub Read** and **Hub Send**. It's the event hub that will be targeted. If it's already on the connection string, this parameter will be ignored.
3. **skipGetRuntimeInfo**: Optional with default of '```false```'. Only affects **Hub Read** and **Hub Send**, because they're the ones that use EventHub and when true, will skip the routine to get information about the connection and the available partitions. Nothing really critical or life-saving is displayed, but I like it.
4. **badgerConfig**: Object containing base configurations specific for the database (BadgerDb).
5. **inboundConfig**: Object with configurations for reading messages.
6. **outboundConfig**: Object with configurations for writing messages.


### InboundConfig Section
1. **consumerGroup**: Required even if you want to use ```$Default```.
2. **partitionId**: Optional with default of ```0```. Which partition will be used while reading messages.
3. **inboundFolder**: Optional with default of ```.inbound```. This is where any messages read will be saved.
4. **readToFile**: Optional with default of ```false```. If true, will save the message to disk as soon as it's read.
5. **ignoreCheckpoint**: Optional with default of ```false```. If true, will process every available message every time the tool is executed. By default, every message will be processed only once.
6. **dumpContentOnly**: Optional with default of ```false```. If true, will save to disk only the content of the message. By default, event information is also saved.
7. **dumpFilter**: Optional with default of null/empty. If the message content contains **any** of the strings passed in this array, it will be saved to disk. The comparison is case-insensitive.


### OutboundConfig Section
1. **outboundFolder**: Optional with default of ```.outbound```. This is there the files that will be sent as messages lives.
2. **buffered**: Optional with default of ```false```. If true, will buffer (save to database) all messages first and then send them.
3. **justSendBuffered**: Optional with default of ```false```. If true, will skip buffering and just send the messages that were already saved to the buffer.
4. **ignoreStatus**: Optional with default of ```false```. If true, will send every message again. By default, once a message is sent, **Hub Send** will not try to send it again.


# Attention Mac users!
You might encounter an error saying "<application> cannot be opened because the developer cannot be verified".
 
If you face that problem, follow the steps bellow to solve it:
1. Open **Finder**.
2. Locate the app you’re trying to open.
3. **Control+Click** the app.
4. Select **Open**.
5. Click **Open**.
6. The app should be saved as an exception in your security settings, allowing you to open it in the future.

# Todo
1. Create tests
2. Refactor
