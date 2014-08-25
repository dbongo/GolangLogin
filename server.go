package main

import (
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
	// bson "gopkg.in/mgo.v2/bson"
	"log"
	"runtime"
	"time"
)

var (
	adress  string
	port    string
	debug   bool
	version = "1.0"
	session *mgo.Session
	DBname  = "login"
	// Will move this to a config file later.
	jwt_secret = "mySuperDuperMegaSecret"
)

func init() {
	flag.StringVar(&adress, "adress", "localhost", "Adress on which the server should be running on")
	flag.StringVar(&port, "port", "3000", "Port on which the server should be running on")
	flag.BoolVar(&debug, "debug", false, "Enables debug mode")
	flag.Parse()

	var err error
	session, err = mgo.Dial("localhost")
	if err != nil {
		log.Fatalln(err.Error())
	}
	// Use maximum number of cores for optimal performance.
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	r := gin.Default()

	public := r.Group("/api")
	// Private route group with jwt token middleware attached.
	private := r.Group("/api/authorized", tokenMiddleWare(jwt_secret))

	/*
		Public route section
		----------------------------------------------------------------
	*/

	//Local login route.
	public.POST("/login", func(c *gin.Context) {
		user := new(User)
		if c.Bind(user) { /*
				Correct data were sent.
				Get user data from DB here.
			*/
			// Use this session to get data from DB. This will save memory in the long run.
			mgo_session := session.Clone()
			defer mgo_session.Close()
			if debug {
				log.Printf("Email: %s Password: %s", user.Email, user.Password)
			}
			// JWT section
			claims := make(map[string]interface{})
			claims["ID"] = "100"
			claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
			token, err := generateToken([]byte(jwt_secret), claims)
			if err != nil {
				c.Fail(500, err)
			}
			c.JSON(200, gin.H{"token": token})
		}
	})

	// Login route with the choosen provider.
	public.GET("/login/:provider", func(c *gin.Context) {
		provider := c.Params.ByName("provider")
		mgo_session := session.Clone()
		defer mgo_session.Close()

		switch provider {
		case "facebook":
			// Handle facebook login here
		case "google":
			// Handle google login here
		default:
			/*
				The user did not choose any appropriate provider.
				Send status code 400 back to client.
			*/
			c.Fail(400, errors.New("Unknown provider"))
		}
	})

	// Get the current API version.
	public.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{"version": version, "author": "Christopher Lillthors"})
	})

	/*
		Private route section
		----------------------------------------------------------------
	*/

	// Get a list of all the users.
	private.GET("/users", func(c *gin.Context) {
		mgo_session := session.Clone()
		defer mgo_session.Close()
	})

	private.GET("/users/:id", func(c *gin.Context) {
		// id := c.Params.ByName("id")
		mgo_session := session.Clone()
		defer mgo_session.Close()
	})

	// Update an existing user in the list of users.
	private.PATCH("/users/:id", func(c *gin.Context) {
		// id := c.Params.ByName("id")
		mgo_session := session.Clone()
		defer mgo_session.Close()
	})

	// Add a user to the list of users.
	private.POST("/users", func(c *gin.Context) {
		mgo_session := session.Clone()
		defer mgo_session.Close()
	})
	r.Run(adress + ":" + port)
}
