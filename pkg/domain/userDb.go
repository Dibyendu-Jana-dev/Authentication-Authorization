package domain

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dibyendu/Authentication-Authorization/lib/constants"
	"github.com/dibyendu/Authentication-Authorization/lib/errs"
	"github.com/dibyendu/Authentication-Authorization/lib/logger"
	"github.com/dibyendu/Authentication-Authorization/lib/utility"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryDb struct {
	client     *mongo.Client
	redisClient *redis.Client
	database   string
	collection map[string]string
}

func NewUserRepositoryDb(dbClient *mongo.Client, redisClient *redis.Client, database string, collection map[string]string) UserRepositoryDb {
	return UserRepositoryDb{
		client:      dbClient,
		redisClient: redisClient,
		database:    database,
		collection:  collection,
	}
}

func(n UserRepositoryDb) CreateUser(ctx context.Context, request CreateUserRequest) (*CreateUserResponse, *errs.AppError){
	var(
		filter = bson.M{}
		data CreateUserResponse
	)

	filter = bson.M{
		"name": request.Name,
		"email": request.Email,
	}
	password, err := utility.HashPassword(request.Password)
	if err != nil {
		logger.Error("password hashing failed"+ err.Error())
		return nil, errs.NewValidationError("password hashing failed"+ err.Error())
	}
	request.Password = password
	err = n.client.Database(n.database).Collection(n.collection["user"]).FindOne(ctx, filter).Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If document doesn't exist, insert it
			result, err := n.client.Database(n.database).Collection(n.collection["user"] ).InsertOne(ctx, request)
			if err != nil {
				logger.Error("error inserting user log: " + err.Error())
				return nil,  errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
			}
			if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
				data.Id = oid
				data.Name = request.Name
				data.Role = request.Role
				data.Email = request.Email
			}
			value, jErr := json.Marshal(data)
			if jErr != nil {
				logger.Error("error marshalling")
				return nil, errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
			}
			_, err = n.redisClient.Set(ctx, data.Id.Hex(), value, time.Duration(5)*time.Minute).Result()
			if err != nil {
				logger.Error("error in set value to redis server for create user: "+err.Error())
				return nil, errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
			}
			return &data, nil
		}
	}
	return nil, &errs.AppError{
		Code:    http.StatusConflict,
		Message: constants.USER_ALREADY_EXISTS,
	}
}

func (n UserRepositoryDb) SignIn(ctx context.Context, request CreateUserRequest) (*CreateUserResponse, *errs.AppError) {
	var (
		filter = bson.M{}
		data   CreateUserResponse
	)

	filter = bson.M{
		"name":  request.Name,
		"email": request.Email,
	}
	password, err := utility.HashPassword(request.Password)
	if err != nil {
		logger.Error("password hashing failed" + err.Error())
		return nil, errs.NewValidationError("password hashing failed" + err.Error())
	}
	request.Password = password
	err = n.client.Database(n.database).Collection(n.collection["user"]).FindOne(ctx, filter).Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If document doesn't exist, insert it
			result, err := n.client.Database(n.database).Collection(n.collection["user"]).InsertOne(ctx, request)
			if err != nil {
				logger.Error("error inserting user log: " + err.Error())
				return nil, errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
			}
			if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
				data.Id = oid
				data.Name = request.Name
				data.Role = request.Role
				data.Email = request.Email
			}
			return &data, nil
		}
	}
	return nil, &errs.AppError{
		Code:    http.StatusConflict,
		Message: constants.USER_ALREADY_EXISTS,
	}
}

func (d UserRepositoryDb) IsEmailExists(ctx context.Context, email string) (*CreateUserResponse, *errs.AppError) {
	var (
		filter     = bson.M{}
		userDetail CreateUserResponse
	)
	filter["email"] = email
	result := d.client.Database(d.database).Collection(d.collection["user"]).FindOne(ctx, filter).Decode(&userDetail)
	if result != nil {
		if result == mongo.ErrNoDocuments {
			logger.Warn("there is not exists this email: " + result.Error())
			return nil, errs.NewNotFoundError("not found this email")
		}
		logger.Error("error fetching user log: " + result.Error())
		return nil, errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
	}

	return &userDetail, nil
}

func(d UserRepositoryDb) GetUser(ctx context.Context, req GetUserRequest) (*GetUserResponse, *errs.AppError){
	var(
		filter = bson.M{}
		userDetail GetUserResponse
	)
	val, err := d.redisClient.Get(ctx, req.Id).Result()
	if err != nil {
		logger.Error("error fetching user from redis: " + err.Error())
		return nil, errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
	}
	if val == "" {
		objectId, err := primitive.ObjectIDFromHex(req.Id)
		if err != nil {
			logger.Error("error converting dtring to objectid: " +err.Error())
			return nil, errs.NewValidationError("unable to convert id to objectId")
		}
		filter["_id"] = objectId
		result := d.client.Database(d.database).Collection(d.collection["user"]).FindOne(ctx, filter).Decode(&userDetail)
		if result != nil {
			if result == mongo.ErrNoDocuments {
				logger.Warn("there is not exists the user detail with this id: " + result.Error())
				return nil, errs.NewNotFoundError("not found user details for this id")
			}
			logger.Error("error fetching user log: " + result.Error())
			return nil, errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
		}
		return &userDetail, nil
	} else {
		var redisUser GetUserResponse
		if err := json.Unmarshal([]byte(val), &redisUser); err != nil {
			logger.Error("error unmarshalling user: " +err.Error())
			return nil, errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
		}
		return &redisUser, nil
	}
	//return &userDetail, nil
}
