package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/rhyeen/shardedcards-service/internal/api"
	"github.com/rhyeen/shardedcards-service/internal/api/handlers"
	"github.com/rhyeen/shardedcards-service/internal/services/gameservice"
	"github.com/rhyeen/shardedcards-service/internal/stores/mongo"
	"github.com/rhyeen/shardedcards-service/internal/util/env"
	"github.com/rs/cors"
	mgo "gopkg.in/mgo.v2"
)

const (
	localURL               = "http://127.0.0.1:8081/"
	defaultMongoTimeout    = "20"
	defaultMongoAddr       = "..."
	defaultMongoDatabase   = "shardedcards-service"
	defaultMongoUsername   = "shardedcards-dev"
	defaultMongoPassword   = "DEFAULT_PASSWORD"
	defaultAdminAuthSecret = "DEFAULT_SECRET"
	defaultPort            = "8082"
	defaultStaticPath      = "../../static"
	defaultDatacenter      = "LOCAL"
)

// Indexer sets up indices for the appropriate data store
type Indexer interface {
	EnsureIndices() error
}

func getHTTPServerAddr() string {
	port := env.Get("PORT", defaultPort)
	return ":" + port
}

func getHTTPServerReadTimeout() time.Duration {
	return 10 * time.Second
}

func getHTTPServerWriteTimeout() time.Duration {
	return 10 * time.Second
}

func getHTTPServerMaxHeaderBytes() int {
	return 1 << 20
}

func getAPIPath() string {
	return "api"
}

func getStaticPath() string {
	return env.Get("STATIC_PATH", defaultStaticPath)
}

func getDatacenter() string {
	return env.Get("DATACENTER", defaultDatacenter)
}

func main() {
	mongoDB, err := setupMongoDB()
	if err != nil {
		log.Fatal(err)
	}
	defer mongoDB.Session.Close()
	apiPath := getAPIPath()
	staticPath := getStaticPath()
	datacenter := getDatacenter()
	handler, err := setupHandler(apiPath, staticPath, datacenter, mongoDB)
	if err != nil {
		log.Fatal(err)
	}
	handler, err = setupCors(datacenter, handler)
	if err != nil {
		log.Fatal(err)
	}
	s := &http.Server{
		Addr:           getHTTPServerAddr(),
		Handler:        handler,
		ReadTimeout:    getHTTPServerReadTimeout(),
		WriteTimeout:   getHTTPServerWriteTimeout(),
		MaxHeaderBytes: getHTTPServerMaxHeaderBytes(),
	}
	log.Fatal(s.ListenAndServe())
}

func setupMongoDB() (*mgo.Database, error) {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{getMongoAddress()},
		Timeout:  getMongoTimeout(),
		Database: getMongoDatabase(),
		Username: getMongoUsername(),
		Password: getMongoPassword(),
	}
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		return nil, err
	}
	// See: https://godoc.org/labix.org/v2/mgo#Session.SetSafe
	session.SetSafe(&mgo.Safe{WMode: "majority"})
	// See: https://godoc.org/labix.org/v2/mgo#Session.SetMode
	session.SetMode(mgo.Strong, true)
	err = session.Ping()
	if err != nil {
		return nil, err
	}
	return session.DB(getMongoDatabase()), nil
}

func getMongoAddress() string {
	return env.Get("MONGO_ADDRESS", defaultMongoAddr)
}

func getMongoTimeout() time.Duration {
	mongoTimeout := env.Get("MONGO_TIMEOUT", defaultMongoTimeout)
	timeout, err := strconv.Atoi(mongoTimeout)
	if err != nil {
		timeout, _ = strconv.Atoi(defaultMongoTimeout)
	}
	return time.Duration(timeout) * time.Second
}

func getMongoDatabase() string {
	return env.Get("MONGO_DATABASE", defaultMongoDatabase)
}

func getMongoUsername() string {
	return env.Get("MONGO_USERNAME", defaultMongoUsername)
}

func getMongoPassword() string {
	return env.Get("MONGO_PASSWORD", defaultMongoPassword)
}

func setupHandler(apiPath, staticPath, datacenter string, mongoDB *mgo.Database) (http.Handler, error) {
	var handler http.Handler
	gameStore := mongo.NewGameStore(mongoDB)
	deckStore := mongo.NewDeckStore(mongoDB)
	userStore := mongo.NewUserStore(mongoDB)
	err := ensureStoreIndices(gameStore, deckStore, userStore)
	if err != nil {
		return handler, err
	}
	gameService := gameservice.GameService{
		GameStore: gameStore,
		DeckStore: deckStore,
		UserStore: userStore,
	}
	routerHandlers := api.RouterHandlers{
		GameHandler: handlers.GameHandler{
			GameService: gameService,
		},
	}
	router := api.NewRouter(apiPath, staticPath, routerHandlers)
	authN, authZ, err := getAuths(apiPath, datacenter)
	if err != nil {
		return handler, err
	}
	return &api.Handler{
		AuthN:      authN,
		AuthZ:      authZ,
		Router:     router,
		Datacenter: datacenter,
		APIPath:    apiPath,
	}, nil
}

func ensureStoreIndices(stores ...Indexer) error {
	skipIndices := env.Get("SKIP_INDICES", "")
	if skipIndices == "" {
		return nil
	}
	for _, store := range stores {
		err := store.EnsureIndices()
		if err != nil {
			return err
		}
	}
	return nil
}

func getAuths(apiPath, datacenter string) (api.AuthN, api.AuthZ, error) {
	adminAuthSecret, err := getAdminAuthSecret(datacenter)
	if err != nil {
		return api.AuthN{}, api.AuthZ{}, err
	}
	authN := api.AuthN{
		Datacenter:      datacenter,
		AdminAuthSecret: adminAuthSecret,
	}
	authZ := api.AuthZ{
		APIPath: apiPath,
	}
	return authN, authZ, nil
}

func getAdminAuthSecret(datacenter string) (string, error) {
	if datacenter != api.LocalEnv {
		return env.Require("ADMIN_AUTH_SECRET")
	}
	return env.Get("ADMIN_AUTH_SECRET", defaultAdminAuthSecret), nil
}

func setupCors(datacenter string, handler http.Handler) (http.Handler, error) {
	if datacenter != api.LocalEnv {
		return handler, nil
	}
	c := cors.New(cors.Options{
		AllowedOrigins: []string{localURL},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"X-AUTH-TOKEN", "Content-Type"},
	})
	return c.Handler(handler), nil
}
