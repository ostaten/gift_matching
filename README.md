# Gift Matching
Given a set of people who are to give each other gifts in a "Secret Santa" style, provides potential matches based on avoiding giving to people recently given to, and not giving to people within their own subgroup.

## Basic use case:
You're arranging a Secret Santa between several different families.  Each person is to give to one other person.  You do not want anyone giving to someone in their own family, and you want to give to someone you haven't given to recently.  

### Uses: GoLang

## What you need to do:
1) Enter your people in the `groups.json` file just like the example provided.  The subgroups are a way to ensure that certain people from certain groups cannot be matched with each other.  If you have no such restrictions, just create one giant subgroup.
2) If you have any restrictions, etc, you need to enter those in the `history.json` file in the format shown in the example file.  Each set of restrictions indicates an interation of gift giving.  Otherwise <b>remove example data while keeping base structure</b>.
3) Make sure you have GoLang installed, then run `go run run_simulation.go` from terminal inside base directory.
