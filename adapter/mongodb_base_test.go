package adapter

import (
	"context"
	mim "github.com/ONSdigital/dp-mongodb-in-memory"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBPersistenceSuite struct {
	suite.Suite
	server                *mim.Server
	client                *mongo.Client
	database              *mongo.Database
	credentialsCollection *mongo.Collection
	userCollection        *mongo.Collection
	credentialsRepository MongoDbCredentialsRepository
	userRepository        MongoDbUserRepository
}

func (suite *MongoDBPersistenceSuite) SetupSuite() {
	testCtx := context.Background()

	server, err := mim.Start(testCtx, "5.0.2")
	manageTestError(err)
	suite.server = server

	client, err := mongo.Connect(testCtx, options.Client().ApplyURI(server.URI()))
	manageTestError(err)
	//Use client as needed
	err = client.Ping(testCtx, nil)
	manageTestError(err)
	suite.client = client
	suite.database = client.Database("auth")
}

func (suite *MongoDBPersistenceSuite) TearDownSuite() {
	ctx := context.TODO()
	defer suite.server.Stop(ctx)
	err := suite.client.Disconnect(ctx)
	manageTestError(err)
}

func (suite *MongoDBPersistenceSuite) setupCredentialsCollection() {
	db := suite.database
	log.Debug().Msgf("Creating collection '%s'", CredentialsCollection)
	err := db.CreateCollection(context.TODO(), CredentialsCollection)
	manageTestError(err)

	collection := db.Collection(CredentialsCollection)

	idIdx := mongo.IndexModel{
		Keys: bson.M{
			"id": 1, // index in ascending order
		}, Options: &options.IndexOptions{Unique: boolPointer(true)},
	}
	s, err := collection.Indexes().CreateOne(context.TODO(), idIdx)
	manageTestError(err)
	log.Debug().Msg(s)

	suite.credentialsCollection = collection
	suite.credentialsRepository = NewMongoDbCredentialsRepository(suite.credentialsCollection)
}

func (suite *MongoDBPersistenceSuite) setupUserCollection() {
	db := suite.database
	log.Debug().Msgf("Creating collection '%s'", UserCollection)
	err := db.CreateCollection(context.TODO(), UserCollection)
	manageTestError(err)

	collection := db.Collection(UserCollection)

	idIdx := []mongo.IndexModel{
		{
			Keys: bson.M{
				"id": 1, // index in ascending order
			}, Options: &options.IndexOptions{Unique: boolPointer(true)},
		},
		{
			Keys: bson.M{
				"email": 1, // index in ascending order
			}, Options: &options.IndexOptions{Unique: boolPointer(true)},
		},
		{
			Keys: bson.M{
				"phone_number": 1, // index in ascending order
			}, Options: &options.IndexOptions{Unique: boolPointer(true)},
		},
	}
	s, err := collection.Indexes().CreateMany(context.TODO(), idIdx)
	manageTestError(err)
	for _, str := range s {
		log.Debug().Msg(str)
	}

	suite.userCollection = collection
	suite.userRepository = NewMongoDbUserRepository(suite.userCollection)
}

func insertOne(coll *mongo.Collection, ctx context.Context, obj interface{}) {
	log.Debug().Msgf("Inserting %v", obj)
	_, err := coll.InsertOne(ctx, obj)
	manageTestError(err)
}

func findOne(coll *mongo.Collection, ctx context.Context, property string, value, entity interface{}) {
	log.Debug().Msgf("Finding object from collection '%s'", coll.Name())
	err := coll.FindOne(ctx, bson.M{
		property: value,
	}, options.FindOne().SetSort(bson.M{})).Decode(entity)
	manageTestError(err)
}

func deleteAll(coll *mongo.Collection, ctx context.Context) {
	log.Debug().Msgf("Deleting all documents in collection '%s'", coll.Name())
	_, err := coll.DeleteMany(ctx, bson.D{})
	manageTestError(err)
}

func manageTestError(err error) {
	if err != nil {
		log.Error().Msg(err.Error())
		panic(err)
	}
}
func boolPointer(b bool) *bool {
	return &b
}
