#!/bin/sh

## CHANGE THESE VALUES based on your application
APP_DESCRIPTION="my test application"       # brief description of the application
APP_URL="YOUR APP.azureiotcentral.com"      # application url
ID_SCOPE="YOUR SCOPE ID HERE"               # application DPS ID scope for the application (get it from admin/device connection)
MASTER_KEY="YOUR MASTER KEY"                # application master key (get it from Admin/Device Connection)
API_TOKEN="YOUR API TOKEN"                  # API token with admin role (get it from Admin/API Token)

## Add an IoT Central application in which the simulated devices will reside.
# You can add multiple applications similar to the one below.
curl --location --request PUT 'http://localhost:6001/api/target' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": "app1",
    "name": "'"$APP_DESCRIPTION"'",
    "provisioningUrl": "global.azure-devices-provisioning.net",
    "idScope": "'"$ID_SCOPE"'",
    "masterKey": "'"$MASTER_KEY"'",
    "appUrl": "'"$APP_URL"'",
    "appToken": "'"$API_TOKEN"'"
}'

## Add a device model.
# These device models are used for simulation.
curl --location --request PUT 'http://localhost:6001/api/model' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": "brewer",
    "name": "brewer",
    "capabilityModel":
[
  {
    "@id": "dtmi:CoffeeCo:Brewer;1",
    "@type": "Interface",
    "contents": [
      {
        "@id": "dtmi:CoffeeCo:Brewer:Temperature;1",
        "@type": "Telemetry",
        "displayName": {
          "en": "Temperature"
        },
        "name": "Temperature",
        "schema": "double"
      },
      {
        "@id": "dtmi:CoffeeCo:Brewer:StartBrew;1",
        "@type": "Command",
        "commandType": "synchronous",
        "displayName": {
          "en": "StartBrew"
        },
        "name": "StartBrew"
      },
      {
        "@id": "dtmi:CoffeeCo:Brewer:Wakeup;1",
        "@type": "Command",
        "displayName": {
          "en": "Wakeup"
        },
        "name": "Wakeup"
      },
      {
        "@id": "dtmi:CoffeeCo:Brewer:SerialNumber;1",
        "@type": "Property",
        "displayName": {
          "en": "SerialNumber"
        },
        "name": "SerialNumber",
        "schema": "string",
        "writable": false
      },
      {
        "@id": "dtmi:CoffeeCo:Brewer:Model;1",
        "@type": "Property",
        "displayName": {
          "en": "Model"
        },
        "name": "Model",
        "schema": "string",
        "writable": false
      },
      {
        "@id": "dtmi:CoffeeCo:Brewer:WakeupTime;1",
        "@type": "Property",
        "displayName": {
          "en": "WakeupTime"
        },
        "name": "WakeupTime",
        "schema": "time",
        "writable": true
      }
    ],
    "displayName": {
      "en": "Brewer"
    },
    "extends": [
      "dtmi:CoffeeCo:BrewerBoiler;1",
      "dtmi:CoffeeCo:BrewerDisplay;1"
    ],
    "@context": [
      "dtmi:iotcentral:context;2",
      "dtmi:dtdl:context;2"
    ]
  },
  {
    "@context": [
      "dtmi:iotcentral:context;2",
      "dtmi:dtdl:context;2"
    ],
    "@id": "dtmi:CoffeeCo:BrewerBoiler;1",
    "@type": "Interface",
    "contents": [
      {
        "@id": "dtmi:CoffeeCo:BrewerBoiler:Pressure;1",
        "@type": [
          "Telemetry",
          "Pressure"
        ],
        "displayName": {
          "en": "Pressure"
        },
        "name": "Pressure",
        "schema": "double"
      },
      {
        "@id": "dtmi:CoffeeCo:BrewerBoiler:BoilerTemperature;1",
        "@type": [
          "Telemetry",
          "Temperature"
        ],
        "displayName": {
          "en": "BoilerTemperature"
        },
        "name": "BoilerTemperature",
        "schema": "double"
      },
      {
        "@id": "dtmi:CoffeeCo:BrewerBoiler:SteamFlowRate;1",
        "@type": [
          "Telemetry",
          "Velocity"
        ],
        "displayName": {
          "en": "SteamFlowRate"
        },
        "name": "SteamFlowRate",
        "schema": "double"
      },
      {
        "@id": "dtmi:CoffeeCo:BrewerBoiler:Clean;1",
        "@type": "Command",
        "displayName": {
          "en": "Clean"
        },
        "name": "Clean"
      },
      {
        "@id": "dtmi:CoffeeCo:BrewerBoiler:ManufacturedDate;1",
        "@type": "Property",
        "displayName": {
          "en": "ManufacturedDate"
        },
        "name": "ManufacturedDate",
        "schema": "date",
        "writable": false
      },
      {
        "@id": "dtmi:CoffeeCo:BrewerBoiler:NextCleaningDate;1",
        "@type": "Property",
        "displayName": {
          "en": "NextCleaningDate"
        },
        "name": "NextCleaningDate",
        "schema": "dateTime",
        "writable": true
      },
      {
        "@id": "dtmi:CoffeeCo:BrewerBoiler:CleaningCycleFrequency;1",
        "@type": "Property",
        "displayName": {
          "en": "CleaningCycleFrequency"
        },
        "name": "CleaningCycleFrequency",
        "schema": "integer",
        "writable": true
      }
    ],
    "displayName": {
      "en": "Boiler"
    }
  },
  {
    "@context": [
      "dtmi:iotcentral:context;2",
      "dtmi:dtdl:context;2"
    ],
    "@id": "dtmi:CoffeeCo:BrewerDisplay;1",
    "@type": "Interface",
    "contents": [
      {
        "@id": "dtmi:CoffeeCo:BrewerDisplay:BrightnessLevel;1",
        "@type": "Telemetry",
        "displayName": {
          "en": "BrightnessLevel"
        },
        "name": "BrightnessLevel",
        "schema": "double"
      },
      {
        "@id": "dtmi:CoffeeCo:BrewerDisplay:Message;1",
        "@type": "Property",
        "displayName": {
          "en": "Message"
        },
        "name": "Message",
        "schema": "string",
        "writable": true
      },
      {
        "@id": "dtmi:CoffeeCo:BrewerDisplay:ImageLocation;1",
        "@type": "Property",
        "displayName": {
          "en": "ImageLocation"
        },
        "name": "ImageLocation",
        "schema": "string",
        "writable": true
      },
      {
        "@id": "dtmi:CoffeeCo:BrewerDisplay:Restart;1",
        "@type": "Command",
        "commandType": "synchronous",
        "displayName": {
          "en": "Restart"
        },
        "name": "Restart"
      }
    ],
    "displayName": {
      "en": "Display"
    }
  }
]
}'

## Add the device model to the target application
# You can add multiple device models in the same models array, after the "brewer".
curl --location --request PUT 'http://localhost:6001/api/target/app1/models' \
--header 'Content-Type: application/json' \
--data-raw '[
    "brewer"
]
'

## Add a simulation.
# This simulation is configured to distribute the devices into two wave groups.
# Each wave group will send data 5 seconds apart.
# Telemetry is sent once in a minute and each time it sends a batch of 2 messages.
# Reported properties are sent once an hour.
# Devices are never disconnected. To simulate an occasionally connected device,
# you can change the disconnectBehavior to 'telemetry' to disconnect the device after sending telemetry.
# A telemetryFormat 'default' sends typical json messages. If it is set to 'opcua', json messages with opcua envelopes will be sent.
curl --location --request PUT 'http://localhost:6001/api/simulation' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": "sim1",
    "name": "sim1",
    "targetId": "app1",
    "status": "created",
    "waveGroupCount": 2,
    "waveGroupInterval": 5,
    "telemetryBatchSize": 1,
    "telemetryInterval": 60,
    "reportedPropertyInterval": 3600,
    "disconnectBehavior": "never",
    "telemetryFormat": "default"
}'

## Configure number of devices for the simulation.
# You can have multiple devices per simulation by adding a new device config just like the one below.
curl --location --request PUT 'http://localhost:6001/api/simulation/sim1/deviceConfig' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": "brewer",
    "modelId": "brewer",
    "deviceCount": 10
}'
