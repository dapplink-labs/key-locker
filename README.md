<!--
parent:
  order: false
-->

<div align="center">
  <h1> Key Locker Repo </h1>
</div>

<div align="center">
  <a href="https://github.com/savour-labs/key-locker/releases/latest">
    <img alt="Version" src="https://img.shields.io/github/tag/savour-labs/key-locker.svg" />
  </a>
  <a href="https://github.com/savour-labs/key-locker/blob/main/LICENSE">
    <img alt="License: Apache-2.0" src="https://img.shields.io/github/license/savour-labs/key-lockersvg" />
  </a>
  <a href="https://pkg.go.dev/github.com/savour-labs/key-locker">
    <img alt="GoDoc" src="https://godoc.org/github.com/savour-labs/key-locker?status.svg" />
  </a>
</div>

key-locker is a key manager for social recovery wallet and mpc wallet private key.

**Tips**: need [Go 1.18+](https://golang.org/dl/)

## Install

### Install dependencies
```bash
go mod tidy
```
### build
```bash
go build or go install locker
```

### start 

#### 1. create database

```
CREATE DATABASE keylocker DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_general_ci;
```

#### 2. build project

```
go build
```

#### 3. migrate database table

```
./locker init
```

#### 4. start web server

```bash
./locker start web
```

#### 5. start rpc server
```bash
./locker start rpc
```

## Contribute

### 1.fork repo

fork key-locker to your github

### 2.clone repo

```bash
git@github.com:guoshijiang/key-locker.git
```

### 3. create new branch and commit code

```bash
git branch -C xxx
git checkout xxx

coding

git add .
git commit -m "xxx"
git push origin xxx
```

### 4.commit PR

Have a pr on your github and submit it to the key-locker repository

### 5.review 

After the key-locker code maintainer has passed the review, the code will be merged into the key-locker repo. At this point, your PR submission is complete
