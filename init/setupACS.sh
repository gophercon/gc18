#!/bin/bash
set -e

RG=gophercon
LOC=southcentralus


az group create --name $RG --location $LOC

az acs create --orchestrator-type kubernetes --resource-group $RG --name $RG --generate-ssh-keys
