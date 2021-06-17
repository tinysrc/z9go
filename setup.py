#!/usr/bin/python
# -*- coding:utf-8 -*-

import os, shlex, subprocess

work_dir = os.path.dirname(os.path.realpath(__file__))

list = [
    "tools/z9",
    "tools/z9protoc"
]
for entry in list:
    cmd = "go install"
    cwd = os.path.join(work_dir, entry)
    print("install " + entry)
    p = subprocess.Popen(shlex.split(cmd), cwd=cwd)
    p.wait()
