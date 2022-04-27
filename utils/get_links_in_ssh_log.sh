#! /bin/bash

cat ../ssh.log | grep -oiahE "https?://[^\"\\'> ]+" | tr -d ';' | sort | uniq
