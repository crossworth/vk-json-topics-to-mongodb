##  VK JSON Topics to MongoDB

This is a companion tool for [VK Group Backup](https://github.com/crossworth/vk-group-backup).
It will read a folder with VK Topics in JSON format and save/update a MongoDB server.


### Usage

This is a command line app.

**Windows**
```bash
json-to-mongo-windows-amd64.exe -folder=backup -mongo=mongodb://user:password@127.0.0.1 -database=my-db -collection=topics
```

**Linux**
```shell script
json-to-mongo-linux-amd64 -folder=backup -mongo=mongodb://user:password@127.0.0.1 -database=my-db -collection=topics
```
