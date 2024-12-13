terraform {
  required_providers {
    clickhouse = {
      source = "hashicorp.com/edu/clickhouse"
    }
  }
}

provider "clickhouse" {
  host     = "localhost:9000"
  username = "admin"
  password = "test"
}

data "clickhouse_databases" "example" {}

data "clickhouse_users" "example" {}

data "clickhouse_roles" "roles" {}

resource "clickhouse_user" "test" {
  username = "manish"
  password = "test"
}

# resource "clickhouse_database" "test" {
#   database = "test_ddatabase"
# }

output "databases" {
  value = data.clickhouse_databases.example.databases
}

output "users" {
  value = data.clickhouse_users.example.users
}

output "roles" {
  value = data.clickhouse_roles.roles.roles
}
