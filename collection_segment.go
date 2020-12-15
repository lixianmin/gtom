package mop

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollectionSegment struct {
	min         interface{}
	max         interface{}
	splitKey    string
	splits      []map[string]interface{}
	subSegments []*CollectionSegment
}

func (cs *CollectionSegment) shrinkTo(next interface{}) {
	cs.max = next
}

func (cs *CollectionSegment) toSelector() bson.M {
	sel, doc := bson.M{}, bson.M{}
	if cs.min != nil {
		doc["$gte"] = cs.min
	}
	if cs.max != nil {
		doc["$lt"] = cs.max
	}
	if len(doc) > 0 {
		sel[cs.splitKey] = doc
	}
	return sel
}

func (cs *CollectionSegment) divide() {
	if len(cs.splits) == 0 {
		return
	}
	ns := &CollectionSegment{
		splitKey: cs.splitKey,
		min:      cs.min,
		max:      cs.max,
	}
	cs.subSegments = nil
	for _, split := range cs.splits {
		ns.shrinkTo(split[cs.splitKey])
		cs.subSegments = append(cs.subSegments, ns)
		ns = &CollectionSegment{
			splitKey: cs.splitKey,
			min:      ns.max,
			max:      cs.max,
		}
	}
	ns = &CollectionSegment{
		splitKey: cs.splitKey,
		min:      cs.splits[len(cs.splits)-1][cs.splitKey],
	}
	cs.subSegments = append(cs.subSegments, ns)
}

func (cs *CollectionSegment) init(c *mongo.Collection) (err error) {
	opts := &options.FindOneOptions{}
	opts.SetSort(bson.M{cs.splitKey: 1})
	doc := make(map[string]interface{})
	if err = c.FindOne(context.Background(), nil, opts).Decode(&doc); err != nil {
		return
	}
	cs.min = doc[cs.splitKey]
	opts = &options.FindOneOptions{}
	opts.SetSort(bson.M{cs.splitKey: -1})
	doc = make(map[string]interface{})
	if err = c.FindOne(context.Background(), nil, opts).Decode(&doc); err != nil {
		return
	}
	cs.max = doc[cs.splitKey]
	return
}