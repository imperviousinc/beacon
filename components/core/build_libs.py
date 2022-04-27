#!/usr/bin/env/python3

# Copyright 2021 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

# See BUILD.gn in this directory for an explanation of what this script is for.

import argparse
import os
import stat
import sys
import shutil
import subprocess
import re

def main():
  parser = argparse.ArgumentParser("build_libs.py")
  parser.add_argument("--clang-base-dir",
                      help="Path to Clang binaries",
                      required=False),
  parser.add_argument("--output",
                      help="Path to library output",
                      required=True)
  parser.add_argument("--depfile", help="Path to write depfile", required=False)

  args = parser.parse_args()

  run_build(args)
      
  # Generate dep file  
  if not args.depfile:
      return

  # root_project/src/beacon/core   
  beacon_libs_dir = os.path.dirname(os.path.realpath(__file__))

  # Ninja depfile Format:
  # <output path to lib>: <path to each Go file>
  with open(args.depfile, 'w') as depfile:
    depfile.write("%s:" % args.output)

    for root, dirs, files in os.walk(beacon_libs_dir):    
      for f in files: 
        if f.endswith('.go'):
          infile = os.path.join(root, f)
          depfile.write(" %s" % infile)               
    depfile.write("\n")

def create_build_env(clang_base_dir):
  env = os.environ.copy()
  if not clang_base_dir or sys.platform == 'win32':
    return env

  clang_bin_dir = os.path.join(clang_base_dir, "bin")
  if "PATH" in env:
    env["PATH"] = clang_bin_dir + ":" + env["PATH"]
  else:
    env["PATH"] = clang_bin_dir

  env["CC"] = os.path.join(clang_bin_dir, "clang")
  return env 

# Removes dylibs so darwin target doesn't link
# against them favoring static libs
def remove_dylibs(d):
  if sys.platform != 'darwin':
    return
  l = os.listdir(d)
  for item in l:
    if item.endswith(".dylib"):
      os.remove(os.path.join(d, item))

def must_exist(p):
  if not os.path.exists(p):
      raise Exception("path " + p + " doesn't exist")

def to_msys_path(p):
   if os.name != 'nt':
       return p 

   return p.replace("C:", "/c").replace("\\", "/")

def call(command, env):
  # On Windows, we need to compile hnsd using MSYS64  
  if sys.platform == 'win32':
    env["MSYSTEM"] = "MINGW64"
    cwd = '"{0}"'.format(to_msys_path(os.getcwd()))
    command = [r"C:\msys64\usr\bin\bash", "-l" , "-c", "cd " + cwd + " && " + ' '.join(command)]

  try:
    subprocess.check_call(command,env=env)
  except subprocess.CalledProcessError as e:
    print(e.output)
    raise e

# Builds this project must be called after build_hsk
def build_core(out, env):
  goBin = 'go'

  # We use MSYS on windows to compile CGO code
  # Note: TDM GCC may still be needed as a depedency.
  if sys.platform == 'win32':
    goBin = '"/c/Program Files/Go/bin/go.exe"'
    out = to_msys_path(out)

  call_args = [goBin,'build','-trimpath','-buildmode=c-shared', '-o', out]
   
  if sys.platform == 'win32':
    call(call_args, env)
    return
    
  # keep in sync with the same min mac osx version used
  # by chromium
  if sys.platform == 'darwin':
    env['CGO_CFLAGS'] = '-mmacosx-version-min=10.11.0'
    env['CGO_LDFLAGS'] = '-mmacos-version-min=10.11.0' 
  try:
    subprocess.check_call(call_args, env=env)
  except subprocess.CalledProcessError as e:
    print(e.output)
    raise e

  # call install tool to update loader path
  # this way it can be found by chromium app bundle
  # we currently dlopen the lib anyways so this
  # isn't necessary but my still be useful.
  if sys.platform == 'darwin':
    call(['install_name_tool','-id', '@loader_path/Libraries/libbeacon.dylib', out], env)      

# Builds libhsk must be called from hnsd root directory
def build_hsk(hsk_out_path, env):    
  # start a clean build
  try: 
    call(['make', 'clean'], env)
  except Exception:
    pass

  call(['./autogen.sh'], env)

  hsk_out_path = to_msys_path(hsk_out_path)
  call(['./configure', '--without-daemon', '--prefix', hsk_out_path],env)

  if sys.platform == 'darwin':
    call(['make', 'CFLAGS=-mmacosx-version-min=10.11.0', '-j'], env)
  elif sys.platform == 'linux' or sys.platform == 'linux2':
    call(['make', 'CFLAGS=-fPIC', '-j'], env)
  else:
    call(['make', '-j'], env)
  
  # installs to hsk_out_path
  call(['make', 'install'], env) 

  # remove dylibs
  remove_dylibs(os.path.join(hsk_out_path, "lib"))

def run_build(args):
  env = create_build_env(args.clang_base_dir)

  # root_project/src/beacon/core
  #+root_project/src/beacon/components/core
  this_dir = os.path.dirname(os.path.realpath(__file__))
  
  # root_project/src/beacon
  beacon_dir = os.path.dirname(os.path.dirname(this_dir)) 
  third_party_dir = os.path.join(beacon_dir, "third_party")

  libhsk_src = os.path.join(beacon_dir, "third_party", "hnsd")
  hsq_src = os.path.join(beacon_dir, "third_party", "hnsquery")
  libhsk_out = os.path.join(hsq_src, "build")

  # Sanity check
  must_exist(hsq_src)
  must_exist(libhsk_src)

  # Starting libhsk build
  # don't build if lib already exists
  if not os.path.exists(os.path.join(libhsk_out, "lib")):
    os.chdir(libhsk_src)
    build_hsk(libhsk_out, env)

  # Building this project  
  os.chdir(this_dir)
  build_core(args.output, env)
 

if __name__ == '__main__':
  sys.exit(main())
