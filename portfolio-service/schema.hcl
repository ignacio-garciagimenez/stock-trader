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

  index "idx_name" {
    columns = [column.name]
    unique = true
  }
}

table "event_journal" {
  schema = schema.portfolio
  column "id" {
    null = false
    type = varchar(36)
  }
  column "timestamp" {
    null = false
    type = datetime(6)
  }
  column "name" {
    null = false
    type = varchar(256)
  }
  column "event_data" {
    null = false
    type = json
  }
  column "sent" {
    null = false
    type = boolean
    default = false
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_sent_x_timestamp" {
    columns = [
      column.sent,
      column.timestamp
    ]
  }
}

schema "portfolio" {
  charset = "utf8mb4"
  collate = "utf8mb4_0900_ai_ci"
}
