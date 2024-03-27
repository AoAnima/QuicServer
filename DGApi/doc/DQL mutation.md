DQL mutation - DQL
Dgraph Query Language (DQL) is Dgraph’s proprietary language to add, modify, delete and fetch data.

Fetching data is done through DQL Queries. Adding, modifying or deleting data is done through DQL Mutations.

This overview explains the structure of DQL Mutations and provides links to the appropriate DQL reference documentation.

DQL mutations support JSON or RDF format.

set block
In DQL, you add data using a set mutation, identified by the set keyword.

   {
   "set": [
     {
       "name":"Star Wars: Episode IV - A New Hope",
       "release_date": "1977-05-25",
       "director": {
         "name": "George Lucas",
         "dgraph.type": "Person"
       },
       "starring" : [
         {
           "name": "Luke Skywalker"
         },
         {
           "name": "Princess Leia"
         },
         {
           "name": "Han Solo"
         }
       ]
     },
     {
       "name":"Star Trek: The Motion Picture",
       "release_date": "1979-12-07"
     }
   ]
 }  
Copy
{
  set {
    # triples in here
    _:n1 <name> "Star Wars: Episode IV - A New Hope" .
    _:n1 <release_date>  "1977-05-25" .
    _:n1 <director> _:n2 .
    _:n2 <name> "George Lucas" .

  }
}
Copy
triples are in RDF format.

Node reference
A mutation can include a blank nodes as an identifier for the subject or object, or a known UID.

{
  set {
    # triples in here
    <0x632ea2> <release_date>  "1977-05-25" .
  }
}
Copy
will add the release_date information to the node identified by UID 0x632ea2.

language support
{
  set {
    # triples in here
    <0x632ea2> <name>  "Star Wars, épisode IV : Un nouvel espoir"@fr .
  }
}
Copy
delete block
A delete mutation, identified by the delete keyword, removes triples from the store.

For example, if the store contained the following:

<0xf11168064b01135b> <name> "Lewis Carrol"
<0xf11168064b01135b> <died> "1998"
<0xf11168064b01135b> <dgraph.type> "Person" .
Copy
Then, the following delete mutation deletes the specified erroneous data, and removes it from any indexes:

{
  delete {
     <0xf11168064b01135b> <died> "1998" .
  }
}
Copy
Wildcard delete
In many cases you will need to delete multiple types of data for a predicate. For a particular node N, all data for predicate P (and all corresponding indexing) is removed with the pattern S P *.

{
  delete {
     <0xf11168064b01135b> <author.of> * .
  }
}
Copy
The pattern S * * deletes all the known edges out of a node, any reverse edges corresponding to the removed edges, and any indexing for the removed data.

Note For mutations that fit the S * * pattern, only predicates that are among the types associated with a given node (using dgraph.type) are deleted. Any predicates that don’t match one of the node’s types will remain after an S * * delete mutation.

{
  delete {
     <0xf11168064b01135b> * * .
  }
}
Copy
If the node S in the delete pattern S * * has only a few predicates with a type defined by dgraph.type, then only those triples with typed predicates are deleted. A node that contains untyped predicates will still exist after a S * * delete mutation.

Note The patterns * P O and * * O are not supported because it’s inefficient to store and find all the incoming edges.

Deletion of non-list predicates
Deleting the value of a non-list predicate (i.e a 1-to-1 relationship) can be done in two ways.

Using the wildcard delete (star notation) mentioned in the last section.
Setting the object to a specific value. If the value passed is not the current value, the mutation will succeed but will have no effect. If the value passed is the current value, the mutation will succeed and will delete the non-list predicate.
For language-tagged values, the following special syntax is supported:

{
  delete {
    <0x12345> <name@es> * .
  }
}
Copy
In this example, the value of the name field that is tagged with the language tag es is deleted. Other tagged values are left untouched.

upsert block
Upsert is an operation where:

A node is searched for, and then
Depending on if it is found or not, either:
Updating some of its attributes, or
Creating a new node with those attributes.
The upsert block allows performing queries and mutations in a single request. The upsert block contains one query block and mutation blocks.

The structure of the upsert block is as follows:

upsert {
  query <query block>
  mutation <mutation block 1>
  [mutation <mutation block 2>]
  ...
}
Copy
Execution of an upsert block also returns the response of the query executed on the state of the database before mutation was executed. To get the latest result, you have to execute another query after the transaction is committed.

Variables defined in the query block can be used in the mutation blocks using the uid and val functions.

conditional upsert
The upsert block also allows specifying conditional mutation blocks using an @if directive. The mutation is executed only when the specified condition is true. If the condition is false, the mutation is silently ignored. The general structure of Conditional Upsert looks like as follows:

upsert {
  query <query block>
  [fragment <fragment block>]
  mutation [@if(<condition>)] <mutation block 1>
  [mutation [@if(<condition>)] <mutation block 2>]
  ...
}
Copy
The @if directive accepts a condition on variables defined in the query block and can be connected using AND, OR and NOT.

Example of Conditional Upsert
Let’s say in our previous example, we know the company1 has less than 100 employees. For safety, we want the mutation to execute only when the variable v stores less than 100 but greater than 50 UIDs in it. This can be achieved as follows:

curl -H "Content-Type: application/rdf" -X POST localhost:8080/mutate?commitNow=true -d  $'
upsert {
  query {
    v as var(func: regexp(email, /.*@company1.io$/))
  }

  mutation @if(lt(len(v), 100) AND gt(len(v), 50)) {
    delete {
      uid(v) <name> * .
      uid(v) <email> * .
      uid(v) <age> * .
    }
  }
}' | jq
Copy
We can achieve the same result using json dataset as follows:

curl -H "Content-Type: application/json" -X POST localhost:8080/mutate?commitNow=true -d '{
  "query": "{ v as var(func: regexp(email, /.*@company1.io$/)) }",
  "cond": "@if(lt(len(v), 100) AND gt(len(v), 50))",
  "delete": {
    "uid": "uid(v)",
    "name": null,
    "email": null,
    "age": null
  }
}' | jq
Copy