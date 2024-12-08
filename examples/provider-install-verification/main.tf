terraform {
  required_providers {
    clickhouse = {
      source = "hashicorp.com/edu/clickhouse"
    }
  }
}

provider "clickhouse" {
  host     = "test.clickhouse.org"
  username = "xyz"
  password = "cddcd"
}

# data "clickhouse_coffees" "example" {}
