+++
title = "Better In Ruby - 02"
description = "Hidden Ruby gems: uncovering little known tips."
template = "post.html"
date = 2020-08-30
tags = ["ruby"]
generate_feed = false

[extra]
comment = false
reaction = true
toc = true
copy = true
img = "/img/post-cover.webp"
+++

Some tips you probably don't know in Ruby - 02

In this article, I want to show you some nice Ruby features that you may know or not. Anyway, it’s a quick read and it’s always interesting to learn new stuff, right?!

## Create a hash from a list of values

<a style="padding:0px, margin:0px" name="assigning-rest-array" ></a>

You can create a hash from a list of values by using Hash[...]. It will create a hash like below:

```ruby
Hash['key1', 'value1', 'key2', 'value2']

# => {"key1"=>"value1", "key2"=>"value2"}
```

## Lambda Literal ->

<a style="padding:0px, margin:0px" name="word-array"></a>

Lambda Literal is a anonymous function like lambda in python, `=>` in JS or C#, allows you to create lambda easily.

```ruby
a = -> { 1 + 1 }
a.call
# => 2

a = -> (v) { v + 1 }
a.call(2)
# => 3
```

## Double star (\*\*)

<a style="padding:0px, margin:0px" name="concate-array"></a>

The double star is a neat little trick in Ruby. See the following method:

```ruby
def my_method(a, *b, **c)
  return a, b, c
end
```

`a` is a regular parameter. `*b` will take all the parameters passed after the first one and put them in an array. `c` will take any parameter given in the format key: value at the end of the method call.

See the following examples:

One parameter

```ruby
my_method(1)
# => [1, [], {}]
```

More than one parameter

```ruby
my_method(1, 2, 3, 4)
# => [1, [2, 3, 4], {}]
```

More than one parameter + hash-style parameters

```ruby
my_method(1, 2, 3, 4, a: 1, b: 2)
# => [1, [2, 3, 4], {:a=>1, :b=>2}]
```

## Handle single object and array in the same way

<a style="padding:0px, margin:0px" name="format-decimal"></a>

Sometimes you might want to give the option to either accept a single object or an array of objects. Instead of checking for the type of object you’ve received, you could use `[*something] or Array(something)`.

Let’s assign two variables. The first one is a single digit and the second one is an array of digits.

```ruby
stuff = 1
stuff_arr = [1, 2, 3]
```

In the following example, I use `[*...]` to loop through whatever is give

```ruby
[*stuff].each { |s| s }
[*stuff_arr].each { |s| s }
```

Same in this one but using Array(...).

```ruby
Array(stuff).each { |s| s }
Array(stuff_arr).each { |s| s }
```

## Double Pipe Equals ||=

<a style="padding:0px, margin:0px" name="remove-folder"></a>

The Double Pipe Equals is a great tool to write concise code.

It’s actually equivalent to the following:

```ruby
a || a = b # Correct
```

And not this one, as a lot of people think:

```ruby
a = a || b # Wrong
```

The second one doesn’t make sense because there is no point reassigning a if we already have it!

This operator can be used to create methods like this in your classes. I love to use it for calculations.

```ruby
def total
  @total ||= (1..100000000).to_a.inject(:+)
end
```

Now you could have other method calling total to get the total value but it will only be calculated the first time.

## Mandatory hash parameters

<a style="padding:0px, margin:0px" name="massive-assignment"></a>

This one was introduced in Ruby 2.0. Instead of just defining a method that takes a hash in parameters like this:

```ruby
def my_method({})
end
```

You can specify the keys that you are waiting for and even define default values for them! a and b are mandatory keys.

```ruby
def my_method(a:, b:, c: 'default')
  return a, b, c
end
```

We can try to call it without giving a value for b but it won’t work.

```ruby
my_method(a: 1)
# => ArgumentError: missing keyword: b
```

Since c has a default value, we can just call the method with a and b.

```ruby
my_method(a: 1, b: 2)
# => [1, 2, "default"]
```

Or with all of them.

```ruby
my_method(a: 1, b: 2, c: 3)
# => [1, 2, 3]
```

All we are doing is passing a hash and using some visual shortcuts but obviously, you can also pass a hash like this:

```ruby
hash = { a: 1, b: 2, c: 3 }
my_method(hash)
# => [1, 2, 3]
```

## Generate array of alphabet or numbers

<a style="padding:0px, margin:0px" name="deep-copy"></a>

You might want to generate a list of numbers or put the entire alphabet inside an array. Well, you can use ruby ranges to do this.

A to Z

```ruby
('a'..'z').to_a
# => ["a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"]
```

1 to 10

```ruby
(1..10).to_a
# => [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
```

## Tap

<a style="padding:0px, margin:0px" name="random-array"></a>

Tap is a nice little method that improves code readability. Let’s take the following class as an example.

```ruby
class User
  attr_accessor :a, :b, :c
end
```

Now let’s say you want to instantiate a new user and assign a value to each of its attributes. You could do it like this:

```ruby
def my_method
  o = User.new
  o.a = 1
  o.b = 2
  o.c = 3
  o
end

```

Or you could use tap to do it like this.

```ruby
def my_method
  User.new.tap do |o|
    o.a = 1
    o.b = 2
    o.c = 3
  end
end

```

Basically, the tap method yields the calling object to the block and returns it.

---

Thanks you for reading! Stay tunned to update the next chapter.
