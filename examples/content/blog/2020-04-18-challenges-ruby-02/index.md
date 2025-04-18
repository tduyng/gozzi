+++
title = "Simple coding challenges with Ruby - 02"
description = "Dive into fun coding challenges in Ruby. Get practical tips and solutions to improve your coding skills."
template = "post.html"
date = 2020-04-18
tags = ["ruby", "algorithms"]
generate_feed = false

[extra]
comment = false
reaction = true
toc = true
copy = true
img = "/img/post-cover.webp"
outdate_alert = true
outdate_alert_days = 365
+++

Doing coding challenges is an excellent way to imporve your programming language & problem-solvings skills.

<sub>_I'am learning and trying to improve my English. So I'm really sorry if my grammar mistakes annoy you. Thank you for your understanding._</sub>

## Kata challenges

In [the last article](/blog/challenges-ruby-01/), we have done the very simple challenges. Now we will challenge in the next level with the 6kyu challenges.

### 1. Replace with alphabet position (6kyu)

All the next challenges have distributed random by codewars.

:bell: Task: [Replace with alphabet position](https://www.codewars.com/kata/546f922b54af40e1e90001da/train/ruby)

> In this kata you are required to, given a string, replace every letter with its position in the alphabet. If anything in the text isn't a letter, ignore it and don't return it.
> "a" = 1, "b" = 2, etc

For example

```
alphabet_position("The sunset sets at twelve o' clock.")
```

Should return

```
Should return "20 8 5 19 21 14 19 5 20 19 5 20 19 1 20 20 23 5 12 22 5 15 3 12 15 3 11" (as a string)
```

---

:arrow_right: My solution:

- This challenge is always handle with the string
- One thing that we need pay attention here is how to convert a character to the ascii caracter.

{% highlight ruby %}
def alphabet_position(text)
text=text.downcase.gsub(/[^a-z]/,'')
text.bytes.map{|x| x-96}.join(' ')
end
{% endhighlight %}

---

:heavy_check_mark: Best solution on codewars

{% highlight ruby %}
def alphabet_position(text)
text.gsub(/[^a-z]/i, '').chars.map{ |c| c.downcase.ord - 96 }.join(' ')
end
{% endhighlight %}

This challenge is still not difficult for you, is it?
We will see the next one.

### 2. Array.diff (6kyu)

:bell: Task: [Replace with alphabet position](https://www.codewars.com/kata/546f922b54af40e1e90001da/train/ruby)

> Your goal in this kata is to implement a difference function, which subtracts one list from another and returns the result. It should remove all values from list a, which are present in list b.

For example

```
array_diff([1,2],[1]) == [2]
array_diff([1,2],[1]) == [2]
```

---

:arrow_right: My solution:

- Very simple here if we know already the set operation method of array in ruby

```
union = A | B
isect = A & B
diff1  = A - B
diff2  = B - A
sdiff = (A - B) | (B - A)
```

```ruby
def array_diff(a, b)
  return (a - b)
end
```

Wow Ruby !!!

### 3. Replace with alphabet position (6kyu)

The last challenge today.

:bell: Task: [Replace with alphabet position](https://www.codewars.com/kata/546f922b54af40e1e90001da/train/ruby)

> Write a function toWeirdCase (weirdcase in Ruby) that accepts a string, and returns the same string with all even indexed characters in each word upper cased, and all odd indexed characters in each word lower cased. The indexing just explained is zero based, so the zero-ith index is even, therefore that character should be upper cased.

> The passed in string will only consist of alphabetical characters and spaces(' '). Spaces will only be present if there are multiple words. Words will be separated by a single space(' ').

For example

```
weirdcase( "String" )#=> returns "StRiNg"
weirdcase( "Weird string case" );#=> returns "WeIrD StRiNg CaSe"
```

---

:arrow_right: My solution:

- The key to resolve this kata is just detect the even index and odd index of caracter.
- The method "each_with_index" will help us to do it.
- I use ".strip" to remove the space caracter in the last string

```ruby
def weirdcase string
  result =""
  arr = string.split(' ').each do |word|
    result  += word.chars.each_with_index.map{|x,i| (i.even?) ? x.upcase: x.downcase}.join('') + " "
  end
  return result.strip
end
```

---

:heavy_check_mark: Best solution on codewars

```ruby
def weirdcase string
  string.gsub(/(\w{1,2})/) { |s| $1.capitalize }
end
```

You can see what the people can do, always one line of code. That's awesome!
I need to learn more about using the block code of the gsub like them.

This article is very short but I need to finish this article here. I don't have so much time to do more the others challenges today. See you in the next chapter.

## Reference

You can find the code in my github repository <i class="fab fa-github"></i> [Codewars Ruby](https://github.com/tduyng/coding-challenges/tree/master/codewars-ruby)

You can also check out my another repository that contains all the programs day by day I worked with The Hacking project <i class="fab fa-github"></i> [thehackingproject](https://github.com/tduyng/thehackingproject)

The Hacking Project website [thehackingproject.org](https://www.thehackingproject.org/)

Codeware website [coderwars.com](https://codewars.com)

## Series Ruby Coding Challenges

1. [Ruby Coding Challenges with TheHackingProject - 01](/blog/challenges-ruby-01)
2. [Ruby Coding Challenges with TheHackingProject - 02](/blog/challenges-ruby-02)
3. [Ruby Coding Challenges with TheHackingProject - 03](/blog/challenges-ruby-03)

**I hope it make more sense or at least it is helpful for you.**
