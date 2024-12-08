terraform {
  required_providers {
    clickhouse = {
      source = "hashicorp.com/edu/clickhouse"
    }
  }
}

provider "clickhouse" {
  host     = "localhost:9000"
  username = "test_user"
  password = "oops"
}

data "clickhouse_databases" "example" {}

output "databases" {
  value = data.clickhouse_databases.example.databases
}