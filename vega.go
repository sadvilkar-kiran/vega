package vega

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/dgraph-io/badger/v4"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/sadvilkar-kiran/vega/cache"
	"github.com/sadvilkar-kiran/vega/mailer"
	"github.com/sadvilkar-kiran/vega/render"
	"github.com/sadvilkar-kiran/vega/session"
)

const version = "1.0.0"

var myRedisCache *cache.RedisCache
var myBadgerCache *cache.BadgerCache
var redisPool *redis.Pool
var badgerConn *badger.DB

// Vega is the overall type for the Vega package. Members that are exported in this type
// are available to any application that uses it.
type Vega struct {
	AppName       string
	Debug         bool
	Version       string
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	RootPath      string
	Routes        *chi.Mux
	Render        *render.Render
	Session       *scs.SessionManager
	DB            Database
	JetViews      *jet.Set
	config        config
	EncryptionKey string
	Cache         cache.Cache
	Scheduler     *cron.Cron
	Mail          mailer.Mail
	Server        Server
}

// Server holds server configuration
type Server struct {
	ServerName string
	Port       string
	Secure     bool
	URL        string
}

type config struct {
	port        string
	renderer    string
	cookie      cookieConfig
	sessionType string
	database    databaseConfig
	redis       redisConfig
}

// New reads the .env file, creates our application config, populates the Vega type with settings
// based on .env values, and creates necessary folders and files if they don't exist
func (v *Vega) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "mail", "data", "public", "tmp", "logs", "middleware"},
	}

	err := v.Init(pathConfig)
	if err != nil {
		return err
	}

	err = v.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	// read .env
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	// create loggers
	infoLog, errorLog := v.startLoggers()

	// connect to database
	if os.Getenv("DATABASE_TYPE") != "" {
		db, err := v.OpenDB(os.Getenv("DATABASE_TYPE"), v.BuildDSN())
		if err != nil {
			errorLog.Println(err)
			os.Exit(1)
		}
		v.DB = Database{
			DataType: os.Getenv("DATABASE_TYPE"),
			Pool:     db,
		}
	}

	scheduler := cron.New()
	v.Scheduler = scheduler

	if os.Getenv("CACHE") == "redis" || os.Getenv("SESSION_TYPE") == "redis" {
		myRedisCache = v.createClientRedisCache()
		v.Cache = myRedisCache
		redisPool = myRedisCache.Conn
	}

	if os.Getenv("CACHE") == "badger" {
		myBadgerCache = v.createClientBadgerCache()
		v.Cache = myBadgerCache
		badgerConn = myBadgerCache.Conn

		_, err = v.Scheduler.AddFunc("@daily", func() {
			_ = myBadgerCache.Conn.RunValueLogGC(0.7)
		})
		if err != nil {
			return err
		}
	}

	v.InfoLog = infoLog
	v.ErrorLog = errorLog
	v.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	v.Version = version
	v.RootPath = rootPath
	v.Mail = v.createMailer()
	v.Routes = v.routes().(*chi.Mux)

	v.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifetime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSISTS"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESSION_TYPE"),
		database: databaseConfig{
			database: os.Getenv("DATABASE_TYPE"),
			dsn:      v.BuildDSN(),
		},
		redis: redisConfig{
			host:     os.Getenv("REDIS_HOST"),
			password: os.Getenv("REDIS_PASSWORD"),
			prefix:   os.Getenv("REDIS_PREFIX"),
		},
	}

	secure := true
	if strings.ToLower(os.Getenv("SECURE")) == "false" {
		secure = false
	}

	v.Server = Server{
		ServerName: os.Getenv("SERVER_NAME"),
		Port:       os.Getenv("PORT"),
		Secure:     secure,
		URL:        os.Getenv("APP_URL"),
	}

	// create session
	sess := session.Session{
		CookieLifetime: v.config.cookie.lifetime,
		CookiePersist:  v.config.cookie.persist,
		CookieName:     v.config.cookie.name,
		SessionType:    v.config.sessionType,
		CookieDomain:   v.config.cookie.domain,
	}

	switch v.config.sessionType {
	case "redis":
		sess.RedisPool = myRedisCache.Conn
	case "mysql", "postgres", "mariadb", "postgresql":
		sess.DBPool = v.DB.Pool
	}

	v.Session = sess.InitSession()
	v.EncryptionKey = os.Getenv("KEY")

	if v.Debug {
		var views = jet.NewSet(
			jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
			jet.InDevelopmentMode(),
		)
		v.JetViews = views
	} else {
		var views = jet.NewSet(
			jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
		)
		v.JetViews = views
	}

	v.createRenderer()
	go v.Mail.ListenForMail()

	return nil
}

// Init creates necessary folders for our Vega application
func (v *Vega) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.folderNames {
		// create folder if it doesn't exist
		err := v.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}

// ListenAndServe starts the web server
func (v *Vega) ListenAndServe() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     v.ErrorLog,
		Handler:      v.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	if v.DB.Pool != nil {
		defer v.DB.Pool.Close()
	}

	if redisPool != nil {
		defer redisPool.Close()
	}

	if badgerConn != nil {
		defer badgerConn.Close()
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
		v.InfoLog.Println("PORT not set, using default: 4000")
	}
	
	v.InfoLog.Printf("Listening on port %s", port)
	srv.Addr = fmt.Sprintf(":%s", port)
	err := srv.ListenAndServe()
	if err != nil {
		v.ErrorLog.Fatal(err)
	}
}

func (v *Vega) checkDotEnv(path string) error {
	err := v.CreateFileIfNotExists(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}

func (v *Vega) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

func (v *Vega) createRenderer() {
	myRenderer := render.Render{
		Renderer:   v.config.renderer,
		RootPath:   v.RootPath,
		Secure:     v.Server.Secure,
		Port:       v.Server.Port,
		ServerName: v.Server.ServerName,
		Version:    v.Version,
		JetViews:   v.JetViews,
		Session:    v.Session,
	}
	v.Render = &myRenderer
}

func (v *Vega) createMailer() mailer.Mail {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	m := mailer.Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Templates:   v.RootPath + "/mail",
		Host:        os.Getenv("SMTP_HOST"),
		Port:        port,
		Username:    os.Getenv("SMTP_USERNAME"),
		Password:    os.Getenv("SMTP_PASSWORD"),
		Encryption:  os.Getenv("SMTP_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
		Jobs:        make(chan mailer.Message, 20),
		Results:     make(chan mailer.Result, 20),
		API:         os.Getenv("MAILER_API"),
		APIKey:      os.Getenv("MAILER_KEY"),
		APIUrl:      os.Getenv("MAILER_URL"),
	}
	return m
}

func (v *Vega) createClientRedisCache() *cache.RedisCache {
	cacheClient := cache.RedisCache{
		Conn:   v.createRedisPool(),
		Prefix: v.config.redis.prefix,
	}
	return &cacheClient
}

func (v *Vega) createClientBadgerCache() *cache.BadgerCache {
	cacheClient := cache.BadgerCache{
		Conn: v.createBadgerConn(),
	}
	return &cacheClient
}

func (v *Vega) createRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   10000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",
				v.config.redis.host,
				redis.DialPassword(v.config.redis.password))
		},

		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
}

func (v *Vega) createBadgerConn() *badger.DB {
	db, err := badger.Open(badger.DefaultOptions(v.RootPath + "/tmp/badger"))
	if err != nil {
		return nil
	}
	return db
}

// BuildDSN builds the datasource name for our database, and returns it as a string
func (v *Vega) BuildDSN() string {
	var dsn string

	switch os.Getenv("DATABASE_TYPE") {
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"))

		// we check to see if a database password has been supplied, since including "password=" with nothing
		// after it sometimes causes postgres to fail to allow a connection.
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("%s password=%s", dsn, os.Getenv("DATABASE_PASS"))
		}

	default:

	}

	return dsn
}

