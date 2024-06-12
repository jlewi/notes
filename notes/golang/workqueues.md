# workqueues

* K8s [workqueues package](https://pkg.go.dev/k8s.io/client-go/util/workqueue) is very useful
  for implementing rate limiting and delaying queues.

* **pitfall* The queue items should be concrete values not pointers.
  * The items you enqueue are used as the keys in a map
  * So if you use pointers; then the key ends up being the memory address not the value
  * As a result if you have two different pointers whose value is the same then they will not be treated as the same
    key and you will have lots of problems.