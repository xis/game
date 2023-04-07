# Game

## Actions in the Service

1. Login

The login action is used to authenticate the user. It returns a JWT token that is used to authorize the user in the other actions.

2. Register

The register action is used to create a new user.

3. Get Leaderboard

The get leaderboard action is used to get the latest leaderboard of the game.

4. Submit User Score

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
