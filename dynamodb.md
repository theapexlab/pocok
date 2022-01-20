# Dyanmo basics
- dynamo db database modeling is different
- do not fake relational modeling
- must identify the access patterns before table design
- most applications need 1 table
- get everything in single query
- add EntityType to fields

## Single table design steps
1. draw entity diagram
2. identify relations in diagram
3. list all access patterns for entitys (crud, filters)
4. identify the primary key for each entity (hash + sort)
5. identify secondary indexes (LocalSI, GlobalSI)

# Key naming
- pk - Primary Key
- sk - Secondary Key
- gsi1pk - global secondary index 1, primary key
- gsi1sk - global secondary index 1, sort key
- lsi2sk - local secondary index 2, sort key
- etc...

# Access Patterns

## Organization
- CRUD

### Keys
- PK: ORG#{orgId}
- SK: #ANY#{orgId}

## Invoice
- CRUD
- Find invoices by status
- Find latest invoice by customerName

### Keys
- PK: ORG#{orgId}
- SK: INVOICE#{invId}
- LSI1SK: STATUS#{status}
    - all attributes projected
- LSI2SK: CUSTOMER#{customerName}#DATE#{createdAt}
    - customerEmail projected

