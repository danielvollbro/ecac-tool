host "test" {
  host = "192.168.50.111"
  port = "22"
  username = "daniel"
  password = "123qweASD"
}

plugin "apt" {
  source = "github.com/danielvollbro/ecac-plugin-apt"
  version = "v0.0.2"
}

task "web" {
  plugin = "apt"
  config {
    packages = [ "nginx", "curl" ]
  }
}
