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

func NewEntityController (session *mgo.Session, logger util.Logger)(*EntityController){
	return &EntityController{
		Session: session,
		Logger: logger,
	}
}

func(self *EntityController) CreateEntity(params CreateEntityParams)(error){
	//Field Validation
	if (!util.IsValidEmail(params.Email)){
		return errors.New("InvalidEmailAddress")
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

	/***
	todo: find out how (and where) to create the pricebook
	mongo creates collections upon data insertion, so should I:
		- create it when creating the entity
		- create when first trade is made
	--going with second one for now
	***/

	return nil
}
