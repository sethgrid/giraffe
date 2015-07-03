# giraffe
Giraffe is an in-memory directional graph and key value store. Giraffe is a work in progress. I want to add the ability to persist and load graphs, to later make giraffe an actual datastore with an exposed tcp or rest api, and eventually make it distributed.

You can check out the full documentation on [godoc.org](https://godoc.org/github.com/sethgrid/giraffe)

# What can giraffe do?
The current API supports:
- creating a graph object (with or without constraints like preventing duplicate keys or circular relationships)
- adding / deleting nodes
- adding / removing relationships between nodes
- assigning a node a key and value
- searching nodes
- creating an HTML/Javascript view of the graph leveraging visjs.org
- safe to use concurrently

![](/giraffe/giraffe.png)
