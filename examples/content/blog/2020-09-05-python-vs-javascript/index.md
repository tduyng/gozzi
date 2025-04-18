+++
title = "Syntax Python vs Javascript"
description = "Discover the key differences and unique syntax nuances between Python and JavaScript in this engaging comparison."
template = "post.html"
date = 2020-09-05
tags = ["python", "javascript"]
generate_feed = false

[extra]
comment = false
reaction = true
toc = true
copy = true
img = "/img/post-cover.webp"
+++

Programming in a new language can feel unfamiliar, intimidating, and tricky to navigate at first. There are new rules, syntax, and nuances to learn.

# Python vs. JavaScript Language Basics

JavaScript and Python are interpreted programming languages, meaning their runtime environments use an interpreter (or engine) that parses and executes code one statement at a time.

The two languages are also “object-based” — everything is (or can be treated as) an object: strings, numbers, data structures, functions, etc.

## Primitive Types

First up, JavaScript and Python have similar built-in data types. For example, both use numeric data types (integers and floats), strings and Booleans.

```js
// JavaScript data types
const pi = 3.14
const age = 31
const greeting = 'good morning'
const isAdmin = true
```

```py
# Python data types
pi = 3.14
age = 13
greeting = "good morning"
is_admin = True
```

## Type Checking and Conversion

Python and JavaScript are “dynamically typed” languages, which means you do not have to set the type of a variable explicitly. The data type is set when you assign a value to a variable.

In JavaScript, you use the typeof operator to verify the data type of a variable. Python provides a similar built-in function, type().

```js
// JavaScript
const greeting = 'good morning'
typeof greeting // "string"
```

```py
# Python
pi = 3.14
type(pi) # float
```

You can convert from one type to another, like a string to a number, in Python with the int() and float() functions:

```py
# Python
input = input('Enter a number: ') #'2'

# convert string to int
int(input) # 2
```

JavaScript includes the methods parseInt() and parseFloat() for the same purpose:

```js
// JavaScript
const input = prompt('Enter a number:') // '4'
parseInt(input) // 4
```

## Built-in String Methods

To convert cased characters in a string from uppercase to lowercase (and the reverse), use Python’s upper() and lower() functions:

```py
# Python
user_name = input('What is your name? ') # GUIL
user_name.lower() # 'guil'

greeting = "good evening"
greeting.upper() # 'GOOD EVENING'
```

JavaScript supplies the toUpperCase() and toLowerCase() methods to convert strings:

```js
// JavaScript
const greeting = 'good evening'
const userName = prompt('What is your name?') // "gUiL"

greeting.toUpperCase() // "GOOD EVENING"
userName.toLowerCase() // "guil"
```

**String Interpolation**
Template literals in JavaScript let you replace ${} placeholders with values inside of a string literal. This process is called string interpolation:

```js
// JavaScript
const greeting = 'Good evening'
const name = 'Guil'
console.log(`${greeting}, ${name}!`) // Good evening, Guil!
```

The Python string format() method inserts values into a template string containing {} replacement fields. You pass the method the values to interpolate. For example:

```py
# Python strings
greeting = "Good evening"
name = "Guil"
print( "{}, {}!".format(greeting, name) ) # Good evening, Guil!
```

Each set of curly braces gets replaced with the values passed to format() in sequential order.

Python’s formatted string literal (f-String) offers a more concise syntax to accomplish the same. It looks like a regular string that’s prepended by the character f, and you include the value to interpolate directly inside the string.

```py
# Python strings
greeting = "Good evening"
name = "Guil"
print(f"{greeting}, {name}!") # Good evening, Guil!
```

# Python vs. JavaScript Data Structures

JavaScript and Python give you comparable structures to store and organize your data.

## Arrays and Lists

Like a JavaScript array, a Python list stores a collection of values in a single container. The values can be different data types like strings, integers, Booleans, etc.

```py
# Python list
students = ['Lee', 'Toni', 'Marie', 'Agata']
# return length of list
len(students) # 4
students[2] # 'Marie'
```

```js
// JavaScript array
const students = ['Lee', 'Toni', 'Marie', 'Agata']
students.length // 4
students[0] // Lee
```

Notice how both languages have similar ways of returning the length of a list and retrieving a value by index.

## Array and List Methods

Since arrays and lists are considered objects in their respective language, there are various properties and methods you can use on them. For example, a common way to add elements to the end of an array in JavaScript is with the push() method:

```js
// JavaScript
const instruments = ['piano', 'drums', 'trumpet']
instruments.push('guitar')
// ['piano', 'drums', 'trumpet', 'guitar']
```

You add an item to the end of a Python list with the append() method:

```py
# Python
instruments = ['piano', 'drums', 'trumpet']
instruments.append('guitar') # ['piano', 'drums', 'trumpet', 'guitar']
```

JavaScript arrays and Python lists have a pop() method for removing and returning items. Calling pop() on a JavaScript array removes the last element:

```js
// JavaScript arrays
const instruments = ['piano', 'drums', 'trumpet']
instruments.pop()
// "trumpet"
```

While JavaScript’s pop() method does not accept arguments, Python’s does! For instance, pass pop() the index of the item you want to remove:

```py
# Python lists
instruments = ['piano', 'drums', 'trumpet']
instruments.pop(1) # 'drums'
```

## Spreading and Unpacking

In JavaScript, you use the spread operator (...) to copy, combine and manipulate arrays:

```js
// JavaScript arrays
const studentsA = ['Lee', 'Toni', 'Marie']
const studentsB = ['Meg', 'Jesse', 'Anwar']
const students = [...studentsA, ...studentsB]
// ['Lee', 'Toni', 'Marie', 'Meg', 'Jesse', 'Anwar']
```

In a similar way, you can “unpack” items from one Python list into another using an asterisks (\*):

```py
# Python lists
students_a = ['Lee', 'Toni', 'Marie']
students_b = ['Meg', 'Jesse', 'Anwar']

students = [*students_a, *students_b]
 # ['Lee', 'Toni', 'Marie', 'Meg', 'Jesse', 'Anwar']
```

Both create a copy of a list or array, preserving the original values.

One handy aspect of the JavaScript spread operator is that you can pass arrays as arguments to functions.

```js
// JavaScript
const numbers = [10, 20, 30, 40]
Math.max(...numbers) // 40
```

Likewise, you have the ability to unpack items in a Python list for function calls:

```py
# Python
numbers = [1, 10]
list( range(*numbers) ) # [1, 2, 3, 4, 5, 6, 7, 8, 9]
```

Both expand an array or list into separate arguments.

## Objects and Dictionaries

If you’re familiar with JavaScript objects, you’ll recognize Python dictionaries. You write both using curly brackets holding related data in the form of key/value pairs.

```js
// JavaScript object
const pet = {
    name: 'Joey',
    type: 'Dog',
    breed: 'Australian Shepherd',
    age: 8,
}
```

```py
# Python dictionary
pet = {
  'name': 'Joey',
  'type': 'Dog',
  'breed': 'Australian Shepherd',
  'age': 8
}
```

Lists and dictionaries, like arrays and objects are mutable, which means that you can change the data inside them without changing their identity. Once you create an object, its type and identity (or the address in memory it’s pointing to) does not change.

**Copying/Merging Objects and Dictionaries**

JavaScript’s spread operator copies key/value pairs from one object literal to another. It’s comparable to the double asterisks (\*\*) operator in Python, which copies and merges dictionaries:

```js
// JavaScript objects
const name = {
    firstName: 'Reggie',
    lastName: 'Williams',
}

const developer = {
    ...name, // place the 'name' key/values here
    title: 'Software developer',
    skills: ['JavaScript', 'HTML', 'CSS'],
}
```

```py
# Python dictionaries
name = {
  'firstName': 'Reggie',
  'lastName': 'Williams'
}

developer = {
  **name, # place the 'name' key/values here
  'title': 'Software developer',
  'skills': ['JavaScript', 'HTML', 'CSS']
}
```

## Functions

Both languages take full advantage of functions for code reuse. Python uses the def keyword compared to function in JavaScript.

```py
# Python function
def add(a, b=10):
    val a + b
    return val
```

```js
// JavaScript function
function add(a, b = 10) {
    const val = a + b
    return val
}
```

Notice how both use the return keyword to return a value, and you’re able to specify default parameters in each function definition.

## Single Line Functions

Arrow functions in JavaScript offer a concise syntax for creating functions. More so, if your function body is only one line of code, you can omit the return keyword and place everything on one line:

```js
// JavaScript arrow function
const add = (a, b) => a + b
```

In Python, the lambda keyword provides syntactic sugar for defining functions as single-line expressions:

```py
# Python lambda
add = lambda a,b : a + b
```

These single-line functions are common when you want to pass an anonymous function as an argument to another (higher-order) function. For example, they’re used with the built-in Python and JavaScript functions map(), filter(), and reduce().

```js
// JavaScript arrow function
const states = ['ca', 'fl', 'hi', 'ny']
states.map((s) => s.toUpperCase())
// ["CA", "FL", "HI", "NY"]
```

```py
# Python lambda
states = ['ca', 'fl', 'hi', 'ny']
list( map(lambda s: s.upper(), states) )
 # ['CA', 'FL', 'HI', 'NY']
```

## Conditional Statements

Python’s flow control statements also look and work similarly to the if/else you know from JavaScript:

```js
// JavaScript conditional
let score = 4

if (score === 5) {
    console.log('Gold Medal!')
} else if (score >= 3) {
    console.log('Silver Medal')
} else if (score >= 1) {
    console.log('Bronze Medal')
} else {
    console.log('No Medal :(')
}
```

The most significant difference besides the absence of curly braces and parentheses around the condition is the elif clause, which is short for “else if”.

```py
# Python conditional
score = 5

if score == 5:
    print("Gold Medal!")
elif score >= 3:
    print("Silver Medal")
elif score >= 1:
    print("Bronze Medal")
else:
    print("No Medal :(")
```

Like JavaScript’s else if clause, you can specify any number of elif clauses, and the optional else clause should appear last.

## Loops and Iteration

Lastly, Python has a while loop, which looks and works almost the same as its JavaScript counterpart:

```
# Python
password = input("Enter the secret password: ")
while password != 'sesame':
    password = input("Invalid password. Try again: ")
```

```
// JavaScript
let password = prompt("Enter the secret password:");
while (password !== 'sesame') {
  password = prompt("Invalid password. Try again: ");
}
```

Data types like strings, lists, and dictionaries are also iterable objects in Python; you use a for loop to iterate over them:

```js
# Python for loop
students = ['Lee', 'Toni', 'Marie', 'Jesse', 'Anwar']
for student in students:
    print(student)
```

The above loop seems more elegant as opposed to JavaScript’s verbose for loop. It’s comparable to the for...of loop introduced in ES2015.

```js
// JavaScript for...of loop
const students = ['Lee', 'Toni', 'Marie', 'Jesse', 'Anwar']
for (let student of students) {
    console.log(student)
}
```

You also use the break keyword in either to exit (or break out of) a while and for loop.

```py
# Python
scores = [50, 20, 30, 0, 10, 15, 35]
for score in scores:
    print(f"Score: {score}");
    if score == 0:
        print("You may not continue if you have a 0 score.")
        break
```

```js
// JavaScript
while (true) {
    let response = prompt("Type 'exit' to make this stop.")
    if (response === 'exit') {
        break
    }
}
```

---

Happy learning!
