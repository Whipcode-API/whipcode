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

read -p "Enter master key: " master_key
read -p "Enter port: " port

declare -A langs=(
    [1]='print("Success!")'
    [2]='console.log("Success!");'
    [3]='echo "Success!"'
    [4]='print "Success!";'
    [5]='print("Success!")'
    [6]='puts "Success!"'
    [7]='#include <stdio.h>\nint main() { printf("Success!"); return 0;}'
    [8]='#include <iostream>\nint main() { std::cout << "Success!"; return 0; } '
    [9]='fn main() { println!("Success!"); }'
    [10]='program hello\n  print *, "Success!"\nend program hello'
    [11]='main = putStrLn "Success!"'
    [12]='public class HelloWorld { public static void main(String[] args) { System.out.println("Success!"); } }'
    [13]='package main; import "fmt"; func main() { fmt.Println("Success!") }'
    [14]='let message: string = "Success!";console.log(message);'
    [15]='(write-line "Success!")'
    [16]='#lang racket\n"Success!"'
    [17]='Console.WriteLine("Success!");'
)

for i in "${!langs[@]}"; do
    code=$(echo -en "${langs[$i]}" | base64 | tr -d '\n')
    result=$(curl -s -X POST -H "Content-Type: application/json" -H "X-Master-Key: $master_key" -d '{"language_id":"'"$i"'","code":"'"$code"'"}' "http://0.0.0.0:$port/run")
    echo "$i $result"
done
