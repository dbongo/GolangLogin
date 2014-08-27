package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	// Cost of bcrypt hash.
	currentCost = 10
	// Sets how many hours until token expire.
	tokenExpire = 24
)

/*
	This function checks if user is in database and returns an token if so.
*/
func getToken(email, password string, session *mgo.Session) (string, error) {
	user := new(LocalUser)
	c := session.DB(DBname).C(UserCollection)
	session.SetMode(mgo.Monotonic, true)

	// Find one user with the email.
	err := c.Find(bson.M{"email": email}).One(user)
	if err != nil {
		// The user could not be found in DB.
		return "", err
	}
	if err = validPass(password, user.Password); err != nil {
		return "", err
	}
	claims := make(map[string]interface{})
	claims["ID"] = user.UserID.Hex()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token, err := generateToken([]byte(jwt_secret), &claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

/*
	This function will create an new user and put it in DB.
	On success it will return an token.
*/
func createUserAndGetToken(email, password string, session *mgo.Session) (string, error) {
	hashedPass, err := hashPass(password)
	if err != nil {
		return "", err
	}
	user := &LocalUser{
		UserID:   bson.NewObjectId(),
		Email:    email,
		Password: hashedPass,
	}
	c := session.DB(DBname).C(UserCollection)
	session.SetMode(mgo.Monotonic, true)
	if err := c.Insert(user); err != nil {
		return "", err
	}
	claims := make(map[string]interface{})
	claims["ID"] = user.UserID.Hex()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token, err := generateToken([]byte(jwt_secret), &claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

func getUserWithID(id string, session *mgo.Session) (*ShowLocalUser, error) {
	user := new(ShowLocalUser)
	c := session.DB(DBname).C(UserCollection)
	session.SetMode(mgo.Monotonic, true)
	err := c.FindId(bson.ObjectIdHex(id)).One(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func getAllUsers(session *mgo.Session) ([]ShowLocalUser, error) {
	var users []ShowLocalUser
	c := session.DB(DBname).C(UserCollection)
	session.SetMode(mgo.Monotonic, true)
	err := c.Find(bson.M{}).All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func updateUser(id string, user *UpdateLocalUser, session *mgo.Session) error {
	c := session.DB(DBname).C(UserCollection)
	session.SetMode(mgo.Monotonic, true)
	var err error
	if user.Password == "" {
		// LocalUser did not post any password.
		err = c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"email": user.Email}})
		if err != nil {
			return err
		}
	}
	if user.Email == "" {
		// LocalUser did not post any email.
		hashedPass, err := hashPass(user.Password)
		if err != nil {
			return err
		}
		err = c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"password": hashedPass}})
		if err != nil {
			return err
		}
	}
	hashedPass, err := hashPass(user.Password)
	if err != nil {
		return err
	}
	err = c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"email": user.Email, "password": hashedPass}})
	if err != nil {
		return err
	}
	return nil
}

/*
	Helper function.
	Will return an hashed string of the password.
*/
func hashPass(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), currentCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

/*
	Helper function.
	Compare password and the hashed password.
	Will return nil on success, otherwise an error.
*/
func validPass(pass, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		return err
	}
	return nil
}
