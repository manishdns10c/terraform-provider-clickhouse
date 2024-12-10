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

resource "clickhouse_user" "test" {
  username = "test"
  password = "test"
}

output "databases" {
  value = data.clickhouse_databases.example.databases
}

output "users" {
  value = data.clickhouse_users.example.users
}
