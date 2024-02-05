package adapter

import (
	"context"
	mim "github.com/ONSdigital/dp-mongodb-in-memory"
	"github.com/brianvoe/gofakeit"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

type MongoDBRepositorySuite struct {
	suite.Suite
	server                *mim.Server
	client                *mongo.Client
	database              *mongo.Database
	credentialsCollection *mongo.Collection
	credentialsRepository MongoDbCredentialsRepository
}

func (suite *MongoDBRepositorySuite) SetupSuite() {
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

	setupCredentialsCollection(suite.database)
	suite.credentialsCollection = suite.database.Collection(CredentialsCollection)
	suite.credentialsRepository = NewMongoDbCredentialsRepository(suite.credentialsCollection)
}

func setupCredentialsCollection(db *mongo.Database) {
	log.Debug().Msgf("Creating collection \"%s\"", CredentialsCollection)
	err := db.CreateCollection(context.TODO(), CredentialsCollection)
	manageTestError(err)

	collection := db.Collection(CredentialsCollection)

	idIdx := mongo.IndexModel{
		Keys: bson.M{
			"id": 1, // index in ascending order
		}, Options: nil,
	}
	log.Debug().Msgf("Creating index whit field \"id\"")
	s, err := collection.Indexes().CreateOne(context.TODO(), idIdx)
	manageTestError(err)
	log.Debug().Msg(s)
}
func (suite *MongoDBRepositorySuite) TearDownSuite() {
	ctx := context.TODO()
	defer suite.server.Stop(ctx)
	err := suite.client.Disconnect(ctx)
	manageTestError(err)
}

func TestCredentialsRepositorySuite(t *testing.T) {
	suite.Run(t, new(MongoDBRepositorySuite))
}

func (suite *MongoDBRepositorySuite) TestFindCredentialsById() {
	var credentials CredentialsDB
	gofakeit.Struct(&credentials)
	gofakeit.Struct(&credentials.LastPasswords)
	insertOne(suite.credentialsCollection, context.TODO(), credentials)

	result, err := suite.credentialsRepository.FindById(context.TODO(), credentials.Id)

	suite.Assert().NoError(err)
	suite.Assert().NotEmpty(result)
	suite.Assert().Equal(&credentials, result)
}

func insertOne(coll *mongo.Collection, ctx context.Context, obj interface{}) {
	log.Debug().Msgf("Inserting %v", obj)
	_, err := coll.InsertOne(ctx, obj)
	log.Debug().Msgf("%v inserted", obj)
	manageTestError(err)
}

func manageTestError(err error) {
	if err != nil {
		log.Error().Msg(err.Error())
		panic(err)
	}
}
