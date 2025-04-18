+++
title = "How to switch to another Rails version"
description = "When you need to use multiple versions of Rails on your computer, how do you switch between them?"
template = "post.html"
date = 2020-06-28
tags = ["ruby", "rails"]
generate_feed = false

[extra]
comment = false
reaction = true
toc = true
copy = true
img = "/img/post-cover.webp"
+++

When you need to use multiple versions of Rails on your computer, how do you switch between them?

## Installation rails

Before I talk about how to switch a version of rails, I will go quickly how to install the different versions of rails. The ruby package manager that I'am using is `rbenv` and I will use it for this tuto.

If you are using another ruby package manager like `rvm`, you should look at `rbenv` and its advantage. `rbenv` is more lighter and faster than the `rvm` and `rbenv` is very recommended by a a lot of developers.

There is an article on [dev.to](https://dev.to/) help you understand more clearly: [Why and How I Replaced RVM with RBENV](https://dev.to/krtb/why-and-how-i-replaced-rvm-with-rbenv-23ad).

Ok, now go on our article.

To install whatever version(s) of rails you want:

`gem install rails -v 6.0.3.2`

To know all rails versions released, you check out on [all rails versions](https://rubygems.org/gems/rails/versions).

You can also do it with your terminal with the command line `gem list rails --remote --all | grep "^rails "`.

I you want to install rails on another version of ruby, you need to switch version of ruby and re-install version of rails.

`rbenv install 2.5.1`

`rbenv global 2.5.1`

You can check list gems to see if you installed succesfullly:

`gem list --local` # result: rails(5.2.4.2,5.2.3)

now you tap `rails version` to see the current version of rails using, but it shows only

When you install multi version of rails on your machine, it will always take the lastest version of rails. So the question is how can create a new project rails with the version wanted?

## Create a project Ruby on Rails with a specific version of rails

I founded two ways to solve it:

- Using ruby package manager(rbenv)
- Using `gem bundler`

### Using ruby package manager (rbenv)

```bash
rbenv exec rails _<version_rails>_ new my_project
```

for example:

if i want to create a project with rails 5.2.4.2

`rbenv exec rails _5.2.4.2_ new my_project`

You can check the version of rails in the `Gemfile` when the new app has been created:

```ruby
ruby '2.7.1'
# Bundle edge Rails instead: gem 'rails', github: 'rails/rails'
gem 'rails', '~> 5.2.4', '>= 5.2.4.2' #version of rails using to create the project
```

### Using bundler

The second way is using `bundler`.

You just do like another gem dependencies. That means you will specify version of rails in your `Gemfile`.

for example:

```ruby
source 'https://rubygems.org'
gem 'rails', '6.0.3'
```

Very simple, isn't it?

Make sure you don't forget `bundle install` to install all the dependencies from your Gemfile.

And then, what we need to do, just finish with:

`rails new new_app`

So that is two ways I used. I hope it will be helpful for you.
