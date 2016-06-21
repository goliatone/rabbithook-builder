# rhbuilder Upstart and SysVinit configuration file

#######################################################
# Replace with valid values.
# This file should be placed in your filesystem as:
# /usr/local/opt/rhbuilder/.rhbuilder
# ENSURE THAT THE USER EXECUTING THE LAUNCH
# SCRIPT HAS PERMISSIONS TO EXECTUTE THIS FILE.
#######################################################

# Customize location of rhbuilder binary (especially for development testing).
#RHBUILDER="/usr/local/bin/rhbuilder"

# Required AMQP URL
RHBUILDER_AMQP_URL=amqp://guest:guest@localhost:5672/

# AMQP Exchange
RHBUILDER_AMQP_EXCHANGE=rabbithook

# AMPQ Topic, by default one of rabbithook.[dockerhub|github|travis]
RHBUILDER_AMQP_TOPIC=rabbithook.dockerhub

# Path where "rhbuilder" looks for job files
RHBUILDER_JOBS_PATH=/usr/local/opt/rhbuilder
