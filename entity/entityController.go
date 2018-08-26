package entity

import (
	"gopkg.in/mgo.v2"
	"financialExchange/util"
	"errors"
	"financialExchange/mongo"
	"gopkg.in/mgo.v2/bson"
	"time"
	"financialExchange/security"
)

type EntityController struct{
	Session		*mgo.Session
	Logger 		util.Logger
}

func NewEntityController (session *mgo.Session, logger *util.Logger)(*EntityController){
	return &EntityController{
		Session: session,
		Logger: *logger,
	}
}

func(self *EntityController) CreateEntity(params CreateEntityParams)(error){
	//Field Validation
	if (!util.IsValidEmail(params.Email)){
		return errors.New("InvalidEmailAddress")
	}

	if len(params.Symbol) < 1 || len(params.Symbol) > 4{
		return errors.New("InvalidSymbolLength")
	}

	if(len(params.Password) < 6){
		return errors.New("PasswordTooShort")
	}

	if (params.Password != params.PasswordVerify){
		return errors.New("PasswordsDontMatch")
	}

	//Verify Entity has not signed up with this email before
	entityCollection := mongo.GetEntityCollection(mongo.GetDataBase(self.Session))

	var findEntity Entity

	err := entityCollection.Find(bson.M{ "email": params.Email}).One(&findEntity)

	if (err != mgo.ErrNotFound){
		return errors.New("EntityExists")
	}
	//Verify symbol has not been used before
	err = entityCollection.Find(bson.M{"symbol" : params.Symbol}).One(&findEntity)

	if (err != mgo.ErrNotFound){
		return errors.New("SymbolExists")
	}

	//Hash Password
	passHash, err := util.HashPassword(params.Password)
	if (err != nil){
		return err
	}

	entityId := bson.NewObjectId()
	securityId := bson.NewObjectId()

	//Create entity and insert to database
	newEntity := Entity{
		Id: entityId,
		Name: params.Name,
		Email: params.Email,
		PassHash: passHash,
		Security: securityId,
		Created: time.Now(),
		Deleted: time.Unix(0,0),
	}

	insertErr := entityCollection.Insert(newEntity)

	if insertErr != nil{
		return err
	}

	if params.CreateSecurity{
		go func() {
			//Create security and insert to database
			securityCollection := mongo.GetSecurityCollection(mongo.GetDataBase(self.Session))

			newSecurity := security.Security{
				Id: securityId,
				Entity: entityId,
				Symbol: params.Symbol,
			}

			insertErr = securityCollection.Insert(newSecurity)
			if insertErr != nil {
				return err
			}
		}()
	}

	/***
	todo: find out how (and where) to create the pricebook
	mongo creates collections upon data insertion, so should I:
		- create it when creating the entity
		- create when first trade is made
	--going with second one for now--
	***/

	return nil
}

func (self *EntityController) SignIn(SignInParams SignInEntityParams)(string, error){
	if (!util.IsValidEmail(SignInParams.Email)){
		return "", errors.New("MissingEmailField")
	}
	if(len(SignInParams.Password) == 0){
		return "", errors.New("MissingPasswordField")
	}

	var verifyUser Entity

	userCollection := mongo.GetEntityCollection(mongo.GetDataBase(self.Session))

	//Find if user is registered in database
	findErr := userCollection.Find(bson.M{"email" : SignInParams.Email}).One(&verifyUser)

	if (findErr != nil){
		return "", errors.New("NonexistentEntity")
	}

	//Check if password provided matches hash on file
	passwordMatch := util.CheckPasswordHash(SignInParams.Password, verifyUser.PassHash)
	if (!passwordMatch){
		return "", errors.New("InvalidPassword")
	}

	//Return access token for usage
	accessToken, err := util.GetAccessToken(verifyUser.Id.Hex())
	if err != nil{
		return "", err
	}

	return accessToken, nil
}