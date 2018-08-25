package customer

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"errors"
	"financialExchange/util"
	"financialExchange/mongo"
	"golang.org/x/net/html"
	"financialExchange/portfolio"
	"financialExchange/api"
)

type CustomerController struct {
	Session 		*mgo.Session
	Logger 			*util.Logger
}

func NewCustomerController(session *mgo.Session, logger *util.Logger)(*CustomerController) {
	return &CustomerController{session, logger}
}

func (self *CustomerController)SignUp(SignUpParams CustomerSignUpParams)(error){
	if (!util.IsValidEmail(SignUpParams.Email)){
		return errors.New("InvalidEmailAddress")
	}

	if(len(SignUpParams.Password) < 6){
		return errors.New("PasswordTooShort")
	}

	if (SignUpParams.Password != SignUpParams.PasswordVerify){
		return errors.New("PasswordsDontMatch")
	}

	customerCollection := mongo.GetCustomerCollection(mongo.GetDataBase(self.Session))

	var findUser Customer

	//Verify user has not signed up with this email before
	err := customerCollection.Find(bson.M{ "email": SignUpParams.Email}).One(&findUser)

	if (err != mgo.ErrNotFound){
		return err
	}

	passHash, err := util.HashPassword(SignUpParams.Password)
	if (err != nil){
		return err
	}
	//First insert user
	newUser := Customer{
		Id			: bson.NewObjectId(),
		FirstName	: SignUpParams.FirstName,
		LastName 	: SignUpParams.LastName,
		PassHash	: passHash,
		Email 		: html.EscapeString(SignUpParams.Email),
		Portfolio 	: bson.NewObjectId(),
	}

	insertErr := customerCollection.Insert(newUser)

	if (insertErr != nil){
		return insertErr
	}

	//Next insert new user portfolio
	newUserPortfolio := portfolio.Portfolio{
		Id: newUser.Portfolio,
		User: newUser.Id,
		Value: api.NewMoneyObject(0.0),
		StockValue: api.NewMoneyObject(0.0),
		CashValue: api.NewMoneyObject(0.0),
		WithdrawableFunds: api.NewMoneyObject(0.0),
		OwnedShares: make(map[bson.ObjectId]int, 0),
		Orders: []bson.ObjectId{},
		Transactions: []bson.ObjectId{},
	}

	portfolioCollection := mongo.GetPortfolioCollection(mongo.GetDataBase(self.Session))

	insertErr = portfolioCollection.Insert(newUserPortfolio)
	if (insertErr != nil){
		return insertErr
	}

	return nil
}

func (self *CustomerController)SignIn(SignInParams CustomerSignInParams)(string, error){
	if (!util.IsValidEmail(SignInParams.Email)){
		return "", errors.New("MissingEmailField")
	}
	if(len(SignInParams.Password) == 0){
		return "", errors.New("MissingPasswordField")
	}

	var verifyUser Customer

	userCollection := mongo.GetCustomerCollection(mongo.GetDataBase(self.Session))

	//Find if user is registered in database
	findErr := userCollection.Find(bson.M{"email" : SignInParams.Email}).One(&verifyUser)

	if (findErr != nil){
		return "", errors.New("NonexistentUser")
	}

	//Check if password provided matches hash on file
	passwordMatch := util.CheckPasswordHash(SignInParams.Password, verifyUser.PassHash)
	if (!passwordMatch){
		return "", errors.New("InvalidPassword")
	}

	//Return access token for usage
	accessToken, err := GetAccessToken(verifyUser.Id.Hex())
	if err != nil{
		return "", err
	}

	return accessToken, nil
}

func (self *CustomerController)GetCurrUser(accessToken string)(*Customer, error){
	uid, err := VerifyAccessToken(accessToken)

	if (err != nil){
		return nil, err
	}
	if (!bson.IsObjectIdHex(uid)){
		return nil, errors.New("InvalidBSONId")
	}

	userCollection := mongo.GetCustomerCollection(mongo.GetDataBase(self.Session))

	var returnUser Customer

	findErr := userCollection.Find(bson.M{ "_id" : bson.ObjectIdHex(uid)}).One(&returnUser)

	if (findErr != nil){
		return nil, findErr
	}

	return &returnUser, nil
}

//func (self *CustomerController)GetUsers(params string)(*[]Customer, error){
//	dec := json.NewDecoder(strings.NewReader(params))
//	var userIds UserIdList
//	decErr := dec.Decode(&userIds)
//
//	if (decErr != nil){
//		return nil, errors.New("CANT DECODE")
//	}
//
//	userCollection := mongo.GetCustomerCollection(mongo.GetDataBase(self.Session))
//
//	var userIdListBson []bson.ObjectId
//
//	for _, uid := range userIds.IdList{
//		if (!bson.IsObjectIdHex(uid)){
//			return nil, errors.New("InvalidBSONId")
//		}else{
//			userIdListBson = append(userIdListBson, bson.ObjectIdHex(uid))
//		}
//	}
//
//	var userList []Customer
//
//	findErr := userCollection.Find(bson.M{"_id" : bson.M{"$in" : userIdListBson}}).All(&userList)
//
//	if (findErr != nil){
//		return nil, errors.New("CANT FIND USERS FROM BSON")
//	}
//
//	return &userList, nil
//}
//
//func(self *CustomerController)DeleteUser(uid bson.ObjectId)(error){
//	userCollection := mongo.GetCustomerCollection(mongo.GetDataBase(self.Session))
//	goalCollection := mongo.GetGoalCollection(mongo.GetDataBase(self.Session))
//	postCollection := mongo.GetPostCollection(mongo.GetDataBase(self.Session))
//
//	var findUser Customer
//	findUserErr := userCollection.Find(bson.M{"_id":uid}).One(&findUser)
//
//	if (findUserErr != nil){
//		return findUserErr
//	}
//
//	if (len(findUser.Goals)==0){
//		return errors.New("UserHasNoGoals")
//	}
//
//	removePostErr := postCollection.Remove(bson.M{"owner":uid})
//
//	if (removePostErr != nil){
//		return removePostErr
//	}
//
//	removeGoalsErr := goalCollection.Remove(bson.M{"owner":uid})
//	if (removeGoalsErr != nil){
//		return removeGoalsErr
//	}
//
//	removeUserErr := userCollection.Remove(bson.M{"_id":uid})
//	if (removeUserErr != nil){
//		return removeUserErr
//	}
//
//	return nil
//}
