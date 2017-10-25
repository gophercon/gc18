#!/bin/bash
set -e

RG=gophercon
LOC=westus2

az group create --name $RG --location $LOC
az aks create --resource-group $RG --name gopherconk8s --agent-count 2 --generate-ssh-keys
