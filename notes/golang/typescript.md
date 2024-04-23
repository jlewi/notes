## Nullables

I think you can `as` to force a nullable type to a non nullable e.g. something like

```
let stringOrNullable: string | null = "hello";
let notNullabel: string = stringOrNullable as string; 
```