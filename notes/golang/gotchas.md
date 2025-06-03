# Shadowing

* Variable in the inner scobe hides the variable in the outer scope
* Here's an example

```sh
offset := 0

for ;; {
  read, offset, err := readFromOffset(offset)

   if read == 0 {
     break
   }
}
```

* Here offset declared in the inner scope is shadowing the value in the outer scope
* So the intent here is each time readFromOffset is called it should read from the last offset
* But since offset is shadowed, it will always read from 0

```sh
echo "hello"
```