#!/bin/bash

#######################################################
# Replace with valid values.
# This file should be placed in your filesystem as:
# /usr/local/opt/rhbuilder/.rhbuilder
# ENSURE THAT THE USER EXECUTING THE LAUNCH
# SCRIPT HAS PERMISSIONS TO EXECTUTE THIS FILE.
#######################################################

RHBUILDER_AMQP_URL=<RHBUILDER_AMQP_URL>
RHBUILDER_AMQP_EXCHANGE=<RHBUILDER_AMQP_EXCHANGE>
RHBUILDER_AMQP_TOPIC=rabbithook.dockerhub
RHBUILDER_JOBS_PATH=/usr/local/opt/rhbuilder
