#!/bin/sh

# Stop the simulation 'sim1'
curl --location --request POST 'http://localhost:6001/api/simulation/sim1/stop'
