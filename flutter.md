# Flutter

## Using GoogleAPIs

```
flutter pub add googleapis
flutter pub add googleapis_auth
```

Looks like there is one package for all apis.

## How to Handle unittesting of dart libraries

How would you run a unittest that ends up pulling in a library
that depends on `dart:js` or other packages not available on native?
Even if the actual code being test never invokes those functions
it will not be allowed to compile.

I think one pattern is to try to organize your code so that platform independent
code doesn't pull in platform dependent code.

Some threads (no clear answer)

* [SO Question](https://stackoverflow.com/questions/22153770/dart-unit-testing-classes-in-a-html-dependant-library)
*[Groups Chat 2014](https://groups.google.com/a/dartlang.org/g/misc/c/pacB66gnVcg)



### Initialization/Dependency Injection Pattern 

A pattern that seems to work well is to try to define interfaces that can then be implemented
by any classes that are platform specific. Code which consumes the interface only depends on the interface
definition and therefore doesn't pull in any platform specific libraries.

In `main.dart` you can initialize the actual implementation and then inject them into the code.

### Plugins and Packages

Linter suggests using [plugins or packages](https://docs.flutter.dev/development/packages-and-plugins/developing-packages). Need to look into that as a solution.

## Flutter/Dart and VSCode


This [doc](https://dartcode.org/docs/launch-configuration/) explains launch configurations
for dart.

### To run a specific dart test file

To run a specific dart test just set the `program` to the path
of the test e.g.

```
	{
      "name": "markdown-test",
      "request": "launch",
      "type": "dart",
      "program": "test/markdown_test.dart",
    },
```

To run a specific test in that file set the test_name

```
{
  "name": "openai_chat_test",
  "request": "launch",
  "type": "dart",
  "program": "test/openai_chat_test.dart",
  "args": [
    "--name=jq-test",
  ],
},
```

## Assertions in tests

You can use the [equality](https://stackoverflow.com/questions/71578026/how-do-i-compare-two-list-in-expect-function-of-flutter-unit-testing) 
function to compare lists for equality. But I'm not sure how you would print out a diff. 


## Table Driven Testing

Here's an example of what this looks like

```
import 'package:test/test.dart';

void main() {
  group('My function', () {
    // Define a list of test cases as tuples containing the inputs and expected outputs.
    final testCases = [
      ['input1', 'input2', 3],
      ['input3', 'input4', 7],
      ['input5', 'input6', 2],
    ];

    // Iterate over the list of test cases and run a separate test for each case.
    for (final testCase in testCases) {
      final input1 = testCase[0];
      final input2 = testCase[1];
      final expectedOutput = testCase[2];

      test('returns the correct output for inputs $input1 and $input2', () {
        final output = myFunction(input1, input2);
        expect(output, equals(expectedOutput));
      });
    }
  });
}

int myFunction(String input1, String input2) {
  // Implementation of the function being tested
  return input1.length + input2.length;
}
```