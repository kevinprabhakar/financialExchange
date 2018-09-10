package entity

import (
	"financialExchange/util"
	"financialExchange/sql"
	"financialExchange/model"
	"time"

	"errors"
	gosql "database/sql"

)

type EntityController struct{
	Database		*sql.MySqlDB
	Logger 			*util.Logger
}

func NewEntityController (database *sql.MySqlDB, logger *util.Logger)(*EntityController){
	return &EntityController{
		Database: database,
		Logger: logger,
	}
}

func(self *EntityController) CreateEntity(params model.CreateEntityParams)(error){
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
	_, err := self.Database.CheckIfEntityInTable(params.Email)
	if err != nil {
		if err == gosql.ErrNoRows{
		} else{
			return errors.New("EntityExists")
		}
	}

	//Verify symbol has not been used before
	_, err = self.Database.CheckIfSecurityInTable(params.Symbol)
	if err != nil {
		if err == gosql.ErrNoRows{
		} else{
			return errors.New("SecurityNameExists")
		}
	}

	//Hash Password
	passHash, err := util.HashPassword(params.Password)
	if (err != nil){
		return err
	}

	creationTime := time.Now().Unix()

	//Create entity and insert to database
	newEntity := model.Entity{
		Name: params.Name,
		Email: params.Email,
		PassHash: passHash,
		Security: -1,
		Created: creationTime,
		Deleted: creationTime,
	}

	entityId, insertErr := self.Database.InsertEntityIntoTable(newEntity)

	if insertErr != nil{
		return err
	}

	if params.CreateSecurity{
		//Create security and insert to database
		newSecurity := model.Security{
			Entity: entityId,
			Symbol: params.Symbol,
		}

		securityId, insertErr := self.Database.InsertSecurityIntoTable(newSecurity)

		if insertErr != nil {
			return err
		}

		err := self.Database.UpdateEntityWithSecurityID(entityId, securityId)
		if err != nil{
			return err
		}
	}

	return nil
}

func (self *EntityController) SignIn(SignInParams model.SignInEntityParams)(string, error){
	if (!util.IsValidEmail(SignInParams.Email)){
		return "", errors.New("MissingEmailField")
	}
	if(len(SignInParams.Password) == 0){
		return "", errors.New("MissingPasswordField")
	}

	findEntity, findErr := self.Database.GetEntityByEmail(SignInParams.Email)

	if (findErr != nil){
		return "", errors.New("NonexistentEntity")
	}

	//Check if password provided matches hash on file
	passwordMatch := util.CheckPasswordHash(SignInParams.Password, findEntity.PassHash)
	if (!passwordMatch){
		return "", errors.New("InvalidPassword")
	}

	//Return access token for usage
	accessToken, err := util.GetAccessToken(findEntity.Email)
	if err != nil{
		return "", err
	}

	return accessToken, nil
}