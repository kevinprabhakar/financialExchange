package customer

import (
	"errors"
	"financialExchange/util"
	"golang.org/x/net/html"
	"financialExchange/model"
	"financialExchange/sql"
	gosql "database/sql"
)

type CustomerController struct {
	Logger 			*util.Logger
	Database 		*sql.MySqlDB
}

func NewCustomerController(logger *util.Logger, db *sql.MySqlDB)(*CustomerController) {
	return &CustomerController{logger, db}
}

func (self *CustomerController)SignUp(SignUpParams model.CustomerSignUpParams)(error){
	if (!util.IsValidEmail(SignUpParams.Email)){
		return errors.New("InvalidEmailAddress")
	}

	if(len(SignUpParams.Password) < 6){
		return errors.New("PasswordTooShort")
	}

	if (SignUpParams.Password != SignUpParams.PasswordVerify){
		return errors.New("PasswordsDontMatch")
	}

	_, doesCustomerExist := self.Database.CheckIfCustomerInTable(SignUpParams.Email)

	if doesCustomerExist != nil {
		if doesCustomerExist == gosql.ErrNoRows{
		} else{
			return doesCustomerExist
		}
	}

	passHash, err := util.HashPassword(SignUpParams.Password)
	if (err != nil){
		return err
	}
	//First insert user
	newUser := model.Customer{
		FirstName	: SignUpParams.FirstName,
		LastName 	: SignUpParams.LastName,
		PassHash	: passHash,
		Email 		: html.EscapeString(SignUpParams.Email),
		Portfolio	: 0,
	}

	userId, insertErr := self.Database.InsertCustomerIntoTable(newUser)

	if (insertErr != nil){
		return insertErr
	}

	//Next insert new user portfolio
	//To start out with, each user will get $100.00 (Set at value and Cash Value)
	//This is solely for testing purposes
	//In the future, stripe integration will let us draw money from credit cards
	newUserPortfolio := model.Portfolio{
		Customer: userId,
		//Value = (Stock + Cash + Withdrawables)
		Value: model.NewMoneyObject(100.0),
		StockValue: model.NewMoneyObject(0.0),
		CashValue: model.NewMoneyObject(100.0),
		WithdrawableFunds: model.NewMoneyObject(0.0),
	}

	portfolioId, err := self.Database.InsertPortfolioToTable(newUserPortfolio)

	if err != nil{
		return err
	}

	//Now update user to link it to the just created Portfolio
	updateErr := self.Database.AttachPortfolioToCustomer(userId, portfolioId)

	if updateErr != nil{
		return err
	}

	return nil
}

func (self *CustomerController)SignIn(SignInParams model.CustomerSignInParams)(string, error){
	if (!util.IsValidEmail(SignInParams.Email)){
		return "", errors.New("MissingEmailField")
	}
	if(len(SignInParams.Password) == 0){
		return "", errors.New("MissingPasswordField")
	}

	userEmail, err := self.Database.CheckIfCustomerInTable(SignInParams.Email)

	//Find if user is registered in database
	if (err != nil){
		return "", errors.New("NonexistentUser")
	}

	userObj, err := self.Database.GetCustomerByEmail(userEmail)
	if err != nil{
		return "", err
	}

	//Check if password provided matches hash on file
	passwordMatch := util.CheckPasswordHash(SignInParams.Password, userObj.PassHash)
	if (!passwordMatch){
		return "", errors.New("InvalidPassword")
	}

	//Return access token for usage
	accessToken, err := util.GetAccessToken(userObj.Email)
	if err != nil{
		return "", err
	}

	return accessToken, nil
}

func (self *CustomerController)GetCurrUser(accessToken string)(*model.Customer, error){
	uid, err := util.VerifyAccessToken(accessToken)

	if (err != nil){
		return nil, err
	}
	userObj, findErr := self.Database.GetCustomerByEmail(uid)

	if (findErr != nil){
		return nil, findErr
	}

	return userObj, nil
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
