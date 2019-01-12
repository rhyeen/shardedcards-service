package mongo

import mgo "gopkg.in/mgo.v2"

// EnsureIndicesOnCollection ensures the given collection has the given indices
func EnsureIndicesOnCollection(mongoDB *mgo.Database, collection string, indices []mgo.Index) error {
	database := mongoDB.Session.Copy().DB(mongoDB.Name)
	defer database.Session.Close()
	database.C(collection).Create(&mgo.CollectionInfo{})
	for _, index := range indices {
		err := database.C(collection).EnsureIndex(index)
		if err != nil {
			return err
		}
	}
	return nil
}
