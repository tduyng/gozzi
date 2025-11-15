# What is Gozzi?

**Gozzi** is a static site generator written in Go. I built it to learn how static site generators work, and now I use it for my personal website [tduyng.com](https://tduyng.com).

## Why I Built This

I wanted to understand how tools like Hugo and Zola work under the hood. Instead of just reading about them, I built one. It started as a learning project and evolved into something I actually use daily.

## What Can You Do With It?

- Build a personal blog
- Create a documentation site
- Make a portfolio website
- Organize notes and learning

## Is Gozzi For You?

Gozzi might be a good fit if you:

- **Want to learn Go** by using a real tool
- **Like simplicity** over feature overload
- **Prefer clear configuration** in TOML
- **Enjoy understanding** how your tools work

Gozzi is **not trying to replace Hugo or Zola**. Those are excellent, mature tools. Use them if they fit your needs!

## Core Features

- **Fast builds** - Sub-second for most sites
- **Live reload** - See changes instantly
- **Markdown content** - Write in markdown, organize in folders
- **Go templates** - Use Go's HTML templates
- **Built-in features** - Tags, pagination, RSS, sitemap

## How It Works

1. Write content in Markdown files
2. Create HTML templates with Go
3. Configure your site in `config.toml`
4. Run `gozzi build` to generate your site
5. Deploy the `public/` folder anywhere

## Philosophy

- **Simple > Complex** - Clear configuration over magic
- **Explicit > Implicit** - You control what happens
- **Learning first** - Code is readable, concepts are clear

## What's Next?

Ready to try it out?

- [Installation](/guide/installation) - Get Gozzi on your machine
- [Getting Started](/guide/getting-started) - Build your first site in 5 minutes
- [Real Example](/examples/real-world) - See how I use it for tduyng.com

Or browse the guides to learn more about how it works.
