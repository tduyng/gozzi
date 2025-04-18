+++
title = "Dynamic Github profile with Bun and Typescript"
description = "Learn how to make your GitHub profile dynamic using Bun and TypeScript as an alternative to Python in my previous article."
template = "post.html"
date = 2024-05-15
tags = ["github", "readme", "typescript"]
generate_feed = true

[extra]
comment = true
reaction = true
toc = true
copy = true
show_ended_words = true
featured = true
img = "img/tduyng.webp"
+++

[In a recent article](https://tduyng.github.io/blog/dynamic-github-profile-readme/), I discussed transforming a GitHub profile into an interactive space with icons, badges, and dynamically updated blog posts using Python scripts and GitHub actions.

After sharing my work in the "r/javascript" community on Reddit, I received serveral positive feedback and upvotes. I chose to post in this community due to my familiarity with JavaScript and TypeScript. After several hours, moderators removed my post because it wasn't directly related to JavaScript. They made the right decision, and I appreciate their careful moderation of the community. This made me rethink my approach and think about using a language that would better connect with the audience, given the simplicity of the script's principle: fetching feed URL, parsing it, and writing the response to the README.md file, along with automatic updates using a cron job on GitHub action.

I decided to explore using JavaScript or TypeScript for this idea. While JavaScript would have been sufficient, I opted for TypeScript because of its enhanced safety and the overall enjoyment it offers during development—for me, it's more fun to work with than JavaScript.

Choosing a runtime was the next consideration. With options like [Node.js](https://nodejs.org/en), [Deno](https://deno.com/), and [Bun](https://bun.sh/) available, I choose on Bun. Because the runtime mainly executes the script, Bun seemed like the perfect choice. It's fast, simple, and allows running JavaScript or TypeScript without needing compilation steps.

**Let's set up and implement the code to see how it works.**

Firstly, ensure that Bun is installed on your machine. There are several ways to install it, all of which are outlined on the Bun website's [Installation](https://bun.sh/docs/installation) page.

Bun undergoes improvements regularly and changes rapidly. You can ensure you're using the latest version of Bun by executing the command: `bun upgrade`.

In the project, you can initialize Bun by using the command: `bun init`, which will generate some files for you. Here is an example of the folder structure after running `bun init`:

```json
❯ tree
├── bun.lockb
├── node_modules
├── package.json
├── README.md
├── src
│  └── index.ts
└── tsconfig.json
```

In my case, I prefer to add it manually for more control.

As I'm using TypeScript, I need to add `@types/bun` to work with TypeScript typings. Run the command:

```bash
bun add -d @types/bun
```

The `tsconfig.json` file is also necessary for TypeScript. You can either use the recommended `tsconfig.json` provided by Bun or generate one. It's minimal enough and follows best practices. Here is my `tsconfig.json`:

```json
{
    "compilerOptions": {
        "lib": ["ESNext"],
        "target": "ESNext",
        "module": "ESNext",
        "moduleDetection": "force",
        "allowJs": true,
        "moduleResolution": "bundler",
        "allowImportingTsExtensions": true,
        "verbatimModuleSyntax": true,
        "noEmit": true,
        "strict": true,
        "skipLibCheck": true,
        "noFallthroughCasesInSwitch": true,
        "noUnusedLocals": true,
        "noUnusedParameters": true,
        "noPropertyAccessFromIndexSignature": true
    },
    "include": ["src/**/*.ts"]
}
```

I specify `"include": ["src/**/*.ts"]` to tell TypeScript to compile all `.ts` files within the `src` folder, simplifying TypeScript's server checking and compilation processes.

Next, I'll use the **[rss-parser](https://www.npmjs.com/package/rss-parser)** library from npm to parse the feed `atom.xml` from my blog. The principle of this library is similar to the one I used in Python, called [feedparser](https://pypi.org/project/feedparser/).

```bash
bun add rss-parser
```

Here's how my `package.json` looks:

```json
{
    "type": "module",
    "scripts": {
        "readme": "bun src/feed.ts"
    },
    "dependencies": {
        "rss-parser": "^3.13.0"
    },
    "devDependencies": {
        "@types/bun": "^1.1.2"
    }
}
```

I use `"type": "module"` to write the code in [ESM](https://developer.mozilla.org/en-US/docs/Glossary/ECMAScript).

The script `"readme": "bun src/feed.ts"` runs the TypeScript script I will write later.

Don’t forget to add `node_modules` to `.gitignore`."

**Now it's time to write the script `feed.ts`.**

I'll use `rss-parser` to fetch the feeds from `atom.xml.`

```tsx
import { promises as fs } from 'fs'
import rssParser from 'rss-parser'

const DEFAULT_N = 5

type Entry = {
    title?: string
    link?: string
    isoDate?: string
}

/**
 * Fetches and parses the feed from the provided URL.
 * @param {string} url - The URL of the feed.
 * @returns {Promise<string[]>} An array of formatted feed entries.
 */
const fetchFeed = async (url: string): Promise<string[]> => {
    try {
        const parser = new rssParser()
        const response = await parser.parseURL(url)
        let feeds = []

        for (const item of response.items) {
            if (item.title && item.link) feeds.push(formatFeedEntry(item))
            if (feeds.length === DEFAULT_N) break
        }

        return feeds
    } catch (error) {
        console.error('Error fetching or parsing the feed:', error)
        return []
    }
}

/**
 * Formats a feed entry into a string.
 * @param {Entry} entry - The feed entry to format.
 * @returns {string} The formatted feed entry.
 */
const formatFeedEntry = ({ title, link, isoDate }: Entry): string => {
    const date = isoDate ? new Date(isoDate).toISOString().slice(0, 10) : ''
    return date ? `[${title}](${link}) - ${date}` : `[${title}](${link})`
}
```

Parsing the URL with `rss-parser` is quite simple:

```tsx
const parser = new rssParser()
const response = await parser.parseURL(url)
```

It returns a response in JSON format. Here's an example of the response:

```json
{
    "items": [
        {
            "title": "How I made my GitHub profile README dynamic",
            "link": "https://tduyng.github.io/blog/dynamic-github-profile-readme/",
            "pubDate": "2024-05-13T00:00:00.000Z",
            "author": "Z",
            "summary": "Explore the process of making your GitHub profile README dynamic with automated updates of your latest blog posts using GitHub Actions and Python scripting",
            "id": "https://tduyng.github.io/blog/dynamic-github-profile-readme/",
            "isoDate": "2024-05-13T00:00:00.000Z"
        },
        {
            "title": "New home for my website",
            "link": "https://tduyng.github.io/blog/new-home-for-my-website/",
            "pubDate": "2024-05-11T00:00:00.000Z",
            "author": "Z",
            "summary": "Discover why I switched my website and blog from Jekyll to Zola.",
            "id": "https://tduyng.github.io/blog/new-home-for-my-website/",
            "isoDate": "2024-05-11T00:00:00.000Z"
        },
        {
            "title": "Start a new journey",
            "link": "https://tduyng.github.io/blog/start-a-new-journey/",
            "pubDate": "2021-05-01T00:00:00.000Z",
            "author": "Z",
            "summary": "I share my journey from being a BIM engineer to becoming a full-time backend developer",
            "id": "https://tduyng.github.io/blog/start-a-new-journey/",
            "isoDate": "2021-05-01T00:00:00.000Z"
        }
    ],
    "link": "https://tduyng.github.io/atom.xml",
    "feedUrl": "https://tduyng.github.io/atom.xml",
    "title": "~/Z",
    "lastBuildDate": "2024-05-14T00:00:00+00:00"
}
```

This response is already sorted by the date of the post. With this response, I can easily filter and format it according to our requirements.

The remaining code reads the current [README.md](http://readme.md/) file, writes the feeds found above, and inserts them into comments in the README file.

```markdown
<!-- blog start -->
<!-- blog end -->
```

Here is the code:

```tsx
const replaceChunk = (
    content: string,
    marker: string,
    chunk: string,
    inline: boolean = false,
): string => {
    const startMarker = `<!-- ${marker} start -->`
    const endMarker = `<!-- ${marker} end -->`

    const pattern = new RegExp(`${startMarker}[\\s\\S]*${endMarker}`, 'g')

    if (!inline) {
        chunk = `\n${chunk}\n`
    }

    return content.replace(pattern, `${startMarker}${chunk}${endMarker}`)
}

const updateReadme = async (): Promise<void> => {
    const url = 'https://tduyng.github.io/atom.xml'
    const feeds = await fetchFeed(url)

    try {
        const readmePath = `${process.cwd()}/README.md`
        let readmeContent = await fs.readFile(readmePath, 'utf-8')
        readmeContent = replaceChunk(readmeContent, 'blog', feeds.join('\n\n'))
        await fs.writeFile(readmePath, readmeContent, 'utf-8')
        console.log('README.md updated successfully!')
    } catch (error) {
        console.error('Error updating README.md:', error)
    }
}

// Since I'm using ESM syntax, I can use the top-level "await" here.
await updateReadme()
```

Alright, the script is complete. You can find the full code [here](https://github.com/tduyng/tduyng/blob/bun/src/feed.ts).

**Next, let's move on to the final step: Configuring the GitHub action.**

It's essentially the same process as I did in [the previous article with Python](https://tduyng.github.io/blog/dynamic-github-profile-readme/). The only difference is replacing the Python implementation with Bun. Here's the code:

```yaml
name: Fetch latest posts from blog for README

on:
    push:
    workflow_dispatch:
    schedule:
        - cron: '0 0 * * *'

jobs:
    build:
        runs-on: ubuntu-latest
        steps:
            - uses: actions/checkout@v4
            - uses: oven-sh/setup-bun@v1
            - uses: actions/cache@v4
              name: Configure bun caching
              with:
                  path: ~/.bun/install/cache
                  key: ${{ runner.os }}-${{ matrix.bun }}-bun-${{ hashFiles('**/bun.lockb') }}
                  restore-keys: |
                      ${{ runner.os }}-${{ matrix.bun }}-bun-
            - run: bun install
            - run: bun run readme
            - name: Commit and push if changed
              run: |-
                  git diff
                  git config --global user.email "${{ vars.USER_EMAIL }}"
                  git config --global user.name "${{ vars.USER_NAME }}"
                  git add -A
                  git commit -m "chore: update blog posts" || exit 0
                  git push
```

In this scenario, I might not require the caching step for Bun.

```yaml
- uses: actions/cache@v4
      name: Configure bun caching
      with:
          path: ~/.bun/install/cache
          key: ${{ runner.os }}-${{ matrix.bun }}-bun-${{ hashFiles('**/bun.lockb') }}
          restore-keys: |
            ${{ runner.os }}-${{ matrix.bun }}-bun-
```

`bun install` is already sufficiently fast. I'm including this for best practice.

The GitHub Action will execute the command `bun run readme` to update the README and create a new commit to reflect the changes.

In summary, I've demonstrated how to enhance your GitHub profile using Bun and TypeScript. This straightforward [alternative to Python](https://tduyng.github.io/blog/dynamic-github-profile-readme/) provides a simple method for adding dynamism to your profile. You can find all the code from this article in [this branch](https://github.com/tduyng/tduyng/tree/bun)

For now, I've decided to keep the Python version on my main branch because it's easier to set up.

I hope it helps. Feel free to read the comments or provide feedback on [my medium post](https://tduyng.medium.com/dynamic-github-profile-with-bun-and-typescript-fe61127acabb).

Happy coding!
