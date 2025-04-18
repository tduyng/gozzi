+++
title = "Better In Ruby - 01"
description = "Hidden Ruby gems: uncovering little known tips."
template = "post.html"
date = 2020-08-01
tags = ["ruby"]
generate_feed = false

[extra]
comment = false
reaction = true
toc = true
copy = true
img = "/img/post-cover.webp"
+++

Some tips you probably don't know in Ruby.

Today, I will make a series that I learned about tips and tricks in ruby.

## Assigning the rest of an array to a variable

<a style="padding:0px, margin:0px" name="assigning-rest-array" ></a>

When destructuring an array, you can unpack and assign the remaining part of it to a variable using the rest pattern

```ruby
  a = [1, 2, 3]
  b, *rest = *a

  b    # => 1
  rest # => [2, 3]
  a    # => [1, 2, 3]
```

## Word array

<a style="padding:0px, margin:0px" name="word-array"></a>

When we want to add a separator in the string

```ruby
%w{This is a string, I want to separate by a comma} * ", "
# "This, is, a, string,, I, want, to, separate, by, a, comma"
```

## Concate array

<a style="padding:0px, margin:0px" name="concate-array"></a>

```ruby
[1, 2, 3] * 3 == [1, 2, 3, 1, 2, 3, 1, 2, 3] # true
new_array = Array.new([0], 5) # [0,0,0,0,0]
```

## Format decimal

<a style="padding:0px, margin:0px" name="format-decimal"></a>

```ruby
number = 9.5
"%.2f" % number # => "9.50"
number.round(2)
```

## Remove a folder

<a style="padding:0px, margin:0px" name="remove-folder"></a>

This is a relatively common job of developer. There a many different ways to delete a folder and this is one of the shortest and fastest way to do:

```ruby
require 'fileutils'
FileUtils.rm_r 'somedir'
```

## Massive assignment

<a style="padding:0px, margin:0px" name="massive-assignment"></a>

Massive assignment allow us to declare many variables at the same time

```ruby
a,b,c,d = 1,2,3,4
```

This feature will be very useful in the method with many parameters

```ruby
def my_method(*args)
  a,b,c,d = args
end
```

## Deep copy

<a style="padding:0px, margin:0px" name="deep-copy"></a>

When we copy an object that contains others objects inside, array for example, you only copy references to those objects

```ruby
food = %w( bread milk orange )
food.map(&:object_id)
=> [69612840, 69612820, 69612800]
[122] pry(main)> food.clone.map(&:object_id)
=> [69612840, 69612820, 69612800]
```

Using <code>Marshal</code> (usually use for [serialization](https://en.wikipedia.org/wiki/Serialization)), you can create a deep copy.

```ruby
def deep_copy(obj)
  Marshal.load(Marshal.dump(obj))
end
```

See the result:

```ruby
ObjectSpace.each_object(String).select{|x| x == food.first}.count
=> 54
deep_copy(food).map &:object_id
=> [77564660, 77564420, 77564300]
ObjectSpace.each_object(String).select{|x| x == food.first}.count
=> 56

```

**Bonus**: Difference between <code>clone</code> and <code>dup</code>(They are almost the same, but there are two points that **dup** doesn't have)

- Clone also singleton class of copied object
- Keep the frozen state of copied object

For example:

- Singleton methods

    ```ruby
    a = Object.new
    def a.foo; :foo end
    p a.foo
    => :foo
    b = a.dup
    p b.foo
    => undefined method `foo' for #<Object:0x007f8bc395ff00> (NoMethodError)
    c = a.clone
    c.foo
    => :foo
    ```

- Frozen state

    ```ruby
    a = Object.new
    a.freeze
    p a.frozen?
    => true
    b = a.dup
    p b.frozen?
    => false
    c = a.clone
    p c.frozen?
    => true
    ```

## Create a an array with random values

<a style="padding:0px, margin:0px" name="random-array"></a>

```ruby
Array.new(10) { rand 300 }
```

It will create an array with 10 item that have the random values from 0 to 299

---

Thanks you for reading! Stay tunned to update the next chapter.
