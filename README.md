# Game

This is a simple game service that is used to authenticate users, get the leaderboard and submit user scores.


1. [Game](#game)
2. [Actions in the Service](#actions-in-the-service)
   1. [Login](#1-login)
   2. [Register](#2-register)
   3. [Get Leaderboard](#3-get-leaderboard)
   4. [Submit User Score](#4-submit-user-score)
3. [Running the Service](#running-the-service)
   1. [Clone the repository](#1-clone-the-repository)
      1. [Using SSH](#using-ssh)
      2. [Using HTTPS](#using-https)
   2. [Install dependencies](#2-install-dependencies)
   3. [Prepare the environment](#3-prepare-the-environment)
   4. [Run the service](#4-run-the-service)

## Actions in the Service

## 1. `Login`
The login action is used to authenticate the user. It returns a JWT token that is used to authorize the user in the other actions.

## 2. `Register`
The register action is used to create a new user.

## 3. `Get Leaderboard`
The get leaderboard action is used to get the latest leaderboard of the game.

## 4. `Submit User Score`
The submit user score action is used to submit the user score to the game. Triggered when a match is finished. If the user score is higher than the previous score, the user score is updated. If not the user score is not updated.

## Running the Service

### 1. Clone the repository

#### Using SSH
```bash
git clone git@github.com:xis/game.git
```

#### Using HTTPS
```bash
git clone https://github.com/xis/game.git
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Prepare the environment

You can find the required environment variables in the `.env.example` file.

### 3. Run the service

```bash
go run cmd/main.go
```
