package hclmodels

type Host struct {
	Name     string `hcl:"name,label"`
	Host     string `hcl:"host"`
	Port     int    `hcl:"port"`
	UserName string `hcl:"username,optional"`
	Password string `hcl:"password,optional"`
	SSH_Key  string `hcl:"ssh-key,optional"`
}
