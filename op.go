package mop

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type Op struct {
	Id                interface{}            `json:"_id"`
	Operation         string                 `json:"operation"`
	Namespace         string                 `json:"namespace"`
	Data              map[string]interface{} `json:"data,omitempty"`
	Timestamp         primitive.Timestamp    `json:"timestamp"`
	Source            QuerySource            `json:"source"`
	Doc               interface{}            `json:"doc,omitempty"`
	UpdateDescription map[string]interface{} `json:"updateDescription,omitempty`
}

func (my *Op) IsDrop() bool {
	if _, drop := my.IsDropDatabase(); drop {
		return true
	}
	if _, drop := my.IsDropCollection(); drop {
		return true
	}
	return false
}

func (my *Op) IsDropCollection() (string, bool) {
	if my.IsCommand() {
		if my.Data != nil {
			if val, ok := my.Data["drop"]; ok {
				return val.(string), true
			}
		}
	}
	return "", false
}

func (my *Op) IsDropDatabase() (string, bool) {
	if my.IsCommand() {
		if my.Data != nil {
			if _, ok := my.Data["dropDatabase"]; ok {
				return my.GetDatabase(), true
			}
		}
	}
	return "", false
}

func (my *Op) IsCommand() bool {
	return my.Operation == "c"
}

func (my *Op) IsInsert() bool {
	return my.Operation == "i"
}

func (my *Op) IsUpdate() bool {
	return my.Operation == "u"
}

func (my *Op) IsDelete() bool {
	return my.Operation == "d"
}

func (my *Op) IsSourceOplog() bool {
	return my.Source == OplogQuerySource
}

func (my *Op) IsSourceDirect() bool {
	return my.Source == DirectQuerySource
}

func (my *Op) ParseNamespace() []string {
	return strings.SplitN(my.Namespace, ".", 2)
}

func (my *Op) GetDatabase() string {
	return my.ParseNamespace()[0]
}

func (my *Op) GetCollection() string {
	if _, drop := my.IsDropDatabase(); drop {
		return ""
	} else if col, drop := my.IsDropCollection(); drop {
		return col
	} else {
		return my.ParseNamespace()[1]
	}
}

