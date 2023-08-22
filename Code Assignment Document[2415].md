
**Code Assignment**

You work in an e-commerce platform as a software engineer. A new campaign module which manipulates prices according to demand, is brought to you as an important business requirement.

* You have **products**, **orders** and **campaigns** in your domain.

* You create a product with **product code**, **price** and **stock**.

* You create an order with **product code** and **quantity**. Price can be assumed the current price of product at that moment.
* You create a campaign with a **name**, a **product code**, **duration**, **price manipulation limit** and **target sales count**.
* Campaign starts after creating and ends after given duration.

* Duration is given in hours.

* A price manipulation limit is the maximum percentage that you can increase or decrease the price of product according to demand.
* Target sales count is the product quantity you want to sell during the campaign.
* You will simulate time in your system. Time will start with 00:00 and it will be increased by you in any amount of **hour**.
* You are free to design your algorithm for how to calculate demand and how to increase and decrease the price during the campaign.



You will have scenario files in the context of this assignment (Sample scenario files are sent you along with this assignment document. **Please use them** to see if your program is running as expected).

* Scenario files have **commands** for the operations defined in business requirements.

* Your program will **read scenario file** and **produce output** for each command.

* You are free to choose the programming language (e.g. C#, Java, C++, Go)

The table given below defines all commands which should be recognized by your system.


|<p></p><p>**Command**</p>|<p></p><p>**What it does**</p>|
| :- | :- |
|create\_product PRODUCTCODE PRICE STOCK|Creates product in your system with given product information.|
|get\_product\_info PRODUCTCODE|Prints product information for given product code.|
|create\_order PRODUCTCODE QUANTITY|Creates order in your system with given information.|
|create\_campaign NAME PRODUCTCODE DURATION PMLIMIT TARGETSALESCOUNT|Creates campaign in your system with given information|
|get\_campaign\_info NAME|Prints campaign information for given campaign name|
|increase\_time HOUR|Increases time in your system.|


The table given below shows sample input and outputs for all possible commands in your system (Given outputs assumes commands are executed successfully)




|<p></p><p>**Command**</p>|<p></p><p>**Sample Output**</p>|
| :- | :- |
|create\_product P1 100 1000|Product created; code P1, price 100, stock 1000|
|get\_product\_info P1|Product P1 info; price 100, stock 1000|
|create\_order P1 3|Order created; product P1, quantity 3|
|create\_campaign C1 P1 10 20 100|Campaign created; name C1, product P1, duration 10, limit 20, target sales count 100|
|get\_campaign\_info C1|Campaign C1 info; Status Active, Target Sales 100, Total Sales 50, Turnover 5000, Average Item Price 100|
|increase\_time 1|Time is 01:00|

The last table shows steps in an example scenario file and their outputs (Price changing strategy in this example follows a linear pattern, you do not have to implement exactly same pattern.)


|**Steps in Example Input File**|**Output**|
| :- | :-: |
|create\_product P1 100 1000|Product created; code P1, price 100, stock 1000|
|create\_campaign C1 P1 5 20 100|Campaign created; name C1, product P1, duration 10, limit 20, target sales count 100|
|get\_product\_info P1|Product P1 info; price 100, stock 1000|
|increase\_time 1|Time is 01:00|
|get\_product\_info P1|Product P1 info; price 95, stock 1000|
|increase\_time 1|Time is 02:00|
|get\_product\_info P1|Product P1 info; price 90, stock 1000|
|increase\_time 1|Time is 03:00|
|get\_product\_info P1|Product P1 info; price 85, stock 1000|
|increase\_time 1|Time is 04:00|
|get\_product\_info P1|Product P1 info; price 80, stock 1000|
|increase\_time 2|Time is 06:00|
|get\_product\_info P1|Product P1 info; price 100, stock 1000|
|get\_campaign\_info C1|Campaign C1 info; Status Ended, Target Sales 100, Total Sales 0, Turnover 0, Average Item Price –|


Following criterias will be considered during the evaluation of the code assignment:

* Code should run as expected.
* Code Quality (Clean Code, SOLID, Applying Patterns “if necessary”, and other Software Craftsmanship techniques)
* Readability
* Unit testing. TDD approach will be favored.
* Packaging (how easy it is to run the code)
* Domain Modeling

You have 5 days to fulfill the assignment. You should submit this assignment by pushing into your repository on GitHub. Plz, keep in touch with us regarding to any enquiries w/ the assignment

Good luck :)
