package helpers

import (
	"Task/constants"

	"go.mongodb.org/mongo-driver/mongo/options"
)

//
func ValidatePageAndSizeDetails(details map[string]interface{}) *options.FindOptions {
	options := options.Find()
	if details == nil {
		return options
	}
	page := details["page"].(int64)
	size := details["size"].(int64)
	if page <= 0 || size <= 0 {
		return options
	}
	options = options.SetLimit(size)
	options = options.SetSkip((page - 1) * size)
	return options

}

//
var ErrorStatusMap map[string]int

func init() {
	ErrorStatusMap = make(map[string]int)
	ErrorStatusMap[constants.InternalServerError] = 500
	ErrorStatusMap[constants.InvalidRecordID] = 400
	ErrorStatusMap[constants.NoRecords] = 404
}
