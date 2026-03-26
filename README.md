# Prerequisites
## What you need installed in your system to run this program:
* Go >=1.25
* Postgres >=16
* sqlc >=1.30

## Database setup
* postgres service needs to be up and running
* set up a user with writing permissions for you
* create database 'bloggator'

## Config file
* create '.gatorconfig.json' in your home directory
* insert '{"db_url":"postgres://<dbuser>:postgres@localhost:5432/bloggator?sslmode=disable"}'
* replace <dbuser> with your postgres user name
* replace ip address and port number if your setup differs from this

## Dependencies
* run 'sqlc generate'
* run 'go install' in root directory
* run 'go build' to create executable 'bloggator'

# Run
## Register user
`bloggator register <username>`

`bloggator users`

## Login
`bloggator login <username>`

## Add RSS feed
`bloggator addfeed <name> <url>`

`bloggator feeds`

## Follow a feed
`bloggator follow <url>`

`bloggator unfollow <url>`

`bloggator following`

## Scrape posts from feeds you are following
(command runs indefinetely)

`bloggator agg <time_between_lookups>`

## Browse latest posts collected from feeds
`bloggator browse <number_of_posts>`
