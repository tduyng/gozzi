+++
title = "Rails 5 or Rails 6 ?"
description = "Trying to Decide Between Rails 5 and Rails 6? Let's Compare!"
template = "post.html"
date = 2020-07-05
tags = ["ruby", "rails"]
generate_feed = false

[extra]
comment = false
reaction = true
toc = true
copy = true
img = "/img/post-cover.webp"
+++

What version of rails you use in your projects, today?

At [THP](https://www.thehackingproject.org/), when we start learning rails, we are required to install rails version 5.2.3. But I when I see that the new version of rails has been released is already 6.0. So we have posed a question for our technical manager of [THP](https://www.thehackingproject.org/):

Why we don't install directly the version rails 6.0?

And he said: "One word: `webpack`. I always don't understand so much for the "webpack", and I need to search more on the internet. And with this article (this is not mine), we will know some advantages and new features of rails 6 .

I hear that: Rails 5 is still just fine
Rails 5.0.Z is still on the severe patches list, and Rails 5.2.Z is still receiving regular security updates; so you're still OK on Rails 5 (especially 5.2 or higher)

Now We will see 4 amazing features of rails 6.

## Parralled Testing

Test’s performance is finally going to be improved (a lot!). Now you can use cores to your advantage of running big tests much faster. Each testing worker runs in its own thread - it should be reflected in the CPU monitor. Thanks to Eileen Uchitelle and Aaron Patterson, parallel-testing will land in the upcoming, final version of Rails framework.

As you probably already know, the way of implementing tests is very important. We try to implement as many suitable tests as we can. It’s great news for us as more tests will no longer have that much of an impact on the execution time. (Stay tuned for our guide on how to write great tests in Rails!)

“The default number [of workers running in parallel] is the actual core count on the machine you are on, but can be changed by the number passed to the parallelize method.” To enable executing tests parallel just set the number of workers in `rails_helper.rb`.

## Native Webpacker Support

Webpacker was introduced a while back with the 5.1 version of Ruby on Rails. It makes using JavaScript pre-processor easier.

We usually use it for JavaScript code, it works really well but it can also be used for CSS, images, fonts and assets as well. From now on, Sprockets is being replaced and Webpacker is the default JavaScript bundler for Rails through the new app/javascript directory.

Currently, in modern applications, using traditional rails views is not so often anymore. User Interfaces are very interactive, there is a lot of dynamic elements which have to respond really fast. Now, the application setup with Rails on the backend and React/Vue on the frontend will be an easy-peasy pleasure :)

## Multi Database Support

Another great feature of upcoming Rails 6.0 is the support for multiple simultaneous database connections. It’s a new, simple API for making that happen without the need to reach deeply into Active Record.

Official RoR blog suggests using two databases to for example split the workflow between two replica databases for a performance boost or records segmentation into databases for scaling. We can definitely see multiple database support as an improvement in building microservices architecture.

One of the real life examples where this feature could be really useful is our in-house project - Artinfo. In this project, there are two database connections, one for old users and one for those who just signed up. Now with Rails 6, this implementation would be much easier.

## Zeitwerk

With this version of Rails a new code loader was introduced - Zeitwerk. It promises to load your project’s classes and modules on demand, no need to write `require` calls. Zeitwerk uses absolute file names making the loader more efficient. Your classes and modules are available everywhere. The team behind the loader claims that it is thread-safe and matches Ruby’s semantics for constants.

Still curious for more? You can find more new stuff over at [rubyonrails.org](https://weblog.rubyonrails.org/2019/8/15/Rails-6-0-final-release/).

### So to reponse the question of article?

I don't know yet. All the new features of rails 6 is still so much compliated for me. At this moment, I will use the rails 5 for the project training and I will discover more about how to use the features of rails 6 and apply it in the new projects.

**But other developer said:**

If you are starting a new Rails project today - I would definitely start with Rails 6. It's officially supported now - and if you look at some of the CHANGELOGS above, you'll see that many of the subsystems have been stable for many months.

If you have a project on Rails 5.2 or higher - then you have some time before you have to do anything. If you're not interested in any of the major announced updates then you'll still be getting security updates for awhile.

You could still try running rails app:update though, just to see how much has changed for you particular app - which will give you an idea of how much work it will be when you do decide to upgrade.

If you have a project on Rails 5.0 or lower - then I would definitely consider upgrading sooner rather than later. It will take more work to upgrade these older projects (especially Rails 4 and below) - but security updates are already stopped for Rails < 4.2 - so now (right after a major release) is a good time to update (so that you don't have to do it again for awhile).
