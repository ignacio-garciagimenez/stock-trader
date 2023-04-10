# Golang Stock Trader


> An implementation of [IBM's Stock Trader](https://developer.ibm.com/blogs/introducing-stocktrader/) example, built with Golang and the following software architecture and design styles: **Microservices Architecture**, **Vertical Slice Architecture** , **CQRS**, **Domain Driven Design (DDD)**, **Event Driven Architecture**.

💡 The focus of this application is mostly on the technical side. It is not intended to be a fully fledged stock trading application. The features set will be small and simple.

🌀 This Application is `in-progress` and I will add new features and technologies over time. 

## Bounded Contexts

- Portfolio
- Broker
- Wire Transfers

## Commands
### Portfolio 
- [ ] Open Portfolio
- [ ] Close Portfolio
- [ ] Place Order
- [ ] Process Trade
- [ ] Acknoledge Order Failure
- [ ] Receive Funds
- [ ] Send Funds
- [ ] Accept Refund

### Wire Transfers
- [ ] Request Funds
- [ ] Send Funds
- [ ] Refund Sender

### Broker
- [ ] Place Order
- [ ] Complete Order
- [ ] Cancel Order


## Features and Technologies
- ✅ Using `Vertical Slice Architecture` as a high level architecture
- ✅ Using `Event Driven Architecture` 
- ✅ Using `CQRS Pattern`
- ✅ Using [Echo](https://github.com/labstack/echo) framework
- `To be Continued`

## Roadmap
- [ ] 🚧 Implement commands following `Domain Driven Design's tactical patterns` 
- [ ] 🚧 Add DBs
- [ ] 🚧 Choose Message Broker
- [ ] 🚧 Implement message relay (CDC or custom relay implementation)
- [ ] 🚧 Dockerize and include Compose (or batect)
- [ ] 🚧 Add `Identity Management` and `OAuth`


