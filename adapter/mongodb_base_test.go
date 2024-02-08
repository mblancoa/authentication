package adapter

import (
	"context"
	mim "github.com/ONSdigital/dp-mongodb-in-memory"
	"github.com/mblancoa/authentication/tools"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDBPersistenceSuite struct {
	suite.Suite
	server                *mim.Server
	client                *mongo.Client
	database              *mongo.Database
	credentialsCollection *mongo.Collection
	userCollection        *mongo.Collection
	credentialsRepository MongoDbCredentialsRepository
	userRepository        MongoDbUserRepository
}

func (suite *mongoDBPersistenceSuite) SetupSuite() {
	testCtx := context.Background()

	server, err := mim.Start(testCtx, "5.0.2")
	tools.ManageTestError(err)
	suite.server = server

	client, err := mongo.Connect(testCtx, options.Client().ApplyURI(server.URI()))
	tools.ManageTestError(err)
	//Use client as needed
	err = client.Ping(testCtx, nil)
	tools.ManageTestError(err)
	suite.client = client
	suite.database = client.Database("auth")
	suite.Assert()
}

func (suite *mongoDBPersistenceSuite) TearDownSuite() {
	ctx := context.TODO()
	defer suite.server.Stop(ctx)
	err := suite.client.Disconnect(ctx)
	tools.ManageTestError(err)
}

func (suite *mongoDBPersistenceSuite) setupCredentialsCollection() {
	db := suite.database
	log.Debug().Msgf("Creating collection '%s'", CredentialsCollection)
	err := db.CreateCollection(context.TODO(), CredentialsCollection)
	tools.ManageTestError(err)

	collection := db.Collection(CredentialsCollection)

	idIdx := mongo.IndexModel{
		Keys: bson.M{
			"_id": 1, // index in ascending order
		},
	}
	s, err := collection.Indexes().CreateOne(context.TODO(), idIdx)
	tools.ManageTestError(err)
	log.Debug().Msg(s)

	suite.credentialsCollection = collection
	suite.credentialsRepository = NewMongoDbCredentialsRepository(suite.credentialsCollection)
}

func (suite *mongoDBPersistenceSuite) setupUserCollection() {
	db := suite.database
	log.Debug().Msgf("Creating collection '%s'", UserCollection)
	err := db.CreateCollection(context.TODO(), UserCollection)
	tools.ManageTestError(err)

	collection := db.Collection(UserCollection)

	idIdx := []mongo.IndexModel{
		{
			Keys: bson.M{
				"_id": 1, // index in ascending order
			},
		},
		{
			Keys: bson.M{
				"email": 1, // index in ascending order
			}, Options: &options.IndexOptions{Unique: tools.BoolPointer(true)},
		},
		{
			Keys: bson.M{
				"phone_number": 1, // index in ascending order
			}, Options: &options.IndexOptions{Unique: tools.BoolPointer(true)},
		},
	}
	s, err := collection.Indexes().CreateMany(context.TODO(), idIdx)
	tools.ManageTestError(err)
	for _, str := range s {
		log.Debug().Msg(str)
	}

	suite.userCollection = collection
	suite.userRepository = NewMongoDbUserRepository(suite.userCollection)
}

func insertOne(coll *mongo.Collection, ctx context.Context, obj interface{}) {
	log.Debug().Msgf("Inserting %v", obj)
	_, err := coll.InsertOne(ctx, obj)
	tools.ManageTestError(err)
}

func findOne(coll *mongo.Collection, ctx context.Context, property string, value, entity interface{}) {
	log.Debug().Msgf("Finding object from collection '%s'", coll.Name())
	err := coll.FindOne(ctx, bson.M{
		property: value,
	}, options.FindOne().SetSort(bson.M{})).Decode(entity)
	tools.ManageTestError(err)
}

func deleteAll(coll *mongo.Collection, ctx context.Context) {
	log.Debug().Msgf("Deleting all documents in collection '%s'", coll.Name())
	_, err := coll.DeleteMany(ctx, bson.D{})
	tools.ManageTestError(err)
}

func count(coll *mongo.Collection, ctx context.Context) int64 {
	c, err := coll.CountDocuments(ctx, bson.D{})
	tools.ManageTestError(err)
	return c
}
