# Generics

* Reciever methods can't have genercis [reference](https://blog.streamelements.com/an-introduction-to-generics-in-go-cc8cdae15ef2)

## Using generic to instantiate concrete classes

Can we use generics to instantiate a concrete version of an interface?

e.g. suppose we have an interface `DBDoc` and different types that implement it
e.g. `Doc1` and `Doc2`.

Now suppose we want to define functions like the following 


```
func NewDoc1() DBDoc {
    return &Doc1{}
}

func NewDoc2() DBDoc {
    return &Doc2{}
}
```

We'd like to define that using a generic the following works

```
func NewDBDoc[Doc any]() *Doc {
	d := new(Doc)
	return d
}
```

But it doesn't define a typeconstraint. As a result suppose we have the following

```
type DocCreator fun() DBDoc

func List(creator DocCreator) {
    ....
}
```

The following won't compile

```
List(NewDBDoc[NewDoc1])
```

because the return type of NewDBDoc is `*Doc1`.

The following won't compile

```
func NewDBDocConstrained[Doc DBDoc]() DBDoc{
	d := new(Doc)
	return d
}
```
* I think the problem is that suppose *Doc1 implements DBDoc
* Then Doc is type *Doc1; so new(Doc) creates **Doc1 which doesn't implement DBDoc
* Put another way we'd have to do `NewDBDocConstrained[*Doc1]`
* If Doc is a pointer type how do we instantiate a new version of it?
* NewDBDoc works because the type constraint is any so we can do `NewDBDoc[Doc1]`