package domain

import (
	"github.com/dibyendu/Authentication-Authorization/lib/constants"
	"github.com/dibyendu/Authentication-Authorization/lib/errs"
	"github.com/dibyendu/Authentication-Authorization/lib/utility"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

type UserRepositoryDb struct {
	client     *mongo.Client
	database   string
	collection map[string]string
}

func NewUserRepositoryDb(dbClient *mongo.Client, database string, collection map[string]string) UserRepositoryDb {
	return UserRepositoryDb{
		client:     dbClient,
		database:   database,
		collection: collection,
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
		return nil, errs.NewValidationError("password hashing failed"+ err.Error())
	}
	request.Password = password
	err = n.client.Database(n.database).Collection(n.collection["user"]).FindOne(ctx, filter).Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If document doesn't exist, insert it
			result, err := n.client.Database(n.database).Collection(n.collection["user"] ).InsertOne(ctx, request)
			if err != nil {
				log.Println("error inserting user log: " + err.Error())
				return nil,  errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
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

func(n UserRepositoryDb) SignIn(ctx context.Context, request CreateUserRequest) (*CreateUserResponse, *errs.AppError){
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
		return nil, errs.NewValidationError("password hashing failed"+ err.Error())
	}
	request.Password = password
	err = n.client.Database(n.database).Collection(n.collection["user"]).FindOne(ctx, filter).Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If document doesn't exist, insert it
			result, err := n.client.Database(n.database).Collection(n.collection["user"]).InsertOne(ctx, request)
			if err != nil {
				log.Println("error inserting user log: " + err.Error())
				return nil,  errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
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

func(d UserRepositoryDb)IsEmailExists(ctx context.Context, email string) (*CreateUserResponse, *errs.AppError){
	var(
		filter = bson.M{}
		userDetail CreateUserResponse
	)
	filter["email"] = email
	result := d.client.Database(d.database).Collection(d.collection["user"]).FindOne(ctx, filter).Decode(&userDetail)
	if result != nil {
		if result == mongo.ErrNoDocuments{
			log.Println("there is not exists this email: " + result.Error())
			return nil,  errs.NewNotFoundError("not found this email")
		}
		log.Println("error fetching user log: " + result.Error())
		return nil,  errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
	}

	return &userDetail, nil
}

func(d UserRepositoryDb) GetUser(ctx context.Context, req GetUserRequest) (*GetUserResponse, *errs.AppError){
	var(
		filter = bson.M{}
		userDetail GetUserResponse
	)
	objectId , err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, errs.NewValidationError("unable to convert id to objectId")
	}
	filter["_id"] = objectId
	result := d.client.Database(d.database).Collection(d.collection["user"]).FindOne(ctx, filter).Decode(&userDetail)
	if result != nil {
		if result == mongo.ErrNoDocuments{
			log.Println("there is not exists the user detail with this id: " + result.Error())
			return nil,  errs.NewNotFoundError("not found user details for this id")
		}
		log.Println("error fetching user log: " + result.Error())
		return nil,  errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
	}

	return &userDetail, nil
}