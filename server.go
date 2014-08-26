package main

import (
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
	"log"
	"runtime"
)

var (
	adress  string
	port    string
	debug   bool
	version = "1.0"
	session *mgo.Session

	/*
		Will move this section to a config file later.
	*/
	DBname         = "GollyLogin"
	UserCollection = "Users"
	jwt_secret     = "mySuperDuperMegaSecret"
)

func init() {
	// Will move this section to a config file later.
	flag.StringVar(&adress, "adress", "localhost", "Adress on which the server should be running on")
	flag.StringVar(&port, "port", "3000", "Port on which the server should be running on")
	flag.BoolVar(&debug, "debug", false, "Enables debug mode")
	flag.Parse()

	var err error
	session, err = mgo.Dial("127.0.0.1")
	if err != nil {
		log.Fatalln(err.Error())
	}
	// Use maximum number of cores for optimal performance.
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	r := gin.Default()
	// Public route group.
	public := r.Group("/api")
	// Private route group with jwt token middleware attached.
	private := r.Group("/api/authorized", tokenMiddleWare(jwt_secret))

	/*
		Public routes section
		----------------------------------------------------------------
	*/

	// Local login route.
	public.POST("/login", func(c *gin.Context) {
		user := new(User)
		if c.Bind(user) {
			mgo_session := session.Copy()
			defer mgo_session.Close()
			if debug {
				log.Printf("Email: %s Password: %s", user.Email, user.Password)
			}
			token, err := getToken(user.Email, user.Password, mgo_session)
			if err != nil {
				c.Fail(400, err)
				return
			}
			c.JSON(200, gin.H{"token": token})
		}
	})

	// Login route with the choosen provider.
	public.GET("/login/:provider", func(c *gin.Context) {
		provider := c.Params.ByName("provider")
		mgo_session := session.Copy()
		defer mgo_session.Close()

		switch provider {
		case "facebook":
			// Handle facebook login here.
		case "google":
			// Handle google login here.
		default:
			/*
				The user did not choose any appropriate provider.
				Send status code 400 back to client.
			*/
			c.Fail(400, errors.New("Unknown provider"))
			return
		}
	})

	public.POST("/create", func(c *gin.Context) {
		user := new(User)
		if c.Bind(user) {
			/*
				Correct data were sent.
				Insert user data into DB here.
			*/

			// Use this session to get data from DB. This will save memory in the long run.
			mgo_session := session.Copy()
			defer mgo_session.Close()
			if debug {
				log.Printf("Email: %s Password: %s", user.Email, user.Password)
			}
			token, err := createUserAndGetToken(user.Email, user.Password, mgo_session)
			if err != nil {
				c.Fail(500, err)
				return
			}
			c.JSON(200, gin.H{"token": token})
		}
	})

	// Get the current API version.
	public.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{"version": version, "author": "Christopher Lillthors"})
	})

	/*
		Private routes section
		----------------------------------------------------------------
	*/

	// Get a list of all the users.
	private.GET("/users", func(c *gin.Context) {
		mgo_session := session.Clone()
		defer mgo_session.Close()
		users, err := getAllUsers(session)
		if err != nil {
			c.Fail(500, err)
			return
		}
		c.JSON(200, users)
	})

	private.GET("/users/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		mgo_session := session.Clone()
		defer mgo_session.Close()
		user, err := getUserWithID(id, mgo_session)
		if err != nil {
			c.Fail(500, err)
			return
		}
		c.JSON(200, user)
	})

	// Update an existing user in the list of users.
	private.PATCH("/users/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		user := new(User)
		if c.Bind(user) {
			mgo_session := session.Clone()
			defer mgo_session.Close()
			err := updateUser(id, user, mgo_session)
			if err != nil {
				c.Fail(500, err)
				return
			}
			c.JSON(200, gin.H{"Status": "ok"})
		}
	})
	r.Run(adress + ":" + port)
}
