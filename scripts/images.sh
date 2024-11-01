#!/usr/bin/env bash
#
#  Copyright 2024 whipcode.app (AnnikaV9)
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#          http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing,
#  software distributed under the License is distributed on
#  an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
#  either express or implied. See the License for the specific
#  language governing permissions and limitations under the License.
#

if ! command -v podman &> /dev/null
then
    echo "Podman is not installed."
    exit 1
fi

if [ ! -d scripts/extra_setup ]; then
    echo "Run this script from the root of the repository."
    exit 1
fi

declare -A langs=(
    # [<lang>]='<apk_deps>'
    # Extra setup should be done in scripts/extra_setup/<lang>.sh
    [bash]='bash'
    [nodejs]='nodejs'
    [c]='gcc'
    [cpp]='g++'
    [fortran]='gfortran'
    [go]='gcc-go'
    [haskell]='ghc'
    [java]='openjdk21'
    [lue]='lua'
    [perl]='perl'
    [python]='python3'
    [ruby]='ruby'
    [rust]='rust'
    [typescript]='npm'
    [list]='sbcl'
    [racket]='racket'
    [csharp]='dotnet8-sdk'
)

HEADER="FROM docker.io/alpine:latest"
PREFIX="RUN apk update --no-cache && apk upgrade --no-cache && apk add --no-cache libc-dev musl-dev "
SUFFIX="apk --purge del apk-tools && rm -rf /var/cache/apk /var/lib/apk /lib/apk /etc/apk /sbin/apk /usr/share/apk /usr/lib/apk /usr/sbin/apk /usr/local/apk /usr/bin/apk /usr/local/bin/apk /usr/local/sbin/apk /usr/local/lib/apk /usr/local/share/apk /usr/local/libexec/apk /usr/local/etc/apk"

trap "rm -f TEMP_CONTAINERFILE; exit" INT

declare -i i=0
for lang in "${!langs[@]}"; do
    i+=1
    header=${HEADER}
    suffix=${SUFFIX}
    prefix=${PREFIX}
    if [ -f scripts/extra_setup/${lang,,}.sh ]; then
        header+="\nCOPY scripts/extra_setup/${lang,,}.sh /tmp/setup.sh\n"
        suffix="sh /tmp/setup.sh && rm -f /tmp/setup.sh && ${suffix}"
    fi
    langs["$lang"]=$(echo -e "${header}\n${prefix}${langs[$lang]} && ${suffix}")
    echo "${langs[$lang]}" > TEMP_CONTAINERFILE
    podman build -t whipcode-${lang,,} -f TEMP_CONTAINERFILE . | while read line; do echo "[$i/16] [${lang,,}] $line"; done
done

rm -f TEMP_CONTAINERFILE
