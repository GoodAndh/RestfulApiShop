package auth

import "golang.org/x/crypto/bcrypt"

//hash the password
func HashPassword(password string)(string,error)  {
	hashed,err:=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed),nil
}

//check hashed password
func CheckHashedPassword(hashed string,password []byte)error  {
	if err:=bcrypt.CompareHashAndPassword([]byte(hashed),password);err!=nil{
		return err
	}
	return nil
}
