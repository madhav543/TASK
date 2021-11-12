package repo

import (
	"Task/constants"
	"Task/database"
	"Task/helpers"
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/go-chassis/openlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//
type Repository struct {
}

//
func (r *Repository) CreateRecord(record interface{}) (interface{}, error) {
	collection := database.Client.Database("Students").Collection("Records")
	res, err := collection.InsertOne(context.TODO(), record)
	if err != nil {
		openlog.Error(err.Error())
		return nil, errors.New(constants.InternalServerError)
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	outputData, outputErr := r.FetchRecordByID(id)
	return outputData, outputErr
}

//
func (r *Repository) FetchRecordByID(id string) (interface{}, error) {
	record := make(map[string]interface{})
	collection := database.Client.Database("Students").Collection("Records")
	docID, convIDErr := primitive.ObjectIDFromHex(id)
	if convIDErr != nil {
		openlog.Error(convIDErr.Error())
		return nil, errors.New(constants.InvalidRecordID)
	}
	err := collection.FindOne(context.TODO(), bson.M{"_id": docID}).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return record, errors.New(constants.NoRecords)
		}
		openlog.Error(err.Error())
		return record, errors.New(constants.InternalServerError)
	}
	return record, err
}

//
func (r *Repository) UpdateRecordByID(id string, update map[string]interface{}) (interface{}, error) {
	collection := database.Client.Database("Students").Collection("Records")
	docID, convIDErr := primitive.ObjectIDFromHex(id)
	if convIDErr != nil {
		openlog.Error(convIDErr.Error())
		return nil, errors.New(constants.InvalidRecordID)
	}
	res := collection.FindOneAndUpdate(context.TODO(), bson.M{"_id": docID}, bson.M{"$set": update})
	if res.Err() != nil {
		openlog.Error(res.Err().Error())
		return nil, errors.New(constants.InternalServerError)
	}
	result, err := r.FetchRecordByID(id)
	return result, err
}

//
func (r *Repository) DeleteRecordByID(id string) (interface{}, error) {
	collection := database.Client.Database("Students").Collection("Records")
	docID, convIDErr := primitive.ObjectIDFromHex(id)
	if convIDErr != nil {
		openlog.Error(convIDErr.Error())
		return nil, errors.New(constants.InvalidRecordID)
	}
	record := make(map[string]interface{})
	err := collection.FindOneAndDelete(context.TODO(), bson.M{"_id": docID}).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(constants.NoRecords)
		}
		openlog.Error(err.Error())
		return nil, errors.New(constants.InternalServerError)
	}
	return record, err
}

//
func (r *Repository) FetchAllRecords(filter, sort, page, size string) (interface{}, error) {
	filterMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(filter), &filterMap)
	if err != nil {
		openlog.Error(err.Error())
	}
	sortMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(sort), &sortMap)
	if err != nil {
		openlog.Error(err.Error())
	}
	pageNo, _ := strconv.ParseInt(page, 10, 64)
	sizeLen, _ := strconv.ParseInt(size, 10, 64)
	pageAndSizeDetails := map[string]interface{}{"page": pageNo, "size": sizeLen}
	options := helpers.ValidatePageAndSizeDetails(pageAndSizeDetails)
	options = options.SetSort(sortMap)
	records := make([]map[string]interface{}, 0)
	collection := database.Client.Database("Students").Collection("Records")
	// fmt.Println("|||||||||||||||||||", filterMap, sortMap, options)
	cur, err := collection.Find(context.TODO(), filterMap, options)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(constants.NoRecords)
		}
		openlog.Error(err.Error())
		// fmt.Println("||||||||||||||||||||||||||||||||", err)
		return nil, errors.New(constants.InternalServerError)
	}
	err = cur.All(context.TODO(), &records)
	if err != nil {
		openlog.Error(err.Error())
		return nil, errors.New(constants.InternalServerError)
	}
	return records, nil
}
