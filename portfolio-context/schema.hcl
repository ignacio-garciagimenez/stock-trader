table "portfolios" {
  schema = schema.portfolio
  column "id" {
    null = false
    type = varchar(36)
  }
  column "name" {
    null = false
    type = varchar(30)
  }
  primary_key {
    columns = [column.id]
  }
}

schema "portfolio" {
  charset = "utf8mb4"
  collate = "utf8mb4_0900_ai_ci"
}
