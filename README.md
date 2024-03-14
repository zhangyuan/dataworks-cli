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

### Get the file list for SQL scripts only

```bash
dataworks-helper files list -o files.json
```

### Get the file list of all files

```bash
dataworks-helper files list-all  -o all.json
```

### Download the files from the file list

```bash
files download -i files.json -o ./code
```
