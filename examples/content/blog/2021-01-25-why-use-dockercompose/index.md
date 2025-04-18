+++
title = "Why use Docker compose"
description = "Discover the advantages and reasons behind utilizing Docker Compose for managing multi-container Docker applications."
template = "post.html"
date = 2021-01-25
tags = ["docker", "compose"]
generate_feed = false

[extra]
comment = false
reaction = true
toc = true
copy = true
img = "/img/post-cover.webp"
+++

## What is **Docker compose**?

If you read this article, I think you maybe already understand what is [Docker](https://www.docker.com/why-docker).

Docker is a tool which helps to create, deploy, and run applications by using containers. Docker provides developers and operators with a friendly interface to build, ship, and run containers on any environment.

With the docker-cli, we can run and execute the Docker container via [Dockerfile](https://docs.docker.com/engine/reference/builder/#:~:text=A%20Dockerfile%20is%20a%20text,can%20use%20in%20a%20Dockerfile%20.). But with Dockerfile, we can work only with a container, but in the project, we normally work with multiple containers. In that case, we need to use **Docker compose**.

**Docker Compose** define the services that make up your app in `docker-compose.yml` so they can be run together in an isolated environment. It get an app running in one command by just running `docker-compose up`.

**Docker compose** uses the `Dockerfile` if one add the build command to your project’s `docker-compose.yml`. Your Docker workflow should be to build a suitable `Dockerfile` for each image you wish to create, then use compose to assemble the images using the build command.

We will create a demo: simple app nodejs using Redis. You will see the problem with only **Dockerfile** and why use **Docker compose**.

## Project demo

### Installation

Make sure you have docker & docker-compose installed in your machine:

- [How to install docker](https://docs.docker.com/get-docker/)
- [How to install docker-compose](https://docs.docker.com/compose/install/)

### Create Docker Images without Docker compose

- Create project folder

    ```bash
    mkdir node-redis-demo
    cd node-redis-demo
    ```

- Create `package.json` file

    In this file, we will declare the dependencies we use for our application. This step is similar as you create the Nodejs application (without **npm install** or **yarn add** dependencies).

    We consider our computer don't have **npm** or **yarn** installed. We just use Docker to create the develop environment.

    ```bash
    touch package.json # or yarn init -y
    ```

    Paste the following code:

    ```javascript
    //package.json
      {
      "name": "node-redis-app",
      "version": "1.0.0",
      "main": "index.js",
      "license": "MIT",
      "scripts": {
        "start": "node index.js",
      },
      "dependencies": {
        "redis": "*", //We don't care about the version of dependencies
        "express": "*", //Docker will fetch the latest version of them
      }
    }
    ```

- Create `index.js` file to config `express` and run server

    ```javascript
    //index.js
    const express = require('express');
    const redis = require('redis');
    const app = express();

    const client = redis.createClient();
    client.on('error', function (error) {
      console.error(error);
    });
    client.set('visits', 0);

    app.get('/', (req, res) => {
      client.get('visits', (err, visits) => {
        res.send('Number of visits is ' + visits);
        client.set('visits', parseInt(visits) + 1);
      });
    });

    const port = 8080;
    app.listen(port, () => {
      console.log('App Redis is running on port 8080');
    });

    const
    ```

- Create **Dockerfile**

    ```yml
    # Dockerfile
    FROM node:12.18-alpine #Version nodejs we want to use

    WORKDIR /app #Define the work directory of project in Docker container

    COPY package*.json ./ #Make sure always copy package.json and package-lock.json to build folder
    # Copy package.json to avoid reinstall all the dependencies for each build

    RUN yarn install #Run install and update dependencies from package.json

    COPY . . #Copy all other files to the build folder

    EXPOSE 8080 #Optional: Port running on Docker container
    CMD ["node", "index.js"] #Command line to run project

    ```

    You can also check an article on my blog for more information of Dockerfile: [Create simple project Nodejs with Docker](https://tduyng.github.io/blog/create-simple-project-nodejs-with-docker)

- Create `.dockerignore`

    ```yml
    node_modules
    npm-debug.log
    yarn-debug.log
    ```

- Build Docker Image
  OK, everything is done. Now, we will build this image:

    ```bash
    docker build -t node-redis-app .
    ```

    Don't forget the point **.** in the end.

    You can check all images available in your computer:

    ```bash
    docker images
    ```

- Run Docker container from Docker image
  We have already an Docker image node-redis-app. The thing to do now is run container from this image to start project.

    ```bash
    docker run -p 8080:8080 -it node-redis-app
    ```

    **-p 8080:8080** flag: declare port using (port on Docker container and port on your machine)
    **-it** flag: allow access inside the terminal of Docker container
    **node-redis-app**: name image of your app

- Problem
  When you run with the command above, you will see an error failed connection with Redis-server:

    ```bash
    $ docker run -p 8080:8080 -it node-redis-app

    App Redis is running on port 8080
    Error: Redis connection to 127.0.0.1:6379 failed - connect ECONNREFUSED 127.0.0.1:6379
        at TCPConnectWrap.afterConnect [as oncomplete] (net.js:1141:16) {
      errno: 'ECONNREFUSED',
      code: 'ECONNREFUSED',
      syscall: 'connect',
      address: '127.0.0.1',
      port: 6379
    }

    ```

    **What is the problem ?**
    The problem is: You want to run 2 containers in same times.
    ==> **redis-server** && **nodejs**
    But when you run your Dockerfile, you run only Nodejs app container, not to redis container.
    You can also create other redis-server container from [redis](https://hub.docker.com/_/redis) image

    There are no any link between 2 this containers. Your Nodejs application can't connect to **redis-server**(from **redis** Docker image).

    **How to resolve this error?**
    The answer is **Docker compose**. In this case, we need to use Docker compose to run and connect multiple Docker container.

### Create Docker image using Docker compose

- To make sure the project above running, we need simply create a `docker-compose.yml` file.
  It is a config file for docker-compose. It allows us to define the services (container) using in our app.

    ```bash
    touch docker-compose.yml
    ```

- Add the following code:

    ```yml
    version: '3.8' #version of docker-compose
    services: # declare containers used
        redis-server: #image of redis
            image: 'redis'
        node-app: #images of your node app (using Dockerfile, package.json ...)
            restart: always #If this container stops for any reason, always attempt to restart it
            build: .
            ports:
                - '8080:8080'
    ```

    Check [docker compose](https://docs.docker.com/compose/gettingstarted/) for more details about configuration.

- Update redis services to `index.js` file

    ```javascript
    // index.js
    const express = require('express')
    const redis = require('redis')
    const app = express()
    const process = require('process')

    const client = redis.createClient({
        host: 'redis-server',
        port: '6379',
    })
    client.on('error', function (error) {
        console.error(error)
    })
    client.set('visits', 0)

    app.get('/', (req, res) => {
        process.exit(0)
        client.get('visits', (err, visits) => {
            res.send('Number of visits is ' + visits)
            client.set('visits', parseInt(visits) + 1)
        })
    })

    const port = 8080
    app.listen(port, () => {
        console.log('App Redis is running on port 8080')
    })
    ```

- Build and run your app with docker-compose

    In the root project:

  - Build images:

                            ```bash
                            docker-compose build
                            ```

                            When you use docker-compose to build the image, the image name is always going to be `<project>_<service>`, where `<service>` in this example is node-redis-app_node-app.

  - Run containers with docker-compose

                        ```bash
                        docker-compose up # all services declared in docker-compose.yml will be started
                        docker-compose up -d #detach mode
                        ```

  - Stop containers with docker-compose

                        ```bash
                        docker-compose down # Stop all containers declared in docker-compose.yml file
                        ```

That's it. Test the code on your machine and check it in then [localhost8080](http://localhost:8080) to see the result.

## Resume

Docker compose:

- Separate CLI that gets installed along with Docker
- Used to startup multiple Docker containers at the same time
- Automate some of the long-winded arguments we were passing to `docker run`
- Docker compose will start and connect all the containers Docker that you declare in the `docker-compose.yml` file. We will use it usually when we work with Docker.

## Reference

- [Code demo GitHub](https://github.com/tduyng/try-docker/tree/master/01.dive-into-docker/simple-with-redis)
- [Previous post: create simple Nodejs app with Docker](https://tduyng.github.io/blog/create-simple-project-nodejs-with-docker)
- [Docker official](https://docs.docker.com/)
- [Learn about Docker compose](https://docs.docker.com/compose/)
- [Docker compose cheatsheet](https://jstobigdata.com/docker-compose-cheatsheet/)
