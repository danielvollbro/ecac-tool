host "test" {
  host = "192.168.50.111"
  port = "22"
  username = "daniel"
  password = "123qweASD"
}

plugin "apt" {
  source = "github.com/example/apt-plugin"
  version = "v1.2.0"
}

task "web" {
  plugin = "apt"
  config {
    packages = [ "nginx", "curl" ]
  }
}
