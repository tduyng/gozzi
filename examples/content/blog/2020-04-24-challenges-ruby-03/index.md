+++
title = "Simple coding challenges with Ruby - 03"
description = "Dive into fun coding challenges in Ruby. Get practical tips and solutions to improve your coding skills."
template = "post.html"
date = 2020-04-24
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

My level on codewars now is 5kyu. You can check my [Codewars](https://www.codewars.com/users/tduyng). So, today when I take the new challenges with the option rank up, Codewars gave me only the 5kyu challenges.

Let's take a look together!

### 1. Rot13 (5kyu)

The first challenge today is still handle string. There are too many string algorithms. Maybe in the next day, I will choose the challenges by subject. For example, I will learn the OOP soon. So I think I will choose the katas in this subject.

Now, we will go to solve the challenge.

:bell: Task: [Rot13](https://www.codewars.com/kata/530e15517bc88ac656000716/train/ruby)

> ROT13 is a simple letter substitution cipher that replaces a letter with the letter 13 letters after it in the alphabet. ROT13 is an example of the Caesar cipher.

> Create a function that takes a string and returns the string ciphered with Rot13. If there are numbers or special characters included in the string, they should be returned as they are. Only letters from the latin/english alphabet should be shifted, like in the original Rot13 "implementation".

For example

```
"grfg" == rot13("test")
"Grfg" == rot13("Test")
```

---

:arrow_right: My solution:

- It becomes more complicated with the 5kyu level. I take me more time to reach the best solution what I can.
- To do this challenge

```ruby
def rot13(str)
  lower=*(97..122)
  upper=*(65..90)
  str_ascii = str.bytes.map do |num|
    if lower.include?(num)
      (num-13) >= 97 ? (num-13) : (num-13+122-96)
    elsif upper.include?(num)
      (num-13) >= 65 ? (num-13) : (num-13+90-64)
    else
      num
    end
  end
  return str_ascii.map(&:chr).join
end
```

---

:heavy_check_mark: Best solution on codewars

```ruby
def rot13(string)
  string.tr("A-Za-z", "N-ZA-Mn-za-m")
end
```

```ruby
def rot13(string)
  origin = ('a'..'z').to_a.join + ('A'..'Z').to_a.join
  cipher = ('a'..'z').to_a.rotate(13).join + ('A'..'Z').to_a.rotate(13).join
  string.tr(origin, cipher)
end
```

Oh la la, how smart are they!!!
I don't think the people can have a lot of super clever idea to resolve this challenge like that.

### 2. Greed is Good (5kyu)

:bell: Task: [Greed is Good](https://www.codewars.com/kata/5270d0d18625160ada0000e4/train/ruby)

> Greed is a dice game played with five six-sided dice. Your mission, should you choose to accept it, is to score a throw according to these rules. You will always be given an array with five six-sided dice values.

dice values

```
  Three 1's => 1000 points
  Three 6's =>  600 points
  Three 5's =>  500 points
  Three 4's =>  400 points
  Three 3's =>  300 points
  Three 2's =>  200 points
  One   1   =>  100 points
  One   5   =>   50 point
```

For example

```
Throw       Score
---------   ------------------
5 1 3 4 1   250:  50 (for the 5) + 2 * 100 (for the 1s)
1 1 1 3 1   1100: 1000 (for three 1s) + 100 (for the other 1)
2 4 4 5 4   450:  400 (for three 4s) + 50 (for the 5)
```

---

:arrow_right: My solution:

- I struggled to find an short and better solution for this challenge.
- Finally I did a quite long solution.
- So my idea is create an dictionary rules contains all case points of numbers. I convert it to string. For example: 31 means three 1's.
- Then I took some steps to convert the input array to an array with the same format of rules: new_array
- Loop through new_array, check with the rules to solve the sum

```ruby
def score(dice)
  rules = {"31"=>1000,
    "36"=>600,
    "35"=>500,
    "34"=>400,
    "33"=>300,
    "32"=>200,
    "11"=>100,
    "21"=>200,
    "41"=>1100,
    "51"=>1200,
    "15"=>50,
    "25"=>100,
    "45"=>550,
    "55"=>600
  }

  arr_uniq=dice.uniq

  # Create a hash {"1" => 0, "2" =>0 ...}
  check = arr_uniq.zip(Array.new(arr_uniq.size,0)).to_h

  # Count the number of dice array and set it in the check hash
  arr = dice.each do |num|
    if check.key?(num)
      check[num] += 1
    end
  end

  # Convert the hash to an array with the same format of rules
  new_arr =[]
  check.each {|k,v|
    new_arr.push(v.to_s+k.to_s)
  }

  # solve the sum
  sum = 0
  new_arr.each do |item|
    if rules.key?(item)
      sum += rules[item]
    end
  end
  return sum
end

```

---

:heavy_check_mark: Best solution on codewars

```ruby
SCORE_MAP = {
1 => [0, 100, 200, 1000, 1100, 1200, 2000],
2 => [0, 0, 0, 200, 200, 200, 400],
3 => [0, 0, 0, 300, 300, 300, 600],
4 => [0, 0, 0, 400, 400, 400, 800],
5 => [0, 50, 100, 500, 550, 600, 1000],
6 => [0, 0, 0, 600, 600, 600, 1200]
}

def score( dice )
  (1..6).inject(0) do |score, die|
    score += SCORE_MAP[die][dice.count(die)]
  end
end
```

```ruby
def score( dice )
[ dice.count(1) / 3 * 1000,
  dice.count(6) / 3 * 600,
  dice.count(5) / 3 * 500,
  dice.count(4) / 3 * 400,
  dice.count(3) / 3 * 300,
  dice.count(2) / 3 * 200,
  dice.count(1) % 3 * 100,
  dice.count(5) % 3 * 50 ].reduce(:+)
end
```

```ruby
def score( dice )
  m = {1 => 100, 5 => 50}
  (1..6).reduce(0) do|res, i|
    count = dice.count(i)
    res + count/3 * i * (i==1 ? 1000 : 100) + count%3*(m[i].to_i)
  end
end

```

I take more 2 solutions of coderwars for this challenge because each solution give an each different way to solve this interesting kata.

I love the first solution because it is the same idea as I thought. But they use an excelent method. The code is very elegent and readable.

### 3. String incrementer (5kyu)

:bell: Task: [String incrementer](https://www.codewars.com/kata/54a91a4883a7de5d7800009c/train/ruby)

> Your job is to write a function which increments a string, to create a new string. If the string already ends with a number, the number should be incremented by 1. If the string does not end with a number. the number 1 should be appended to the new string.

For example

```
foo -> foo1

foobar23 -> foobar24

foo0042 -> foo0043

foo9 -> foo10

foo099 -> foo100
```

---

:arrow_right: My solution:

- Return "1" if input string is empty
- Split array to two parts: only last numbers and the first part
- With the last numbers, we will increase one and keep their same format
- Join the first part and the new numbers and we are the result

```ruby
def increment_string(input)
  if input.empty?
    return "1"
  end
  arr = input.split(/(\d*$)/)
  num = (arr[1].to_i+1).to_s
  new_num = (arr.size > 1) && num.length < arr[1].length ? ("0"*(arr[1].length-num.length)+num) : num
  return arr[0]+new_num
end
```

---

:heavy_check_mark: Best solution on codewars

```ruby
def increment_string(input)
  input.sub(/\d*$/) { |n| n.empty? ? 1 : n.succ }
end
```

```ruby
def increment_string(input)
  input[/\d\z/] ? input.sub(/(\d+)\z/) { $1.next } : input + '1'
end
```

```ruby
def increment_string(input)
  input.gsub(/\d+/) { |match| match.succ } if Integer(input.chars.last) rescue input << "1"
end
```

You can see what the people can do, always one line of code. That's awesome!
With these solutions, I can learn more about a method ".succ" and ".next" and more about using the block code of the gsub like them.

## Reference

You can find the code in my github repository <i class="fab fa-github"></i> [Codewars Ruby](https://github.com/tduyng/coding-challenges/tree/master/codewars-ruby)

You can also check out my another repository that contains all the programs day by day I worked with The Hacking project <i class="fab fa-github"></i> [thehackingproject](https://github.com/tduyng/thehackingproject)

The Hacking Project website [thehackingproject.org](https://www.thehackingproject.org/)

Codeware website [coderwars.com](https://codewars.com)

## Series Ruby Coding Challenges

1. [Ruby Coding Challenges with TheHackingProject - 01](/blog/challenges-ruby-01)
2. [Ruby Coding Challenges with TheHackingProject - 02](/blog/challenges-ruby-02)
3. [Ruby Coding Challenges with TheHackingProject - 03](/blog/challenges-ruby-03)
