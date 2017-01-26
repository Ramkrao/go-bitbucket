package bitbucket

// User is the sub struct of Client
type User struct {
	c *Client
}

// Profile is getting the user data
func (u *User) Profile() interface{} {
	url := GetApiBaseURL() + "/user/"
	return u.c.execute("GET", url, "")
}

// Emails is getting user's emails
func (u *User) Emails() interface{} {
	url := GetApiBaseURL() + "/user/emails"
	return u.c.execute("GET", url, "")
}
