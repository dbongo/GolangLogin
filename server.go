package main

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/golang/oauth2"
	mgo "gopkg.in/mgo.v2"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
)

var (
	adress  string = "localhost"
	port    string = "3000"
	debug   bool
	version = "1.0"
	session *mgo.Session

	/*
		Will move this section to a config file later.
	*/
	DBname         = "GollyLogin"
	UserCollection = "Users"
	jwt_secret     = "I<3Unicorns" // I mean... Who doesn't?
	facebookConf   = new(oauth2.Config)
	googleConf     = new(oauth2.Config)
	// linkedinConf   = new(oauth2.Config)
)

func init() {
	flag.BoolVar(&debug, "debug", false, "Enables debug mode")
	flag.Parse()

	var err error
	session, err = mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	facebookConf, err = newFacebookConf()
	if err != nil {
		panic(err)
	}

	googleConf, err = newGoogleConf()
	if err != nil {
		panic(err)
	}

	// linkedinConf, err = newLinkedInConf()
	// if err != nil {
	// 	panic(err)
	// }
	// Use maximum number of cores for optimal performance.
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	r := gin.Default()
	//Use CORS as a global middleware.
	r.Use(CORSMiddleware())

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
		user := new(LocalUser)
		if c.Bind(user) {
			mgo_session := session.Copy()
			defer mgo_session.Close()
			if debug {
				log.Printf("Email: %s Password: %s", user.Email, user.Password)
			}
			token, err := getToken(user.Email, user.Password, mgo_session)
			if err != nil {
				c.Fail(http.StatusInternalServerError, err)
				return
			}
			c.JSON(http.StatusOK, gin.H{"token": token})
		}
	})

	/*
		Login route with one of the choosen following providers.
		* Facebook
		* Google
		* Linkedin
	*/
	public.GET("/login/:provider", func(c *gin.Context) {
		provider := c.Params.ByName("provider")
		mgo_session := session.Copy()
		defer mgo_session.Close()

		switch provider {
		case "facebook":
			c.Redirect(http.StatusTemporaryRedirect, facebookConf.AuthCodeURL("unicorn", "offline", "auto"))
		case "google":
			c.Redirect(http.StatusTemporaryRedirect, googleConf.AuthCodeURL("unicorn", "offline", "auto"))
		case "linkedin":
			// c.Redirect(http.StatusTemporaryRedirect, linkedinConf.AuthCodeURL("unicorn", "offline", "auto"))
		case "kth":
			c.Redirect(http.StatusTemporaryRedirect, "https://login.kth.se/login?service=http://localhost:3000/api/callback/kth")
		default:
			/*
				The user did not choose any appropriate provider.
				Send status code 404 back to client.
			*/
			c.Fail(http.StatusNotFound, errors.New("Unknown provider"))
			return
		}
	})

	public.POST("/create", func(c *gin.Context) {
		user := new(LocalUser)
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
			token, err := createLocalUserAndGetToken(user.Email, user.Password, mgo_session)
			if err != nil {
				c.Fail(http.StatusInternalServerError, err)
				return
			}
			c.JSON(http.StatusOK, gin.H{"token": token})
		}
	})

	/*
				Callback for several social media providers.
		----------------------------------------------------------------
	*/
	public.GET("/callback/facebook", func(c *gin.Context) {
		code := c.Request.URL.Query().Get("code")
		mgo_session := session.Copy()
		defer mgo_session.Close()

		t, err := facebookConf.NewTransportWithCode(code)
		if err != nil {
			c.Fail(http.StatusInternalServerError, err)
			return
		}
		client := http.Client{Transport: t}
		res, err := client.Get("https://graph.facebook.com/me")
		if err != nil {
			c.Fail(http.StatusInternalServerError, err)
			return
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			c.Fail(http.StatusInternalServerError, err)
			return
		}
		var user map[string]interface{}
		if err := json.Unmarshal(body, &user); err != nil {
			c.Fail(http.StatusInternalServerError, err)
			return
		}
		c.JSON(200, user)
		// token, err := lookUpProviderUserAndGetToken(user, session)
		// if err != nil {
		// 	c.Fail(http.StatusInternalServerError, err)
		// 	return
		// }
		// c.JSON(http.StatusOK, gin.H{"token": token})
	})

	public.GET("/callback/google", func(c *gin.Context) {
		mgo_session := session.Copy()
		defer mgo_session.Close()
	})

	public.GET("/callback/linkedin", func(c *gin.Context) {
		mgo_session := session.Copy()
		defer mgo_session.Close()
	})

	public.GET("/callback/kth", func(c *gin.Context) {
		// ticket := c.Request.URL.Query().Get("ticket")

		// c.JSON(200, gin.H{"user": string(data)})
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
			c.Fail(http.StatusInternalServerError, err)
			return
		}
		c.JSON(200, users)
	})

	// Get a user with id.
	private.GET("/users/:id", func(c *gin.Context) {
		mgo_session := session.Clone()
		defer mgo_session.Close()
		id := c.Params.ByName("id")
		user, err := getUserWithID(id, mgo_session)
		if err != nil {
			c.Fail(http.StatusInternalServerError, err)
			return
		}
		c.JSON(200, user)
	})

	// Update an existing user with id in the list of users.
	private.PATCH("/users/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		user := new(UpdateLocalUser)
		if c.Bind(user) {
			mgo_session := session.Clone()
			defer mgo_session.Close()
			err := updateUser(id, user, mgo_session)
			if err != nil {
				c.Fail(http.StatusInternalServerError, err)
				return
			}
			c.JSON(200, gin.H{"Status": "ok"})
		}
	})
	r.Run(adress + ":" + port)
}
