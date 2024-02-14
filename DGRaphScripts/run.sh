#!/bin/bash

# dgraph zero  & dgraph alpha 
./dgraph zero --config ./zero_auth.json & 
./dgraph alpha --config ./alpha_auth.json &
./ratel &