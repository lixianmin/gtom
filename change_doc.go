package mop

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ChangeDoc struct {
	DocKey            map[string]interface{} "documentKey"
	Id                interface{}            "_id"
	Operation         string                 "operationType"
	FullDoc           map[string]interface{} "fullDocument"
	Namespace         ChangeDocNs            "ns"
	Timestamp         primitive.Timestamp    "clusterTime"
	UpdateDescription map[string]interface{} "updateDescription"
}

func (cd *ChangeDoc) docId() interface{} {
	return cd.DocKey["_id"]
}

func (cd *ChangeDoc) mapTimestamp() primitive.Timestamp {
	if cd.Timestamp.T > 0 {
		// only supported in version 4.0
		return cd.Timestamp
	} else {
		// for versions prior to 4.0 simulate a timestamp
		now := time.Now().UTC()
		return primitive.Timestamp{
			T: uint32(now.Unix()),
			I: uint32(now.Nanosecond()),
		}
	}
}

func (cd *ChangeDoc) mapOperation() string {
	if cd.Operation == "insert" {
		return "i"
	} else if cd.Operation == "update" || cd.Operation == "replace" {
		return "u"
	} else if cd.Operation == "delete" {
		return "d"
	} else if cd.Operation == "invalidate" || cd.Operation == "drop" || cd.Operation == "dropDatabase" {
		return "c"
	} else {
		return ""
	}
}

func (cd *ChangeDoc) hasUpdate() bool {
	return cd.UpdateDescription != nil
}

func (cd *ChangeDoc) hasDoc() bool {
	return (cd.mapOperation() == "i" || cd.mapOperation() == "u") && cd.FullDoc != nil
}

func (cd *ChangeDoc) isInvalidate() bool {
	return cd.Operation == "invalidate"
}

func (cd *ChangeDoc) isDrop() bool {
	return cd.Operation == "drop"
}

func (cd *ChangeDoc) isDropDatabase() bool {
	return cd.Operation == "dropDatabase"
}

func (cd *ChangeDoc) mapNs() string {
	if cd.Namespace.Collection != "" {
		return cd.Namespace.Database + "." + cd.Namespace.Collection
	} else {
		return cd.Namespace.Database + ".cmd"
	}
}