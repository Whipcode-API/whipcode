REST service for executing untrusted code with Podman.

Implemented as a scalable stateless microservice with
no user management or authentication, minimizing damage
from potential zero-day breakouts.

See https://whipcode.app/#docs for API documentation.


Supported languages:
    Python                 1
    Javascript (Node.js)   2
    Bash                   3
    Perl                   4
    Lua	                   5
    Ruby                   6
    C                      7
    C++                    8
    Rust                   9
    Fortran               10
    Haskell               11
    Java                  12
    Go                    13
    TypeScript            14
    Common Lisp           15
    Racket                16
    C#                    17


Please set this up on a SELinux-enabled system.
Tested on Fedora Server 41.


Environment setup:
    Go:
        $ wget https://go.dev/dl/go1.<version>.linux-<arch>.tar.gz
        $ sudo rm -rf /usr/local/go
        $ sudo tar -C /usr/local -xzf go1.<version>.linux-<arch>.tar.gz
        $ sudo echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
        $ source /etc/profile

    go-task:
        $ wget https://github.com/go-task/task/releases/latest/download/task_linux_<arch>.rpm
        $ sudo dnf install task_linux_<arch>.rpm

    Podman:
        $ sudo dnf install podman
        $ sudo echo "$USER:100000:65536" | sudo tee /etc/subuid /etc/subgid
        $ podman system reset

        Ensure that `podman info | grep graphDriverName` returns `overlay`.

    SELinux:
        Ensure that `sudo getenforce` returns `Enforcing`.

        $ sudo dnf install container-selinux
        $ sudo semodule -i selinux/whipcode.cil selinux/base_container.cil


Build:
    To build everything:
        $ task all

    To build only the service:
        $ task build

    Or just the container images:
        $ task build-images

    To rebuild images at a later time:
        $ task rebuild-images

    Update (git pull), build service and images:
        $ task update


Start the service:

    # # # # ! WARNING ! WARNING ! WARNING ! WARNING ! WARNING ! # # # #
    #                                                                 #
    #    DO NOT run this service without a reverse proxy or API       #
    #    gateway in front of it. While whipcode does have a           #
    #    `--standalone` mode for per IP rate limiting, it is          #
    #    not meant to be used in production. Use an API gateway       #
    #    like Kong, Tyk and WSO2 to enforce rate limits, policies,    #
    #    and authentication. Configure your gateway to add a          #
    #    `X-Master-Key` header to every request with the secret       #
    #    defined below. DO NOT host the gateway on the same system.   #
    #                                                                 #
    # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

    DO NOT run whipcode as root, or with SELinux disabled/permissive.

    Save your master key's hash to .masterkey:
        $ export HISTFILE=/dev/null
        $ task key -- <KEY>

    Copy the configuration template:
        $ cp config.toml.example config.toml

    Default port 8000:
        $ task run

    Port 6060 with /ping enabled:
        $ task run -- --ping --port 6060

    The endpoint will be available at /run.

    Test the service:
        $ task test

        If you see 16 responses with "Success!" in the `stdout` field,
        the service is working correctly.


CLI options:
    -p, --port  PORT
        The port to listen on. May not always work with authbind
        when attempting to bind to ports < 1024. (default: 8000)

    -m, --max  BYTES
        The maximum size of the request body in bytes. Requests
        larger than this will be rejected. (default: 1000000)

    -t, --timeout  SECONDS
        The maximum time allowed for code execution. Should be
        set lower than the server's write timeout, which is 20
        seconds. (default: 10)

    -k, --key  FILE
        Path to the file containing the master key's argon2 hash
        and salt. (default: .masterkey)

    --proxy  ADDR
        The address of the reverse proxy or API gateway in front
        of whipcode. Requests not originating from this address
        will be rejected. (default: none)

    --cache
        Enables an LRU cache for code executions. This will speed
        up responses for repeated requests. (default: false)

        Note: The cache is not persistent and will be lost on
        restart. While this feature is intended to reduce server
        load and latency, in some situations it may end up
        worsening it. Memory usage will also increase.

    --tls
        Enables TLS. Requires tls/cert.pem and tls/key.pem to be
        present.

    --ping
        Enables the /ping endpoint. Replies with "pong".

    --standalone
        Enables per IP rate limiting, without the need for a
        reverse proxy or API gateway. This is NOT RECOMMENDED in
        production. (default: false)

    --burst  COUNT     (Requires --standalone)
        The number of requests allowed in a burst. (default: 3)

    --refill  SECONDS  (Requires --standalone)
        The number of seconds for each request to refill in the
        burst bucket. (default: 1)


License:
    Copyright 2024 whipcode.app (AnnikaV9)

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

            http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on
    an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
    either express or implied. See the License for the specific
    language governing permissions and limitations under the License.
