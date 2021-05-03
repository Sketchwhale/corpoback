#ABOUT

Corpoback lets you view, edit, and create companies, and also add beneficial owners to these.

#HOW

Adding a new company is done with a simple cURL command:

<blockquote>
curl -X POST -H "Content-Type: application/json" -d '{
               "Id": "3",
               "Name": "Mitsubishi",
               "Address": "Thetion for my new post",
               "City": "my articles content",
	"Country": "snaps",
	"Phone": "4343"
           }' http://localhost:10000/company
</blockquote>

Similarly, adding an owner to a company is handled thusly:

<blockquote>
curl -X POST -H "Content-Type: application/json" -d '{
                              "Owner": "SonyJapan",
                              "Owned": "SonyEurope"
                          }' http://localhost:10000/addOwnership
</blockquote>

Note that Id, Name, Address, City, and Country are required.

The end points are as follows:

* create new company
	* /company
* retrieve all companies
	* /all
* retrieve or edit specific company
	* /company/{Id}
* add beneficial owner
	* /addOwnership

Ownerships are relationships between companies. They are detailed when retrieving a single company, but are not part of the company object.
