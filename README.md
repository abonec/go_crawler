# Crawler for go word

Simple crawl all web pages given in stdin and search for `go` word

# Features

* Parallel
* Have upper limit for number of workers
* Workers spawns only if necessary and shutdown after idle timeout
* Handling sigint (`ctrl+c`) for interrupt and print out the current result

# Usage

* `go build`
* `cat text_file_with_links.txt | ./go_crawler`
