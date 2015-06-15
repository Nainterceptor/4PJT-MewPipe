Database structure and optimization
===================================

## Which Database Engine ?

Mewpipe must be ready to scale quickly and with ease. With that in mind, we had to make the choice of NoSQL. 
MongoDB is an open-source document database engine, and the leading NoSQL database. 
It's a good choice because this engine is agile and scalable through many ways: Replica sets and Sharding.

## Our schema

MongoDB permits evolution of our schema over time. Some keys must be omitted when empty.

Our collections:

* media
* media.chunks
* media.files
* media.shareCounts
* media.thumbnails.chunks
* media.thumbnails.files
* media.views
* system.indexes
* users

### Users

```JavaScript
{
	"_id" : ObjectId("555a076a2fd06c1891000001"),
	"createdAt" : ISODate("2015-06-13T14:46:40.556Z"), //Bad time, fault on my old
	"name" : {
		"firstname" : "Foo",
		"lastname" : "Bar",
		"nickname" : "FooBar"
	},
	"email" : "foo@bar.com",
	"roles" : [ ], //List of roles (strings), available : admin
	"hashedpassword" : "$2a$10$LD9dA75m2U6jx8cRRKmc.u5DQYMHACKym.4OZkJ0T91qSVJCZHvU2", // Bcrypt encoded
	"usertokens" : [
		{
			"token" : ObjectId("557cab510000000000000000"),
			"expireat" : ISODate("2015-06-13T23:14:41.330Z")
		}
	],
	"twitter" : {
		"userId" : "99999999"
	}
}
```

HashedPassword uses [bcrypt](http://codahale.com/how-to-safely-store-a-password/).
When User tokens are cleaned, all older tokens are removed

### Media

When we store a media, it's splitted between many collections :

* media, which store some metadata like title, summary, etc... And do the link with other media metadata
* media.files and media.chunks, which store binary (chunks) and metadata of the file (md5, content type...) 
through GridFS
* media.thumbnails.files and media.thumbnails.chunks, which store binary (chunks) and metadata of the thumbnail 
(md5, content type...) through GridFS

```JavaScript
{
	"_id" : ObjectId("557c42502fd06c127000001c"),
	"createdAt" : ISODate("2015-06-13T14:46:40.556Z"),
	"title" : "Penguins_of_Madagascar.webm",
	"summary" : "My amazing video",
	"publisher" : { //Embed the user who publish the video
		"_id" : ObjectId("555a076a2fd06c1891000002"),
		"name" : {
			"firstname" : "Admin",
			"lastname" : "Admin",
			"nickname" : "Admin"
		},
		"email" : "admin@admin.com"
	},
	"file" : ObjectId("557c42502fd06c127000001d"), //ObjectID of the file
	"thumbnail" : ObjectId("557c42502fd06c1270000023"), //ObjectID of the thumbnail
	"scope" : "public",
	"views" : 2, //Cache of the number of views, calculated from media.views 
	"shares" : 2 //Cache of the number of shares, calculated from media.shareCounts, 
}
```

### Shares & Views

This collection groups views by user and media, the aim is to be able to calculate some statistics, 
like the most viewed media of a user.

```JavaScript
{
	"_id" : ObjectId("556e3ba94d96f23c00b302be"),
	"media" : ObjectId("556dda352fd06c0fb6000002"),
	"user" : ObjectId("555a076a2fd06c1891000001"),
	"count" : 4
}
```

An aggregation query calculates the number of views mongo-side and cache the result in the media.

## Our strategy

* Media and thumbnails are split in files of 16Mb, they are concatenated on the fly by MongoS (MongoDB Shard). 
This is a powerful and fast system.
* The user is embedded inside the media. This method facilitates the access to the user's info (like his nickname). 
During update, those user-specific data will be updated on all related media.

## Optimizations

Some indexes are useful for speeding up a MongoDB database and also in adding some controls.
In our case we have, for example, indexes on the user's email and his Twitter ID.

All indexes are in the collection _system.indexes_