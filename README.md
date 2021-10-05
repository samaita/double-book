# Flash Sale Force Cancel & Refund Purchace

## Overview

Flash sale especially as big as 12.12 like in some e-commerce attract massive people. Since it scheduled within a short time range, people trying to get their favorite item since the first second. Massive access since first second is inevitable.

High amount of user making concurrent requests to the service. People might do add to cart, checkout, and pay at the very same time.

An Issue might arise from those concurrent requests and make several bad user experience. One of them: Multiple payment refunded even though the user make valid purchase during flash sale. Stock is displayed as available, they able to ATC, Checkout, and Pay.

Payment refunded will surely trigger user to make bad reviews. Stock mismatch is the issue to be fixed. How come the stock/inventory during flash sale goes mismatch with actual condition?

## Negative Inventory Issue Possibility

There is no issue if total inventory purchased is still lower than the inventory available. But when the total inventory purchased is higher, even by 1 item, this will trigger disappointed customer & inevitable bad review.

There is several possibility that would make inventory issue goes negative:

1. **PRD Incomplete: There is no inventory count** - The most fatal situation.
    - All product in flash sale treated as unlimited.
    - Anyone can buy it anytime during flash sale, even though usual Flash Sale has limited stock.
2. **PRD Incomplete: There is no inventory check during ATC, Checkout, and/or Payment**
    - Inventory information goes to waste, it exist but never used.
    - Another example of with unlimited stock flash sale.
3. **Testing Incomplete: Inventory management has bugs**
    - Bad logic implementation that make bugs, For examples:
        - Common usage inventory query update is using this math logic `inventory = inventory - 1`, without validate the remaining inventory quantity.
        ```
        UPDATE db_inventory SET stock = stock - 1
        ```
        - Inventory check exist, but remaining inventory quantity is never reduced. Another example of unlimited stock flash sale.
4. **Human Error**
    - PIC for flash sale give wrong input, for example: product A only prepared 10 item, but the PIC accidentally put 100.
5. **Architectural Issues: Inventory check exist, done to Slave/Secondary but DB met Replication Lag**
    - Replication Lag causing information in slave/secondary outdated compared to master in a split millisecond. Even 20-50ms lag can make a huge difference!
    - Invalid information due to replication lag will lead to inventory quantities invalid. For example:
        - Master already got 2 update query so the total stock now 0.
        - But due to replication lag, Slave just ACK 1 update so the total still 1.
        - In 20-50ms time window, another user that query to Slave will able to do another update stock - 1 to Master.
    - This condition might get worse since DB load much higher during flash sale
6. **Architectural Issues: Inventory check exist, done to cache first then DB, but cache has outdated value**
    - Cache exist to prevent too many load to DB
    - But, cache invalidation must be properly implemented to prevent such as: Cache says 10 item is still in stock, but actually only 9 left in DB.
    - This too will make inventory quantities invalid.
    - Using cache also put some disadvantages such as the backend service has no connection pool left or redis got too many connection.
7. **Architectural Issues: Inventory check & update to DB Master/Primary but done in separated process.**
    - This will open to race condition and make inventory quantities invalid.
    - For example:
        - 2 process done in paralel
        - Both check to DB with SELECT
        - Both got info that stock remaining is 1
        - both doing update to DB
        - Now DB stock remaining is -1

# Preliminary
There is two solution mentioned: easiest & proposed. Both focusing on how to address the ultimate problem:
> What is the most guaranteed way to ensure the amount of purchased inventory is not negative?

In the same time, the proposed solution will also address some possible architectural solution during implementation, such:
1. Keeping DB load low
2. Avoid Replication lag
3. Avoid exhausting cache resource
4. How to do simple test with 10.000 concurrency

But possibility such human error, will not be addressed here. That is one of the possible root cause of the inventory mismatch but let's focus on technical aspect.

The solution describe here will use RDBMS as DB choice since data consistency is the most important matter in this issue. 
All solution will also always assume flash sale condition with high traffic.

# Easy Solution

Judging from the possible causes of negative inventory, **inventory check & update** is the vital process.

As summary, to perform inventory check & update for flash sale condition, the easiest possible way is to use:

**Query to DB Master/Primary in a single process/transaction**

By using a query such this:
```
UPDATE db_inventory SET stock = stock - 1 WHERE stock > 0
```
- Pros
    - Avoid replication lag
    - Only update when the stock is really available or > 0
    - Require DB only
- Cons
    - Updating same cell will lead to DB Locking.
    - Thus more process queued in DB and cause higher load in DB.
    - In the end this approach will endanger entire business as other table in same DB Master might not be operating smoothly.
    - Many query queued due to lock. When the stock is 0, all that queue will have 0 rows affected. Wasted resource!
- Mitigation
    - DB for inventory will have to have a isolated DB infrastructure to ensure no other table in the same DB Master affected by flash sale event.
    - Cache might be added if stock is empty to prevent any other request goes to DB, with risk of cache invalid value

## But, is it the only way?

Of course not, the main idea for this possible way is to prepare the DB Master to handle the traffic, helped by cache. There might be solution such as multiple table, multiple DB but the idea is same: **preparing DB to take the traffic**.

As stated before **inventory check & update** is the vital process to prevent negative inventory. That is why, DB hold the most important role one here.

While that statement correct, there is another thing to consider:
1. **Timing of inventory check & update** - Is it during ATC, Checkout, Payment, or all the time? 1 time check such as ATC only need considerable amount of confidence.
2. **Cost DB** - using the easiest solution will put all the load to DB Master, it need to have large pool of resource, more RAM, more CPU to satisfy incoming traffic
3. **Cost Dedicated Cache** - If choosing dedicated cache such as Redis to prevent unnecessary traffic stock check to DB Master, a cache cluster required to met the traffic

# Proposed Solution

## Main Idea

The proposed solution main idea is:
> How to prevent traffic abnormality overload the most important role in the entire backend: DB.

Putting the load to DB is requiring considerable resource and **that is not cheap**. Meanwhile, **DB is the most difficult one to scale**. The effort to increase DB resource require much time & possibly downtime of service.

It is much easier to adding more instance to cache cluster. By only adding instance for cache value that perishable.

It even much easier to adding more instance to backend service. Terraforming can handle it at ease with concept: copy the binary & run it!

So, how this solution able to prevent traffic abnormality?

**By forcing it into a queue with Message Queue** - So backend service & MQ will act as enforcer, prevent most of the traffic get to DB.

![ilustration](https://1.bp.blogspot.com/-StCX890PlF0/YVw9SMWb9tI/AAAAAAAATCU/NJNAeLE8nKAuKk518Mxnivllr92aolvtACLcBGAsYHQ/s2719/Evermos%2BE-Commerce%2BDiagram%25282%2529.png, "Ilustration main idea ")

So basically, backend will only relay the request to MQ. MQ will let MQ Consumer to process the request 1 by 1. It is expected for MQ to build depth as the process consuming much lower than MQ depth.

This will be done during ATC process, to prevent people doing ATC if the stock is empty. It doesn't make sense if user able to perform ATC a product that already empty, right?

## Architecture Diagram

![poc architecture diagram](https://1.bp.blogspot.com/-ZS5mQVdzlv4/YVtXwnAPWGI/AAAAAAAATB4/OLLHu5WWd5EDijipZO7JvzoKYxqjm8gNgCLcBGAsYHQ/s2869/Evermos%2BArchitecture%2BDiagram.png, "Architecture Diagram")

To perform this solution, the stack used are:
- RDBMS - SQLite for simplicity of PoC
- MQ - NSQ
- Go - for main logic, consist of
    - main app + main consumer, both under 1 binary
    - consumer ATC + cron stock evaluation, both under 1 binary

## Logic App

![Logic](https://1.bp.blogspot.com/-Nee8-1Dyjl4/YVthhuHUmdI/AAAAAAAATCA/NdmPBCWyjf8eWjusFWYyb7U8JXgt7Cg2wCLcBGAsYHQ/s2048/Evermos%2BFlow.png)

- Mass traffic request hit the backend service.
- Backend service handler publish each request into a ATC message to a MQ ATC
- Backend service handler set to wait for confirmation listen to a channel with go routine
- MQ ATC set to 1 topic & N channel per flash sale, where N is the amount of product listed during flash sale. Each channel only have 1 consumer 1 goroutine to process the ATC & stock decrement.
- Consumer ATC discard all irrelevant message received, for example: Consumer assigned to process ATC for product 555 only, if any message inbound contain ATC for product 123, this consumer will discard it.
- Consumer ATC discard all ATC message if a process success query but no row affected.
- Consumer ATC will back processing message if there is back in stock message from cron.
- Consumer ATC will then publish an information to MQ to be consumed by Main consumer success or not.
- Again, Main Consumer discard all irrelevant message received by only process message with his IP.
- Main Consumer then gave information to Backend service handler to return the response.
- A cron run to evaluate Cart & Invoice. Any cart or invoice with flash sale product will have to be paid under 5 minutes or the product cart will be removed and invoice will be cancelled. Thus the stock will be +1

## Database Design

![e-commerce basic diagram](https://1.bp.blogspot.com/-baDJZkV6aDo/YVw5i1ONsoI/AAAAAAAATCM/YJ6nECtIOfQnem_c_C017uTwCvpwvlvkACLcBGAsYHQ/s2048/Evermos%2BE-Commerce%2BDiagram%25281%2529.png "Database Design")

The type data used here are designed for SQLite3, text for every `create_time` & `update_time` should be changed into timestamp instead.

for every `id`, could be use specific data such as `int64` for fast growth entity.

There are 10 Table:
- `User` contain all registered user
- `Shop` contain all created shop, 1 shop owned by 1 user, generated by default
- `Product` contain all product owned by a shop, 1 shop owned multiple product
- `Stock` contain product information regarding stock & it remaining, 1 product 1 stock info. Flashsale product & non flash sale product should have different product ID.
- `Cart` contain every cart an user own, 1 user only allowed to have 1 cart at the same time.
- `Cart Detail` contain every cart detail for each cart ID.
- `Invoice` is a copy of cart, copied if an user checkout & the status change is to indicate the paid / unpaid status.
- `Invoice Detail` contain every item purchased under 1 invoice ID
- `Flash Sale` hold information about flash sale management such as schedule
- `Flash Sale Detail` contain all product id that will be in flash sale.

## How to Run the Solution
```
docker compose
```
## How to test the Solution

Hit API X to create a new flash sale
```

```
Hit API Y to add product with predefined stock
```

```
Hit API Z to add product to flash sale
```

```

Note flash sale ID & product ID
in folder functional_test run
```
-flashsale=123 -productID=456 -h your_sdocker or http://....
```
wait for bot to struggling in 30s, ahh feels like Squid Game!
```

```
The script hit API A to get current stock for the product ID

The script also hit API B to get all invoice with product ID 

## Conclusion
There is no perfect solution, this one might be one of the alternative out there. Hopefully this solution can be your consideration to prevent double book or inventory negative in your service.

Pros:
 - People will fight in ATC, once it already in their cart, they can proceed at ease.
 - No DB lock, each cell of flash sale stock can only write by 1 go routine, not by all request.
 - Relatively lower DB load.

Cons:
 - MQ message depth will inevitable especially in channel of a product. But, disk is cheaper than CPU & Memory.
 - Bad implementation could lead to panic mmeory since using channel to communicate between goroutine
