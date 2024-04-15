package main

type creds struct {
	login, password string
}

func (x creds) GetCredentials(int) (string, string, error) {
	return x.login, x.password, nil
}
