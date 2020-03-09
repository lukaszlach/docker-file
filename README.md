# Docker File

![Version](https://img.shields.io/badge/version-1.0-lightgrey.svg?style=flat)

**Docker File** allows to validate `Dockerfile`, `docker-compose.yml`, `.env` and `.dockerignore` files and returns their contents in JSON format for further analysis. **Docker File** is distributed as a standalone binary and therefore does not need Docker or any other container runtime installed or running.

## Building

```bash
git clone https://github.com/lukaszlach/docker-file.git
cd docker-file
make clean-build
```

## Running

```
$ docker-file -help

NAME:
   docker-file - Tools to handling Docker-related files

USAGE:
   docker-file [global options] command [command options] [arguments...]

COMMANDS:
   compose, c       Docker Compose file handlers
   dockerfile, d    Dockerfile handlers
   env, e           .env file handlers
   dockerignore, i  .dockerignore file handlers
   help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)

COPYRIGHT:
   2020 Łukasz Lach https://lach.dev
```

### `dockerfile parse FILE`

Parses the `Dockerfile` using [BuildKit](https://github.com/moby/buildkit) parser and converts it to JSON format with [dockerfile-json](https://github.com/keilerkonzept/dockerfile-json).

```
$ docker-file dockerfile parse ./Dockerfile
{"MetaArgs":[{"Key":"B","DefaultValue":"999","ProvidedValue":null,"Value":"999"}],"Stages":[{"Name":"test1","BaseName":"nginx:${B}","SourceCode":"FROM nginx:${B} AS test1","Platform":"","As":"test1","From":{"Image":"nginx:${B}"},"Commands":[{"CmdLine":["apt-get update \u0026\u0026     apt-get install -y jq"],"Name":"run","PrependShell":true},{"Name":"workdir","Path":"/app"}]},{"Name":"test2","BaseName":"nginx:latest","SourceCode":"FROM nginx:latest AS test2","Platform":"","As":"test2","From":{"Image":"nginx:latest"},"Commands":[{"Chown":"","From":"","Name":"copy","SourcesAndDest":[".","."]},{"Chown":"","Name":"add","SourcesAndDest":["a1","a2"]},{"CmdLine":["/app/binary","run"],"Name":"cmd","PrependShell":false},{"CmdLine":["go test ."],"Name":"run","PrependShell":true},{"Env":[{"Key":"A","Value":"1234"}],"Name":"env"},{"Env":[{"Key":"A","Value":"123"}],"Name":"env"},{"Key":"B","Name":"arg","Value":"2"},{"Name":"volume","Volumes":["/dir1"]},{"Name":"volume","Volumes":["/dir2"]}]}]}
```

When BuildKit fails to parse the file the actual error message is returned with a non-zero exit code.

```
$ docker-file dockerfile parse ./Dockerfile.invalid
Dockerfile parse error line 2: COPY requires at least two arguments, but only one was provided. Destination could not be determined.
```

### `compose parse FILE`

Parses the docker-compose YAML file with docker-cli Compose Loader used when triggering `docker stack deploy` command. 
Returns the file contents in JSON format.

```
$ docker-file compose parse ./docker-compose.yml
{"services":{"web":{"image":"nginx:latest"}},"version":"3"}
```

When docker-cli fails to parse the file the actual error message is returned with a non-zero exit code.

```
$ docker-file compose parse ./docker-compose-invalid.yml
services.web.ports must be a list
```

### `env parse FILE`

Parses the `.env` file with docker-cli, return the contents in JSON format.

```
$ docker-file env parse ./.env
{"A":"1","B":"2","C":"123"}
```

When docker-cli fails to parse the file the actual error message is returned with a non-zero exit code.

```
$ docker-file env parse ./.env-invalid
poorly formatted environment: no variable name on line '=3'
```

### `dockerignore parse FILE`

Parses the `.dockerignore` file with Docker Builder parser used when triggering `docker build` command. 
Returns the file contents in JSON format.

```
$ docker-file dockerignore parse ./.dockerignore
[".idea",".git","*.md"]
```

### `dockerignore check IGNORE_FILE FILE`

Checks if the file matches rules from pointed `.dockerignore` file.
Returns `true` or `false`.

```
$ docker-file dockerignore check ./.dockerignore README.md
true
```

## License

MIT License

Copyright (c) 2020 Łukasz Lach <llach@llach.pl>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
