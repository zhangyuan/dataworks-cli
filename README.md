# dataworks-helper

A tool to help data engineers and data architects to know about the stuffs locked in Dataworks more clearly and efficiently.

**FYI: Dataworks is a service provided by Aliyun. This project is not associated with the Dataworks nor Aliyun.**

## Usage

### Prepare the `.env` file

Add the following env variables to `.env` file:

```env
ACCESS_KEY_ID=
ACCESS_KEY_SECRET=
DATAWORKS_PROJECT_ID=
DATAWORKS_ENDPOINT=dataworks.cn-beijing.aliyuncs.com # or other endpoint
```

### Get the file list by file types

```bash
dataworks-helper files list -t 10,23  -o files/manifest.json
```


### Download the files from the file list

```bash
dataworks-helper files list -t 10,23  -o files/manifest.json
```

### Get the table list by data source

```bash
dataworks-helper files download -i files.json -o ./files
```

### Get the DI sync tasks

```bash
dataworks-helper di list-sync-tasks -s data-source-name -o di-sync-tasks.json
```
