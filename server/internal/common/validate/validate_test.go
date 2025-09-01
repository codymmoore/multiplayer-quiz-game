package validate

import (
	"testing"
)

const (
	ValidUsername = "test-username"
	ValidEmail    = "test@email.com"
	ValidPassword = "testPassword1234#?!@$%^&*-"
)

func TestUsername_Success(t *testing.T) {
	username := ValidUsername
	if err := Username(username); err != nil {
		t.Errorf(`Username("%s") = "%v", expected "<nil>"`, username, err)
	}
}

func TestUsername_Short(t *testing.T) {
	username := "uh"
	err := Username(username)
	if err == nil {
		t.Errorf(
			`Username("%s") = "%v", expected "username must be between %d and %d characters"`,
			username,
			err,
			MinUsernameLength,
			MaxUsernameLength,
		)
	}
}

func TestUsername_Long(t *testing.T) {
	username := "testInvalidUsername"
	err := Username(username)
	if err == nil {
		t.Errorf(
			`Username("%s") = "%v", expected "username must be between %d and %d characters"`,
			username,
			err,
			MinUsernameLength,
			MaxUsernameLength,
		)
	}
}

func TestUsername_IllegalCharacter(t *testing.T) {
	username := "username*"
	err := Username(username)
	if err == nil {
		t.Errorf(
			`Username("%s") = "%v", expected "illegal character. username must contain only letters, numbers, underscores, and hyphens"`,
			username,
			err,
		)
	}
}

func TestEmail_Success(t *testing.T) {
	email := ValidEmail
	if err := Email(email); err != nil {
		t.Errorf(`Email("%s") = "%v", expected "<nil>"`, email, err)
	}
}

func TestEmail_Format(t *testing.T) {
	email := "invalidEmail"
	err := Email(email)
	if err == nil {
		t.Errorf(`Email("%s") = "%v", expected "invalid email format"`, email, err)
	}
}

func TestPassword_Success(t *testing.T) {
	password := ValidPassword
	if err := Password(password); err != nil {
		t.Errorf(`Password("%s") = "%v", expected "<nil>"`, password, err)
	}
}

func TestPassword_Short(t *testing.T) {
	password := "Invalid123456!"
	err := Password(password)
	if err == nil {
		t.Errorf(
			`Password("%s") = "%v", expected "password must be between %d and %d characters"`,
			password,
			err,
			MinPasswordLength,
			MaxPasswordLength,
		)
	}
}

func TestPassword_Long(t *testing.T) {
	password := "InvalidPassword-1234567890123457890123456789012345678901234567-1!"
	err := Password(password)
	if err == nil {
		t.Errorf(
			`Password("%s") = "%v", expected "password must be between %d and %d characters"`,
			password,
			err,
			MinPasswordLength,
			MaxPasswordLength,
		)
	}
}

func TestPassword_MissingUpper(t *testing.T) {
	password := "testpassword1234!"
	err := Password(password)
	if err == nil {
		t.Errorf(
			`Password("%s") = "%v", expected "invalid password. Must contain at least one of the following: upper case English character, lower case English character, number, special character"`,
			password,
			err,
		)
	}
}

func TestPassword_MissingLower(t *testing.T) {
	password := "TESTPASSWORD1234!"
	err := Password(password)
	if err == nil {
		t.Errorf(
			`Password("%s") = "%v", expected "invalid password. Must contain at least one of the following: upper case English character, lower case English character, number, special character"`,
			password,
			err,
		)
	}
}

func TestPassword_MissingNumber(t *testing.T) {
	password := "testPassword#?!@$%^&*-"
	err := Password(password)
	if err == nil {
		t.Errorf(
			`Password("%s") = "%v", expected "invalid password. Must contain at least one of the following: upper case English character, lower case English character, number, special character"`,
			password,
			err,
		)
	}
}

func TestPassword_MissingSymbol(t *testing.T) {
	password := "testPassword12345"
	err := Password(password)
	if err == nil {
		t.Errorf(
			`Password("%s") = "%v", expected "invalid password. Must contain at least one of the following: upper case English character, lower case English character, number, special character"`,
			password,
			err,
		)
	}
}

func TestPassword_IllegalCharacter(t *testing.T) {
	password := ValidPassword + "😘"
	err := Password(password)
	if err == nil {
		t.Errorf(
			`Password("%s") = "%v", expected "password contains illegal characters"`,
			password,
			err,
		)
	}
}
