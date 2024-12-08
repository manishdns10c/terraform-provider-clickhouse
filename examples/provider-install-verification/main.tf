terraform {
  required_providers {
    clickhouse = {
      source = "hashicorp.com/edu/clickhouse"
    }
  }
}

provider "clickhouse" {}

data "clickhouse_coffees" "example" {}
