# Time Tracker

This is a fun little command line tool to track what you've done through out the day.

## Usage
```
$ timetrack --help
Usage of timetrack:
  -print
    	Prints out records in the database
  -since string
    	Limits by date the records to print [only applies when -print is enabled] (default "May-29")
```

## Here's how it works

### Enter things you do
```
$ timetrack 
5m20s since your last entry
Enter a task description:

> Wrote a timetracking app

Saved /Users/natebosscher/.timetrack.json
$ timetrack 
12m20s since your last entry
Enter a task description:

> Deployed my cool web app

Saved /Users/natebosscher/.timetrack.json
```

### Report on what you did
```
$ timetrack --print
+--------------+----------+--------------------------------+
|     TIME     | DURATION |          DESCRIPTION           |
+--------------+----------+--------------------------------+
| May-29 07:00 |        - | Good morning                   |
| May-29 07:05 |  5m20s   | Wrote a timetracking app       |
| May-29 07:17 | 12m20s   | Deployed my cool web app       |
+--------------+----------+--------------------------------+
```
