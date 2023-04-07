package main

import (
	"context"
	"net"
	"time"

	"github.com/caarlos0/env/v8"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	grpccontroller "game/internal/controllers/grpc"
	"game/internal/domain"
	bcryptpasswordhasher "game/internal/passwordhashers/bcrypt"
	leaderboard "game/internal/proto/leaderboard/proto"
	user "game/internal/proto/user/proto"
	usermongo "game/internal/repositories/user/mongo"
	userscoreredis "game/internal/repositories/userscore/redis"
	service "game/internal/services"

	jwttokenmanager "game/internal/tokenmanagers/jwt"
)

type EnvironmentVariables struct {
	MongoURI                 string `env:"MONGO_URI,required"`
	MongoDatabaseName        string `env:"MONGO_DATABASE_NAME,required"`
	MongoUsersCollectionName string `env:"MONGO_USERS_COLLECTION_NAME,required"`
	JWTSecretKey             string `env:"JWT_SECRET_KEY,required"`
	RedisAddr                string `env:"REDIS_ADDR,required"`
	GrpcServerPort           string `env:"GRPC_SERVER_PORT,required"`
	JWTTokenTTLInHours       int    `env:"JWT_TOKEN_TTL_IN_HOURS,required"`
}

func main() {
	logger := logrus.New()

	logrus.Info("starting server...")

	environments := EnvironmentVariables{}

	if err := env.Parse(&environments); err != nil {
		logger.Fatal("failed to parse environment variables", err)
	}

	jwtTokenManager := jwttokenmanager.NewJWTTokenManager(jwttokenmanager.JWTTokenCreatorDependencies{
		SecretKey: environments.JWTSecretKey,
		TokenTTL:  time.Duration(environments.JWTTokenTTLInHours) * time.Hour,
	})

	bcryptPasswordHasher := bcryptpasswordhasher.NewBcryptPasswordHasher()

	mongoClient, err := connectToMongoDB(environments.MongoURI)
	if err != nil {
		logger.Fatal("failed to connect to MongoDB", err)
	}

	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			logger.Fatal("failed to disconnect from MongoDB", err)
		}
	}()

	usersCollection := mongoClient.
		Database(environments.MongoDatabaseName).
		Collection(environments.MongoUsersCollectionName)

	mongoUserRepository := usermongo.NewMongoUserRepository(usermongo.MongoUserRepositoryDependencies{
		UsersCollection: usersCollection,
	})

	redisUserScoreRepository, err := createRedisUserScoreRepository(environments.RedisAddr, mongoUserRepository)
	if err != nil {
		logger.Fatal("failed to create redis user score repository", err)
	}

	userService := service.NewUserService(service.UserServiceDependencies{
		UserRepository: mongoUserRepository,
		TokenManager:   jwtTokenManager,
		PasswordHasher: bcryptPasswordHasher,
	})

	userController := grpccontroller.NewUserController(grpccontroller.UserControllerDependencies{
		UserService: userService,
		Logger:      logger,
	})

	leaderboardService := service.NewLeaderboardService(service.LeaderboardServiceDependencies{
		UserRepository:      mongoUserRepository,
		UserScoreRepository: redisUserScoreRepository,
	})

	leaderboardController := grpccontroller.NewLeaderboardController(grpccontroller.LeaderboardControllerDependencies{
		LeaderboardService: leaderboardService,
		Logger:             logger,
	})

	unaryInterceptor := grpccontroller.NewUnaryInterceptor(grpccontroller.UnaryInterceptorDependencies{
		TokenManager: jwtTokenManager,
		AuthorizedMethodNames: []string{
			"/leaderboard.LeaderboardService/SubmitUserScore",
			"/leaderboard.LeaderboardService/GetLeaderboard",
		},
	})

	server := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor.Intercept),
	)

	user.RegisterUserServiceServer(server, userController)
	leaderboard.RegisterLeaderboardServiceServer(server, leaderboardController)

	listener, err := net.Listen("tcp", ":"+environments.GrpcServerPort)
	if err != nil {
		logger.Fatal("failed to listen", err)
	}

	logger.Info("server started")

	err = server.Serve(listener)
	if err != nil {
		logger.Fatal("failed to serve", err)
	}
}

func connectToMongoDB(
	mongoURI string,
) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func createRedisUserScoreRepository(
	redisAddr string,
	userRepository domain.UserRepository,
) (domain.UserScoreRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return userscoreredis.NewRedisUserScoreRepository(
		userscoreredis.RedisUserScoreRepositoryDependencies{
			Client:         client,
			UserRepository: userRepository,
		},
	), nil
}
