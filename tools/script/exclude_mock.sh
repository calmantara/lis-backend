#!/bin/sh
cat coverage.txt | grep -v mock > coverage.final.txt
mv coverage.final.txt coverage.txt